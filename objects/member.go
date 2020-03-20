package objects

import (
	"github.com/Dot-Rar/gdl/utils"
	"time"
)

type Member struct {
	User     *User                   `json:"user"`
	Nick     string                  `json:"nick"`
	Roles    utils.Uint64StringSlice `json:"roles,string"`
	JoinedAt time.Time               `json:"joined_at"`
	Deaf     bool                    `json:"deaf"`
	Mute     bool                    `json:"mute"`
}
