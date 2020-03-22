package routes

type RestRouteManager struct {
	channels map[uint64]*ChannelRoute
}

func NewRestRouteManager() RestRouteManager {
	return RestRouteManager{
		channels: make(map[uint64]*ChannelRoute, 0),
	}
}

func (rrm *RestRouteManager) GetChannelRoute(id uint64) *ChannelRoute {
	route, ok := rrm.channels[id]; if ok {
		return route
	}

	route = NewChannelRoute(id)
	rrm.channels[id] = route
	return route
}

var RouteManager = NewRestRouteManager()
