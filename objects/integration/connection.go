package integration

type Visibility int

const (
	VisibilityNone     Visibility = 0
	VisibilityEveryone Visibility = 1
)

type Connection struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"` // youtube, twitch, etc.
	Revoked      bool          `json:"revoked"`
	Integrations []Integration `json:"integrations"`
	Verified     bool          `json:"verified"`
	FriendSync   bool          `json:"friend_sync"`
	ShowActivity bool          `json:"show_activity"`
	Visibility   Visibility    `json:"visibility"`
}
