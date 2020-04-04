package message

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyId string              `json:"party_id"`
}
