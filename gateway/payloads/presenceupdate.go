package payloads

type PresenceUpdate struct {
	Opcode int      `json:"op"`
	Data   Presence `json:"d"`
}

func NewPresenceUpdate(data Presence) PresenceUpdate {
	return PresenceUpdate{
		Opcode: 3,
		Data:   data,
	}
}
