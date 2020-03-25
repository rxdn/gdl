package events

import "github.com/rxdn/gdl/objects"

type Ready struct {
	GatewayVersion  int              `json:"v"`
	User            *objects.User    `json:"user"`
	PrivateChannels []uint64         `json:"private_channels,string"` // Note: This slice will always be empty
	Guilds          []*objects.Guild `json:"guilds"`
	SessionId       string           `json:"session_id"`
	Shard           []int            `json:"shard"` // Slice of [shard_id, num_shards]
}
