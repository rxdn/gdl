package objects

type MessageActivity struct {
	Type    int
	PartyId uint64 `json:",string"`
}
