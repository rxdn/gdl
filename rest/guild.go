package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/utils"
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
func CreateGuild(token string, data CreateGuildData) (guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    "/guilds",
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuild, 0),
		RateLimiter: nil,
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, data, &guild)
	return guild, err
}

func GetGuild(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) (guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d?with_counts=true", guildId), // TODO: Allow users to specify whether they want with_counts
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuild, guildId),
		RateLimiter: rateLimiter,
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, nil, &guild)
	return guild, err
}

func GetGuildPreview(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) (guild.GuildPreview, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/preview", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildPreview, guildId),
		RateLimiter: rateLimiter,
	}

	var preview guild.GuildPreview
	err, _ := endpoint.Request(token, nil, &preview)
	return preview, err
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

func ModifyGuild(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data ModifyGuildData) (guild.Guild, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuild, guildId),
		RateLimiter: rateLimiter,
	}

	var guild guild.Guild
	err, _ := endpoint.Request(token, data, &guild)
	return guild, err
}

func DeleteGuild(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGuild, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildChannels(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildChannels, guildId),
		RateLimiter: rateLimiter,
	}

	var channels []channel.Channel
	err, _ := endpoint.Request(token, nil, &channels)
	return channels, err
}

type CreateChannelData struct {
	Name                 string                        `json:"name"`
	Type                 channel.ChannelType           `json:"type"`
	Topic                string                        `json:"topic,omitempty"`
	Bitrate              int                           `json:"bitrate,omitempty"`
	UserLimit            int                           `json:"user_limit,omitempty"`
	RateLimitPerUser     int                           `json:"rate_limit_per_user"`
	Position             int                           `json:"position,omitempty"`
	PermissionOverwrites []channel.PermissionOverwrite `json:"permission_overwrites"`
	ParentId             uint64                        `json:"parent_id,string,omitempty"`
	Nsfw                 bool                          `json:"nsfw,omitempty"`
}

func CreateGuildChannel(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data CreateChannelData) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildChannel, guildId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	err, _ := endpoint.Request(token, data, &channel)
	return channel, err
}

type Position struct {
	ChannelId uint64 `json:"id,string"`
	Position  int    `json:"position"`
}

func ModifyGuildChannelPositions(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, positions []Position) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/channels", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildChannelPositions, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, positions, nil)
	return err
}

func GetGuildMember(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64) (member.Member, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildMember, guildId),
		RateLimiter: rateLimiter,
	}

	var member member.Member
	err, _ := endpoint.Request(token, nil, &member)
	return member, err
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

func ListGuildMembers(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data ListGuildMembersData) ([]member.Member, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/members?%s", guildId, data.Query()),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteListGuildMembers, guildId),
		RateLimiter: rateLimiter,
	}

	var members []member.Member
	err, _ := endpoint.Request(token, nil, &members)
	return members, err
}

type ModifyGuildMemberData struct {
	Nick      string                   `json:"nick,omitempty"`
	Roles     *utils.Uint64StringSlice `json:"roles,omitempty"`
	Mute      *bool                    `json:"mute,omitempty"`
	Deaf      *bool                    `json:"deaf,omitempty"`
	ChannelId uint64                   `json:"channel_id,string,omitempty"` // id of channel to move user to (if they are connected to voice)
}

func ModifyGuildMember(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64, data ModifyGuildMemberData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildMember, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, data, nil)
	return err
}

func ModifyCurrentUserNick(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, nick string) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/@me/nick", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyCurrentUserNick, guildId),
		RateLimiter: rateLimiter,
	}

	data := map[string]interface{}{
		"nick": nick,
	}

	err, _ := endpoint.Request(token, data, nil)
	return err
}

func AddGuildMemberRole(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d/roles/%d", guildId, userId, roleId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteAddGuildMemberRole, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func RemoveGuildMemberRole(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d/roles/%d", guildId, userId, roleId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteRemoveGuildMemberRole, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func RemoveGuildMember(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteRemoveGuildMember, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildBans(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]guild.Ban, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildBans, guildId),
		RateLimiter: rateLimiter,
	}

	var bans []guild.Ban
	err, _ := endpoint.Request(token, nil, &bans)
	return bans, err
}

func GetGuildBan(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64) (guild.Ban, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildBan, guildId),
		RateLimiter: rateLimiter,
	}

	var ban guild.Ban
	err, _ := endpoint.Request(token, nil, &ban)
	return ban, err
}

type CreateGuildBanData struct {
	DeleteMessageDays int    `json:"delete_message_days,omitempty"` // 1 - 7
	Reason            string `json:"reason,omitempty"`
}

func CreateGuildBan(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64, data CreateGuildBanData) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildBan, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, data, nil)
	return err
}

