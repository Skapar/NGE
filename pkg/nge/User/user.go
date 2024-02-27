package user

import (
	"github.com/Skapar/NGE/common/database/initializers"
	"github.com/Skapar/NGE/models"

	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeDatabase(database *gorm.DB) {
	db = initializers.DB
	db.AutoMigrate(&models.User{}) // Make sure our database schema is updated
}



func GetAllUsers() ([]models.User) {
    var students []models.User
    return students
}

// func GetUserByID(id string) (*Student, error) {
// 	for _, student := range studentList {
// 		if student.ID == id {
// 			return &student, nil
// 		}
// 	}
// 	return nil, errors.New("student not found")
// }