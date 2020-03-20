package events

import "github.com/Dot-Rar/gdl/objects"

// Sent when a bot removes all instances of a given emoji from the reactions of a message.
type MessageReactionRemoveEmoji struct {
	ChannelId uint64         `json:"channel_id,string"`
	GuildId   uint64         `json:"guild_id,string"`
	MessageId uint64         `json:"message_id,string"`
	Emoji     *objects.Emoji `json:"emoji,string"` // Partial emoji object; https://discordapp.com/developers/docs/resources/emoji#emoji-object-gateway-reaction-standard-emoji-example
}
