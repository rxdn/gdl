package guild

import "github.com/rxdn/gdl/objects/guild/emoji"

type GuildPreview struct {
	Id                       uint64         `json:"id,string"`
	Name                     string         `json:"name"`
	Icon                     string         `json:"icon"`
	Splash                   string         `json:"splash"`
	DiscoverySplash          string         `json:"discovery_splash"`
	Emojis                   []emoji.Emoji  `json:"emojis"`
	Features                 []GuildFeature `json:"features"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
	Description              string         `json:"description"`
}
