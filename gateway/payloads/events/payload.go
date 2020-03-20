package events

import "encoding/json"

type Payload struct {
	Name string          `json:"t"`
	Raw  json.RawMessage `json:"d"`
	Data interface{}     `json:"-"`
}
