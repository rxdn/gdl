package member

import (
	"github.com/rxdn/gdl/objects/user"
	"time"
)

type CachedMember struct {
	Nick         string     `db:"nick"`
	Roles        []uint64   `db:"roles"`
	JoinedAt     time.Time  `db:"joined_at"`
	PremiumSince *time.Time `db:"premium_since"` // when the user started boosting the guild
	Deaf         bool       `db:"deaf"`
	Mute         bool       `db:"mute"`
}

func (m *CachedMember) ToMember(user user.User) Member {
	return Member{
		User:         user,
		Nick:         m.Nick,
		Roles:        m.Roles,
		JoinedAt:     m.JoinedAt,
		PremiumSince: m.PremiumSince,
		Deaf:         m.Deaf,
		Mute:         m.Mute,
	}
}
