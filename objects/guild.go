package objects

import "time"

type Guild struct {
	Id                          uint64        `json:"id,string"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon"`
	Splash                      string        `json:"splash"`
	Owner                       bool          `json:"owner"`
	OwnerId                     uint64        `json:"owner_id,string"`
	Permissions                 int           `json:"permissions"`
	Region                      string        `json:"region"`
	AfkChannelId                uint64        `json:"afk_channel_id,string"`
	AfkTimeout                  int           `json:"afk_timeout"`
	EmbedEnabled                bool          `json:"embed_enabled"`
	EmbedChannelId              uint64        `json:"embed_channel_id,string"`
	VerificationLevel           int           `json:"verification_level"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Roles                       []*Role       `json:"roles"`
	Emojis                      []*Emoji      `json:"emojis"`
	Features                    []string      `json:"features"`
	MfaLevel                    int           `json:"mfa_level"`
	ApplicationId               uint64        `json:"application_id,string"`
	WidgetEnabled               bool          `json:"widget_enabled"`
	WidgetChannelId             uint64        `json:"widget_channel_id,string"`
	SystemChannelId             uint64        `json:"system_channel_id,string"`
	JoinedAt                    time.Time     `json:"joined_at"`
	Large                       bool          `json:"large"`
	Unavailable                 bool          `json:"unavailable"`
	MemberCount                 int           `json:"member_count"`
	VoiceStates                 []*VoiceState `json:"voice_state"`
	Members                     []*Member     `json:"members"`
	Channels                    []*Channel    `json:"channels"`
	Presences                   []*Presence   `json:"presences"`
	MaxPresences                int           `json:"max_presences"`
	MaxMembers                  int           `json:"max_members"`
	VanityUrlCode               string        `json:"vanity_url_code"`
	Description                 string        `json:"description"`
	Banner                      string        `json:"banner"`
}
