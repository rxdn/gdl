package guild

import (
	"fmt"
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"strings"
	"time"
)

type Guild struct {
	Id                          uint64                    `json:"id,string"`
	Name                        string                    `json:"name"`
	Icon                        string                    `json:"icon"`
	Splash                      string                    `json:"splash"`
	Owner                       bool                      `json:"owner"`
	OwnerId                     uint64                    `json:"owner_id,string"`
	Permissions                 uint64                    `json:"permissions,string"`
	Region                      string                    `json:"region"`
	AfkChannelId                objects.NullableSnowflake `json:"afk_channel_id"`
	AfkTimeout                  int                       `json:"afk_timeout"`
	VerificationLevel           int                       `json:"verification_level"`
	DefaultMessageNotifications int                       `json:"default_message_notifications"`
	ExplicitContentFilter       int                       `json:"explicit_content_filter"`
	Roles                       []Role                    `json:"roles"`
	Emojis                      []emoji.Emoji             `json:"emojis"`
	Features                    []GuildFeature            `json:"features"`
	MfaLevel                    int                       `json:"mfa_level"`
	ApplicationId               objects.NullableSnowflake `json:"application_id"`
	WidgetEnabled               bool                      `json:"widget_enabled"`
	WidgetChannelId             objects.NullableSnowflake `json:"widget_channel_id"`
	SystemChannelId             objects.NullableSnowflake `json:"system_channel_id"`
	SystemChannelFlags          uint16                     `json:"system_channel_flags"`
	RulesChannelId              objects.NullableSnowflake `json:"rules_channel_id,omitempty"`
	JoinedAt                    time.Time                 `json:"joined_at"`
	Large                       bool                      `json:"large"`
	Unavailable                 *bool                     `json:"unavailable"`
	MemberCount                 int                       `json:"member_count"`
	VoiceStates                 []VoiceState              `json:"voice_state"`
	Members                     []member.Member           `json:"members"`
	Channels                    []channel.Channel         `json:"channels"`
	Threads                     []channel.Channel         `json:"threads"`
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
	ApproximateMemberCount      int                       `json:"approximate_member_count"`   // Returned on GET /guild/:id
	ApproximatePresenceCount    int                       `json:"approximate_presence_count"` // Returned on GET /guild/:id
	WelcomeScreen               WelcomeScreen             `json:"welcome_screen"`
	Nsfw                        bool                      `json:"nsfw"`
}

func (g *Guild) IconUrl() string {
	if g.Icon == "" {
		return ""
	}

	extension := "png"
	if strings.HasPrefix(g.Icon, "a_") {
		extension = "gif"
	}

	return fmt.Sprintf("https://cdn.discordapp.com/icons/%d/%s.%s", g.Id, g.Icon, extension)
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
		VerificationLevel:           g.VerificationLevel,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		ExplicitContentFilter:       g.ExplicitContentFilter,
		Features:                    g.Features,
		MfaLevel:                    g.MfaLevel,
		ApplicationId:               g.ApplicationId,
		WidgetEnabled:               g.WidgetEnabled,
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

	for _, role := range g.Roles {
		cached.Roles = append(cached.Roles, role.Id)
	}

	for _, emoji := range g.Emojis {
		cached.Emojis = append(cached.Emojis, emoji.Id.Value)
	}

	for _, channel := range g.Channels {
		cached.Channels = append(cached.Channels, channel.Id)
	}

	return
}
