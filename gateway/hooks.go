package gateway

type Hooks struct {
	ReconnectHook func(*Shard)
	IdentifyHook  func(*Shard)
	RestHook      func(url string)
}
