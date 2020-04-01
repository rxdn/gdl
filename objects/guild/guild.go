package guild

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"time"
)

type ExplicitContentFilterLevel int

const (
	DISABLED              ExplicitContentFilterLevel = 0
	MEMBERS_WITHOUT_ROLES ExplicitContentFilterLevel = 1
	ALL_MEMBERS           ExplicitContentFilterLevel = 2
)

type DefaultMessageNotificationLevel int

const (
	DefaultMessageNotificationLevelAllMessages  DefaultMessageNotificationLevel = 0
	DefaultMessageNotificationLevelOnlyMengions DefaultMessageNotificationLevel = 1
)

type GuildFeature string

const (
	GuildFeatureInviteSplash   GuildFeature = "INVITE_SPLASH"
	GuildFeatureVipRegions     GuildFeature = "VIP_REGIONS" // guild has access to set 384kbps bitrate in voice (previously VIP voice servers)
	GuildFeatureVanityUrl      GuildFeature = "VANITY_URL"
	GuildFeatureVerified       GuildFeature = "VERIFIED"
	GuildFeaturePartnered      GuildFeature = "PARTNERED"
	GuildFeaturePublic         GuildFeature = "PUBLIC"
	GuildFeatureCommerce       GuildFeature = "COMMERCE"
	GuildFeatureNews           GuildFeature = "NEWS"
	GuildFeatureDiscoverable   GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable     GuildFeature = "FEATURABLE"
	GuildFeatureBanner         GuildFeature = "BANNER"
	GuildFeaturePublicDisabled GuildFeature = "PUBLIC_DISABLED"
)

type Guild struct {
	Id                          uint64             `json:"id,string"`
	Name                        string             `json:"name"`
	Icon                        string             `json:"icon"`
	Splash                      string             `json:"splash"`
	Owner                       bool               `json:"owner"`
	OwnerId                     uint64             `json:"owner_id,string"`
	Permissions                 int                `json:"permissions"`
	Region                      string             `json:"region"`
	AfkChannelId                uint64             `json:"afk_channel_id,string"`
	AfkTimeout                  int                `json:"afk_timeout"`
	EmbedEnabled                bool               `json:"embed_enabled"`
	EmbedChannelId              uint64             `json:"embed_channel_id,string"`
	VerificationLevel           int                `json:"verification_level"`
	DefaultMessageNotifications int                `json:"default_message_notifications"`
	ExplicitContentFilter       int                `json:"explicit_content_filter"`
	Roles                       []*Role            `json:"roles"`
	Emojis                      []*emoji.Emoji     `json:"emojis"`
	Features                    []string           `json:"features"`
	MfaLevel                    int                `json:"mfa_level"`
	ApplicationId               uint64             `json:"application_id,string"`
	WidgetEnabled               bool               `json:"widget_enabled"`
	WidgetChannelId             uint64             `json:"widget_channel_id,string"`
	SystemChannelId             uint64             `json:"system_channel_id,string"`
	JoinedAt                    time.Time          `json:"joined_at"`
	Large                       bool               `json:"large"`
	Unavailable                 *bool              `json:"unavailable"`
	MemberCount                 int                `json:"member_count"`
	VoiceStates                 []*VoiceState      `json:"voice_state"`
	Members                     []*member.Member   `json:"members"`
	Channels                    []*channel.Channel `json:"channels"`
	Presences                   []*user.Presence   `json:"presences"`
	MaxPresences                int                `json:"max_presences"`
	MaxMembers                  int                `json:"max_members"`
	VanityUrlCode               string             `json:"vanity_url_code"`
	Description                 string             `json:"description"`
	Banner                      string             `json:"banner"`
}
