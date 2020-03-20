package objects

type Presence struct {
	User         *User
	Roles        []uint64 `json:",string"`
	Game         Activity
	GuildId      uint64 `json:",string"`
	Status       string
	Activities   []*Activity
	ClientStatus ClientStatus
}
