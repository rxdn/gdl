package events

import (
	"github.com/rxdn/gdl/objects/user"
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
	user.Presence
	PremiumSince time.Time `json:"premium_since"` // When the user started boosting the guild
	Nick         string    `json:"nick"`
}
