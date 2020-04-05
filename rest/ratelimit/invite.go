package ratelimit

import "fmt"

type InviteRoute struct {
	Id string
}

func NewInviteRoute(inviteCode string) *InviteRoute {
	return &InviteRoute{
		Id: inviteCode,
	}
}

func (e *InviteRoute) Endpoint() string {
	return fmt.Sprintf("/invites/%s", e.Id)
}
