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
	db, err := config.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	config.Migrate()
	app := App{DB: db}
	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheckHandler)

	//r.HandleFunc("/GetUserById/{id}", app.GetUserHandler)

	r.HandleFunc("/events", app.AddEventHandler).Methods("POST")
	r.HandleFunc("/events/{id}", app.GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{id}", app.DeleteEventHandler).Methods("DELETE")
	r.HandleFunc("/events/{id}", app.UpdateEventHandler).Methods("PUT")

	//r.HandleFunc("/createUser", app.createUserHandler).Methods("POST")

	// r.HandleFunc("/CreateEvent", app.createEventHandler).Methods("POST")

	// r.HandleFunc("/signup", signupHandler)
	// r.HandleFunc("/signin", signinHandler)

	// r.HandleFunc("/students", getAllUsersHandler)
	// r.HandleFunc("/student/{id}", getUserByIDHandler)

	fmt.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
