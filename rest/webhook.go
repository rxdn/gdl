package rest

import (
	"bytes"
	"fmt"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"io"
	"mime/multipart"
	"net/textproto"
	"strconv"
	"strings"
)

type WebhookData struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func CreateWebhook(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data WebhookData) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/webhooks", channelId),
		Bucket:      ratelimit.NewChannelBucket(channelId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(token, data, &webhook)
	return webhook, err
}

func GetChannelWebhooks(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/webhooks", channelId),
		Bucket:      ratelimit.NewChannelBucket(channelId),
		RateLimiter: rateLimiter,
	}

	var webhooks []guild.Webhook
	err, _ := endpoint.Request(token, nil, &webhooks)
	return webhooks, err
}

func GetGuildWebhooks(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/webhooks", guildId),
		Bucket:      ratelimit.NewGuildBucket(guildId),
		RateLimiter: rateLimiter,
	}

	var webhooks []guild.Webhook
	err, _ := endpoint.Request(token, nil, &webhooks)
	return webhooks, err
}

func GetWebhook(token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(token, nil, &webhook)
	return webhook, err
}

// does not return a User object
func GetWebhookWithToken(webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request("", nil, &webhook)
	return webhook, err
}

type ModifyWebhookData struct {
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	ChannelId uint64 `json:"channel_id,string,omitempty"`
}

func ModifyWebhook(token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, data ModifyWebhookData) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(token, data, &webhook)
	return webhook, err
}

func ModifyWebhookWithToken(webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, data WebhookData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request("", data, nil)
	return err
}

func DeleteWebhook(token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func DeleteWebhookWithToken(webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Bucket:      ratelimit.NewWebhookBucket(webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request("", nil, nil)
	return err
}

type WebhookBody struct {
	Content         string                 `json:"content,omitempty"`
	Username        string                 `json:"username,omitempty"`
	AvatarUrl       string                 `json:"avatar_url,omitempty"`
	Tts             bool                   `json:"tts"`
	File            *File                  `json:"file,omitempty"`
	Embeds          []*embed.Embed         `json:"embeds,omitempty"`
	PayloadJson     string                 `json:"payload_json"`
	AllowedMentions message.AllowedMention `json:"allowed_mentions,omitempty"`
}

func (d WebhookBody) EncodeMultipartFormData() ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if d.File != nil {
		fileName := d.File.Name
		fileName = strings.Replace(fileName, "\\", "\\\\", -1)
		fileName = strings.Replace(fileName, "\"", "\\\"", -1)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
		h.Set("Content-Type", d.File.ContentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			return body.Bytes(), writer.Boundary(), err
		}

		if _, err := io.Copy(part, d.File.Reader); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if d.Username != "" {
		if err := writer.WriteField("username", d.Username); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if d.AvatarUrl != "" {
		if err := writer.WriteField("avatar_url", d.AvatarUrl); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if err := writer.WriteField("tts", strconv.FormatBool(d.Tts)); err != nil {
		return body.Bytes(), writer.Boundary(), err
	}

	if d.PayloadJson != "" {
		if err := writer.WriteField("payload_json", d.PayloadJson); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	return []byte(string(body.Bytes()) + "\r\n--" + writer.Boundary() + "--"), writer.Boundary(), nil
}

// if wait=true, a message object will be returned
func ExecuteWebhook(webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, wait bool, data WebhookBody) (*message.Message, error) {
	var endpoint request.Endpoint

	if data.File == nil {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s?wait=%t", webhookId, webhookToken, wait),
			Bucket:      ratelimit.NewWebhookExecuteBucket(webhookId),
			RateLimiter: rateLimiter,
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s?wait=%t", webhookId, webhookToken, wait),
			Bucket:      ratelimit.NewWebhookExecuteBucket(webhookId),
			RateLimiter: rateLimiter,
		}

	}
	if wait {
		var message *message.Message
		err, _ := endpoint.Request("", data, message)
		return message, err
	} else {
		err, _ := endpoint.Request("", data, nil)
		return nil, err
	}
}
