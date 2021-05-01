package guild

import "github.com/rxdn/gdl/objects"

type WelcomeScreen struct {
	Description     string                 `json:"description"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels"`
}

type WelcomeScreenChannel struct {
	ChannelId   uint64                    `json:"channel_id,string"`
	Description string                    `json:"description"`
	EmojiId     objects.NullableSnowflake `json:"emoji_id"`
	EmojiName   *string                   `json:"emoji_name"`
}
