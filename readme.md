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
        CacheFactory: cache.BoltCacheFactory(cache.CacheOptions{
            Guilds:      true,
            Users:       true,
            Members:     true,
            Channels:    true,
            Roles:       true,
            Emojis:      true,
            VoiceStates: true,
        }, cache.BoltOptions{
            ClearOnRestart: false, // Should we clear the cache on start-up? (Do not use on WSL)
            Path:           "bolt.db",
            FileMode:       600,
            Options:        nil, // Additional bolt options
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

Other examples are available in the examples package

# Events
View the following package for a list of events: [gateway/payloads/events](https://github.com/rxdn/gdl/tree/master/gateway/payloads/events)

Gateway events are also available to listen on: [gateway/payloads](https://github.com/rxdn/gdl/tree/master/gateway/payloads)

# Commands
GDL comes with a built-in command handler, however, feel free to build your own.

## Example
```go
sm := gateway.NewShardManager(token, shardOptions)
ch := command.NewCommandHandler(sm, "!", "-") // register a new command handler with the prefixes ! and -

// create a new "hello" command, with no aliases, that executes in a goroutine
myCommand := command.NewCommand("hello", nil, true, func(ctx command.CommandContext) {
    _, _ = ctx.Shard.CreateMessage(ctx.ChannelId, fmt.Sprintf("Hello, %s!", ctx.Author.Username))
})

// create a subcommand (i.e. when the user runs !hello world), with an alias of "alias" (so !hello alias), that doesn't execute in a goroutine
myCommand.RegisterSubCommand(command.NewCommand("world", []string{"alias"}, false, func(ctx command.CommandContext) {
    _ = ctx.Shard.CreateReaction(ctx.ChannelId, ctx.Message.Id, "👍")
}))

// register our hello command
// subcommands do not need to be registered on the command handler
ch.RegisterCommand(myCommand)
```

![Command handler in action](https://i.imgur.com/eNH0NIb.png)

Take a look at the [examples package](https://github.com/rxdn/gdl/tree/master/examples) for more usage examples
(also feel free to contribute some!).

# Caching
GDL currently offers 2 caches, however, you are free to develop your own:

## PostgreSQL cache (recommended)
The PostgreSQL cache offers the best performance and scalability, as well as multi-instance support. So therefore, it is
recommended that you use it if you have access to a PostgreSQL server.

GDL uses the [PGX library](https://github.com/jackc/pgx) for accessing PostgreSQL. You are responsible for making a
pgx.DB instance and passing it to GDL.

### Example
```go
db, err := pgxpool.Connect(context.Background(), "postgres://user:pwd@localhost/database?pool_max_conns=2")
if err != nil {
	panic(err)
}

c := cache.PgCacheFactory(db, cache.CacheOptions{
    Guilds:      true,
    Users:       true, 
    Members:     true,
    Channels:    true,
    Roles:       true,
    Emojis:      true,
    VoiceStates: true,
})

shardOptions := gateway.ShardOptions{
    ...
    CacheFactory: c,
    ...
}
```

## Bolt cache
If you do not have access to a PostgreSQL database, GDL offers a [Bolt](https://github.com/boltdb/bolt) based cache.
Bolt is a flat-file key-value store, which means that it is not able to offer the best performance or scalability, as
well as being restricted to use by a single instance.

## Example
```go
c := cache.BoltCacheFactory(cache.CacheOptions{
    Guilds:      true,
    Users:       true,
    Members:     true,
    Channels:    true,
    Roles:       true,
    Emojis:      true,
    VoiceStates: true,
}, cache.BoltOptions{
    ClearOnRestart: false, // Should we clear the cache on start-up? (Do not use on WSL)
    Path:           "bolt.db",
    FileMode:       600,
    Options:        nil, // Additional bolt options
})

shardOptions := gateway.ShardOptions{
    ...
    CacheFactory: c,
    ...
}
```

# FAQ  
## I'm getting a pkg-config / zlib error!  
The czlib library that GDL uses for compression requires the C zlib library to be installed. You can install it by:  
  
- Ubuntu: `# apt-get install zlib1g-dev`  
- CentOS: `# yum install zlib-devel`

## I'm getting a panic: invalid page type: 0: 4 when using WSL!
This is a [known issue](https://github.com/microsoft/WSL/issues/3162) with WSL. Luckily, it only happens on the first
run, so you can use Bolt with WSL if you set ClearOnRestart to false.
