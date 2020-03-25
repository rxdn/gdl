package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/rest/routes"
)

func GetInvite(token string, inviteCode string, withCounts bool) (*objects.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/invites/%s?with_counts=%v", inviteCode, withCounts),
	}

	var invite objects.Invite
	err, _ := endpoint.Request(token, &routes.RouteManager.GetInviteRoute(inviteCode).Ratelimiter, nil, &invite)
	return &invite, err
}

func DeleteInvite(token string, inviteCode string) (*objects.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/invites/%s", inviteCode),
	}

	var invite objects.Invite
	err, _ := endpoint.Request(token, &routes.RouteManager.GetInviteRoute(inviteCode).Ratelimiter, nil, &invite)
	return &invite, err
}
