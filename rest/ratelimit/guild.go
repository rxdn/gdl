package ratelimit

import "fmt"

type GuildRoute struct {
	Id uint64
}

func NewGuildRoute(id uint64) *GuildRoute {
	return &GuildRoute{
		Id: id,
	}
}

func (e *GuildRoute) Endpoint() string {
	return fmt.Sprintf("/guilds/%d", e.Id)
}
