package objects

import "github.com/Dot-Rar/gdl/utils"

type Emoji struct {
	Id            uint64                  `json:"id,string"`
	Name          string                  `json:"name"`
	Roles         utils.Uint64StringSlice `json:"roles,string"`
	User          *User                   `json:"user"`
	RequireColons bool                    `json:"require_colons"`
	Managed       bool                    `json:"managed"`
	Animated      bool                    `json:"animated"`
}
