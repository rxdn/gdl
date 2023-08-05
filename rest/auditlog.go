package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/auditlog"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"net/url"
	"strconv"
)

type GetGuildAuditLogData struct {
	UserId     uint64
	ActionType auditlog.AuditLogEvent
	Before     uint64 // audit log entry ID
	Limit      int
}

func (d *GetGuildAuditLogData) Query() string {
	query := url.Values{}

	if d.UserId != 0 {
		query.Set("user_id", strconv.FormatUint(d.UserId, 10))
	}

	if d.ActionType != 0 {
		query.Set("action_type", strconv.Itoa(int(d.ActionType)))
	}

	if d.Before != 0 {
		query.Set("before", strconv.FormatUint(d.Before, 10))
	}

	if d.Limit > 100 || d.Limit < 1 {
		d.Limit = 50
	}
	query.Set("limit", strconv.Itoa(d.Limit))

	return query.Encode()
}

func GetGuildAuditLog(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data GetGuildAuditLogData) (log auditlog.AuditLog, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/audit-logs?%s", guildId, data.Query()),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildAuditLog, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &log)
	return
}
