package guild

import (
	"github.com/rxdn/gdl/objects/user"
)

type Ban struct {
	Reason string     `json:"reason,omitempty"`
	User   *user.User `json:"user"`
}
