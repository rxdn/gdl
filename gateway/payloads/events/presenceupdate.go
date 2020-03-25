package events

import (
	"github.com/rxdn/gdl/objects"
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
	*objects.Presence
	PremiumSince time.Time               `json:"premium_since"` // When the user started boosting the guild
	Nick         string                  `json:"nick"`
}
