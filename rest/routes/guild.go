package routes

import "fmt"

type GuildRoute struct {
	Id          uint64
	Ratelimiter Ratelimiter
}

func NewGuildRoute(id uint64, rrm *RestRouteManager) *GuildRoute {
	return &GuildRoute{
		Id:          id,
		Ratelimiter: NewRatelimiter(rrm),
	}
}

func (e *GuildRoute) Endpoint() string {
	return fmt.Sprintf("/guilds/%d", e.Id)
}

func (e *GuildRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
