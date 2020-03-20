package objects

import (
	"fmt"
	"github.com/Dot-Rar/gdl/utils"
	"github.com/fatih/structs"
	"reflect"
)

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
	JoinedAt                    uint64        `json:"joined_at,string"`
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

func (g *Guild) KeyName() string {
	return fmt.Sprintf("cache:guild:%s", g.Id)
}

func (g *Guild) Serialize() map[string]map[string]interface{} {
	fields := make(map[string]map[string]interface{})
	utils.Initialise(fields, g.KeyName())

	for k, v := range structs.Map(g) {
		if k == "Roles" {
			if len(g.Roles) == 0 {
				continue
			}

			var ids []string
			for _, role := range g.Roles {
				fields = utils.Append(fields, role.Serialize())
				ids = append(ids, role.Id)
			}
			fields[g.KeyName()][k] = ids
		} else if k == "Emojis" {
			if len(g.Emojis) == 0 {
				continue
			}

			var ids []string
			for _, emoji := range g.Emojis {
				fields = utils.Append(fields, emoji.Serialize())
				ids = append(ids, emoji.Id)
			}
			fields[g.KeyName()][k] = ids
		} else if k == "Features" {
			fields[g.KeyName()][k] = g.Features
		} else if k == "VoiceStates" {
			if len(g.VoiceStates) == 0 {
				continue
			}

			var ids []string
			for _, state := range g.VoiceStates {
				fields = utils.Append(fields, state.Serialize())
				ids = append(ids, state.UserId)
			}
			fields[g.KeyName()][k] = ids
		} else if k == "Members" {
			if len(g.Members) == 0 { // Should only happen when lazy loading is used
				continue
			}

			var ids []string
			for _, member := range g.Members {
				fields = utils.Append(fields, member.Serialize(g.Id))
				ids = append(ids, member.User.Id)
			}

			fields[g.KeyName()][k] = ids
		} else if k == "Channels" {
			if len(g.Members) == 0 {
				continue
			}

			var ids []string
			for _, channel := range g.Channels {
				fields = utils.Append(fields, channel.Serialize())
				ids = append(ids, channel.Id)
			}

			fields[g.KeyName()][k] = ids
		} else {
			if !utils.IsZero(reflect.ValueOf(v)) {
				fields[g.KeyName()][k] = v
			}
		}
	}

	return fields
}
