package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint      `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	PasswordHash  string    `json:"password"`
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
