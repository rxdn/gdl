package routes

import "fmt"

type InviteRoute struct {
	Id          string
	Ratelimiter Ratelimiter
}

func NewInviteRoute(inviteCode string, rrm *RestRouteManager) *InviteRoute {
	return &InviteRoute{
		Id:          inviteCode,
		Ratelimiter: NewRatelimiter(rrm),
	}
}

func (e *InviteRoute) Endpoint() string {
	return fmt.Sprintf("/invites/%s", e.Id)
}

func (e *InviteRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
