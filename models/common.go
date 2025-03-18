package models

import (
	"time"

	"gorm.io/gorm"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginationQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type CommonTime struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
