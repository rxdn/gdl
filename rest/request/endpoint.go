package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	BASE_URL             = "https://discord.com/api/v10"
	AuditLogReasonHeader = "X-Audit-Log-Reason"
)

type Endpoint struct {
	RequestType       RequestType
	ContentType       ContentType
	Endpoint          string
	Route             ratelimit.Route
	RateLimiter       *ratelimit.Ratelimiter
	AdditionalHeaders map[string]string
}

type ResponseWithContent struct {
	*http.Response
	Content []byte
}

// figure out a better way to do this
// token, request
var (
	preRequestHooks   []func(string, *http.Request)
	preRequestHooksMu sync.RWMutex

	postRequestHooks   []func(*http.Response, []byte)
	postRequestHooksMu sync.RWMutex
)

func RegisterHook(hook func(string, *http.Request)) {
	RegisterPreRequestHook(hook)
}

func RegisterPreRequestHook(hook func(string, *http.Request)) {
	preRequestHooksMu.Lock()
	preRequestHooks = append(preRequestHooks, hook)
	preRequestHooksMu.Unlock()
}

func RegisterPostRequestHook(hook func(*http.Response, []byte)) {
	postRequestHooksMu.Lock()
	postRequestHooks = append(postRequestHooks, hook)
	postRequestHooksMu.Unlock()
}

// TODO: Allow users to specify custom timeouts
var Client = http.Client{
	Transport: &http.Transport{
		TLSHandshakeTimeout: time.Second * 3,
	},
	Timeout: time.Second * 3,
}

func (e *Endpoint) Request(ctx context.Context, token string, body any, response any) (error, *ResponseWithContent) {
	url := BASE_URL + e.Endpoint

	// Ratelimit
	if e.RateLimiter != nil {
		ch := make(chan error)
		go e.RateLimiter.ExecuteCall(e.Route, ch)

		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context deadline exceeded while waiting for ratelimit"), nil
		case err := <-ch:
			if err != nil {
				return err, nil
			}
		}
	}

	// Create req
	var req *http.Request
	var err error
	if body == nil || e.ContentType == Nil {
		req, err = http.NewRequestWithContext(ctx, string(e.RequestType), url, nil)
	} else {
		contentType := string(e.ContentType)

		// Encode body
		var encoded []byte
		if e.ContentType == ApplicationJson {
			raw, err := json.Marshal(&body)
			if err != nil {
				return err, nil
			}
			encoded = raw
		} else if e.ContentType == ApplicationFormUrlEncoded {
			str, err := qs.Marshal(body)
			if err != nil {
				return err, nil
			}
			encoded = []byte(str)
		} else if e.ContentType == MultipartFormData {
			data, ok := body.(MultipartPayload)
			if !ok {
				return errors.New("Content-Type MultipartFormData specified but EncodeMultipartFormData was missing"), nil
			}

			var boundary string
			encoded, boundary, err = EncodeMultipartFormData(data)
			if err != nil {
				return err, nil
			}

			contentType = fmt.Sprintf("%s; boundary=%s", MultipartFormData, boundary)
		}

		buff := bytes.NewBuffer(encoded)
		req, err = http.NewRequestWithContext(ctx, string(e.RequestType), url, buff)
		if err != nil {
			return err, nil
		}

		req.Header.Set("Content-Type", contentType)
	}

	if err != nil {
		return err, nil
	}

	if token != "" {
		var header string
		if !strings.HasPrefix(token, "Bot ") && !strings.HasPrefix(token, "Bearer ") {
			header += "Bot "
		}

		header += token
		req.Header.Set("Authorization", header)
	}

	for key, value := range e.AdditionalHeaders {
		req.Header.Set(key, value)
	}

	req.Header.Set("User-Agent", "DiscordBot (https://github.com/rxdn/gdl, 1.0.0)")

	var res *http.Response
	var content []byte

	// Execute hooks
	executePreRequestHooks(token, req)
	defer func() {
		executePostRequestHooks(res, content)
	}()

	res, err = Client.Do(req)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()

	if e.RateLimiter != nil {
		e.applyNewRatelimits(res.Header)
	}

	content, err = io.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}

	if res.StatusCode < 200 || res.StatusCode > 226 {
		var parsed ApiV8Error
		if err := json.Unmarshal(content, &parsed); err != nil {
			parsed.Message = string(content)
		}

		err = RestError{
			StatusCode: res.StatusCode,
			ApiError:   parsed,
			Url:        url,
			Raw:        content,
		}

		return err, &ResponseWithContent{
			Response: res,
			Content:  content,
		}
	}

	if response != nil {
		return json.Unmarshal(content, response), &ResponseWithContent{
			Response: res,
			Content:  content,
		}
	} else {
		return nil, &ResponseWithContent{
			Response: res,
			Content:  content,
		}
	}
}

func (e *Endpoint) applyNewRatelimits(header http.Header) {
	// check global limit
	if global, err := strconv.ParseBool(header.Get("X-RateLimit-Global")); err == nil && global {
		if retryAfter, err := strconv.ParseInt(header.Get("Retry-After"), 10, 64); err == nil {
			e.RateLimiter.Store.UpdateGlobalRateLimit(time.Duration(retryAfter) * time.Second)
			return
		}
	}

	// check route limit
	if remaining, err := strconv.Atoi(header.Get("X-Ratelimit-Remaining")); err == nil {
		if resetAfterSeconds, err := strconv.ParseFloat(header.Get("X-Ratelimit-Reset-After"), 32); err == nil {
			if err := e.RateLimiter.Store.UpdateRateLimit(e.Route, remaining, time.Duration(resetAfterSeconds*1000)*time.Millisecond); err != nil {
				logrus.Warnf("Error occurred updating ratelimits: %e", err)
			}
		}
	}
}

func executePreRequestHooks(token string, req *http.Request) {
	preRequestHooksMu.RLock()
	defer preRequestHooksMu.RUnlock()

	for _, hook := range preRequestHooks {
		hook(token, req)
	}
}

func executePostRequestHooks(res *http.Response, body []byte) {
	postRequestHooksMu.RLock()
	defer postRequestHooksMu.RUnlock()

	for _, hook := range postRequestHooks {
		hook(res, body)
	}
}
