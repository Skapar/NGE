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