func RemoveGuildBan(token string, rateLimiter *ratelimit.Ratelimiter, guildId, userId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/bans/%d", guildId, userId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteRemoveGuildBan, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildRoles(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildRoles, guildId),
		RateLimiter: rateLimiter,
	}

	var roles []guild.Role
	err, _ := endpoint.Request(token, nil, &roles)
	return roles, err
}

type GuildRoleData struct {
	Name        string  `json:"name,omitempty"`
	Permissions *uint64 `json:"permissions,omitempty,string"`
	Color       *int    `json:"color,omitempty"`
	Hoist       *bool   `json:"hoist,omitempty"`
	Mentionable *bool   `json:"mentionable,omitempty"`
}

func CreateGuildRole(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data GuildRoleData) (guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildRole, guildId),
		RateLimiter: rateLimiter,
	}

	var role guild.Role
	err, _ := endpoint.Request(token, data, &role)
	return role, err
}

func ModifyGuildRolePositions(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, positions []Position) ([]guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildRolePositions, guildId),
		RateLimiter: rateLimiter,
	}

	var roles []guild.Role
	err, _ := endpoint.Request(token, positions, &roles)
	return roles, err
}

func ModifyGuildRole(token string, rateLimiter *ratelimit.Ratelimiter, guildId, roleId uint64, data GuildRoleData) (guild.Role, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles/%d", guildId, roleId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildRole, guildId),
		RateLimiter: rateLimiter,
	}

	var role guild.Role
	err, _ := endpoint.Request(token, data, &role)
	return role, err
}

func DeleteGuildRole(token string, rateLimiter *ratelimit.Ratelimiter, guildId, roleId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/roles/%d", guildId, roleId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGuildRole, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildPruneCount(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, days int) (int, error) {
	if days < 1 {
		days = 7
	}

	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/prune?days=%d", guildId, days),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildPruneCount, guildId),
		RateLimiter: rateLimiter,
	}

	var res map[string]int
	err, _ := endpoint.Request(token, nil, &res)
	return res["pruned"], err
}

// computePruneCount = whether 'pruned' is returned, discouraged for large guilds
func BeginGuildPrune(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, days int, computePruneCount bool) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/prune?days=%d&compute_prune_count=%t", guildId, days, computePruneCount),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteBeginGuildPrune, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildVoiceRegions(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]guild.VoiceRegion, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/regions", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildVoiceRegions, guildId),
		RateLimiter: rateLimiter,
	}

	var regions []guild.VoiceRegion
	err, _ := endpoint.Request(token, nil, &regions)
	return regions, err
}

func GetGuildInvites(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]invite.InviteMetadata, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/invites", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildInvites, guildId),
		RateLimiter: rateLimiter,
	}

	var invites []invite.InviteMetadata
	err, _ := endpoint.Request(token, nil, &invites)
	return invites, err
}

func GetGuildIntegrations(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]integration.Integration, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildIntegrations, guildId),
		RateLimiter: rateLimiter,
	}

	var integrations []integration.Integration
	err, _ := endpoint.Request(token, nil, &integrations)
	return integrations, err
}

type CreateIntegrationData struct {
	Type string
	Id   uint64 `json:"id,string"`
}

func CreateGuildIntegration(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data CreateIntegrationData) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteCreateGuildIntegration, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, data, nil)
	return err
}

type ModifyIntegrationData struct {
	ExpireBehaviour   integration.IntegrationExpireBehaviour `json:"expire_behavior"`
	ExpireGracePeriod int                                    `json:"expire_grace_period"`
	EnableEmoticons   bool                                   `json:"enable_emoticons"`
}

func ModifyGuildIntegration(token string, rateLimiter *ratelimit.Ratelimiter, guildId, integrationId uint64, data ModifyIntegrationData) error {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d", guildId, integrationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildIntegration, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, data, nil)
	return err
}

func DeleteGuildIntegration(token string, rateLimiter *ratelimit.Ratelimiter, guildId, integrationId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d", guildId, integrationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteDeleteGuildIntegration, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func SyncGuildIntegration(token string, rateLimiter *ratelimit.Ratelimiter, guildId, integrationId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/integrations/%d/sync", guildId, integrationId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteSyncGuildIntegration, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetGuildWidget(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) (widget guild.GuildWidget, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/widget.json", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildWidget, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, &widget)
	return
}

func ModifyGuildEmbed(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data guild.GuildEmbed) (widget guild.GuildEmbed, err error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/guilds/%d/embed", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteModifyGuildWidget, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, data, &widget)
	return
}

// returns invite object with only "code" and "uses" fields
func GetGuildVanityURL(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) (invite invite.Invite, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/vanity-url", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetGuildVanityURL, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(token, nil, &invite)
	return
}
