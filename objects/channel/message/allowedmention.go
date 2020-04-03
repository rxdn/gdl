package message

import "github.com/rxdn/gdl/utils"

type AllowedMention struct {
	Parse []AllowedMentionType      `json:"parse"`
	Roles []utils.Uint64StringSlice `json:"roles"`
	Users []utils.Uint64StringSlice `json:"users"`
}

// Helper
var MentionEveryone = AllowedMention{
	Parse: []AllowedMentionType{EVERYONE},
}
