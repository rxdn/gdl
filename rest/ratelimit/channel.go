package ratelimit

import "fmt"

type ChannelRoute struct {
	Id uint64
}

func NewChannelRoute(id uint64) *ChannelRoute {
	return &ChannelRoute{
		Id: id,
	}
}

func (c *ChannelRoute) Endpoint() string {
	return fmt.Sprintf("/channels/%d", c.Id)
}
