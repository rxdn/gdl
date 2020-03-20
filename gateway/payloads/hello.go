package payloads

import "encoding/json"

type (
	Hello struct {
		Opcode         int        `json:"op"`
		EventName      *string    `json:"t"`
		EventData      *HelloData `json:"d"`
		SequenceNumber *int       `json:"s"`
	}

	HelloData struct {
		Interval int      `json:"heartbeat_interval"`
		Trace    []string `json:"_trace"`
	}
)

func NewHello(raw []byte) (Hello, error) {
	var payload Hello
	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return Hello{}, err
	}

	return payload, nil
}
