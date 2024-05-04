package models

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Filters struct {
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
	Sort         string `json:"sort"`
	SortSafeList []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (f Filters) Validate() error {
	if f.Page <= 0 {
		return fmt.Errorf("page must be greater than 0")
	}
	if f.Page > 10_000_0000 {
		return fmt.Errorf("page must be a maximum of 10 million")
	}
	if f.PageSize <= 0 {
		return fmt.Errorf("page_size must be greater than 0")
	}
	if f.PageSize > 100 {
		return fmt.Errorf("page_size must be a maximum of 100")
	}
	if !isValidSort(f.Sort, f.SortSafeList) {
		return fmt.Errorf("invalid sort value")
	}
	return nil
}

func isValidSort(sort string, safeList []string) bool {
	for _, safeValue := range safeList {
		if sort == safeValue || sort == "-"+safeValue {
			return true
		}
	}
	return false
}

func (f Filters) SortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func Limit(filters Filters) int {
	return filters.PageSize
}

func Offset(filters Filters) int {
	return (filters.Page - 1) * filters.PageSize
}

func FetchPosts(db *gorm.DB, limit, offset int, sortColumn, sortDirection string) ([]Post, error) {
	var posts []Post
	query := fmt.Sprintf("SELECT * FROM posts ORDER BY %s %s LIMIT %d OFFSET %d", sortColumn, sortDirection, limit, offset)

	if err := db.Raw(query).Scan(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}