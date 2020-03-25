package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/rest/routes"
	"net/url"
	"strconv"
)

func GetCurrentUser(token string) (*objects.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    "/users/@me",
	}

	var user objects.User
	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, nil, &user)
	return &user, err
}

func GetUser(token string, userId uint64) (*objects.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/%d", userId),
	}

	var user objects.User
	err, _ := endpoint.Request(token, &routes.RouteManager.GetUserRoute(userId).Ratelimiter, nil, &user)
	return &user, err
}

type ModifyUserData struct {
	Username string `json:"username,omitempty"`
	Avatar   Image  `json:"avatar,omitempty"`
}

func ModifyCurrentUser(token string, data ModifyUserData) (*objects.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    "/users/@me",
	}

	var user objects.User
	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, data, &user)
	return &user, err
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

	if d.Limit > 100 || d.Limit < 1 {
		d.Limit = 100
	}
	query.Set("limit", strconv.Itoa(d.Limit))

	return query.Encode()
}

func GetCurrentUserGuilds(token string, data CurrentUserGuildsData) ([]*objects.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/guilds?%s", data.Query()),
	}

	var guilds []*objects.Guild
	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, nil, &guilds)
	return guilds, err
}

func LeaveGuild(token string, guildId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/guilds/%d", guildId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, nil, nil)
	return err
}

func CreateDM(token string, recipientId uint64) (*objects.Channel, error) {
	endpoint := request.Endpoint{
		RequestType:       request.POST,
		ContentType:       request.ApplicationJson,
		Endpoint:          fmt.Sprintf("/users/@me/channels"),
	}

	body := map[string]interface{}{
		"recipient_id": strconv.FormatUint(recipientId, 10),
	}

	var channel objects.Channel
	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, body, &channel)
	return &channel, err
}

func GetUserConnections(token string) ([]*objects.Connection, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/users/@me/connections"),
	}

	var connections []*objects.Connection
	err, _ := endpoint.Request(token, &routes.RouteManager.GetSelfRoute().Ratelimiter, nil, &connections)
	return connections, err
}
