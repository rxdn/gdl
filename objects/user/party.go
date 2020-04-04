package user

type Party struct {
	Id   string `json:"id"`
	Size []int  `json:"size"`
}

func (p *Party) GetCurrentSize() int {
	if len(p.Size) >= 1 {
		return p.Size[0]
	} else {
		return 0
	}
}

func (p *Party) GetMaxSize() int {
	if len(p.Size) >= 2 {
		return p.Size[1]
	} else {
		return 0
	}
}
