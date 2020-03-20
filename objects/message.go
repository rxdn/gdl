package objects

import "github.com/Dot-Rar/gdl/utils"

type Message struct {
	Id              uint64 `json:",string"`
	ChannelId       uint64 `json:",string"`
	GuildId         uint64 `json:",string"`
	Author          *User
	Member          *Member
	Content         string
	Timestamp       string
	EditedTimestamp string
	Tts             bool
	MentionEveryone bool
	Mentions        []*MessageMentionedUser
	MentionsRoles   utils.Uint64StringSlice `json:",string"`
	Attachments     []Attachment
	Embeds          []Embed
	Reactions       []Reaction
	Nonce           string
	Pinned          bool
	WebhookId       uint64
	Type            int
	Activity        MessageActivity
	Application     MessageApplication
}

// Mentions is an array of users with partial member
type MessageMentionedUser struct {
	*User
	Member Member
}
