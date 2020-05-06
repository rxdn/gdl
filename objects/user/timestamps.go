package user

// uses millis since unix epoch
type Timestamps struct {
	Start uint64 `json:"start,omitempty"`
	End   uint64 `json:"end,omitempty"`
}
