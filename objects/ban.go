package objects

type Ban struct {
	Reason string `json:"reason,omitempty"`
	User   *User  `json:"user"`
}
