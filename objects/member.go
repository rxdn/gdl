package objects

import (
	"github.com/rxdn/gdl/utils"
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

type Attachment struct {
	Id       uint64 `json:",string"`
	Filename string
	Size     int
	url      string
	ProxyUrl string
	height   int
	Width    int
}
