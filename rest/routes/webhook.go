package routes

import "fmt"

type WebhookRoute struct {
	Id          uint64
	Ratelimiter Ratelimiter
}

func NewWebhookRoute(id uint64) *WebhookRoute {
	return &WebhookRoute{
		Id:          id,
		Ratelimiter: NewRatelimiter(),
	}
}

func (e *WebhookRoute) Endpoint() string {
	return fmt.Sprintf("/webhooks/%d", e.Id)
}

func (e *WebhookRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
