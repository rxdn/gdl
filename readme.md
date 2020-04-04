# About  
GDL (Go Discord Library) provides a wrapper around the Discord API with additional helper functions to make life easier  
  
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
		Total:   2,
		Lowest:  0, // Inclusive
		Highest: 2, // Exclusive
	}  
 
	 // create our cache tuned to our liking (remember the zero value of a bool is false)
	 cache := cache.MemoryCacheFactory(cache.CacheOptions{
		 Guilds:      true,
		 Users:       true,
		 Channels:    true,
		 Roles:       true,
		 Emojis:      false,
		 VoiceStates: false, 
	})  
	token := ""
	sm := gateway.NewShardManager(token, shardOptions, cache)  
	sm.GuildSubscriptions = false // Toggle guild subscriptions
	sm.Presence = user.BuildStatus(user.ActivityTypePlaying, "with GDL") // Set the status of the bot  
	sm.RegisterListeners(echoListener)
	sm.Connect()
	sm.WaitForInterrupt()
}  

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