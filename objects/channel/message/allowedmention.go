package message

import "github.com/rxdn/gdl/utils"

type AllowedMention struct {
	Parse       []AllowedMentionType      `json:"parse,omitempty"`
	Roles       []utils.Uint64StringSlice `json:"roles,omitempty"`
	Users       []utils.Uint64StringSlice `json:"users,omitempty"`
	RepliedUser bool                      `json:"replied_user"`
}

// Helper
var MentionEveryone = AllowedMention{
	Parse: []AllowedMentionType{EVERYONE},
}
