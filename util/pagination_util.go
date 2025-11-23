package util

import (
	"math"

	"gorm.io/gorm"
)

type PaginationMeta struct {
	TotalRows  int64 `json:"total_rows"`
	LastPage   int   `json:"last_page"`
	Page       int   `json:"page"`
}

type Pagination struct {
	Data interface{} `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

func Paginate(page int, pagination *Pagination, totalRows int64, size int) func(db *gorm.DB) *gorm.DB {

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size
	pagination.Meta.Page = page
	pagination.Meta.TotalRows = totalRows
	pagination.Meta.LastPage = int(math.Ceil(float64(totalRows) / float64(size)))
	if page > pagination.Meta.LastPage {
		pagination.Meta.Page = 0
		pagination.Meta.TotalRows = 0
		pagination.Meta.LastPage = 0
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(size)
	}
}
