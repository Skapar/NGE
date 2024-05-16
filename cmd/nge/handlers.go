package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
func extractUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no Authorization token provided")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("authorization header format must be Bearer {token}")
	}

	tokenStr := parts[1]

	claims, err := models.ValidateToken(tokenStr)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
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
// Role

func (app *App) CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var newRole models.Role
	if err := json.NewDecoder(r.Body).Decode(&newRole); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	createdRole, err := models.AddRole(app.DB, newRole)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdRole)
}

// User's handler

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
	userID, err := extractUserIDFromToken(r)
	// fmt.Println(userID)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid or missing token"})
		return
	}

	role, err := models.GetUserRole(app.DB, userID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	fmt.Println(role)
	if role != 1 {
		writeJSONResponse(w, http.StatusForbidden, ErrorResponse{"Insufficient permissions"})
		return
	}

	vars := mux.Vars(r)
	userIDToDelete, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid user ID"})
		return
	}

	err = models.DeleteUser(app.DB, uint(userIDToDelete))
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userID, err := strconv.ParseUint(vars["id"], 10, 64)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid user ID"})
// 		return
// 	}

// 	user, err := models.GetUserByID(app.DB, uint(userID))
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
// 		return
// 	}

// 	writeJSONResponse(w, http.StatusOK, user)
// }

// AUTH Handler
// _________________________________________________________

func (app *App) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     int    `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	createdUser, err := models.Signup(app.DB, input.Username, input.Email, input.Password, input.Role)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdUser)
}

func (app *App) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	user, err := models.GetUserByEmail(app.DB, input.Email)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid credentials"})
		return
	}

	token, err := models.GenerateToken(user.ID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"token": token})
}

func (app *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"No Authorization token provided"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Authorization header format must be Bearer {token}"})
			return
		}

		tokenStr := parts[1]

		// Validate the token
		claims, err := models.ValidateToken(tokenStr)
		if err != nil {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid or expired token"})
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

//

// ----------------------------------------------
// FILTER'S HANDLER
func (app *App) FilterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := 1
		pageSize := 10
		sort := "-created_at"

		filters := models.Filters{
			Page:         page,
			PageSize:     pageSize,
			Sort:         sort,
			SortSafeList: []string{"created_at", "-created_at"},
		}

		if err := filters.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit := models.Limit(filters)
		offset := models.Offset(filters)

		posts, err := models.FetchPosts(db, limit, offset, filters.SortColumn(), filters.SortDirection())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		for _, post := range posts {
			w.Write([]byte(post.Text + "\n"))
		}
	}
}

// Companies CRUD
// _____________________________________________________

func (app *App) AddCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Company

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assuming you have a DB connection available as app.DB
	err = models.CreateCompany(app.DB, req.Name, req.Description, req.StartDate, req.EndDate, req.Owners)
	if err != nil {
		http.Error(w, "Failed to add company", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{"result": "Company added successfully"})
}

func (app *App) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteCompany(app.DB, uint(id)); err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "Company deleted successfully"})
}

func (app *App) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.Company // Use the Company struct for decoding request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := models.UpdateCompany(app.DB, uint(id), req.Name, req.Description, req.StartDate, req.EndDate, req.Owners); err != nil {
		http.Error(w, "Failed to update company", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"result": "Company updated successfully"})
}

func (app *App) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	company, err := models.GetCompanyByID(app.DB, uint(id))
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, company)
}
