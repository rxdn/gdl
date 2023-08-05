package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

type WebhookData struct {
	Username string `json:"name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func CreateWebhook(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data WebhookData) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/webhooks", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteCreateWebhook, channelId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(ctx, token, data, &webhook)
	return webhook, err
}

func GetChannelWebhooks(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/webhooks", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetChannelWebhooks, channelId),
		RateLimiter: rateLimiter,
	}

	var webhooks []guild.Webhook
	err, _ := endpoint.Request(ctx, token, nil, &webhooks)
	return webhooks, err
}

func GetGuildWebhooks(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/webhooks", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildWebhooks, guildId),
		RateLimiter: rateLimiter,
	}

	var webhooks []guild.Webhook
	err, _ := endpoint.Request(ctx, token, nil, &webhooks)
	return webhooks, err
}

func GetWebhook(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteGetWebhook, webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(ctx, token, nil, &webhook)
	return webhook, err
}

// does not return a User object
func GetWebhookWithToken(ctx context.Context, webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteGetWebhookWithToken, webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(ctx, "", nil, &webhook)
	return webhook, err
}

type ModifyWebhookData struct {
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	ChannelId uint64 `json:"channel_id,string,omitempty"`
}

func ModifyWebhook(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, data ModifyWebhookData) (guild.Webhook, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteModifyWebhook, webhookId),
		RateLimiter: rateLimiter,
	}

	var webhook guild.Webhook
	err, _ := endpoint.Request(ctx, token, data, &webhook)
	return webhook, err
}

func ModifyWebhookWithToken(ctx context.Context, webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, data WebhookData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteModifyWebhookWithToken, webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, "", data, nil)
	return err
}

func DeleteWebhook(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d", webhookId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteDeleteWebhook, webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func DeleteWebhookWithToken(ctx context.Context, webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", webhookId, webhookToken),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteDeleteWebhookWithToken, webhookId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, "", nil, nil)
	return err
}

type WebhookBody struct {
	Content         string                 `json:"content,omitempty"`
	Username        string                 `json:"username,omitempty"`
	AvatarUrl       string                 `json:"avatar_url,omitempty"`
	Tts             bool                   `json:"tts"`
	Flags           uint                   `json:"flags,omitempty"`
	Embeds          []*embed.Embed         `json:"embeds,omitempty"`
	AllowedMentions message.AllowedMention `json:"allowed_mentions,omitempty"`
	Components      []component.Component  `json:"components,omitempty"`
	Attachments     []request.Attachment   `json:"attachments,omitempty"`
	ThreadName      string                 `json:"thread_name,omitempty"`
}

func (d WebhookBody) GetAttachments() []request.Attachment {
	return d.Attachments
}

// if wait=true, a message object will be returned
func ExecuteWebhook(ctx context.Context, webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId uint64, wait bool, data WebhookBody) (*message.Message, error) {
	var endpoint request.Endpoint

	if len(data.Attachments) == 0 {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s?wait=%t", webhookId, webhookToken, wait),
			Route:       ratelimit.NewWebhookRoute(ratelimit.RouteExecuteWebhook, webhookId),
			RateLimiter: rateLimiter,
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s?wait=%t", webhookId, webhookToken, wait),
			Route:       ratelimit.NewWebhookRoute(ratelimit.RouteExecuteWebhook, webhookId),
			RateLimiter: rateLimiter,
		}

	}

	if wait {
		var message message.Message
		err, _ := endpoint.Request(ctx, "", data, &message)
		return &message, err
	} else {
		err, _ := endpoint.Request(ctx, "", data, nil)
		return nil, err
	}
}

type WebhookEditBody struct {
	Content         string                 `json:"content"`
	Embeds          []*embed.Embed         `json:"embeds"`
	AllowedMentions message.AllowedMention `json:"allowed_mentions"`
	Components      []component.Component  `json:"components"`
	Attachments     []request.Attachment   `json:"attachments"`
}

func (d WebhookEditBody) GetAttachments() []request.Attachment {
	return d.Attachments
}

func EditWebhookMessage(ctx context.Context, webhookToken string, rateLimiter *ratelimit.Ratelimiter, webhookId, messageId uint64, data WebhookEditBody) (msg message.Message, err error) {
	var endpoint request.Endpoint

	if len(data.Attachments) == 0 {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/%d", webhookId, webhookToken, messageId),
			Route:       ratelimit.NewWebhookRoute(ratelimit.RouteEditWebhookMessage, webhookId),
			RateLimiter: rateLimiter,
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/%d", webhookId, webhookToken, messageId),
			Route:       ratelimit.NewWebhookRoute(ratelimit.RouteEditWebhookMessage, webhookId),
			RateLimiter: rateLimiter,
		}
	}

	err, _ = endpoint.Request(ctx, "", data, &msg)
	return
}
