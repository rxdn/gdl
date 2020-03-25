package objects

import "time"

type IntegrationExpireBehaviour int

const (
	RemoveRole IntegrationExpireBehaviour = 0
	Kick       IntegrationExpireBehaviour = 1
)

type Integration struct {
	Id                uint64                     ` json:"id,string"`
	Name              string                     `json:"name"`
	Type              string                     `json:"type"` // twitch, youtube, etc.
	Enabled           bool                       `json:"enabled"`
	Syncing           bool                       `json:"syncing"`
	RoleId            uint64                     `json:"role_id,string"`
	EnableEmoticons   bool                       `json:"enable_emoticons"`
	ExpireBehaviour   IntegrationExpireBehaviour `json:"expire_behavior"`
	ExpireGracePeriod int                        `json:"expire_grace_period"`
	User              *User                      `json:"user"`
	Account           *Account                   `json:"account"`
	SyncedAt          time.Time                  `json:"synced_at"`
}
