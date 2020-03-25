package objects

type AllowedMentionType string

const(
	EVERYONE AllowedMentionType = "everyone"
	USERS AllowedMentionType = "users"
	ROLES AllowedMentionType = "roles"
)

type AllowedMention struct {
	Parse []AllowedMentionType `json:"parse"`
	Roles []string `json:"roles"`
	Users []string `json:"users"`
}

// Helper
var MentionEveryone = AllowedMention{
	Parse: []AllowedMentionType{EVERYONE},
}
