package gateway

type ShardOptions struct {
	Total   int
	Lowest  int // Inclusive
	Highest int // Inclusive
}
