package avatar

import (
	"fmt"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/command"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cacheFactory := cache.BoltCacheFactory(cache.CacheOptions{
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

	registerCommands(sm)

	sm.Connect()
	sm.WaitForInterrupt()
}

func registerCommands(sm *gateway.ShardManager) {
	// register a command handler with prefix !
	ch := command.NewCommandHandler(sm, "!")

	// create a new avatar command
	cmd := command.NewCommand("avatar", []string{"a"}, onCommand)

	// register our command
	ch.RegisterCommand(cmd)
}

func onCommand(ctx command.CommandContext) {
	if len(ctx.Mentions) == 0 {
		_, _ = ctx.Shard.CreateMessage(ctx.ChannelId, "You need to mention a user")
		return
	}

	for _, mention := range ctx.Mentions {
		res, err := http.Get(mention.AvatarUrl(2048)); if err != nil {
			logrus.Warn(err.Error())
			continue
		}
		defer res.Body.Close()

		data := rest.CreateMessageData{
			Content:         fmt.Sprintf("%s's avatar is:", mention.Username),
			File:            &rest.File{
				Name:        "avatar.png",
				ContentType: res.Header.Get("Content-Type"),
				Reader:      res.Body,
			},
		}

		ctx.Shard.CreateMessageComplex(ctx.ChannelId, data)
	}
}
