package message

type MessageActivity struct {
	Type    int
	PartyId uint64 `json:",string"`
}
