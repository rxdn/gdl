package ratelimit

import "fmt"

type UserRoute struct {
	Id uint64
}

func NewUserRoute(id uint64) *UserRoute {
	return &UserRoute{
		Id: id,
	}
}

func (e *UserRoute) Endpoint() string {
	return fmt.Sprintf("/users/%d", e.Id)
}
