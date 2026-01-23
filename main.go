package main

import (
	"log"
	"net/http"

	"backend/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/api/items", handlers.GetItems).Methods("GET")
	r.HandleFunc("/api/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/api/items/{id}", handlers.GetItem).Methods("GET")
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
