package objects

type AllowedMention struct {
	Parse []AllowedMentionType `json:"parse"`
	Roles []string `json:"roles"`
	Users []string `json:"users"`
}

// Helper
var MentionEveryone = AllowedMention{
	Parse: []AllowedMentionType{EVERYONE},
}
