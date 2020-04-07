package emoji

import "github.com/rxdn/gdl/objects/user"

type CachedEmoji struct {
	Name          string   `db:"name"`
	Roles         []uint64 `db:"roles"`
	User          uint64   `db:"user"`
	RequireColons bool     `db:"require_colons"`
	Managed       bool     `db:"managed"`
	Animated      bool     `db:"animated"`
}

func (e *CachedEmoji) ToEmoji(emojiId uint64, user user.User) Emoji {
	return Emoji{
		Id:            emojiId,
		Name:          e.Name,
		Roles:         e.Roles,
		User:          user,
		RequireColons: e.RequireColons,
		Managed:       e.Managed,
		Animated:      e.Animated,
	}
}
