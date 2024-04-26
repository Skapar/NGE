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
// PROFILE HANDLER

func (app *App) addProfile(w http.ResponseWriter, r *http.Request) {
	var newProfile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&newProfile); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if err := models.AddProfile(app.DB, &newProfile); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, newProfile)
}

func (app *App) getProfileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid profile ID"})
		return
	}

	profile, err := models.GetProfileById(app.DB, uint(profileID))
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if profile == nil {
		writeJSONResponse(w, http.StatusNotFound, ErrorResponse{"Profile not found"})
		return
	}

	writeJSONResponse(w, http.StatusOK, profile)
}

func (app *App) updateProfileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid profile ID"})
		return
	}

	var updatedProfile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if err := models.UpdateProfileById(app.DB, uint(profileID), &updatedProfile); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedProfile)
}

func (app *App) deleteProfileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid profile ID"})
		return
	}

	if err := models.DeleteProfileById(app.DB, uint(profileID)); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Profile deleted successfully"})
}

// _________________________________________________________

// FILTER HANDLER

func (app *App) FilterPosts(w http.ResponseWriter, r *http.Request) {
	var filters models.FilterParams
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pagination
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1 // Default to page 1 if not provided or invalid
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10 // Default page size if not provided or invalid
	}
	filters.Page = page
	filters.PageSize = pageSize

	// Sorting
	sortBy := r.URL.Query().Get("sort_by")
	filters.SortBy = sortBy

	posts, err := models.FilterPosts(app.DB, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
