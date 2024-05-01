package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	NGE "github.com/Skapar/NGE/pkg/nge"
	"github.com/gorilla/mux"

	"github.com/Skapar/NGE/pkg/nge/models"
	"github.com/Skapar/NGE/pkg/nge/validator"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
	Check  string `json:"Check"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CustomDB struct {
	*gorm.DB
}

type EventRequest struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
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

// Events CRUD
// _____________________________________________________
func (app *App) AddEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = models.AddEvent(app.DB, req.Date, req.Description)
	if err != nil {
		http.Error(w, "Failed to add event", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{"result": "Event added successfully"})
}

func (app *App) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteEvent(app.DB, uint(id)); err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "Event deleted successfully"})
}

func (app *App) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req EventRequest // Use the EventRequest struct for decoding request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := models.UpdateEvent(app.DB, uint(id), req.Date, req.Description); err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "Event updated successfully"})
}

func (app *App) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	event, err := models.GetEventByID(app.DB, uint(id))
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, event)
}

// _________________________________________________________

// POSTS HANDLER

func (app *App) addPost(w http.ResponseWriter, r *http.Request) {
	var newPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	createdPost, err := models.AddPost(app.DB, newPost)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdPost)
}

func (app *App) updatePostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid post ID"})
		return
	}

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	updatedPost.ID = uint(postID)
	updatedPost, err = models.UpdatePost(app.DB, updatedPost)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedPost)
}

func (app *App) deletePostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid post ID"})
		return
	}

	err = models.DeletePost(app.DB, uint(postID)) // pass app.DB to the function
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

func (app *App) getPostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid post ID"})
		return
	}

	post, err := models.GetPost(app.DB, uint(postID))
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, post)
}

func (app *App) getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts(app.DB)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, posts)
}

// _____________________________________________________________

// User's handler

func (app *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	createdUser, err := models.CreateUser(app.DB, newUser)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdUser)
}

func (app *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid user ID"})
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	updatedUser.ID = uint(userID)
	updatedUser, err = models.UpdateUser(app.DB, updatedUser)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedUser)
}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid user ID"})
		return
	}

	err = models.DeleteUser(app.DB, uint(userID))
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid user ID"})
		return
	}

	user, err := models.GetUserByID(app.DB, uint(userID))
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}

// _____________________________________________________________
// FILTER'S HANDLER
func (app *App) FilterHandler(w http.ResponseWriter, r *http.Request) {
	page := parseIntQueryParam(r, "page", 1)
	pageSize := parseIntQueryParam(r, "page_size", 10)
	sort := r.URL.Query().Get("sort")

	filters := models.Filters{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sort,
		SortSafeList: []string{"created_at", "-created_at"},
	}

	validator := &validator.Validator{}
	models.ValidateFilters(validator, filters)

	limit := models.Limit(filters)
	offset := models.Offset(filters)
	sortColumn := filters.SortColumn()
	sortDirection := filters.SortDirection()
	posts, err := models.FetchPosts(app.DB, limit, offset, sortColumn, sortDirection)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, posts)
}
func parseIntQueryParam(r *http.Request, param string, defaultValue int) int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		return defaultValue
	}
	return parsedValue
}
