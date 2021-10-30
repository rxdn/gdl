package ratelimit

type RouteType uint8

const (
	RouteTypeGuild RouteType = iota
	RouteTypeChannel
	RouteTypeWebhook
	RouteTypeApplication
	RouteTypeOther
)

type RouteId uint16

type Route struct {
	Type      RouteType
	Id        RouteId
	Snowflake uint64
}

func NewRoute(routeType RouteType, id RouteId, snowflake uint64) Route {
	return Route{
		Type:      routeType,
		Id:        id,
		Snowflake: snowflake,
	}
}

func NewGuildRoute(id RouteId, snowflake uint64) Route {
	return NewRoute(RouteTypeGuild, id, snowflake)
}

func NewChannelRoute(id RouteId, snowflake uint64) Route {
	return NewRoute(RouteTypeChannel, id, snowflake)
}

func NewWebhookRoute(id RouteId, snowflake uint64) Route {
	return NewRoute(RouteTypeWebhook, id, snowflake)
}

func NewApplicationRoute(id RouteId, snowflake uint64) Route {
	return NewRoute(RouteTypeApplication, id, snowflake)
}

func NewOtherRoute(id RouteId, snowflake uint64) Route {
	return NewRoute(RouteTypeOther, id, snowflake)
}

const (
	// /guilds/:id/audit-logs
	RouteGetGuildAuditLog RouteId = iota

	// /channels/:id/...
	RouteGetChannel
	RouteModifyChannel
	RouteDeleteChannel
	RouteGetChannelMessages
	RouteGetChannelMessage
	RouteCreateMessage
	RouteCrosspostMessage
	RouteCreateReaction
	RouteDeleteOwnReaction
	RouteDeleteUserReaction
	RouteGetReactions
	RouteDeleteAllReactions
	RouteDeleteAllReactionsForEmoji
	RouteEditMessage
	RouteDeleteMessage
	RouteBulkDeleteMessages
	RouteEditChannelPermissions
	RouteGetChannelInvites
	RouteCreateChannelInvite
	RouteDeleteChannelPermission
	RouteFollowNewsChannel
	RouteTriggerTypingIndicator
	RouteGetPinnedMessages
	RouteAddPinnedChannelMessage
	RouteDeletePinnedChannelMessage
	RouteGetThreadMembers
	RouteStartThreadWithMessage
	RouteStartThreadWithoutMessage
	RouteGetActiveThreads
	RouteGetArchivedPrivateSelfThreads
	RouteGetArchivedPublicThreads
	RouteGetArchivedPrivateThreads
	RouteGroupDMAddRecipient
	RouteGroupDMRemoveRecipient

	// /guilds/:id/emojis
	RouteListGuildEmojis
	RouteGetGuildEmoji
	RouteCreateGuildEmoji
	RouteModifyGuildEmoji
	RouteDeleteGuildEmoji

	// /guilds/:id/...
	RouteCreateGuild
	RouteGetGuild
	RouteGetGuildPreview
	RouteModifyGuild
	RouteDeleteGuild
	RouteGetGuildChannels
	RouteCreateGuildChannel
	RouteModifyGuildChannelPositions
	RouteGetGuildMember
	RouteSearchGuildMembers
	RouteListGuildMembers
	RouteAddGuildMember
	RouteModifyGuildMember
	RouteModifyCurrentUserNick
	RouteAddGuildMemberRole
	RouteRemoveGuildMemberRole
	RouteRemoveGuildMember
	RouteGetGuildBans
	RouteGetGuildBan
	RouteCreateGuildBan
	RouteRemoveGuildBan
	RouteGetGuildRoles
	RouteCreateGuildRole
	RouteModifyGuildRolePositions
	RouteModifyGuildRole
	RouteDeleteGuildRole
	RouteGetGuildPruneCount
	RouteBeginGuildPrune
	RouteGetGuildVoiceRegions
	RouteGetGuildInvites
	RouteGetGuildIntegrations
	RouteCreateGuildIntegration
	RouteModifyGuildIntegration
	RouteDeleteGuildIntegration
	RouteSyncGuildIntegration
	RouteGetGuildWidgetSettings
	RouteModifyGuildWidget
	RouteGetGuildWidget
	RouteGetGuildVanityURL
	RouteGuildWidgetImage

	// /invites/:id
	// Invites seemingly don't have ratelimits, but we need these enums internally
	RouteGetInvite
	RouteDeleteInvite

	// /guilds/templates/:code
	// Also seemingly no ratelimits
	RouteGetTemplate
	RouteCreateGuildTemplate

	// /users/:id/...
	// Again, seemingly no ratelimits but we need these internally
	RouteGetCurrentUser
	RouteGetUser
	RouteModifyCurrentUser
	RouteGetCurrentUserGuilds
	RouteLeaveGuild
	RouteGetUserDMs
	RouteCreateDM
	RouteCreateGroupDM
	RouteGetUserConnections

	// /voice/regions
	RouteListVoiceRegions

	// /channels/:id/webhooks
	RouteCreateWebhook
	RouteGetChannelWebhooks

	// /guilds/:id/webhooks
	RouteGetGuildWebhooks

	// /webhooks/:id/...
	RouteGetWebhook
	RouteGetWebhookWithToken
	RouteModifyWebhook
	RouteModifyWebhookWithToken
	RouteDeleteWebhook
	RouteDeleteWebhookWithToken
	RouteExecuteWebhook
	RouteEditWebhookMessage

	// /applications/:id/...
	RouteGetGlobalCommands
	RouteCreateGlobalCommand
	RouteModifyGlobalCommand
	RouteModifyGlobalCommands
	RouteDeleteGlobalCommand

	// /applications/:id/guild/...
	RouteGetGuildCommands
	RouteCreateGuildCommand
	RouteModifyGuildCommand
	RouteModifyGuildCommands
	RouteDeleteGuildCommand

	// /applications/:id/guild/.../permissions
	RouteGetCommandPermissions
	RouteGetBulkCommandPermissions
	RouteEditCommandPermissions
	RouteEditBulkCommandPermissions

	// /webhooks/:id/:token/...
	RouteGetOriginalInteractionResponse
	RouteEditOriginalInteractionResponse
	RouteDeleteOriginalInteractionResponse
	RouteCreateFollowupMessage
	RouteGetFollowupMessage
	RouteEditFollowupMessage
	RouteDeleteFollowupMessage
)
