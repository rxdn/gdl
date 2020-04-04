package routes

import "sync"

type RestRouteManager struct {
	*sync.RWMutex
	GlobalRetryAfter int64 // unix epoch for when to retry after
	channels         map[uint64]*ChannelRoute
	emojis           map[uint64]*EmojiRoute
	guilds           map[uint64]*GuildRoute
	invites          map[string]*InviteRoute
	users            map[uint64]*UserRoute
	self             *UserRoute
	voice            *VoiceRoute
	webhooks         map[uint64]*WebhookRoute
}

func NewRestRouteManager() RestRouteManager {
	rrm := RestRouteManager{
		RWMutex:  &sync.RWMutex{},
		channels: make(map[uint64]*ChannelRoute, 0),
		emojis:   make(map[uint64]*EmojiRoute, 0),
		guilds:   make(map[uint64]*GuildRoute, 0),
		invites:  make(map[string]*InviteRoute, 0),
		users:    make(map[uint64]*UserRoute, 0),
		webhooks: make(map[uint64]*WebhookRoute, 0),
	}

	rrm.voice = NewVoiceRoute(&rrm)
	rrm.self = NewUserRoute(0, &rrm)

	return rrm
}

var RouteManager = NewRestRouteManager()

func (rrm *RestRouteManager) GetChannelRoute(channelId uint64) *ChannelRoute {
	rrm.RLock()
	route, ok := rrm.channels[channelId]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewChannelRoute(channelId, rrm)
	rrm.Lock()
	rrm.channels[channelId] = route
	rrm.RUnlock()
	return route
}

func (rrm *RestRouteManager) GetEmojiRoute(guildId uint64) *EmojiRoute {
	rrm.RLock()
	route, ok := rrm.emojis[guildId]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewEmojiRoute(guildId, rrm)
	rrm.Lock()
	rrm.emojis[guildId] = route
	rrm.Unlock()
	return route
}

func (rrm *RestRouteManager) GetGuildRoute(guildId uint64) *GuildRoute {
	rrm.RLock()
	route, ok := rrm.guilds[guildId]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewGuildRoute(guildId, rrm)
	rrm.Lock()
	rrm.guilds[guildId] = route
	rrm.Unlock()
	return route
}

func (rrm *RestRouteManager) GetInviteRoute(inviteCode string) *InviteRoute {
	rrm.RLock()
	route, ok := rrm.invites[inviteCode]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewInviteRoute(inviteCode, rrm)
	rrm.Lock()
	rrm.invites[inviteCode] = route
	rrm.Unlock()
	return route
}

func (rrm *RestRouteManager) GetUserRoute(userId uint64) *UserRoute {
	rrm.RLock()
	route, ok := rrm.users[userId]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewUserRoute(userId, rrm)
	rrm.Lock()
	rrm.users[userId] = route
	rrm.Unlock()
	return route
}

func (rrm *RestRouteManager) GetSelfRoute() *UserRoute {
	return rrm.self
}

func (rrm *RestRouteManager) GetVoiceRoute() *VoiceRoute {
	return rrm.voice
}

func (rrm *RestRouteManager) GetWebhookRoute(webhookId uint64) *WebhookRoute {
	rrm.RLock()
	route, ok := rrm.webhooks[webhookId]
	rrm.RUnlock()
	if ok {
		return route
	}

	route = NewWebhookRoute(webhookId, rrm)
	rrm.Lock()
	rrm.webhooks[webhookId] = route
	rrm.Unlock()
	return route
}
