package guild

import (
	"github.com/rxdn/gdl/objects"
	"time"
)

type CachedGuild struct {
	Id                          uint64                    `json:"id"`
	Name                        string                    `json:"name"`
	Icon                        string                    `json:"icon"`
	Splash                      string                    `json:"splash"`
	Owner                       bool                      `json:"owner"`
	OwnerId                     uint64                    `json:"owner_id"`
	Permissions                 uint64                    `json:"permissions"`
	Region                      string                    `json:"region"`
	AfkChannelId                objects.NullableSnowflake `json:"afk_channel_id"`
	AfkTimeout                  int                       `json:"afk_timeout"`
	VerificationLevel           int                       `json:"verification_level"`
	DefaultMessageNotifications int                       `json:"default_message_notifications"`
	ExplicitContentFilter       int                       `json:"explicit_content_filter"`
	Roles                       []uint64                  `json:"-"`
	Emojis                      []uint64                  `json:"-"`
	Features                    []GuildFeature            `json:"features"`
	MfaLevel                    int                       `json:"mfa_level"`
	ApplicationId               objects.NullableSnowflake `json:"application_id"`
	WidgetEnabled               bool                      `json:"widget_enabled"`
	WidgetChannelId             objects.NullableSnowflake `json:"widget_channel_id"`
	SystemChannelId             objects.NullableSnowflake `json:"system_channel_id"`
	SystemChannelFlags          uint8                     `json:"system_channel_flags"`
	RulesChannelId              objects.NullableSnowflake `json:"rules_channel_id,omitempty"`
	JoinedAt                    time.Time                 `json:"joined_at"`
	Large                       bool                      `json:"large"`
	Unavailable                 *bool                     `json:"unavailable"`
	MemberCount                 int                       `json:"member_count"`
	Channels                    []uint64                  `json:"-"`
	Presences                   []uint64                  `json:"-"`
	MaxPresences                int                       `json:"max_presences"`
	MaxMembers                  int                       `json:"max_members"`
	VanityUrlCode               string                    `json:"vanity_url_code"`
	Description                 string                    `json:"description"`
	Banner                      string                    `json:"banner"`
	PremiumTier                 PremiumTier               `json:"premium_tier"`
	PremiumSubscriptionCount    int                       `json:"premium_subscription_count"`
	PreferredLocale             string                    `json:"preferred_locale"`
	PublicUpdatesChannelId      objects.NullableSnowflake `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                       `json:"max_video_channel_users"`
}

func (g *CachedGuild) ToGuild(guildId uint64) Guild {
	return Guild{
		Id:                          guildId,
		Name:                        g.Name,
		Icon:                        g.Icon,
		Splash:                      g.Splash,
		Owner:                       g.Owner,
		OwnerId:                     g.OwnerId,
		Permissions:                 g.Permissions,
		Region:                      g.Region,
		AfkChannelId:                g.AfkChannelId,
		AfkTimeout:                  g.AfkTimeout,
		VerificationLevel:           g.VerificationLevel,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		ExplicitContentFilter:       g.ExplicitContentFilter,
		Features:                    g.Features,
		MfaLevel:                    g.MfaLevel,
		ApplicationId:               g.ApplicationId,
		WidgetEnabled:               bool(g.WidgetEnabled),
		WidgetChannelId:             g.WidgetChannelId,
		SystemChannelId:             g.SystemChannelId,
		SystemChannelFlags:          g.SystemChannelFlags,
		RulesChannelId:              g.RulesChannelId,
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
}
