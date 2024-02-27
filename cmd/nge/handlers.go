package main

import (
	"encoding/json"
	"net/http"

	NGE "github.com/Skapar/NGE/pkg/nge"
	User "github.com/Skapar/NGE/pkg/nge/User"
)
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

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	allStudents := User.GetAllUsers()
	writeJSONResponse(w, http.StatusOK, allStudents)
}

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
