package events

import (
	"github.com/rxdn/gdl/objects/user"
)

type UserUpdate struct {
	*user.User
}
