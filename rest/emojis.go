package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/utils"
)

func ListGuildEmojis(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64) ([]emoji.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis", guildId),
		BaseRoute:   ratelimit.NewEmojiRoute(guildId),
		RateLimiter: rateLimiter,
	}

	var emojis []emoji.Emoji
	err, _ := endpoint.Request(token, nil, &emojis)
	return emojis, err
}

func GetGuildEmoji(token string, rateLimiter *ratelimit.Ratelimiter, guildId, emojiId uint64) (emoji.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
		BaseRoute:   ratelimit.NewEmojiRoute(guildId),
		RateLimiter: rateLimiter,
	}

	var emoji emoji.Emoji
	if err, _ := endpoint.Request(token, nil, &emoji); err != nil {
		return emoji, err
	}

	return emoji, nil
}

type CreateEmojiData struct {
	Name  string
	Image Image
	Roles []uint64 // roles for which this emoji will be whitelisted
}

func CreateGuildEmoji(token string, rateLimiter *ratelimit.Ratelimiter, guildId uint64, data CreateEmojiData) (emoji.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: data.Image.ContentType,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis", guildId),
		BaseRoute:   ratelimit.NewEmojiRoute(guildId),
		RateLimiter: rateLimiter,
	}

	var emoji emoji.Emoji
	imageData, err := data.Image.Encode()
	if err != nil {
		return emoji, err
	}

	body := map[string]interface{}{
		"name":  data.Name,
		"image": imageData,
		"roles": utils.Uint64StringSlice(data.Roles),
	}

	if err, _ := endpoint.Request(token, body, &emoji); err != nil {
		return emoji, err
	}

	return emoji, nil
}

// updating Image is not permitted
func ModifyGuildEmoji(token string, rateLimiter *ratelimit.Ratelimiter, guildId, emojiId uint64, data CreateEmojiData) (emoji.Emoji, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
		BaseRoute:   ratelimit.NewEmojiRoute(guildId),
		RateLimiter: rateLimiter,
	}

	body := map[string]interface{}{
		"name":  data.Name,
		"roles": utils.Uint64StringSlice(data.Roles),
	}

	var emoji emoji.Emoji
	if err, _ := endpoint.Request(token, body, &emoji); err != nil {
		return emoji, err
	}

	return emoji, nil
}

func DeleteGuildEmoji(token string, rateLimiter *ratelimit.Ratelimiter, guildId, emojiId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/guilds/%d/emojis/%d", guildId, emojiId),
		BaseRoute:   ratelimit.NewEmojiRoute(guildId),
		RateLimiter: rateLimiter,
	}

	err, _ := endpoint.Request(token, nil, nil)
	return err
}
