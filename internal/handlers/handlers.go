package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/internal/database"
	"backend/internal/logger"
	"backend/internal/models"
	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rows, err := database.DB.Query("SELECT id, name, description, created_at FROM items ORDER BY created_at DESC")
	if err != nil {
		if logger.APILogger != nil {
			logger.APILogger.LogError(r.Method, r.URL.Path, "Database query failed: "+err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	items := []models.Item{}
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err != nil {
			if logger.APILogger != nil {
				logger.APILogger.LogError(r.Method, r.URL.Path, "Row scan failed: "+err.Error())
			}
			continue
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var item models.Item
	err = database.DB.QueryRow("SELECT id, name, description, created_at FROM items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	json.NewEncoder(w).Encode(item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.Item
	
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		if logger.APILogger != nil {
			logger.APILogger.LogError(r.Method, r.URL.Path, "Invalid JSON in request body: "+err.Error())
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at",
		item.Name, item.Description,
	).Scan(&item.ID, &item.CreatedAt)

	if err != nil {
		if logger.APILogger != nil {
			logger.APILogger.LogError(r.Method, r.URL.Path, "Failed to insert item: "+err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create item"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}
// Client handlers

func GetClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var client models.Client
	err = database.DB.QueryRow("SELECT id, full_name, email, birth_date, country, created_at FROM clients WHERE id = $1", id).
		Scan(&client.ID, &client.FullName, &client.Email, &client.BirthDate, &client.Country, &client.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}

	json.NewEncoder(w).Encode(client)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Parse birth_date if it's provided as string
	var birthDate time.Time
	if !client.BirthDate.IsZero() {
		birthDate = client.BirthDate
	}

	_, err = database.DB.Exec(
		"UPDATE clients SET full_name = $1, email = $2, birth_date = $3, country = $4 WHERE id = $5",
		client.FullName, client.Email, birthDate, client.Country, id,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update client"})
		return
	}

	// Fetch updated client
	err = database.DB.QueryRow("SELECT id, full_name, email, birth_date, country, created_at FROM clients WHERE id = $1", id).
		Scan(&client.ID, &client.FullName, &client.Email, &client.BirthDate, &client.Country, &client.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch updated client"})
		return
	}

	json.NewEncoder(w).Encode(client)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM clients WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete client"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client deleted successfully"})
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var client models.Client
	
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO clients (full_name, email, birth_date, country) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		client.FullName, client.Email, client.BirthDate, client.Country,
	).Scan(&client.ID, &client.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create client"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

// Bank handlers

func GetBanks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rows, err := database.DB.Query("SELECT id, name, type, created_at FROM banks ORDER BY created_at DESC")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	banks := []models.Bank{}
	for rows.Next() {
		var bank models.Bank
		if err := rows.Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt); err != nil {
			continue
		}
		banks = append(banks, bank)
	}

	json.NewEncoder(w).Encode(banks)
}

func GetBank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var bank models.Bank
	err = database.DB.QueryRow("SELECT id, name, type, created_at FROM banks WHERE id = $1", id).
		Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		return
	}

	json.NewEncoder(w).Encode(bank)
}

func CreateBank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bank models.Bank
	
	if err := json.NewDecoder(r.Body).Decode(&bank); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate bank type
	if bank.Type != models.BankTypePrivate && bank.Type != models.BankTypeGovernment {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank type. Must be PRIVATE or GOVERNMENT"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO banks (name, type) VALUES ($1, $2) RETURNING id, created_at",
		bank.Name, bank.Type,
	).Scan(&bank.ID, &bank.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create bank"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bank)
}

func UpdateBank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var bank models.Bank
	if err := json.NewDecoder(r.Body).Decode(&bank); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate bank type
	if bank.Type != models.BankTypePrivate && bank.Type != models.BankTypeGovernment {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank type. Must be PRIVATE or GOVERNMENT"})
		return
	}

	_, err = database.DB.Exec(
		"UPDATE banks SET name = $1, type = $2 WHERE id = $3",
		bank.Name, bank.Type, id,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update bank"})
		return
	}

	// Fetch updated bank
	err = database.DB.QueryRow("SELECT id, name, type, created_at FROM banks WHERE id = $1", id).
		Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch updated bank"})
		return
	}

	json.NewEncoder(w).Encode(bank)
}

func DeleteBank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM banks WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete bank"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bank deleted successfully"})
}
// Credit handlers

