package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"net/url"
	"strconv"
)

func GetCurrentUser(token string, rateLimiter *ratelimit.Ratelimiter) (user.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    "/users/@me",
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteGetCurrentUser, 0),
		RateLimiter: rateLimiter,
	}

	var user user.User
	err, _ := endpoint.Request(token, nil, &user)
	return user, err
}

func GetUser(token string, rateLimiter *ratelimit.Ratelimiter, userId uint64) (user.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/%d", userId),
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteGetUser, userId),
		RateLimiter: rateLimiter,
	}

	var user user.User
	err, _ := endpoint.Request(token, nil, &user)
	return user, err
}

type ModifyUserData struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func ModifyCurrentUser(token string, rateLimiter *ratelimit.Ratelimiter, data ModifyUserData) (user.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    "/users/@me",
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteModifyCurrentUser, 0),
		RateLimiter: rateLimiter,
	}

	var user user.User
	err, _ := endpoint.Request(token, data, &user)
	return user, err
}

type CurrentUserGuildsData struct {
	Before uint64
	After  uint64
	Limit  int
}

func (d *CurrentUserGuildsData) Query() string {
	query := url.Values{}

	if d.Before != 0 {
		query.Set("before", strconv.FormatUint(d.Before, 10))
	}

	if d.After != 0 {
		query.Set("after", strconv.FormatUint(d.After, 10))
	}

	if d.Limit > 200 || d.Limit < 1 {
		d.Limit = 200
	}
	query.Set("limit", strconv.Itoa(d.Limit))

	return query.Encode()
}

func GetCurrentUserGuilds(token string, rateLimiter *ratelimit.Ratelimiter, data CurrentUserGuildsData) ([]guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/guilds?%s", data.Query()),
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteGetCurrentUserGuilds, 0),
		RateLimiter: rateLimiter,
	}

	var guilds []guild.Guild
	err, _ := endpoint.Request(token, nil, &guilds)
	return guilds, err
}

func LeaveGuild(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/guilds/%d", guildId),
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteLeaveGuild, 0),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func CreateDM(token string, rateLimiter *ratelimit.Ratelimiter, recipientId uint64) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/users/@me/channels"),
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteCreateDM, recipientId),
		RateLimiter: rateLimiter,
	}

	body := map[string]interface{}{
		"recipient_id": strconv.FormatUint(recipientId, 10),
	}

	var channel channel.Channel
	err, _ := endpoint.Request(token, body, &channel)
	return channel, err
}

func GetUserConnections(token string, rateLimiter *ratelimit.Ratelimiter) ([]integration.Connection, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/connections"),
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteGetUserConnections, 0),
		RateLimiter: rateLimiter,
	}

	var connections []integration.Connection
	err, _ := endpoint.Request(token, nil, &connections)
	return connections, err
}
