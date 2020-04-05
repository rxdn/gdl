package ratelimit

type Route interface {
	Endpoint() string
}
