package events

import (
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"time"
)

type GuildMemberUpdate struct {
	GuildId      uint64                  `json:"guild_id,string"`
	Roles        utils.Uint64StringSlice `json:"roles,string"`
	User         user.User               `json:"user"`
	Nick         string                  `json:"nick"`
	PremiumSince *time.Time              `json:"premium_since"` // When the user started boosting the guidl
}
