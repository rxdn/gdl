package rest

import (
	"bytes"
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/rest/routes"
	"github.com/rxdn/gdl/utils"
	"image"
	"image/png"
	"net/url"
	"strconv"
)

type CreateGuildData struct {
	Name                        string                                `json:"name"`
	Region                      string                                `json:"region"` // voice region ID TODO: Helper function
	Icon                        string                                `json:"icon"`
	VerificationLevel           guild.VerificationLevel               `json:"verification_level"`
	DefaultMessageNotifications guild.DefaultMessageNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       guild.ExplicitContentFilterLevel      `json:"explicit_content_filter"`
	Roles                       []*guild.Role                         `json:"roles"`    // only @everyone
	Channels                    []*channel.Channel                    `json:"channels"` // channels cannot have a ParentId
	AfkChannelId                uint64                                `json:"afk_channel_id,string"`
	AfkTimeout                  int                                   `json:"afk_timeout"`
	SystemChannelId             uint64                                `json:"system_channel_id"`
}

// only available to bots in < 10 guilds
func CreateGuild(token string, data CreateGuildData) (*guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    "/guilds",
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, nil, data, &guild)
	return &guild, err
}

func GetGuild(token string, guildId uint64) (*guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d", guildId),
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &guild)
	return &guild, err
}

func GetGuildPreview(token string, guildId uint64) (*guild.GuildPreview, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/preview", guildId),
	}

	var preview guild.GuildPreview
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &preview)
	return &preview, err
}

type ModifyGuildData struct {
	Name                        string                                `json:"name"`
	Region                      string                                `json:"region"` // voice region ID TODO: Helper function
	VerificationLevel           guild.VerificationLevel               `json:"verification_level"`
	DefaultMessageNotifications guild.DefaultMessageNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       guild.ExplicitContentFilterLevel      `json:"explicit_content_filter"`
	AfkChannelId                uint64                                `json:"afk_channel_id,string"`
	AfkTimeout                  int                                   `json:"afk_timeout"`
	Icon                        string                                `json:"icon"`
	OwnerId                     uint64                                `json:"owner_id"`
	Splash                      string                                `json:"splash"`
	Banner                      string                                `json:"banner"`
	SystemChannelId             uint64                                `json:"system_channel_id"`
	RulesChannelId              uint64                                `json:"rules_channel_id"`
	PublicUpdatesChannelId      uint64                                `json:"public_updates_channel_id"`
	PreferredLocale             string                                `json:"preferred_locale"`
}

func ModifyGuild(token string, guildId uint64, data ModifyGuildData) (*guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d", guildId),
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, &guild)
	return &guild, err
}

func DeleteGuild(token string, guildId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d", guildId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildChannels(token string, guildId uint64) ([]*channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
	}

	var channels []*channel.Channel
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &channels)
	return channels, err
}

type CreateChannelData struct {
	Name                 string                         `json:"name"`
	Type                 channel.ChannelType            `json:"type"`
	Topic                string                         `json:"topic"`
	Bitrate              int                            `json:"bitrate"`
	UserLimit            int                            `json:"user_limit"`
	RateLimitPerUser     int                            `json:"rate_limit_per_user"`
	Position             int                            `json:"position"`
	PermissionOverwrites []*channel.PermissionOverwrite `json:"permission_overwrites"`
	ParentId             uint64                         `json:"parent_id,string"`
	Nsfw                 bool                           `json:"nsfw"`
}

func CreateGuildChannel(token string, guildId uint64, data CreateChannelData) (*channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
	}

	var channel channel.Channel
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, &channel)
	return &channel, err
}

type Position struct {
	ChannelId uint64 `json:"id,string"`
	Position  int    `json:"position"`
}

func ModifyGuildChannelPositions(token string, guildId uint64, positions []Position) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, positions, nil)
	return err
}

func GetGuildMember(token string, guildId, userId uint64) (*member.Member, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
	}

	var member member.Member
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &member)
	return &member, err
}

// all parameters are optional
type ListGuildMembersData struct {
	Limit int    // 1 - 1000
	After uint64 // Highest user ID in the previous page
}

func (d *ListGuildMembersData) Query() string {
	query := url.Values{}

	if d.Limit != 0 {
		query.Set("limit", strconv.Itoa(d.Limit))
	}

	if d.After != 0 {
		query.Set("after", strconv.FormatUint(d.After, 10))
	}

	return query.Encode()
}

func ListGuildMembers(token string, guildId uint64, data ListGuildMembersData) ([]*member.Member, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/members?%s", guildId, data.Query()),
	}

	var members []*member.Member
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &members)
	return members, err
}

type ModifyGuildMemberData struct {
	Nick      string                   `json:"nick,omitempty"`
	Roles     *utils.Uint64StringSlice `json:"roles,omitempty"`
	Mute      *bool                    `json:"mute,omitempty"`
	Deaf      *bool                    `json:"deaf,omitempty"`
	ChannelId uint64                   `json:"channel_id,string,omitempty"` // id of channel to move user to (if they are connected to voice)
}

func ModifyGuildMember(token string, guildId, userId uint64, data ModifyGuildMemberData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, nil)
	return err
}

