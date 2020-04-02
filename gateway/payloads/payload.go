package payloads

import (
	"encoding/json"
)

type Payload struct {
	Opcode         int                 `json:"op"`
	Data           json.RawMessage `json:"d"`
	SequenceNumber *int                `json:"s"`
	EventName      string              `json:"t"`
}

func NewPayload(raw []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}
