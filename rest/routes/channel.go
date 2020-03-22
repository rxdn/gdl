package routes

import "fmt"

type ChannelRoute struct {
	Id          uint64
	Ratelimiter Ratelimiter
}

func NewChannelRoute(id uint64) *ChannelRoute {
	return &ChannelRoute{
		Id:          id,
		Ratelimiter: NewRatelimiter(),
	}
}

func (c *ChannelRoute) Endpoint() string {
	return fmt.Sprintf("/channels/%d", c.Id)
}

func (c *ChannelRoute) GetRatelimit() *Ratelimiter {
	return &c.Ratelimiter
}
