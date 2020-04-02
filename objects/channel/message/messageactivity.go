package message

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyId uint64              `json:"party_id,string"`
}
