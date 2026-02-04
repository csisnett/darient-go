package main

import (
	"fmt"
	"log"
	"os"

	"backend/internal/database"
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

	// Check what tables exist
	fmt.Println("Checking existing tables...")
	rows, err := database.DB.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name;
	`)
	if err != nil {
		log.Fatal("Failed to query tables:", err)
	}
	defer rows.Close()

	fmt.Println("Existing tables:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		fmt.Printf("- %s\n", tableName)
	}

	// Check if banks table exists and has the right structure
	fmt.Println("\nChecking banks table structure...")
	rows2, err := database.DB.Query(`
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'banks' 
		ORDER BY ordinal_position;
	`)
	if err != nil {
		fmt.Printf("Error checking banks table: %v\n", err)
	} else {
		defer rows2.Close()
		fmt.Println("Banks table columns:")
		for rows2.Next() {
			var columnName, dataType, isNullable string
			if err := rows2.Scan(&columnName, &dataType, &isNullable); err != nil {
				continue
			}
			fmt.Printf("- %s (%s, nullable: %s)\n", columnName, dataType, isNullable)
		}
	}

	// Try to count records in banks table
	fmt.Println("\nChecking banks table data...")
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM banks").Scan(&count)
	if err != nil {
		fmt.Printf("Error counting banks: %v\n", err)
	} else {
		fmt.Printf("Banks table has %d records\n", count)
	}

	// Test a simple query on banks table
	fmt.Println("\nTesting banks query...")
	rows3, err := database.DB.Query("SELECT id, name, type, created_at FROM banks ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		fmt.Printf("Error querying banks: %v\n", err)
	} else {
		defer rows3.Close()
		fmt.Println("Sample banks data:")
		for rows3.Next() {
			var id int
			var name, bankType string
			var createdAt string
			if err := rows3.Scan(&id, &name, &bankType, &createdAt); err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			fmt.Printf("- ID: %d, Name: %s, Type: %s, Created: %s\n", id, name, bankType, createdAt)
		}
	}
}