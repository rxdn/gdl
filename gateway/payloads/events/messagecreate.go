package events

import (
	"github.com/rxdn/gdl/objects/channel/message"
)

type MessageCreate struct {
	*message.Message
}
