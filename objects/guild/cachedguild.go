package guild

import (
	"time"
)

type CachedGuild struct {
	Id                          uint64            `json:"id"`
	Name                        string            `json:"name"`
	Icon                        string            `json:"icon"`
	Splash                      string            `json:"splash"`
	Owner                       bool              `json:"owner"`
	OwnerId                     uint64            `json:"owner_id"`
	Permissions                 int               `json:"permissions"`
	Region                      string            `json:"region"`
	AfkChannelId                uint64            `json:"afk_channel_id"`
	AfkTimeout                  int               `json:"afk_timeout"`
	EmbedEnabled                bool              `json:"embed_enabled"`
	EmbedChannelId              uint64            `json:"embed_channel_id"`
	VerificationLevel           int               `json:"verification_level"`
	DefaultMessageNotifications int               `json:"default_message_notifications"`
	ExplicitContentFilter       int               `json:"explicit_content_filter"`
	Features                    []GuildFeature    `json:"features"`
	MfaLevel                    int               `json:"mfa_level"`
	ApplicationId               uint64            `json:"application_id,string"`
	WidgetEnabled               bool              `json:"widget_enabled"`
	WidgetChannelId             uint64            `json:"widget_channel_id"`
	SystemChannelId             uint64            `json:"system_channel_id"`
	JoinedAt                    time.Time         `json:"joined_at"`
	Large                       bool              `json:"large"`
	Unavailable                 *bool             `json:"unavailable"`
	MemberCount                 int               `json:"member_count"`
	MaxPresences                int               `json:"max_presences"`
	MaxMembers                  int               `json:"max_members"`
	VanityUrlCode               string            `json:"vanity_url_code"`
	Description                 string            `json:"description"`
	Banner                      string            `json:"banner"`
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
	}
}

