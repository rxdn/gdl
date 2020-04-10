package command

import (
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/gateway/payloads/events"
)

type CommandContext struct {
	*events.MessageCreate
	Shard *gateway.Shard
	Args []string
}
