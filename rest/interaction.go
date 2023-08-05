package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

func GetGlobalCommands(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteGetGlobalCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &commands)
	return
}

type CreateCommandData struct {
	Id          uint64                                 `json:"id,omitempty"` // Optional: Use to rename without changing ID
	Name        string                                 `json:"name"`
	Description string                                 `json:"description"`
	Options     []interaction.ApplicationCommandOption `json:"options"`
	Type        interaction.ApplicationCommandType     `json:"type"`
}

func CreateGlobalCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteCreateGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &command)
	return
}

func ModifyGlobalCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, commandId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands/%d", applicationId, commandId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteModifyGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &command)
	return
}

func ModifyGlobalCommands(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data []CreateCommandData) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/commands", applicationId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteModifyGlobalCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &commands)
	return
}

func DeleteGlobalCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, commandId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/commands/%d", applicationId, commandId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteDeleteGlobalCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

func GetGuildCommands(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &commands)
	return
}

func CreateGuildCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &command)
	return
}

func ModifyGuildCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64, data CreateCommandData) (command interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &command)
	return
}

func ModifyGuildCommands(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64, data []CreateCommandData) (commands []interaction.ApplicationCommand, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildCommands, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &commands)
	return
}

func DeleteGuildCommand(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGuildCommand, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

type CommandWithPermissionsData struct {
	Id            uint64                                      `json:"id,string,omitempty"`
	ApplicationId uint64                                      `json:"application_id,string,omitempty"`
	GuildId       uint64                                      `json:"guild_id,string,omitempty"`
	Permissions   []interaction.ApplicationCommandPermissions `json:"permissions"`
}

func GetCommandPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64) (command CommandWithPermissionsData, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d/permissions", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetCommandPermissions, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &command)
	return
}

func GetBulkCommandPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64) (commands []CommandWithPermissionsData, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/permissions", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetBulkCommandPermissions, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &commands)
	return
}

func EditCommandPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId, commandId uint64, data CommandWithPermissionsData) (command CommandWithPermissionsData, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/%d/permissions", applicationId, guildId, commandId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteEditCommandPermissions, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &command)
	return
}

func EditBulkCommandPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, guildId uint64, data []CommandWithPermissionsData) (commands []CommandWithPermissionsData, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/guilds/%d/commands/permissions", applicationId, guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteEditBulkCommandPermissions, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &commands)
	return
}

func GetOriginalInteractionResponse(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64) (msg message.Message, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/@original", applicationId, token),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteGetOriginalInteractionResponse, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", nil, &msg)
	return
}

func EditOriginalInteractionResponse(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data WebhookEditBody) (msg message.Message, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/@original", applicationId, token),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteEditOriginalInteractionResponse, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", data, &msg)
	return
}

func DeleteOriginalInteractionResponse(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/@original", applicationId, token),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteEditOriginalInteractionResponse, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", nil, nil)
	return
}

func CreateFollowupMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data WebhookBody) (msg message.Message, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s", applicationId, token),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteCreateFollowupMessage, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", data, &msg)
	return
}

func GetFollowupMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, messageId uint64) (msg message.Message, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/%d", applicationId, token, messageId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteGetFollowupMessage, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", nil, &msg)
	return
}

func EditFollowupMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, messageId uint64, data WebhookBody) (msg message.Message, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/%d", applicationId, token, messageId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteEditFollowupMessage, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", data, &msg)
	return
}

func DeleteFollowupMessages(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, messageId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/webhooks/%d/%s/messages/%d", applicationId, token, messageId),
		Route:       ratelimit.NewWebhookRoute(ratelimit.RouteDeleteFollowupMessage, applicationId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, "", nil, nil)
	return
}
