package interaction

import (
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
)

type ApplicationCommandCallbackData struct {
	Tts             bool                   `json:"tts"`
	Content         string                 `json:"content,omitempty"`
	Embeds          []*embed.Embed          `json:"embeds,omitempty"`
	AllowedMentions message.AllowedMention `json:"allowed_mentions,omitempty"`
	Flags           uint                   `json:"flags"`
}
