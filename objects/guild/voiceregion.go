package guild

type VoiceRegion struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Vip        bool   `json:"vip"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}
