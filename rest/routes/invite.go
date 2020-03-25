package routes

import "fmt"

type InviteRoute struct {
	Id          string
	Ratelimiter Ratelimiter
}

func NewInviteRoute(inviteCode string) *InviteRoute {
	return &InviteRoute{
		Id:          inviteCode,
		Ratelimiter: NewRatelimiter(),
	}
}

func (e *InviteRoute) Endpoint() string {
	return fmt.Sprintf("/invites/%s", e.Id)
}

func (e *InviteRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
