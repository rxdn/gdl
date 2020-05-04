package events

import (
	"github.com/rxdn/gdl/objects/guild/emoji"
)

// Sent when a bot removes all instances of a given emoji from the reactions of a message.
type MessageReactionRemoveEmoji struct {
	ChannelId uint64      `json:"channel_id,string"`
	GuildId   uint64      `json:"guild_id,string"`
	MessageId uint64      `json:"message_id,string"`
	Emoji     emoji.Emoji `json:"emoji,string"` // Partial emoji object; https://discord.com/developers/docs/resources/emoji#emoji-object-gateway-reaction-standard-emoji-example
}