func ModifyCurrentUserNick(token string, guildId uint64, nick string) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/@me/nick", guildId),
	}

	data := map[string]interface{}{
		"nick": nick,
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, nil)
	return err
}

func AddGuildMemberRole(token string, guildId, userId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d/roles/%d", guildId, userId, roleId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func RemoveGuildMemberRole(token string, guildId, userId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d/roles/%d", guildId, userId, roleId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func RemoveGuildMember(token string, guildId, userId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildBans(token string, guildId uint64) ([]*guild.Ban, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans", guildId),
	}

	var bans []*guild.Ban
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &bans)
	return bans, err
}

func GetGuildBan(token string, guildId, userId uint64) (*guild.Ban, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
	}

	var ban guild.Ban
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &ban)
	return &ban, err
}

type CreateGuildBanData struct {
	DeleteMessageDays int    `json:"delete-message-days,omitempty"` // 1 - 7
	Reason            string `json:"reason,omitempty"`
}

func CreateGuildBan(token string, guildId, userId uint64, data CreateGuildBanData) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, nil)
	return err
}

func RemoveGuildBan(token string, guildId, userId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildRoles(token string, guildId uint64) ([]*guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
	}

	var roles []*guild.Role
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &roles)
	return roles, err
}

type GuildRoleData struct {
	Name        string `json:"name,omitempty"`
	Permissions *int   `json:"permissions,omitempty"`
	Color       *int   `json:"color,omitempty"`
	Hoist       *bool  `json:"hoist,omitempty"`
	Mentionable *bool  `json:"mentionable,omitempty"`
}

func CreateGuildRole(token string, guildId uint64, data GuildRoleData) (*guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
	}

	var role *guild.Role
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, &role)
	return role, err
}

func ModifyGuildRolePositions(token string, guildId uint64, positions []Position) ([]*guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
	}

	var roles []*guild.Role
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, positions, &roles)
	return roles, err
}

func ModifyGuildRole(token string, guildId, roleId uint64, data GuildRoleData) (*guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles/%d", guildId, roleId),
	}

	var role *guild.Role
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, &role)
	return role, err
}

func DeleteGuildRole(token string, guildId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles/%d", guildId, roleId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildPruneCount(token string, guildId uint64, days int) (int, error) {
	if days < 1 {
		days = 7
	}

	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/prune?days=%d", guildId, days),
	}

	var res map[string]int
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &res)
	return res["pruned"], err
}

// computePruneCount = whether 'pruned' is returned, discouraged for large guilds
func BeginGuildPrune(token string, guildId uint64, days int, computePruneCount bool) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/prune?days=%d&compute_prune_count=%t", guildId, days, computePruneCount),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildVoiceRegions(token string, guildId uint64) ([]*guild.VoiceRegion, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/regions", guildId),
	}

	var regions []*guild.VoiceRegion
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &regions)
	return regions, err
}

func GetGuildInvites(token string, guildId uint64) ([]*invite.InviteMetadata, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/regions", guildId),
	}

	var invites []*invite.InviteMetadata
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &invites)
	return invites, err
}

func GetGuildIntegrations(token string, guildId uint64) ([]*integration.Integration, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations", guildId),
	}

	var integrations []*integration.Integration
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &integrations)
	return integrations, err
}

type CreateIntegrationData struct {
	Type string
	Id   uint64 `json:"id,string"`
}

func CreateGuildIntegration(token string, guildId uint64, data CreateIntegrationData) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations", guildId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, nil)
	return err
}

type ModifyIntegrationData struct {
	ExpireBehaviour   integration.IntegrationExpireBehaviour `json:"expire_behavior"`
	ExpireGracePeriod int                                    `json:"expire_grace_period"`
	EnableEmoticons   bool                                   `json:"enable_emoticons"`
}

func ModifyGuildIntegration(token string, guildId, integrationId uint64, data ModifyIntegrationData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d", guildId, integrationId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, nil)
	return err
}

func DeleteGuildIntegration(token string, guildId, integrationId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d", guildId, integrationId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func SyncGuildIntegration(token string, guildId, integrationId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d/sync", guildId, integrationId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	return err
}

func GetGuildEmbed(token string, guildId uint64) (*guild.GuildEmbed, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/embed", guildId),
	}

	var embed guild.GuildEmbed
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &embed)
	return &embed, err
}

func ModifyGuildEmbed(token string, guildId uint64, data guild.GuildEmbed) (*guild.GuildEmbed, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/embed", guildId),
	}

	var embed guild.GuildEmbed
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, data, &embed)
	return &embed, err
}

// returns invite object with only "code" and "uses" fields
func GetGuildVanityURL(token string, guildId uint64) (*invite.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/vanity-url", guildId),
	}

	var invite invite.Invite
	err, _ := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, &invite)
	return &invite, err
}

func GetGuildWidgetImage(token string, guildId uint64, style guild.WidgetStyle) (*image.Image, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/widget.png?style=%s", guildId, string(style)),
	}

	err, res := endpoint.Request(token, &routes.RouteManager.GetGuildRoute(guildId).Ratelimiter, nil, nil)
	if err != nil {
		return nil, err
	}

	image, err := png.Decode(bytes.NewReader(res.Content))
	if err != nil {
		return nil, err
	}

	return &image, err
}
