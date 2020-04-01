package message

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
)

type Message struct {
	Id              uint64 `json:",string"`
	ChannelId       uint64 `json:",string"`
	GuildId         uint64 `json:",string"`
	Author          *user.User
	Member          *member.Member
	Content         string
	Timestamp       string
	EditedTimestamp string
	Tts             bool
	MentionEveryone bool
	Mentions        []*MessageMentionedUser
	MentionsRoles   utils.Uint64StringSlice `json:",string"`
	Attachments     []channel.Attachment
	Embeds          []embed.Embed
	Reactions       []Reaction
	Nonce           string
	Pinned          bool
	WebhookId       uint64
	Type            int
	Activity        MessageActivity
	Application     MessageApplication
}
