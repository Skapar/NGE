package models

import (
	"gorm.io/gorm"
)

type FilterParams struct {
	Text     string `json:"text"`
	SortBy   string `json:"sort_by"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

func FilterPosts(db *gorm.DB, filters FilterParams) ([]Post, error) {
	var posts []Post
	query := db

	if filters.Text != "" {
		query = query.Where("text LIKE ?", "%"+filters.Text+"%")
	}

	if filters.SortBy != "" {
		query = query.Order(filters.SortBy)
	}

	if filters.PageSize > 0 {
		offset := (filters.Page - 1) * filters.PageSize
		limit := filters.PageSize
		query = query.Offset(offset).Limit(limit)
	}
	err := query.Find(&posts).Error
	return posts, err
}
