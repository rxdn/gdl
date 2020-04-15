![GDL](https://all-my-homies-use.aurieh.services/5eNXRKJ.png)
# About  

GDL (Go Discord Library) provides a wrapper around the Discord API with additional helper functions to make life easier for developers.

**GDL keeps in mind the main principles of the Go programming language**:    
Simple and easy to learn.

Most methods and properties are named directly after what the Discord API names them.

# Support
Join our [Discord](https://discord.gg/adVPPGp) and ask in #support if you're having trouble with some code or the library itself.

# Getting Started  
## Installation  
**Note: GDL uses Go modules. Run `go mod init` to initialise your go.mod file.**

**Go modules require Go v1.11 or higher**
  
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

// create a new "hello" command, with no aliases
myCommand := command.NewCommand("hello", nil, func(ctx command.CommandContext) {
    _, _ = ctx.Shard.CreateMessage(ctx.ChannelId, fmt.Sprintf("Hello, %s!", ctx.Author.Username))
})

// create a subcommand (i.e. when the user runs !hello world), with an alias of "alias" (so !hello alias)
myCommand.RegisterSubCommand(command.NewCommand("world", []string{"alias"}, func(ctx command.CommandContext) {
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

# Error Handling
When calling a REST API method, Discord may send an error response. You can tell what kind of error has occurred through
calling `errors.Is` and comparing the error to one of [GDL's error types](https://github.com/rxdn/gdl/blob/master/rest/request/errors.go).

More generally, GDL wraps common errors, such as 404, 403 in either a ClientError or ServerError type. You can then run
`request.IsServerError(err)` and `request.IsClientError(err)` to determine whether an error is a ClientError or
ServerError.

If an error that GDl does not provide a wrapper for occurs (it provides wrappers for all response codes that I have
seen Discord return), an `ErrUnknown` will be returned, and full details will be logged. In this case, you should open
an issue so that I can create the required error wrapper (or PR it!).

The benefit of this is that you are able to do things like this:
```go
ch, err := s.CreateGuildChannel(guildId, data)
if err != nil {
	if errors.Is(err, request.ErrForbidden) {
		_, _ = s.CreateMessage(e.ChannelId, "I do not have permission to create channels!")
	}
	return
}
```

Note: However, in this case it is recommended to use GDL's [permission calculator](https://github.com/rxdn/gdl/blob/master/permission/permissioncalculator.go)
to determine whether your bot has the required permissions for an action, rather than sending a request to Discord that
is guaranteed to fail, as sending 10000 requests in 10 minutes that fail with a 401, 403 or 429 will ban your token
from the API for an entire hour.

# FAQ  
## I'm getting a pkg-config / zlib error!  
The library that GDL uses for compression requires the C zlib library to be installed. You can install it by:  
  
- Ubuntu: `# apt-get install zlib1g-dev`  
- CentOS: `# yum install zlib-devel`

## I'm getting a panic: invalid page type: 0: 4 when using WSL!
This is a [known issue](https://github.com/microsoft/WSL/issues/3162) with WSL. Luckily, it only happens on the first
run, so you can use Bolt with WSL if you set ClearOnRestart to false.
