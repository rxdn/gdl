package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/utils"
	"net/url"
	"strconv"
	"time"
)

func GetChannel(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetChannel, channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(ctx, token, nil, &channel); err != nil {
		return channel, err
	}

	return channel, nil
}

type ModifyChannelData struct {
	Name                 string                        `json:"name,omitempty"`
	Position             int                           `json:"position,omitempty"`
	Topic                string                        `json:"topic,omitempty"`
	Nsfw                 bool                          `json:"nsfw,omitempty"`
	RateLimitPerUser     int                           `json:"rate_limit_per_user,omitempty"`
	Bitrate              int                           `json:"bitrate,omitempty"`
	UserLimit            int                           `json:"user_limit,omitempty"`
	PermissionOverwrites []channel.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentId             uint64                        `json:"parent_id,string,omitempty"`
	*ThreadMetadataModifyData
}

type ThreadMetadataModifyData struct {
	Archived            *bool   `json:"archived,omitempty"`
	AutoArchiveDuration *uint16 `json:"auto_archive_duration,omitempty"`
	Locked              *bool   `json:"locked,omitempty"`
	Invitable           *bool   `json:"invitable,omitempty"`
}

func ModifyChannel(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data ModifyChannelData) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteModifyChannel, channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(ctx, token, data, &channel); err != nil {
		return channel, err
	}

	return channel, nil
}

func DeleteChannel(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteChannel, channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(ctx, token, nil, &channel); err != nil {
		return channel, err
	}

	return channel, nil
}

// The before, after, and around keys are mutually exclusive, only one may be passed at a time.
type GetChannelMessagesData struct {
	Around uint64 // get messages around this message ID
	Before uint64 // get messages before this message ID
	After  uint64 // get messages after this message ID
	Limit  int    // 1 - 100
}

func (o *GetChannelMessagesData) Query() string {
	query := url.Values{}

	if o.Around != 0 {
		query.Set("around", strconv.FormatUint(o.Around, 10))
	}

	if o.Before != 0 {
		query.Set("before", strconv.FormatUint(o.Before, 10))
	}

	if o.After != 0 {
		query.Set("after", strconv.FormatUint(o.After, 10))
	}

	if o.Limit > 100 || o.Limit < 1 {
		o.Limit = 50
	}
	query.Set("limit", strconv.Itoa(o.Limit))

	return query.Encode()
}

func GetChannelMessages(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data GetChannelMessagesData) ([]message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages?%s", channelId, data.Query()),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetChannelMessages, channelId),
		RateLimiter: rateLimiter,
	}

	var messages []message.Message
	if err, _ := endpoint.Request(ctx, token, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetChannelMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) (message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetChannelMessage, channelId),
		RateLimiter: rateLimiter,
	}

	var message message.Message
	if err, _ := endpoint.Request(ctx, token, nil, &message); err != nil {
		return message, err
	}

	return message, nil
}

type CreateMessageData struct {
	Content          string                    `json:"content"`
	Nonce            string                    `json:"nonce,omitempty"`
	Tts              bool                      `json:"tts,omitempty"`
	Embeds           []*embed.Embed            `json:"embeds,omitempty"`
	Flags            uint                      `json:"flags,omitempty"`
	AllowedMentions  message.AllowedMention    `json:"allowed_mentions"`
	MessageReference *message.MessageReference `json:"message_reference,omitempty"`
	Components       []component.Component     `json:"components,omitempty"`
	StickerIds       []uint64                  `json:"sticker_ids,omitempty"`
	Attachments      []request.Attachment      `json:"attachments,omitempty"`
}

func (d CreateMessageData) GetAttachments() []request.Attachment {
	return d.Attachments
}

func CreateMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data CreateMessageData) (message.Message, error) {
	var endpoint request.Endpoint
	if len(data.Attachments) == 0 {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
			Route:       ratelimit.NewChannelRoute(ratelimit.RouteCreateMessage, channelId),
			RateLimiter: rateLimiter,
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
			Route:       ratelimit.NewChannelRoute(ratelimit.RouteCreateMessage, channelId),
			RateLimiter: rateLimiter,
		}
	}

	var message message.Message
	if err, _ := endpoint.Request(ctx, token, data, &message); err != nil {
		return message, err
	}

	return message, nil
}

// emoji is the raw unicode emoji
func CreateReaction(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/@me", channelId, messageId, url.QueryEscape(emoji)),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteCreateReaction, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteOwnReaction(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/@me", channelId, messageId, url.QueryEscape(emoji)),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteOwnReaction, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteUserReaction(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId, userId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/%d", channelId, messageId, url.QueryEscape(emoji), userId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteUserReaction, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

type GetReactionsData struct {
	Before uint64 // get users before this user ID
	After  uint64 // get users after this user ID
	Limit  int    // 1 - 100
}

func (o *GetReactionsData) Query() string {
	query := url.Values{}

	if o.Before != 0 {
		query.Set("before", strconv.FormatUint(o.Before, 10))
	}

	if o.After != 0 {
		query.Set("after", strconv.FormatUint(o.After, 10))
	}

	if o.Limit > 100 || o.Limit < 1 {
		o.Limit = 25
	}
	query.Set("limit", strconv.Itoa(o.Limit))

	return query.Encode()
}

func GetReactions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string, data GetReactionsData) ([]user.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s?%s", channelId, messageId, url.QueryEscape(emoji), data.Query()),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetReactions, channelId),
		RateLimiter: rateLimiter,
	}

	var users []user.User
	if err, _ := endpoint.Request(ctx, token, nil, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteAllReactions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteAllReactions, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func DeleteAllReactionsEmoji(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s", channelId, messageId, url.QueryEscape(emoji)),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteAllReactionsForEmoji, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

type EditMessageData struct {
	Content    string                `json:"content"`
	Embeds     []*embed.Embed        `json:"embeds"`
	Flags      uint                  `json:"flags"` // https://discord.com/developers/docs/resources/channel#message-object-message-flags TODO: Helper function
	Components []component.Component `json:"components"`
}

func EditMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, data EditMessageData) (message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteEditMessage, channelId),
		RateLimiter: rateLimiter,
	}

	var message message.Message
	if err, _ := endpoint.Request(ctx, token, data, &message); err != nil {
		return message, err
	}

	return message, nil
}

func DeleteMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteMessage, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func BulkDeleteMessages(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, messages []uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/bulk-delete", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteBulkDeleteMessages, channelId),
		RateLimiter: rateLimiter,
	}

	body := map[string]interface{}{
		"messages": utils.Uint64StringSlice(messages),
	}

	err, _ := endpoint.Request(ctx, token, body, nil)
	return err
}

func EditChannelPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, updated channel.PermissionOverwrite) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, updated.Id),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteEditChannelPermissions, channelId),
		RateLimiter: rateLimiter,
	}

	updated.Id = 0

	err, _ := endpoint.Request(ctx, token, updated, nil)
	return err
}

func GetChannelInvites(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]invite.InviteMetadata, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetChannelInvites, channelId),
		RateLimiter: rateLimiter,
	}

	var invites []invite.InviteMetadata
	if err, _ := endpoint.Request(ctx, token, nil, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

type CreateInviteData struct {
	MaxAge         int    `json:"max_age"`  // seconds, 0 = never
	MaxUses        int    `json:"max_uses"` // 0 = unlimited
	Temporary      bool   `json:"temporary"`
	Unique         bool   `json:"unique"`
	TargetUser     uint64 `json:"target_user,string,omitempty"`
	TargetUserType int    `json:"target_user_type,omitempty"`
}

func CreateChannelInvite(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data CreateInviteData) (invite.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteCreateChannelInvite, channelId),
		RateLimiter: rateLimiter,
	}

	var invite invite.Invite
	if err, _ := endpoint.Request(ctx, token, data, &invite); err != nil {
		return invite, err
	}

	return invite, nil
}

func DeleteChannelPermissions(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, overwriteId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, overwriteId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeleteChannelPermission, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func TriggerTypingIndicator(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/typing", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteTriggerTypingIndicator, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func GetPinnedMessages(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetPinnedMessages, channelId),
		RateLimiter: rateLimiter,
	}

	var messages []message.Message
	if err, _ := endpoint.Request(ctx, token, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func AddPinnedChannelMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteAddPinnedChannelMessage, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func DeletePinnedChannelMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteDeletePinnedChannelMessage, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}

func JoinThread(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members/@me", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteJoinThread, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

func AddThreadMember(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, userId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members/%d", channelId, userId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteAddThreadMember, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

func LeaveThread(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members/@me", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteLeaveThread, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

func RemoveThreadMember(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, userId uint64) (err error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members/%d", channelId, userId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteRemoveThreadMember, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, nil)
	return
}

func GetThreadMember(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, userId uint64) (member channel.ThreadMember, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members/%d", channelId, userId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetThreadMember, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &member)
	return
}

func ListThreadMembers(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (members []channel.ThreadMember, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/thread-members", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteListThreadMembers, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &members)
	return
}

type StartThreadWithMessageData struct {
	Name                string `json:"name"`
	AutoArchiveDuration uint16 `json:"auto_archive_duration"`
}

func StartThreadWithMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, data StartThreadWithMessageData) (ch channel.Channel, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/threads", channelId, messageId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteStartThreadWithMessage, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &ch)
	return
}

type StartThreadWithoutMessageData struct {
	Name                string              `json:"name"`
	AutoArchiveDuration uint16              `json:"auto_archive_duration"`
	Type                channel.ChannelType `json:"type"`
	Invitable           bool                `json:"invitable"`
}

func StartThreadWithoutMessage(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data StartThreadWithoutMessageData) (ch channel.Channel, err error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/threads", channelId),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteStartThreadWithoutMessage, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, data, &ch)
	return
}

type ThreadsResponse struct {
	Threads []channel.Channel      `json:"threads"`
	Members []channel.ThreadMember `json:"members"`
	HasMore bool                   `json:"has_more"`
}

func ListActiveThreads(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) (threads ThreadsResponse, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/threads/active", guildId),
		Route:       ratelimit.NewGuildRoute(ratelimit.RouteGetActiveThreads, guildId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &threads)
	return
}

type ListThreadsData struct {
	Before time.Time
	Limit  int
}

func (d *ListThreadsData) Query() string {
	query := url.Values{}

	if !d.Before.IsZero() {
		query.Set("before", d.Before.String())
	}

	if d.Limit > 0 {
		query.Set("limit", strconv.Itoa(d.Limit))
	}

	return query.Encode()
}

func ListPublicArchivedThreads(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data ListThreadsData) (threads ThreadsResponse, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/threads/archived/public?%s", channelId, data.Query()),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetArchivedPublicThreads, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &threads)
	return
}

func ListPrivateArchivedThreads(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data ListThreadsData) (threads ThreadsResponse, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/threads/archived/private?%s", channelId, data.Query()),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetArchivedPrivateThreads, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &threads)
	return
}

func ListJoinedPrivateArchivedThreads(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data ListThreadsData) (threads ThreadsResponse, err error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/usrs/@methreads/archived/private?%s", channelId, data.Query()),
		Route:       ratelimit.NewChannelRoute(ratelimit.RouteGetArchivedPrivateThreads, channelId),
		RateLimiter: rateLimiter,
	}

	err, _ = endpoint.Request(ctx, token, nil, &threads)
	return
}
