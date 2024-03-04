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

	r.HandleFunc("/events", app.AddEventHandler).Methods("POST")
	r.HandleFunc("/events/{id}", app.get).Methods("GET")
	r.HandleFunc("/events/{id}", app.DeleteEventHandler).Methods("DELETE")
	r.HandleFunc("/events/{id}", app.UpdateEventHandler).Methods("PUT")

	r.HandleFunc("/addPost", app.add).Methods("POST")
	r.HandleFunc("/getPost/{id}", app.get).Methods("GET")
	r.HandleFunc("/updatePost/{id}", app.update).Methods("PUT")
	r.HandleFunc("/deletePost/{id}", app.delete).Methods("DELETE")

	r.HandleFunc("/user", app.CreateUserHandler).Methods("POST")
	r.HandleFunc("/user/{id}", app.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", app.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", app.DeleteUserHandler).Methods("DELETE")



	fmt.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}