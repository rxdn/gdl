package user

// uses millis since unix epoch
type Timestamps struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}
