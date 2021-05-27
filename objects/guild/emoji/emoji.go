package emoji

import (
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
)

// https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	Id            objects.NullableSnowflake `json:"id"`
	Name          string                    `json:"name"` // if this is not a custom emote, Name will be the unicode emoji, and Id will be 0
	Roles         utils.Uint64StringSlice   `json:"roles,omitempty"`
	User          user.User                 `json:"user"`
	RequireColons bool                      `json:"require_colons"`
	Managed       bool                      `json:"managed"`
	Animated      bool                      `json:"animated"`
}

func (e *Emoji) ToCachedEmoji(guildId uint64) CachedEmoji {
	return CachedEmoji{
		GuildId:       guildId,
		Name:          e.Name,
		Roles:         e.Roles,
		User:          e.User.Id,
		RequireColons: e.RequireColons,
		Managed:       e.Managed,
		Animated:      e.Animated,
	}
}
