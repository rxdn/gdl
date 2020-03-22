package routes

type GenericRoute interface {
	Endpoint() string
	GetRatelimit() *Ratelimiter
}
