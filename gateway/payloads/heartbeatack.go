package payloads

import "encoding/json"

type HeartbeatAck struct {
	Opcode int `json:"op"`
	SequenceNumber int `json:"d"`
}

func NewHeartbeackAck(raw []byte) (Hello, error) {
	var payload Hello
	err := json.Unmarshal(raw, &payload); if err != nil {
		return Hello{}, err
	}

	return payload, nil
}
