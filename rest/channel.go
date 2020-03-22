package rest

import (
	"bytes"
	"fmt"
	"github.com/Dot-Rar/gdl/objects"
	"github.com/Dot-Rar/gdl/rest/request"
	"github.com/Dot-Rar/gdl/rest/routes"
	"github.com/Dot-Rar/gdl/utils"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

func GetChannel(channelId uint64, token string) (*objects.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
	}

	var channel objects.Channel
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

type ModifyChannelData struct {
	Name                 string                        `json:"name,omitempty"`
	Position             int                           `json:"position,omitempty"`
	Topic                string                        `json:"topic,omitempty"`
	Nsfw                 bool                          `json:"nsfw,omitempty"`
	RateLimitPerUser     int                           `json:"rate_limit_per_user,omitempty"`
	Bitrate              int                           `json:"bitrate,omitempty"`
	UserLimit            int                           `json:"user_limit,omitempty"`
	PermissionOverwrites []objects.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentId             uint64                        `json:"parent_id,string,omitempty"`
}

func ModifyChannel(channelId uint64, token string, data ModifyChannelData) (*objects.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
	}

	var channel objects.Channel
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, data, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func DeleteChannel(channelId uint64, token string) (*objects.Channel, error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d", channelId),
	}

	var channel objects.Channel
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// The before, after, and around keys are mutually exclusive, only one may be passed at a time.
type GetChannelMessagesData struct {
	Around *uint64 // get messages around this message ID
	Before *uint64 // get messages before this message ID
	After  *uint64 // get messages after this message ID
	Limit  int     // 1 - 100
}

func (o *GetChannelMessagesData) Query() string {
	query := url.Values{}

	if o.Around != nil {
		query.Set("around", strconv.FormatUint(*o.Around, 10))
	}

	if o.Before != nil {
		query.Set("before", strconv.FormatUint(*o.Before, 10))
	}

	if o.After != nil {
		query.Set("after", strconv.FormatUint(*o.After, 10))
	}

	if o.Limit > 100 || o.Limit < 1 {
		o.Limit = 50
		query.Set("limit", strconv.Itoa(o.Limit))
	}

	return query.Encode()
}

func GetChannelMessages(channelId uint64, token string, data GetChannelMessagesData) ([]objects.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages?%s", channelId, data.Query()),
	}

	var messages []objects.Message
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetChannelMessage(channelId, messageId uint64, token string) (*objects.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
	}

	var message objects.Message
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &message); err != nil {
		return nil, err
	}

	return &message, nil
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
	Embed           *objects.Embed         `json:"embed,omitempty"`
	PayloadJson     string                 `json:"-"` // TODO: Helper method
	AllowedMentions objects.AllowedMention `json:"allowed_mentions"`
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

func CreateMessage(channelId uint64, token string, data CreateMessageData) (*objects.Message, error) {
	var endpoint request.Endpoint
	if data.File == nil {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.ApplicationJson,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
		}
	} else {
		endpoint = request.Endpoint{
			RequestType: request.POST,
			ContentType: request.MultipartFormData,
			Endpoint:    fmt.Sprintf("/channels/%d/messages", channelId),
		}
	}

	var message objects.Message
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// emoji is the raw unicode emoji
func CreateReaction(channelId, messageId uint64, emoji, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s", channelId, messageId, url.QueryEscape(emoji)),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteOwnReaction(channelId, messageId uint64, emoji, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/@me", channelId, messageId, url.QueryEscape(emoji)),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

// emoji is the raw unicode emoji
func DeleteUserReaction(channelId, messageId, userId uint64, emoji, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s/%d", channelId, messageId, url.QueryEscape(emoji), userId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

type GetReactionsData struct {
	Before *uint64 // get users before this user ID
	After  *uint64 // get users after this user ID
	Limit  int     // 1 - 100
}

func (o *GetReactionsData) Query() string {
	query := url.Values{}

	if o.Before != nil {
		query.Set("before", strconv.FormatUint(*o.Before, 10))
	}

	if o.After != nil {
		query.Set("after", strconv.FormatUint(*o.After, 10))
	}

	if o.Limit > 100 || o.Limit < 1 {
		o.Limit = 25
		query.Set("limit", strconv.Itoa(o.Limit))
	}

	return query.Encode()
}

func GetReactions(channelId, messageId uint64, emoji, token string, data GetReactionsData) ([]objects.User, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s?%s", channelId, messageId, url.QueryEscape(emoji), data.Query()),
	}

	var users []objects.User
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteAllReactions(channelId, messageId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions", channelId, messageId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

func DeleteAllReactionsEmoji(channelId, messageId uint64, emoji, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d/reactions/%s", channelId, messageId, url.QueryEscape(emoji)),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

type EditMessageData struct {
	Content string         `json:"content,omitempty"`
	Embed   *objects.Embed `json:"embed,omitempty"`
	Flags   int            `json:"flags,omitempty"` // https://discordapp.com/developers/docs/resources/channel#message-object-message-flags TODO: Helper function
}

func EditMessage(channelId, messageId uint64, token string, data ModifyChannelData) (*objects.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
	}

	var message objects.Message
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func DeleteMessage(channelId, messageId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/%d", channelId, messageId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

func BulkDeleteMessages(channelId uint64, messages []uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/messages/bulk-delete", channelId),
	}

	body := map[string]interface{}{
		"messages": utils.Uint64StringSlice(messages),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, body, nil)
	return err
}

func EditChannelPermissions(channelId uint64, token string, updated objects.Overwrite) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, updated.Id),
	}

	updated.Id = 0

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, updated, nil)
	return err
}

func GetChannelInvites(channelId uint64, token string) ([]objects.InviteMetadata, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
	}

	var invites []objects.InviteMetadata
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

func CreateChannelInvite(channelId uint64, token string, data objects.InviteMetadata) (*objects.Invite, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/invites", channelId),
	}

	var invite objects.Invite
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, data, &invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

func DeleteChannelPermissions(channelId, overwriteId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/channels/%d/permissions/%d", channelId, overwriteId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

func TriggerTypingIndicator(channelId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/typing", channelId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

func GetPinnedMessages(channelId uint64, token string) ([]objects.Message, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins", channelId),
	}

	var messages []objects.Message
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func AddPinnedChannelMessage(channelId, messageId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.PUT,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}

func DeletePinnedChannelMessage(channelId, messageId uint64, token string) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/channels/%d/pins/%d", channelId, messageId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetChannelRoute(channelId).Ratelimiter, nil, nil)
	return err
}
