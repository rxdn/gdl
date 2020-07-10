package guild

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"time"
)

type Guild struct {
	Id                          uint64            `json:"id,string"`
	Name                        string            `json:"name"`
	Icon                        string            `json:"icon"`
	Splash                      string            `json:"splash"`
	Owner                       bool              `json:"owner"`
	OwnerId                     uint64            `json:"owner_id,string"`
	Permissions                 int               `json:"permissions"`
	Region                      string            `json:"region"`
	AfkChannelId                uint64            `json:"afk_channel_id,string"`
	AfkTimeout                  int               `json:"afk_timeout"`
	EmbedEnabled                bool              `json:"embed_enabled"`
	EmbedChannelId              uint64            `json:"embed_channel_id,string"`
	VerificationLevel           int               `json:"verification_level"`
	DefaultMessageNotifications int               `json:"default_message_notifications"`
	ExplicitContentFilter       int               `json:"explicit_content_filter"`
	Roles                       []Role            `json:"roles"`
	Emojis                      []emoji.Emoji     `json:"emojis"`
	Features                    []GuildFeature    `json:"features"`
	MfaLevel                    int               `json:"mfa_level"`
	ApplicationId               uint64            `json:"application_id,string"`
	WidgetEnabled               bool              `json:"widget_enabled"`
	WidgetChannelId             uint64            `json:"widget_channel_id,string"`
	SystemChannelId             uint64            `json:"system_channel_id,string"`
	JoinedAt                    time.Time         `json:"joined_at"`
	Large                       bool              `json:"large"`
	Unavailable                 *bool             `json:"unavailable"`
	MemberCount                 int               `json:"member_count"`
	VoiceStates                 []VoiceState      `json:"voice_state"`
	Members                     []member.Member   `json:"members"`
	Channels                    []channel.Channel `json:"channels"`
	MaxPresences                int               `json:"max_presences"`
	MaxMembers                  int               `json:"max_members"`
	VanityUrlCode               string            `json:"vanity_url_code"`
	Description                 string            `json:"description"`
	Banner                      string            `json:"banner"`
	PremiumTier                 PremiumTier       `json:"premium_tier"`
	PremiumSubscriptionCount    int               `json:"premium_subscription_count"`
	PreferredLocale             string            `json:"preferred_locale"`
	PublicUpdatesChannelId      uint64            `json:"public_updates_channel_id,string"`
	MaxVideoChannelUsers        int               `json:"max_video_channel_users"`
	ApproximateMemberCount      int               `json:"approximate_member_count"`   // Returned on GET /guild/:id
	ApproximatePresenceCount    int               `json:"approximate_presence_count"` // Returned on GET /guild/:id
}

func (g *Guild) ToCachedGuild() (cached CachedGuild) {
	cached = CachedGuild{
		Id:                          g.Id,
		Name:                        g.Name,
		Icon:                        g.Icon,
		Splash:                      g.Splash,
		Owner:                       g.Owner,
		OwnerId:                     g.OwnerId,
		Permissions:                 g.Permissions,
		Region:                      g.Region,
		AfkChannelId:                g.AfkChannelId,
		AfkTimeout:                  g.AfkTimeout,
		EmbedEnabled:                g.EmbedEnabled,
		EmbedChannelId:              g.EmbedChannelId,
		VerificationLevel:           g.VerificationLevel,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		ExplicitContentFilter:       g.ExplicitContentFilter,
		Features:                    g.Features,
		MfaLevel:                    g.MfaLevel,
		ApplicationId:               g.ApplicationId,
		WidgetEnabled:               g.WidgetEnabled,
		WidgetChannelId:             g.WidgetChannelId,
		SystemChannelId:             g.SystemChannelId,
		JoinedAt:                    g.JoinedAt,
		Large:                       g.Large,
		Unavailable:                 g.Unavailable,
		MemberCount:                 g.MemberCount,
		MaxPresences:                g.MaxPresences,
		MaxMembers:                  g.MaxMembers,
		VanityUrlCode:               g.VanityUrlCode,
		Description:                 g.Description,
		Banner:                      g.Banner,
		PremiumTier:                 g.PremiumTier,
		PremiumSubscriptionCount:    g.PremiumSubscriptionCount,
		PreferredLocale:             g.PreferredLocale,
		PublicUpdatesChannelId:      g.PublicUpdatesChannelId,
		MaxVideoChannelUsers:        g.MaxVideoChannelUsers,
	}

	for _, role := range g.Roles {
		cached.Roles = append(cached.Roles, role.Id)
	}

	for _, emoji := range g.Emojis {
		cached.Emojis = append(cached.Emojis, emoji.Id)
	}

	for _, channel := range g.Channels {
		cached.Channels = append(cached.Channels, channel.Id)
	}

	return
}
