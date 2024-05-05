package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Skapar/NGE/pkg/nge/config"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func main() {
	r := mux.NewRouter()

	db, err := config.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	config.Migrate()
	app := App{DB: db}

	r.HandleFunc("/health", healthCheckHandler)

	r.HandleFunc("/signup", app.SignUpHandler).Methods("POST")
	r.HandleFunc("/signin", app.SignInHandler).Methods("POST")

	r.HandleFunc("/deleteUser/{id}", app.AuthMiddleware(app.DeleteUserHandler)).Methods("DELETE")

	r.HandleFunc("/role", app.CreateRoleHandler).Methods("POST")

	r.HandleFunc("/events", app.AddEventHandler).Methods("POST")
	r.HandleFunc("/events/{id}", app.GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{id}", app.DeleteEventHandler).Methods("DELETE")
	r.HandleFunc("/events/{id}", app.UpdateEventHandler).Methods("PUT")

	r.HandleFunc("/post", app.addPost).Methods("POST")
 r.HandleFunc("/post/{id}", app.getPostById).Methods("GET")
 r.HandleFunc("/post/{id}", app.updatePostById).Methods("PUT")
 r.HandleFunc("/post/{id}", app.deletePostById).Methods("DELETE")
 r.HandleFunc("/getAllPosts", app.getAllPosts).Methods("GET")
 r.HandleFunc("/filter", app.FilterHandler(app.DB)).Methods("GET")

	fmt.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
