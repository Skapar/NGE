package models

import "gorm.io/gorm"

type Role struct {
	ID int
	Title  string
}

func AddRole(db *gorm.DB, role Role) (Role, error) {
	err := db.Create(&role).Error
	return role, err
}