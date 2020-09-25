package events

import "github.com/rxdn/gdl/objects/user"

type Status string

const (
	IDLE    Status = "idle"
	DND     Status = "dnd"
	ONLINE  Status = "online"
	OFFLINE Status = "offline"
)

type PresenceUpdate struct {
	User         user.User         `json:"user"`
	GuildId      uint64            `json:"guild_id,string"`
	Status       string            `json:"status"`
	Activities   []user.Activity   `json:"activities"`
	ClientStatus user.ClientStatus `json:"client_status"`
}
