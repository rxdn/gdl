package objects

import "github.com/Dot-Rar/gdl/utils"

type Presence struct {
	User         *User
	Roles        utils.Uint64StringSlice `json:",string"`
	Game         Activity
	GuildId      uint64 `json:",string"`
	Status       string
	Activities   []*Activity
	ClientStatus ClientStatus
}
