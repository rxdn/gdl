package routes

type VoiceRoute struct {
	Ratelimiter Ratelimiter
}

func NewVoiceRoute(rrm *RestRouteManager) *VoiceRoute {
	return &VoiceRoute{
		Ratelimiter: NewRatelimiter(rrm),
	}
}

func (e *VoiceRoute) Endpoint() string {
	return "/voice"
}

func (e *VoiceRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
