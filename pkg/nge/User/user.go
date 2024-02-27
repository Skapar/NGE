package user

import (
	"errors"

	"github.com/Skapar/common/database/initializers"

	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeDatabase(database *gorm.DB) {
	db = initializers.DB
	db.AutoMigrate(&Student{}) // Make sure our database schema is updated
}



func GetAllUsers() []Student {
	return studentList
}

func GetUserByID(id string) (*Student, error) {
	for _, student := range studentList {
		if student.ID == id {
			return &student, nil
		}
	}
	return nil, errors.New("student not found")
}