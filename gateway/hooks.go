package gateway

import "net/http"

type Hooks struct {
	ReconnectHook func(*Shard)
	IdentifyHook  func(*Shard)
	RestHook      func(token string, req *http.Request)
}
