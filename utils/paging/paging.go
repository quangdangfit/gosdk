package paging

import (
	"math"
)

type Paging struct {
	Current   int `json:"current"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Limit     int `json:"limit"`
	Skip      int `json:"skip"`
}

func New(page int, pageSize int, total int) *Paging {
	var pageInfo Paging
	limit := 50

	if pageSize > 0 && pageSize <= limit {
		pageInfo.Limit = pageSize
	} else {
		pageInfo.Limit = limit
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageInfo.Limit)))
	pageInfo.Total = total
	pageInfo.TotalPage = totalPage
	if page < 1 {
		page = 1
	}
	if page > totalPage {
		page = totalPage
	}
	pageInfo.Current = page
	pageInfo.Skip = (page - 1) * pageInfo.Limit
	return &pageInfo
}
