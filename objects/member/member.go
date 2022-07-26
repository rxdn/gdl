package member

import (
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"time"
)

type Member struct {
	User         user.User               `json:"user"`
	Nick         string                  `json:"nick"`
	Roles        utils.Uint64StringSlice `json:"roles"`
	JoinedAt     time.Time               `json:"joined_at"`
	PremiumSince *time.Time              `json:"premium_since"` // when the user started boosting the guild
	Deaf         bool                    `json:"deaf"`
	Mute         bool                    `json:"mute"`
	Permissions  uint64                  `json:"permissions,string"`
}

func (m *Member) HasRole(roleId uint64) bool {
	for _, memberRole := range m.Roles {
		if memberRole == roleId {
			return true
		}
	}
	return false
}

func (m *Member) ToCachedMember() CachedMember {
	return CachedMember{
		Nick:         m.Nick,
		Roles:        m.Roles,
		JoinedAt:     m.JoinedAt,
		PremiumSince: m.PremiumSince,
		Deaf:         m.Deaf,
		Mute:         m.Mute,
	}
}