func GetCredits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rows, err := database.DB.Query(`
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, 
		       credit_type, status, created_at 
		FROM credits 
		ORDER BY created_at DESC
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	credits := []models.Credit{}
	for rows.Next() {
		var credit models.Credit
		if err := rows.Scan(&credit.ID, &credit.ClientID, &credit.BankID, 
			&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
			&credit.CreditType, &credit.Status, &credit.CreatedAt); err != nil {
			continue
		}
		credits = append(credits, credit)
	}

	json.NewEncoder(w).Encode(credits)
}

func GetCredit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var credit models.Credit
	err = database.DB.QueryRow(`
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, 
		       credit_type, status, created_at 
		FROM credits WHERE id = $1
	`, id).Scan(&credit.ID, &credit.ClientID, &credit.BankID, 
		&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
		&credit.CreditType, &credit.Status, &credit.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Credit not found"})
		return
	}

	json.NewEncoder(w).Encode(credit)
}

func CreateCredit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var credit models.Credit
	
	if err := json.NewDecoder(r.Body).Decode(&credit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate credit type
	if credit.CreditType != models.CreditTypeAuto && 
	   credit.CreditType != models.CreditTypeMortgage && 
	   credit.CreditType != models.CreditTypeCommercial {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credit type. Must be AUTO, MORTGAGE, or COMMERCIAL"})
		return
	}

	// Validate status if provided, otherwise default to PENDING
	if credit.Status == "" {
		credit.Status = models.CreditStatusPending
	} else if credit.Status != models.CreditStatusPending && 
	          credit.Status != models.CreditStatusApproved && 
	          credit.Status != models.CreditStatusRejected {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid status. Must be PENDING, APPROVED, or REJECTED"})
		return
	}

	// Validate payment amounts
	if credit.MinPayment <= 0 || credit.MaxPayment <= 0 || credit.MinPayment > credit.MaxPayment {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payment amounts. Min and max must be positive, and min must be <= max"})
		return
	}

	// Validate term months
	if credit.TermMonths <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Term months must be positive"})
		return
	}

	err := database.DB.QueryRow(`
		INSERT INTO credits (client_id, bank_id, min_payment, max_payment, term_months, credit_type, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at
	`, credit.ClientID, credit.BankID, credit.MinPayment, credit.MaxPayment, 
	   credit.TermMonths, credit.CreditType, credit.Status).Scan(&credit.ID, &credit.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create credit"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func UpdateCredit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	var credit models.Credit
	if err := json.NewDecoder(r.Body).Decode(&credit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate credit type
	if credit.CreditType != models.CreditTypeAuto && 
	   credit.CreditType != models.CreditTypeMortgage && 
	   credit.CreditType != models.CreditTypeCommercial {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credit type. Must be AUTO, MORTGAGE, or COMMERCIAL"})
		return
	}

	// Validate status
	if credit.Status != models.CreditStatusPending && 
	   credit.Status != models.CreditStatusApproved && 
	   credit.Status != models.CreditStatusRejected {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid status. Must be PENDING, APPROVED, or REJECTED"})
		return
	}

	// Validate payment amounts
	if credit.MinPayment <= 0 || credit.MaxPayment <= 0 || credit.MinPayment > credit.MaxPayment {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payment amounts. Min and max must be positive, and min must be <= max"})
		return
	}

	// Validate term months
	if credit.TermMonths <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Term months must be positive"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE credits 
		SET client_id = $1, bank_id = $2, min_payment = $3, max_payment = $4, 
		    term_months = $5, credit_type = $6, status = $7 
		WHERE id = $8
	`, credit.ClientID, credit.BankID, credit.MinPayment, credit.MaxPayment, 
	   credit.TermMonths, credit.CreditType, credit.Status, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update credit"})
		return
	}

	// Fetch updated credit
	err = database.DB.QueryRow(`
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, 
		       credit_type, status, created_at 
		FROM credits WHERE id = $1
	`, id).Scan(&credit.ID, &credit.ClientID, &credit.BankID, 
		&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
		&credit.CreditType, &credit.Status, &credit.CreatedAt)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch updated credit"})
		return
	}

	json.NewEncoder(w).Encode(credit)
}

func DeleteCredit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM credits WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete credit"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Credit not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Credit deleted successfully"})
}

func GetCreditsByClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	clientID, err := strconv.Atoi(params["clientId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid client ID"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, 
		       credit_type, status, created_at 
		FROM credits 
		WHERE client_id = $1 
		ORDER BY created_at DESC
	`, clientID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	credits := []models.Credit{}
	for rows.Next() {
		var credit models.Credit
		if err := rows.Scan(&credit.ID, &credit.ClientID, &credit.BankID, 
			&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
			&credit.CreditType, &credit.Status, &credit.CreatedAt); err != nil {
			continue
		}
		credits = append(credits, credit)
	}

	json.NewEncoder(w).Encode(credits)
}

func GetCreditsByBank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bankID, err := strconv.Atoi(params["bankId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, 
		       credit_type, status, created_at 
		FROM credits 
		WHERE bank_id = $1 
		ORDER BY created_at DESC
	`, bankID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	credits := []models.Credit{}
	for rows.Next() {
		var credit models.Credit
		if err := rows.Scan(&credit.ID, &credit.ClientID, &credit.BankID, 
			&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
			&credit.CreditType, &credit.Status, &credit.CreatedAt); err != nil {
			continue
		}
		credits = append(credits, credit)
	}

	json.NewEncoder(w).Encode(credits)
}