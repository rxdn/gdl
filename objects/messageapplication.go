package objects

type MessageApplication struct {
	Id          uint64 `json:",string"`
	CoverImage  string
	Description string
	Icon        string
	Name        string
}
