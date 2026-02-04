package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/models"
	"github.com/gorilla/mux"
)

// Test data for unit tests
var (
	validClient = models.Client{
		FullName:  "John Doe",
		Email:     "john.doe@example.com",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Country:   "USA",
	}

	validBank = models.Bank{
		Name: "Test Bank",
		Type: models.BankTypePrivate,
	}

	validCredit = models.Credit{
		ClientID:   1,
		BankID:     1,
		MinPayment: 100.0,
		MaxPayment: 1000.0,
		TermMonths: 12,
		CreditType: models.CreditTypeAuto,
		Status:     models.CreditStatusPending,
	}

	validItem = models.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}
)

func TestHealthCheckUnit(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Could not parse response")
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

// Test invalid JSON handling
func TestCreateItemInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/items", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateItem)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateClientInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/clients", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateClient)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateBankInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/banks", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateBank)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateBankInvalidType(t *testing.T) {
	invalidBank := models.Bank{
		Name: "Invalid Bank",
		Type: "INVALID_TYPE",
	}
	jsonData, _ := json.Marshal(invalidBank)
	req, _ := http.NewRequest("POST", "/api/banks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateBank)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid bank type: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateCreditInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/credits", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCredit)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateCreditInvalidType(t *testing.T) {
	invalidCredit := validCredit
	invalidCredit.CreditType = "INVALID_TYPE"
	jsonData, _ := json.Marshal(invalidCredit)
	req, _ := http.NewRequest("POST", "/api/credits", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCredit)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid credit type: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateCreditInvalidPaymentAmounts(t *testing.T) {
	invalidCredit := validCredit
	invalidCredit.MinPayment = 1000.0
	invalidCredit.MaxPayment = 100.0 // min > max
	jsonData, _ := json.Marshal(invalidCredit)
	req, _ := http.NewRequest("POST", "/api/credits", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCredit)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid payment amounts: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateCreditInvalidTermMonths(t *testing.T) {
	invalidCredit := validCredit
	invalidCredit.TermMonths = -1
	jsonData, _ := json.Marshal(invalidCredit)
	req, _ := http.NewRequest("POST", "/api/credits", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCredit)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid term months: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Test invalid ID handling
func TestGetItemInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/items/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/items/{id}", GetItem).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetClientInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/clients/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/clients/{id}", GetClient).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetBankInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/banks/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/banks/{id}", GetBank).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetCreditInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/credits/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/credits/{id}", GetCredit).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetCreditsByClientInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/clients/invalid/credits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/clients/{clientId}/credits", GetCreditsByClient).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetCreditsByBankInvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/banks/invalid/credits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/banks/{bankId}/credits", GetCreditsByBank).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}