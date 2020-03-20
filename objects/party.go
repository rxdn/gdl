package objects

type Party struct {
	Id   uint64 `json:",string"`
	Size []int
}
