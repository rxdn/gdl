package request

import "encoding/json"

type RestError struct {
	ErrorCode int
	Message   string // json
}

func (r *RestError) IsClientError() bool {
	return r.ErrorCode >= 400 && r.ErrorCode < 500
}

func (r *RestError) IsServerError() bool {
	return r.ErrorCode >= 500
}

func (r RestError) Error() string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(r.Message), &data); err == nil {
		if msg, ok := data["message"]; ok {
			if s, ok := msg.(string); ok {
				return s
			}
		}
	}

	return r.Message
}
