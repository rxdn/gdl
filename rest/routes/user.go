package routes

import "fmt"

type UserRoute struct {
	Id          uint64
	Ratelimiter Ratelimiter
}

func NewUserRoute(id uint64, rrm *RestRouteManager) *UserRoute {
	return &UserRoute{
		Id:          id,
		Ratelimiter: NewRatelimiter(rrm),
	}
}

func (e *UserRoute) Endpoint() string {
	return fmt.Sprintf("/users/%d", e.Id)
}

func (e *UserRoute) GetRatelimit() *Ratelimiter {
	return &e.Ratelimiter
}
