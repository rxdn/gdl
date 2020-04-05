# About  
GDL (Go Discord Library) provides a wrapper around the Discord API with additional helper functions to make life easier  
  
[![Discord Invite](https://discordapp.com/api/v6//guilds/691411922477908008/widget.png?style=banner2)](https://discord.gg/adVPPGp)
# Getting Started  
## Installation  
Note: GDL uses Go modules. Run `go mod init` to initialise your go.mod file.  
  
To install GDL, run `go get github.com/rxdn/gdl`  
  
## Example
From the following example, you should be able to work out the gist of a basic bot:

```go    
package main   
    
import (  
   "github.com/rxdn/gdl/cache"   
   "github.com/rxdn/gdl/gateway"  
   "github.com/rxdn/gdl/gateway/payloads/events"  
   "github.com/rxdn/gdl/objects/user"  
)    
    
func main() { 
	shardOptions := gateway.ShardOptions{ 
		ShardCount: gateway.ShardCount{  
			Total:   1,  
			Lowest:  0,  // Inclusive
			Highest: 1,  // Exclusive
		},
		RateLimitStore: ratelimit.NewMemoryStore(), // ratelimit.NewRedisStore() is also available
		// We can choose exactly what we want to cache. Remember the zero value of a bool is false!
		CacheFactory: cache.MemoryCacheFactory(cache.CacheOptions{
			Guilds:      true,  
			Users:       true,
			Members:     true,
			Channels:    true,  
			Roles:       true,  
			Emojis:      false,
			VoiceStates: false,
		}),
		GuildSubscriptions: false,
		Presence: user.BuildStatus(user.ActivityTypePlaying, "DM for help | t!help"), // Set the status of the bot
	}

   token := ""  
   sm := gateway.NewShardManager(token, shardOptions)
   sm.RegisterListeners(echoListener)  
   sm.Connect()  
   sm.WaitForInterrupt()  
}    
  
// Example listener that will just echo back any messages sent
func echoListener(s *gateway.Shard, e *events.MessageCreate) {    
    _, _ = s.CreateMessage(e.ChannelId, e.Content)  
}    
```

## Events
View the following package for a list of events: [gateway/payloads/events](https://github.com/rxdn/gdl/tree/master/gateway/payloads/events)

Gateway events are also available to listen on: [gateway/payloads](https://github.com/rxdn/gdl/tree/master/gateway/payloads)

# FAQ  
## I'm getting a pkg-config / zlib error!  
The czlib library that GDL uses for compression requires the C zlib library to be installed. You can install it by:  
  
- Ubuntu: `# apt-get install zlib1g-dev`  
- CentOS: `# yum install zlib-devel`