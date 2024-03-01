package main

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	NGE "github.com/Skapar/NGE/pkg/nge"
	"github.com/Skapar/NGE/pkg/nge/models"
	// Auth "github.com/Skapar/NGE/pkg/nge/Auth"
)

var db *gorm.DB
type HealthCheckResponse struct {
	Status string `json:"status"`
	Check   string `json:"Check"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{"Internal Server Error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, HealthCheckResponse{"ok", NGE.HealthCheck()})
}

// func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	allStudents := User.GetAllUsers()
// 	writeJSONResponse(w, http.StatusOK, allStudents)
// }

// func signupHandler(w http.ResponseWriter, r *http.Request) {
// 	// Call the Signup function from the User package, passing the database connection
// 	// along with the response writer and request. Since you haven't provided the database
// 	// connection in this file, I'm assuming you have it set up elsewhere and you'll pass it here.
// 	Auth.Signup(db, w, r)
// }

// func signinHandler(w http.ResponseWriter, r *http.Request) {
// 	// Similarly, call the Signin function from the User package.
// 	Auth.Signin(db, w, r)
// }


// func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	studentID := vars["id"]

// 	student, err := User.GetUserByID(studentID)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusNotFound, ErrorResponse{"User not found"})
// 		return
// 	}

// 	writeJSONResponse(w, http.StatusOK, student)
// }



func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extracting the user ID from the URL parameter
	vars := mux.Vars(r)
	idStr := vars["id"] // Ensure your route variable is named 'id'
	id, err := strconv.Atoi(idStr) // Converts the ID from string to int
	if err != nil {
		// If there's an error in conversion, return a bad request response
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the user by ID using the GetUserByID function
	user, err := models.GetUserByID(db, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If no user is found, return a not found response
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For any other errors, return an internal server error response
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// If the user is found, encode and return the user as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}