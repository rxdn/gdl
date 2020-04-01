package message

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"regexp"
	"strconv"
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

var channelMentionRegex = regexp.MustCompile(`<#(\d+)>`)
var userMentionRegex = regexp.MustCompile(`<@(\d+)>`)

// no guarantee that these are actual channels
func (m *Message) ChannelMentions() []uint64 {
	mentions := make([]uint64, 0)
	for _, id := range channelMentionRegex.FindStringSubmatch(m.Content) {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil {
			mentions = append(mentions, parsed)
		}
	}
	return mentions
}

// no guarantee that these are actual users
func (m *Message) UserMentions() []uint64 {
	mentions := make([]uint64, 0)
	for _, id := range userMentionRegex.FindStringSubmatch(m.Content) {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil {
			mentions = append(mentions, parsed)
		}
	}
	return mentions
}
