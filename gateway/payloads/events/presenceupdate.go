package events

import (
	"github.com/Dot-Rar/gdl/objects"
	"github.com/Dot-Rar/gdl/utils"
	"time"
)

type Status string

const (
	IDLE    Status = "idle"
	DND     Status = "dnd"
	ONLINE  Status = "online"
	OFFLINE Status = "offline"
)

type PresenceUpdate struct {
	User         *objects.User           `json:"user"`
	Roles        utils.Uint64StringSlice `json:"roles,string"`
	Game         *objects.Activity       `json:"game"`
	GuildId      uint64                  `json:"guild_id,string"`
	Status       Status                  `json:"status"`
	Activities   []*objects.Activity     `json:"activities"`
	ClientStatus *objects.ClientStatus   `json:"client_status"`
	PremiumSince time.Time               `json:"premium_since"` // When the user started boosting the guild
	Nick         string                  `json:"nick"`
}
