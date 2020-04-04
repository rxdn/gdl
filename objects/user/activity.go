package user

import "github.com/rxdn/gdl/objects/guild/emoji"

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	Url           string       `json:"url"`
	Timestamps    Timestamps   `json:"timestamps"`
	ApplicationId uint64       `json:"application_id,string"`
	Details       string       `json:"details"`
	State         string       `json:"state"`
	Emoji         emoji.Emoji  `json:"emoji"`
	Party         Party        `json:"party"`
	Assets        Asset        `json:"assets"`
	Secret        Secret       `json:"secret"`
	Instance      bool         `json:"instance"`
	Flags         int          `json:"flags"` // TODO: Wrap this
}
