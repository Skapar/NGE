package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Skapar/NGE/pkg/nge/config"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


func main() {
	config.Connect()
	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheckHandler)

	r.HandleFunc("/GetUserById", GetUser )



	// r.HandleFunc("/signup", signupHandler)
	// r.HandleFunc("/signin", signinHandler)

	// r.HandleFunc("/students", getAllUsersHandler)
	// r.HandleFunc("/student/{id}", getUserByIDHandler)

	fmt.Println("Server listening on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}