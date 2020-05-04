package payloads

import "github.com/rxdn/gdl/objects/user"

type PresenceUpdate struct {
	Opcode int               `json:"op"`
	Data   user.UpdateStatus `json:"d"`
}

func NewPresenceUpdate(data user.UpdateStatus) PresenceUpdate {
	return PresenceUpdate{
		Opcode: 3,
		Data:   data,
	}
}
