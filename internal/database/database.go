package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func InitSchema() error {
	itemsQuery := `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	clientsQuery := `
	CREATE TABLE IF NOT EXISTS clients (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		birth_date DATE NOT NULL,
		country VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	banksQuery := `
	CREATE TABLE IF NOT EXISTS banks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(20) NOT NULL CHECK (type IN ('PRIVATE', 'GOVERNMENT')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	creditsQuery := `
	CREATE TABLE IF NOT EXISTS credits (
		id SERIAL PRIMARY KEY,
		client_id INTEGER NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
		bank_id INTEGER NOT NULL REFERENCES banks(id) ON DELETE CASCADE,
		min_payment DECIMAL(15,2) NOT NULL,
		max_payment DECIMAL(15,2) NOT NULL,
		term_months INTEGER NOT NULL,
		credit_type VARCHAR(20) NOT NULL CHECK (credit_type IN ('AUTO', 'MORTGAGE', 'COMMERCIAL')),
		status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(itemsQuery)
	if err != nil {
		return err
	}

	_, err = DB.Exec(clientsQuery)
	if err != nil {
		return err
	}

	_, err = DB.Exec(banksQuery)
	if err != nil {
		return err
	}

	_, err = DB.Exec(creditsQuery)
	if err != nil {
		return err
	}

	log.Println("Database schema initialized")
	return nil
}
