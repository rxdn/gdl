package objects

import "time"

type TargetUserType int

const (
	STREAM TargetUserType = 1
)

type Invite struct {
	Code                     string         `json:"code"`
	Guild                    *Guild         `json:"guild"`
	Channel                  *Channel       `json:"channel"`
	Inviter                  *User          `json:"inviter"`
	TargetUser               *User          `json:"target_user"`
	TargetUserType           TargetUserType `json:"target_user_type"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
}

type InviteMetadata struct {
	*Invite
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}
