package objects

type Emoji struct {
	Id            uint64   `json:"id,string"`
	Name          string   `json:"name"`
	Roles         []uint64 `json:"roles,string"`
	User          *User    `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
}
