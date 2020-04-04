package routes

import (
	"github.com/rxdn/gdl/utils"
	"sync"
	"time"
)

type Ratelimiter struct {
	sync.RWMutex

	RouteManager *RestRouteManager
	Limit        int
	Remaining    int
	Reset        int64 // Now + ResetAfter
	Bucket       string
	taskQueue    []chan struct{}
}

func NewRatelimiter(manager *RestRouteManager) Ratelimiter {
	return Ratelimiter{
		RouteManager: manager,
		Limit:        5, // Default to some arbitrary value so that the first request will go through
		Remaining:    5, // Default to some arbitrary value so that the first request will go through
		taskQueue:    make([]chan struct{}, 0),
	}
}

// Seconds
func (r *Ratelimiter) GetWaitTime() int64 {
	specificWaitTime := r.Reset - time.Now().Unix()

	r.RouteManager.RLock()
	globalWaitTime := r.RouteManager.GlobalRetryAfter - utils.GetCurrentTimeMillis()
	r.RouteManager.RUnlock()

	r.RouteManager.RLock()
	if globalWaitTime > 0 && globalWaitTime > specificWaitTime {
		return globalWaitTime
	}

	return specificWaitTime
}

func (r *Ratelimiter) Queue(task chan struct{}) {
	r.Lock()
	if r.Remaining > 0 || r.GetWaitTime() <= 0 {
		r.Remaining -= 1
		r.Unlock()

		task <- struct{}{}
	} else {
		r.taskQueue = append(r.taskQueue, task)

		if len(r.taskQueue) == 1 {
			r.Unlock()
			r.queueNext()
		} else {
			r.Unlock()
		}
	}
}

func (r *Ratelimiter) queueNext() {
	var task chan struct{}

	r.Lock()
	if r.Remaining > 0 || r.GetWaitTime() <= 0 {
		r.Remaining -= 1
		task, r.taskQueue = r.taskQueue[0], r.taskQueue[1:]
		r.Unlock()
		task <- struct{}{}
	} else {
		r.Unlock()
		time.Sleep(time.Duration(r.GetWaitTime()) * time.Second)
		r.Lock()

		r.Remaining -= 1
		task, r.taskQueue = r.taskQueue[0], r.taskQueue[1:]
		r.Unlock()
		task <- struct{}{}
	}
}
