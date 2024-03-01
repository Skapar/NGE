package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	NGE "github.com/Skapar/NGE/pkg/nge"
	"github.com/Skapar/NGE/pkg/nge/models"
)

// HealthCheckResponse represents the response for health check endpoint
type HealthCheckResponse struct {
	Status string `json:"status"`
	Check  string `json:"Check"`
}

// ErrorResponse represents the response for error scenarios
type ErrorResponse struct {
	Error string `json:"error"`
}

// CustomDB embeds *gorm.DB to allow defining methods on it
type CustomDB struct {
	*gorm.DB
}

type EventRequest struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// HealthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, HealthCheckResponse{"ok", NGE.HealthCheck()})
}

// writeJSONResponse writes JSON response to the http.ResponseWriter
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

func (db *CustomDB) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the URL parameter
	vars := mux.Vars(r)
	idStr := vars["id"]            // Ensure your route variable is named 'id'
	id, err := strconv.Atoi(idStr) // Converts the ID from string to int
	if err != nil {
		// If there's an error in conversion, return a bad request response
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the user by ID using the GetUserByID function
	user, err := models.GetUserByID(db.DB, uint(id))
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

func (app *App) AddEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assuming a function in the models package that adds an event
	err = models.AddEvent(app.DB, req.Date, req.Description)
	if err != nil {
		http.Error(w, "Failed to add event", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	writeJSONResponse(w, http.StatusCreated, map[string]string{"result": "Event added successfully"})
}
