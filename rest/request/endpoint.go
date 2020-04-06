package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pasztorpisti/qs"
	"github.com/rxdn/gdl/rest/ratelimit"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const BASE_URL = "https://discordapp.com/api/v6"

type Endpoint struct {
	RequestType       RequestType
	ContentType       ContentType
	Endpoint          string
	Bucket            string
	RateLimiter       *ratelimit.Ratelimiter
	AdditionalHeaders map[string]string
}

type ResponseWithContent struct {
	*http.Response
	Content []byte
}

// figure out a better way to do this
var Hook func(string)

func (e *Endpoint) Request(token string, body interface{}, response interface{}) (error, *ResponseWithContent) {
	url := BASE_URL + e.Endpoint

	if Hook != nil {
		Hook(url)
	}

	// Ratelimit
	if e.RateLimiter != nil {
		ch := make(chan error)
		go e.RateLimiter.ExecuteCall(e.Bucket, ch)
		if err := <-ch; err != nil {
			return err, nil
		}
	}

	// Create req
	var req *http.Request
	var err error
	if body == nil || e.ContentType == Nil {
		req, err = http.NewRequest(string(e.RequestType), url, nil)
	} else {
		contentType := string(e.ContentType)

		// Encode body
		var encoded []byte
		if e.ContentType == ApplicationJson {
			raw, err := json.Marshal(body)
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
			data, ok := body.(MultipartData)
			if !ok {
				return errors.New("Content-Type MultipartFormData specified but EncodeMultipartFormData was missing"), nil
			}

			var boundary string
			encoded, boundary, err = data.EncodeMultipartFormData()
			if err != nil {
				return err, nil
			}

			contentType = fmt.Sprintf("%s; boundary=%s", MultipartFormData, boundary)
		}

		buff := bytes.NewBuffer(encoded)
		req, err = http.NewRequest(string(e.RequestType), url, buff)
		req.Header.Set("Content-Type", contentType)
	}

	if err != nil {
		return err, nil
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bot %s", token))
	}

	req.Header.Set("X-RateLimit-Precision", "millisecond")

	for key, value := range e.AdditionalHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	res, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()

	if e.RateLimiter != nil {
		e.applyNewRatelimits(res.Header)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}

	if res.StatusCode < 200 || res.StatusCode > 226 {
		return errors.New(fmt.Sprintf("http status %d at %s: %s", res.StatusCode, e.Endpoint, string(content))), nil
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
	// TODO: Global limit
	/*// check global limit
	if global, err := strconv.ParseBool(header.Get("X-RateLimit-Global")); err == nil && global {
		if retryAfter, err := strconv.ParseInt(header.Get("Retry-After"), 10, 64); err == nil {
			ratelimiter.RouteManager.Lock()
			ratelimiter.RouteManager.GlobalRetryAfter = utils.GetCurrentTimeMillis() + retryAfter
			ratelimiter.RouteManager.Unlock()
			ratelimiter.Unlock()
			return
		}
	}*/

	if remaining, err := strconv.Atoi(header.Get("X-Ratelimit-Remaining")); err == nil {
		if resetAfterSeconds, err := strconv.ParseFloat(header.Get("X-Ratelimit-Reset-After"), 32); err == nil {
			e.RateLimiter.Store.UpdateRateLimit(e.Bucket, remaining, time.Duration(resetAfterSeconds*1000)*time.Millisecond)
		}
	}
}
