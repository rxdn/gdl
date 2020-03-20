package gateway

import (
	"encoding/json"
	"github.com/Dot-Rar/gdl/gateway/payloads/events"
)

func (s *Shard) ExecuteEvent(eventType events.EventType, data json.RawMessage) {
	switch eventType {
	case events.GUILD_CREATE:
		var decoded events.GuildCreate
		if err := json.Unmarshal(data, &decoded); err == nil {
			for _, listener := range s.ShardManager.EventBus.Listeners {
				if fn, ok := listener.(func(s *Shard, e *events.GuildCreate)); ok {
					fn(s, &decoded)
				}
			}
		}

	case events.MESSAGE_CREATE:
		var decoded events.MessageCreate
		if err := json.Unmarshal(data, &decoded); err == nil {
			for _, listener := range s.ShardManager.EventBus.Listeners {
				if fn, ok := listener.(func(s *Shard, e *events.MessageCreate)); ok {
					fn(s, &decoded)
				}
			}
		}
	}
}