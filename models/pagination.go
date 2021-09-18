package models

type Pagination struct {
	Next struct {
		Page  int `json:"page,omitempty"`
		Limit int `json:"limit,omitempty"`
	} `json:"next"`

	Prev struct {
		Page  int `json:"page,omitempty"`
		Limit int `json:"limit,omitempty"`
	} `json:"prev"`
}

func (p *Pagination) Fill(page int, limit int, startIndex int, endIndex int, total int) {
	if endIndex < total {
		p.Next.Page = page + 1
		p.Next.Limit = limit
	}
	if startIndex > 0 {
		p.Prev.Page = page - 1
		p.Prev.Limit = limit
	}
}
