package routes

type VoiceRoute struct {
	Ratelimiter Ratelimiter
}

func NewVoiceRoute() *VoiceRoute {
	return &VoiceRoute{
		Ratelimiter: NewRatelimiter(),
	}
}

func (e *VoiceRoute) Endpoint() string {
	return "/voice"
}

func (e *VoiceRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
