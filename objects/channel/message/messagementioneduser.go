package message

import (
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

// Mentions is an array of users with partial member
type MessageMentionedUser struct {
	user.User
	Member member.Member
}
