package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

func GetInvite(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, inviteCode string, withCounts bool) (invite.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/invites/%s?with_counts=%v", inviteCode, withCounts),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetInvite, 0), // No ratelimit
		RateLimiter: rateLimiter,
	}

	var invite invite.Invite
	err, _ := endpoint.Request(ctx, token, nil, &invite)
	return invite, err
}

func DeleteInvite(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, inviteCode string) (invite.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/invites/%s", inviteCode),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteInvite, 0), // No ratelimit
		RateLimiter: rateLimiter,
	}

	var invite invite.Invite
	err, _ := endpoint.Request(ctx, token, nil, &invite)
	return invite, err
}
