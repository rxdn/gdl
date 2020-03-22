package routes

type RestRouteManager struct {
	channels map[uint64]*ChannelRoute
	emojis   map[uint64]*EmojiRoute
}

func NewRestRouteManager() RestRouteManager {
	return RestRouteManager{
		channels: make(map[uint64]*ChannelRoute, 0),
		emojis:   make(map[uint64]*EmojiRoute, 0),
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
