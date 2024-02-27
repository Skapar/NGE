package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Skapar/NGE/pkg/ng/models"
	"github.com/Skapar/NGE/pkg/ng/utils"
	"github.com/gorilla/mux"
)

// GetBook handles the HTTP request for getting all books.
func GetBook(w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetBookById handles the HTTP request for getting a book by ID.
func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseUint(bookId, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing book ID", http.StatusBadRequest)
		return
	}
	bookDetails, err := models.GetBookById(uint(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// CreateBook handles the HTTP request for creating a new book.
func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook := &models.Book{}
	utils.ParseBody(r, CreateBook)
	b := CreateBook.CreateBook()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteBook handles the HTTP request for deleting a book.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseUint(bookId, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing book ID", http.StatusBadRequest)
		return
	}
	err = models.DeleteBook(uint(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted successfully"))
}

// UpdateBook handles the HTTP request for updating a book.
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseUint(bookId, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing book ID", http.StatusBadRequest)
		return
	}
	bookDetails, err := models.GetBookById(uint(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	// db := models.UpdateBook(bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
