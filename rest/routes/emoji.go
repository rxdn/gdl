package routes

import "fmt"

type EmojiRoute struct {
	GuildId     uint64
	Ratelimiter Ratelimiter
}

func NewEmojiRoute(guildId uint64, rrm *RestRouteManager) *EmojiRoute {
	return &EmojiRoute{
		GuildId:     guildId,
		Ratelimiter: NewRatelimiter(rrm),
	}
}

func (r *EmojiRoute) Endpoint() string {
	return fmt.Sprintf("/guilds/%d/emojis", r.GuildId)
}

func (r *EmojiRoute) GetRatelimit() *Ratelimiter {
	return &r.Ratelimiter
}
