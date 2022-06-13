package interaction

import (
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

type ResolvedData struct {
	Users       map[objects.Snowflake]user.User          `json:"users"`
	Members     map[objects.Snowflake]member.Member      `json:"members"`
	Roles       map[objects.Snowflake]guild.Role         `json:"roles"`
	Channels    map[objects.Snowflake]channel.Channel    `json:"channels"`
	Messages    map[objects.Snowflake]message.Message    `json:"messages"`
	Attachments map[objects.Snowflake]channel.Attachment `json:"attachments"`
}
