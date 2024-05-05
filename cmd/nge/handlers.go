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
        // Handle error
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userID := r.Context().Value("userID").(uint)

    
    createdPost, err := models.AddPost(app.DB, newPost)
    if err != nil {
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }
	fmt.Println(userID)

    // Respond with the created post
    json.NewEncoder(w).Encode(createdPost)
}

func (app *App) updatePostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid post ID"})
		return
	}

	fmt.Println(postID)

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
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
		Username 	string	`json:"username"`
		Email	 	string	`json:"email"`
		Password 	string	`json:"password"`
		Role 		int		`json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input);
	err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	createdUser, err:= models.Signup(app.DB, input.Username, input.Email, input.Password, input.Role)
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
        // Retrieve the token from the Authorization header
        // Typically, the Authorization header is in the format "Bearer token"
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"No Authorization token provided"})
            return
        }

        // Split the header to separate the prefix from the token value
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Authorization header format must be Bearer {token}"})
            return
        }

        // Extract the JWT token from the parts array
        tokenStr := parts[1]

        // Validate the token
        claims, err := models.ValidateToken(tokenStr)
        if err != nil {
            writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{"Invalid or expired token"})
            return
        }

        // Add the user ID from the token to the context
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)

        // Create a new request with the updated context
        r = r.WithContext(ctx)

        // Call the next handler in the chain with the updated request
        next.ServeHTTP(w, r)
    }
}

// 