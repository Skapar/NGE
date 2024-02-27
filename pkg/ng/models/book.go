package models

import (
	"github.com/Skapar/NGE/pkg/ng/config"
	"gorm.io/gorm"
)

// Assuming config.Connect() returns a *gorm.DB and error, handle the error appropriately.
var db, err = config.Connect()

// Book defines the model for the database table.
type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

// Initialize the package, running database migrations.
func Initialize() {
	if err != nil {
		// Handle the error, maybe log it or panic
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Book{})
}

// CreateBook adds a new book to the database.
func (b *Book) CreateBook() *Book {
	db.Create(b)
	return b
}

// GetAllBooks retrieves all books from the database.
func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

// GetBookById retrieves a book by its ID.
func GetBookById(Id uint) (*Book, error) {
	var getBook Book
	result := db.First(&getBook, Id)
	return &getBook, result.Error
}

// DeleteBook removes a book by its ID.
func DeleteBook(ID uint) error {
	result := db.Delete(&Book{}, ID)
	return result.Error
}
