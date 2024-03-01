package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Skapar/NGE/pkg/nge/config"
	"github.com/Skapar/NGE/pkg/nge/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type App struct {
	DB      *gorm.DB
	DBModel *models.DBModel
}

func main() {
	db, err := config.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	config.Migrate()
	dbModel = &models.DBModel{DB: db}
	app := App{DB: db}
	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheckHandler)

	r.HandleFunc("/GetUserById", GetUser)

	r.HandleFunc("/addPost", app.add).Methods("POST")
	r.HandleFunc("/getPost/{id}", app.get).Methods("GET")
	r.HandleFunc("/updatePost/{id}", app.update).Methods("PUT")
	r.HandleFunc("/deletePost/{id}", app.delete).Methods("DELETE")

	// r.HandleFunc("/signup", signupHandler)
	// r.HandleFunc("/signin", signinHandler)

	// r.HandleFunc("/students", getAllUsersHandler)
	// r.HandleFunc("/student/{id}", getUserByIDHandler)

	fmt.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
