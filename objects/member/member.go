package member

import (
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"time"
)

type Member struct {
	User     *user.User              `json:"user"`
	Nick     string                  `json:"nick"`
	Roles    utils.Uint64StringSlice `json:"roles,string"`
	JoinedAt time.Time               `json:"joined_at"`
	Deaf     bool                    `json:"deaf"`
	Mute     bool                    `json:"mute"`
}

func (m *Member) HasRole(roleId uint64) bool {
	for _, memberRole := range m.Roles {
		if memberRole == roleId {
			return true
		}
	}
	return false
}
