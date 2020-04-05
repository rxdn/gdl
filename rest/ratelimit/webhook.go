package ratelimit

import "fmt"

type WebhookRoute struct {
	Id uint64
}

func NewWebhookRoute(id uint64) *WebhookRoute {
	return &WebhookRoute{
		Id: id,
	}
}

func (e *WebhookRoute) Endpoint() string {
	return fmt.Sprintf("/webhooks/%d", e.Id)
}
