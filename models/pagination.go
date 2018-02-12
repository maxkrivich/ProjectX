package models

import "math"

type Pagination struct {
	Limit   int  `json:"perPage"`
	Page    int  `json:"page"`
	Count   int  `json:"count"`
	HasNext bool `json:"hasNext"`
}

func NewPagination(limit, page int) *Pagination {
	return &Pagination{
		Limit:   int(math.Max(1, math.Min(10000, float64(limit)))),
		Page:    int(math.Max(1, float64(page))),
		Count:   0,
		HasNext: false,
	}
}
