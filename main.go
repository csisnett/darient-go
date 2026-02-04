package main

import (
	"log"
	"net/http"
	"os"

	"backend/internal/database"
	"backend/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	if err := database.Connect(databaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := database.InitSchema(); err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	r := mux.NewRouter()
	
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/api/items", handlers.GetItems).Methods("GET")
	r.HandleFunc("/api/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/api/items/{id}", handlers.GetItem).Methods("GET")
	
	// Client routes
	r.HandleFunc("/api/clients", handlers.CreateClient).Methods("POST")
	r.HandleFunc("/api/clients/{id}", handlers.GetClient).Methods("GET")
	r.HandleFunc("/api/clients/{id}", handlers.UpdateClient).Methods("PUT")
	r.HandleFunc("/api/clients/{id}", handlers.DeleteClient).Methods("DELETE")

	// Bank routes
	r.HandleFunc("/api/banks", handlers.GetBanks).Methods("GET")
	r.HandleFunc("/api/banks", handlers.CreateBank).Methods("POST")
	r.HandleFunc("/api/banks/{id}", handlers.GetBank).Methods("GET")
	r.HandleFunc("/api/banks/{id}", handlers.UpdateBank).Methods("PUT")
	r.HandleFunc("/api/banks/{id}", handlers.DeleteBank).Methods("DELETE")

	// Credit routes
	r.HandleFunc("/api/credits", handlers.GetCredits).Methods("GET")
	r.HandleFunc("/api/credits", handlers.CreateCredit).Methods("POST")
	r.HandleFunc("/api/credits/{id}", handlers.GetCredit).Methods("GET")
	r.HandleFunc("/api/credits/{id}", handlers.UpdateCredit).Methods("PUT")
	r.HandleFunc("/api/credits/{id}", handlers.DeleteCredit).Methods("DELETE")
	r.HandleFunc("/api/clients/{clientId}/credits", handlers.GetCreditsByClient).Methods("GET")
	r.HandleFunc("/api/banks/{bankId}/credits", handlers.GetCreditsByBank).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
