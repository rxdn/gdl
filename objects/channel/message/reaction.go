package message

import "github.com/rxdn/gdl/objects/guild/emoji"

type Reaction struct {
	Count int
	Me    bool
	Emoji emoji.Emoji
}
