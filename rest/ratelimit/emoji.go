package ratelimit

import "fmt"

type EmojiRoute struct {
	GuildId uint64
}

func NewEmojiRoute(guildId uint64) *EmojiRoute {
	return &EmojiRoute{
		GuildId: guildId,
	}
}

func (r *EmojiRoute) Endpoint() string {
	return fmt.Sprintf("/guilds/%d/emojis", r.GuildId)
}
