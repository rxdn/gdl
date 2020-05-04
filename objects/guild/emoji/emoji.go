package emoji

import (
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
)

// https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	Id            uint64                  `json:"id,string"`
	Name          string                  `json:"name"` // if this is not a custom emote, Name will be the unicode emoji, and Id will be 0
	Roles         utils.Uint64StringSlice `json:"roles,string"`
	User          user.User               `json:"user"`
	RequireColons bool                    `json:"require_colons"`
	Managed       bool                    `json:"managed"`
	Animated      bool                    `json:"animated"`
}

func (e *Emoji) ToCachedEmoji() CachedEmoji {
	return CachedEmoji{
		Name:          e.Name,
		Roles:         e.Roles,
		User:          e.User.Id,
		RequireColons: e.RequireColons,
		Managed:       e.Managed,
		Animated:      e.Animated,
	}
}
