package routes

type RestRouteManager struct {
	channels map[uint64]*ChannelRoute
	emojis   map[uint64]*EmojiRoute
	guilds   map[uint64]*GuildRoute
	invites  map[string]*InviteRoute
	users    map[uint64]*UserRoute
	self     *UserRoute
	voice    *VoiceRoute
	webhooks    map[uint64]*WebhookRoute
}

func NewRestRouteManager() RestRouteManager {
	return RestRouteManager{
		channels: make(map[uint64]*ChannelRoute, 0),
		emojis:   make(map[uint64]*EmojiRoute, 0),
		guilds:   make(map[uint64]*GuildRoute, 0),
		invites:  make(map[string]*InviteRoute, 0),
		users:    make(map[uint64]*UserRoute, 0),
		self:     NewUserRoute(0),
		voice:    NewVoiceRoute(),
		webhooks:    make(map[uint64]*WebhookRoute, 0),
	}
}

var RouteManager = NewRestRouteManager()

func (rrm *RestRouteManager) GetChannelRoute(channelId uint64) *ChannelRoute {
	route, ok := rrm.channels[channelId]
	if ok {
		return route
	}

	route = NewChannelRoute(channelId)
	rrm.channels[channelId] = route
	return route
}

func (rrm *RestRouteManager) GetEmojiRoute(guildId uint64) *EmojiRoute {
	route, ok := rrm.emojis[guildId]
	if ok {
		return route
	}

	route = NewEmojiRoute(guildId)
	rrm.emojis[guildId] = route
	return route
}

func (rrm *RestRouteManager) GetGuildRoute(guildId uint64) *GuildRoute {
	route, ok := rrm.guilds[guildId]
	if ok {
		return route
	}

	route = NewGuildRoute(guildId)
	rrm.guilds[guildId] = route
	return route
}

func (rrm *RestRouteManager) GetInviteRoute(inviteCode string) *InviteRoute {
	route, ok := rrm.invites[inviteCode]
	if ok {
		return route
	}

	route = NewInviteRoute(inviteCode)
	rrm.invites[inviteCode] = route
	return route
}

func (rrm *RestRouteManager) GetUserRoute(userId uint64) *UserRoute {
	route, ok := rrm.users[userId]
	if ok {
		return route
	}

	route = NewUserRoute(userId)
	rrm.users[userId] = route
	return route
}

func (rrm *RestRouteManager) GetSelfRoute() *UserRoute {
	return rrm.self
}

func (rrm *RestRouteManager) GetVoiceRoute() *VoiceRoute {
	return rrm.voice
}

func (rrm *RestRouteManager) GetWebhookRoute(webhookId uint64) *WebhookRoute {
	route, ok := rrm.webhooks[webhookId]
	if ok {
		return route
	}

	route = NewWebhookRoute(webhookId)
	rrm.webhooks[webhookId] = route
	return route
}
