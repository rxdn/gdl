package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

func GetGlobalCommands(token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGlobalCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, &commands)
	return
}

type CreateCommandData struct {
	Id          uint64                                 `json:"id,omitempty"` // Optional: Use to rename without changing ID
	Name        string                                 `json:"name"`
	Description string                                 `json:"description"`
	Options     []interaction.ApplicationCommandOption `json:"options"`
}

func CreateGlobalCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &command)
	return
}

func ModifyGlobalCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, commandId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands/%d", applicationId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &command)
	return
}

func ModifyGlobalCommands(token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data []CreateCommandData) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGlobalCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &commands)
	return
}

func DeleteGlobalCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, commandId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/commands/%d", applicationId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, nil)
	return
}

func GetGuildCommands(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, &commands)
	return
}

func CreateGuildCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &command)
	return
}

func ModifyGuildCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &command)
	return
}

func DeleteGuildCommand(token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, nil)
	return
}
