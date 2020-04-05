package rest

import (
	"bytes"
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/utils"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

func GetChannel(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(token, nil, &channel); err != nil {
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
}

func ModifyChannel(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data ModifyChannelData) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(token, data, &channel); err != nil {
		return channel, err
	}

	return channel, nil
}

func DeleteChannel(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) (channel.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var channel channel.Channel
	if err, _ := endpoint.Request(token, nil, &channel); err != nil {
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

func GetChannelMessages(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data GetChannelMessagesData) ([]message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages?%s", channelId, data.Query()),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var messages []message.Message
	if err, _ := endpoint.Request(token, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetChannelMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) (message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var message message.Message
	if err, _ := endpoint.Request(token, nil, &message); err != nil {
		return message, err
	}

	return message, nil
}

type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

type CreateMessageData struct {
	Content         string                 `json:"content"`
	Nonce           string                 `json:"nonce,omitempty"`
	Tts             bool                   `json:"tts,omitempty"`
	File            *File                  `json:"-"`
	Embed           *embed.Embed           `json:"embed,omitempty"`
	PayloadJson     string                 `json:"-"` // TODO: Helper method
	AllowedMentions message.AllowedMention `json:"allowed_mentions"`
}

func (d CreateMessageData) EncodeMultipartFormData() ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if d.File != nil {
		fileName := d.File.Name
		fileName = strings.Replace(fileName, "\\", "\\\\", -1)
		fileName = strings.Replace(fileName, "\"", "\\\"", -1)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
		h.Set("Content-Type", d.File.ContentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			return body.Bytes(), writer.Boundary(), err
		}

		if _, err := io.Copy(part, d.File.Reader); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if d.Content != "" {
		if err := writer.WriteField("content", d.Content); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if err := writer.WriteField("tts", strconv.FormatBool(d.Tts)); err != nil {
		return body.Bytes(), writer.Boundary(), err
	}

	if d.Nonce != "" {
		if err := writer.WriteField("nonce", d.Nonce); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	if d.PayloadJson != "" {
		if err := writer.WriteField("payload_json", d.PayloadJson); err != nil {
			return body.Bytes(), writer.Boundary(), err
		}
	}

	return []byte(string(body.Bytes()) + "\r\n--" + writer.Boundary() + "--"), writer.Boundary(), nil
}

func CreateMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data CreateMessageData) (message.Message, error) {
	var endpoint request.Endpoint
	if data.File == nil {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
			BaseRoute:   ratelimit.NewChannelRoute(channelId),
			RateLimiter: rateLimiter,
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
		}
	}

	var message message.Message
	if err, _ := endpoint.Request(token, data, &message); err != nil {
		return message, err
	}

	return message, nil
}

// emoji is the raw unicode emoji
func CreateReaction(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/@me", channelId, messageId, url.QueryEscape(emoji)),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteOwnReaction(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/@me", channelId, messageId, url.QueryEscape(emoji)),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteUserReaction(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId, userId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/%d", channelId, messageId, url.QueryEscape(emoji), userId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
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

func GetReactions(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string, data GetReactionsData) ([]user.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s?%s", channelId, messageId, url.QueryEscape(emoji), data.Query()),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var users []user.User
	if err, _ := endpoint.Request(token, nil, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteAllReactions(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func DeleteAllReactionsEmoji(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, emoji string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s", channelId, messageId, url.QueryEscape(emoji)),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

type EditMessageData struct {
	Content string       `json:"content,omitempty"`
	Embed   *embed.Embed `json:"embed,omitempty"`
	Flags   int          `json:"flags,omitempty"` // https://discordapp.com/developers/docs/resources/channel#message-object-message-flags TODO: Helper function
}

func EditMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64, data ModifyChannelData) (message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var message message.Message
	if err, _ := endpoint.Request(token, data, &message); err != nil {
		return message, err
	}

	return message, nil
}

func DeleteMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func BulkDeleteMessages(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, messages []uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/bulk-delete", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	body := map[string]interface{}{
		"messages": utils.Uint64StringSlice(messages),
	}

	err, _ := endpoint.Request(token, body, nil)
	return err
}

func EditChannelPermissions(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, updated channel.PermissionOverwrite) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, updated.Id),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	updated.Id = 0

	err, _ := endpoint.Request(token, updated, nil)
	return err
}

func GetChannelInvites(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]invite.InviteMetadata, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var invites []invite.InviteMetadata
	if err, _ := endpoint.Request(token, nil, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

func CreateChannelInvite(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64, data invite.InviteMetadata) (invite.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var invite invite.Invite
	if err, _ := endpoint.Request(token, data, &invite); err != nil {
		return invite, err
	}

	return invite, nil
}

func DeleteChannelPermissions(token string, rateLimiter *ratelimit.Ratelimiter, channelId, overwriteId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, overwriteId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func TriggerTypingIndicator(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/typing", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func GetPinnedMessages(token string, rateLimiter *ratelimit.Ratelimiter, channelId uint64) ([]message.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins", channelId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	var messages []message.Message
	if err, _ := endpoint.Request(token, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func AddPinnedChannelMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}

func DeletePinnedChannelMessage(token string, rateLimiter *ratelimit.Ratelimiter, channelId, messageId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
		BaseRoute:   ratelimit.NewChannelRoute(channelId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}
