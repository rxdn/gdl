package rest

import (
	"encoding/base64"
	"fmt"
	"github.com/Dot-Rar/gdl/objects"
	"github.com/Dot-Rar/gdl/rest/request"
	"github.com/Dot-Rar/gdl/rest/routes"
	"github.com/Dot-Rar/gdl/utils"
	"io"
	"io/ioutil"
)

func ListGuildEmojis(guildId uint64, token string) ([]*objects.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis", guildId),
	}

	var emojis []*objects.Emoji
	err, _ := endpoint.Request(token, &routes.RouteManager.GetEmojiRoute(guildId).Ratelimiter, nil, &emojis)
	return emojis, err
}

func GetGuildEmoji(guildId, emojiId uint64, token string) (*objects.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
	}

	var emoji objects.Emoji
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetEmojiRoute(guildId).Ratelimiter, nil, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

type Image struct {
	ContentType request.ContentType
	ImageReader io.Reader
}

func (i *Image) Encode() (string, error) {
	content, err := ioutil.ReadAll(i.ImageReader)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	return fmt.Sprintf("data:%s;base64,%s", string(i.ContentType), encoded), nil
}

type CreateEmojiData struct {
	Name  string
	Image Image
	Roles []uint64 // roles for which this emoji will be whitelisted
}

func CreateGuildEmoji(guildId uint64, token string, data CreateEmojiData) (*objects.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: data.Image.ContentType,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis", guildId),
	}

	body := map[string]interface{}{
		"name": data.Name,
		"image": data.Image.Encode(),
		"roles": utils.Uint64StringSlice(data.Roles),
	}

	var emoji objects.Emoji
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetEmojiRoute(guildId).Ratelimiter, body, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// updating Image is not permitted
func ModifyGuildEmoji(guildId, emojiId uint64, token string, data CreateEmojiData) (*objects.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
	}

	body := map[string]interface{}{
		"name": data.Name,
		"roles": utils.Uint64StringSlice(data.Roles),
	}

	var emoji objects.Emoji
	if err, _ := endpoint.Request(token, &routes.RouteManager.GetEmojiRoute(guildId).Ratelimiter, body, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

func DeleteGuildEmoji(guildId, emojiId uint64, token string) (error) {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
	}

	err, _ := endpoint.Request(token, &routes.RouteManager.GetEmojiRoute(guildId).Ratelimiter, nil, nil)
	return err
}

