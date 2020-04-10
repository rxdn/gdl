package examples

import (
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/sirupsen/logrus"
)

const(
	StarboardChannelId = 1
	StarThreshold = 3
)

func main() {
	cacheFactory := cache.BoltCacheFactory(cache.CacheOptions{
		Guilds:      true,
		Users:       true,
	}, cache.BoltOptions{
		ClearOnRestart: false,
		Path:           "bolt.db",
		FileMode:       600,
		Options:        nil,
	})

	shardOptions := gateway.ShardOptions{
		ShardCount: gateway.ShardCount{
			Total:   1,
			Lowest:  0,
			Highest: 1,
		},
		RateLimitStore: ratelimit.NewMemoryStore(),
		CacheFactory:       cacheFactory,
		GuildSubscriptions: false,
	}

	token := ""
	sm := gateway.NewShardManager(token, shardOptions)

	sm.RegisterListeners(reactListener)
	sm.Connect()
	sm.WaitForInterrupt()
}

func reactListener(s *gateway.Shard, e *events.MessageReactionAdd) {
	// check the new reaction is a star
	if e.Emoji.Name != "⭐" {
		return
	}
	
	// get people who have reacted with a star so we can check if we've met the threshold
	reactors, err := s.GetReactions(e.ChannelId, e.MessageId, "⭐", rest.GetReactionsData{}); if err != nil {
		logrus.Warn(err)
		return
	}

	if len(reactors) >= StarThreshold {
		// get the message object so we can get the content & sender
		msg, err := s.GetChannelMessage(e.ChannelId, e.MessageId); if err != nil {
			logrus.Warn(err)
			return
		}

		// create an embed
		embed := embed.NewEmbed().
			SetTitle("New Star").
			SetAuthor(msg.Author.Username, "", msg.Author.AvatarUrl(256)).
			SetDescription(msg.Content)

		if _, err := s.CreateMessageEmbed(StarboardChannelId, embed); err != nil {
			logrus.Warn(err)
			return
		}
	}
}
