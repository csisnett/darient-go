package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/models"
	"github.com/gorilla/mux"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// Setup test database
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		fmt.Println("TEST_DATABASE_URL not set, skipping integration tests")
		os.Exit(0)
	}

	if err := database.Connect(testDBURL); err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := database.InitSchema(); err != nil {
		fmt.Printf("Failed to initialize test schema: %v\n", err)
		os.Exit(1)
	}

	// Setup test server
	r := setupRouter()
	testServer = httptest.NewServer(r)
	defer testServer.Close()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestData()
	os.Exit(code)
}

func setupRouter() *mux.Router {
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

	return r
}

func cleanupTestData() {
	database.DB.Exec("DELETE FROM credits")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM banks")
	database.DB.Exec("DELETE FROM items")
}

func TestIntegrationHealthCheck(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	
	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", result["status"])
	}
}

func TestIntegrationItemsCRUD(t *testing.T) {
	// Clean up before test
	database.DB.Exec("DELETE FROM items")

	// Test Create Item
	item := models.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}
	
	jsonData, _ := json.Marshal(item)
	resp, err := http.Post(testServer.URL+"/api/items", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdItem models.Item
	json.NewDecoder(resp.Body).Decode(&createdItem)
	
	if createdItem.Name != item.Name {
		t.Errorf("Expected name %s, got %s", item.Name, createdItem.Name)
	}

	// Test Get Item
	resp2, err := http.Get(fmt.Sprintf("%s/api/items/%d", testServer.URL, createdItem.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp2.StatusCode)
	}

	// Test Get Items
	resp3, err := http.Get(testServer.URL + "/api/items")
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp3.StatusCode)
	}

	var items []models.Item
	json.NewDecoder(resp3.Body).Decode(&items)
	
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
}

func TestIntegrationClientsCRUD(t *testing.T) {
	// Clean up before test
	database.DB.Exec("DELETE FROM clients")

	// Test Create Client
	client := models.Client{
		FullName:  "John Doe",
		Email:     "john.doe@example.com",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Country:   "USA",
	}
	
	jsonData, _ := json.Marshal(client)
	resp, err := http.Post(testServer.URL+"/api/clients", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdClient models.Client
	json.NewDecoder(resp.Body).Decode(&createdClient)
	
	if createdClient.FullName != client.FullName {
		t.Errorf("Expected name %s, got %s", client.FullName, createdClient.FullName)
	}

	// Test Get Client
	resp2, err := http.Get(fmt.Sprintf("%s/api/clients/%d", testServer.URL, createdClient.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp2.StatusCode)
	}

	// Test Update Client
	updatedClient := createdClient
	updatedClient.FullName = "Jane Doe"
	updatedClient.Email = "jane.doe@example.com"
	
	jsonData2, _ := json.Marshal(updatedClient)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/clients/%d", testServer.URL, createdClient.ID), bytes.NewBuffer(jsonData2))
	req.Header.Set("Content-Type", "application/json")
	
	client2 := &http.Client{}
	resp3, err := client2.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp3.StatusCode)
	}

	// Test Delete Client
	req2, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/clients/%d", testServer.URL, createdClient.ID), nil)
	resp4, err := client2.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp4.StatusCode)
	}
}

func TestIntegrationBanksCRUD(t *testing.T) {
	// Clean up before test
	database.DB.Exec("DELETE FROM banks")

	// Test Create Bank
	bank := models.Bank{
		Name: "Test Bank",
		Type: models.BankTypePrivate,
	}
	
	jsonData, _ := json.Marshal(bank)
	resp, err := http.Post(testServer.URL+"/api/banks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdBank models.Bank
	json.NewDecoder(resp.Body).Decode(&createdBank)
	
	if createdBank.Name != bank.Name {
		t.Errorf("Expected name %s, got %s", bank.Name, createdBank.Name)
	}

	// Test Get Bank
	resp2, err := http.Get(fmt.Sprintf("%s/api/banks/%d", testServer.URL, createdBank.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp2.StatusCode)
	}

	// Test Get Banks
	resp3, err := http.Get(testServer.URL + "/api/banks")
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp3.StatusCode)
	}

	var banks []models.Bank
	json.NewDecoder(resp3.Body).Decode(&banks)
	
	if len(banks) != 1 {
		t.Errorf("Expected 1 bank, got %d", len(banks))
	}

	// Test Update Bank
	updatedBank := createdBank
	updatedBank.Name = "Updated Bank"
	updatedBank.Type = models.BankTypeGovernment
	
	jsonData2, _ := json.Marshal(updatedBank)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/banks/%d", testServer.URL, createdBank.ID), bytes.NewBuffer(jsonData2))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp4, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp4.StatusCode)
	}

	// Test Delete Bank
	req2, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/banks/%d", testServer.URL, createdBank.ID), nil)
	resp5, err := client.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp5.Body.Close()

	if resp5.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp5.StatusCode)
	}
}

func TestIntegrationCreditsCRUD(t *testing.T) {
	// Clean up before test
	database.DB.Exec("DELETE FROM credits")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM banks")

	// Create test client and bank first
	client := models.Client{
		FullName:  "John Doe",
		Email:     "john.doe@example.com",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Country:   "USA",
	}
	
	jsonData, _ := json.Marshal(client)
	resp, err := http.Post(testServer.URL+"/api/clients", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var createdClient models.Client
	json.NewDecoder(resp.Body).Decode(&createdClient)

	bank := models.Bank{
		Name: "Test Bank",
		Type: models.BankTypePrivate,
	}
	
	jsonData2, _ := json.Marshal(bank)
	resp2, err := http.Post(testServer.URL+"/api/banks", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	var createdBank models.Bank
	json.NewDecoder(resp2.Body).Decode(&createdBank)

	// Test Create Credit
	credit := models.Credit{
		ClientID:   createdClient.ID,
		BankID:     createdBank.ID,
		MinPayment: 100.0,
		MaxPayment: 1000.0,
		TermMonths: 12,
		CreditType: models.CreditTypeAuto,
		Status:     models.CreditStatusPending,
	}
	
	jsonData3, _ := json.Marshal(credit)
	resp3, err := http.Post(testServer.URL+"/api/credits", "application/json", bytes.NewBuffer(jsonData3))
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp3.StatusCode)
	}

	var createdCredit models.Credit
	json.NewDecoder(resp3.Body).Decode(&createdCredit)
	
	if createdCredit.ClientID != credit.ClientID {
		t.Errorf("Expected client ID %d, got %d", credit.ClientID, createdCredit.ClientID)
	}

	// Test Get Credit
	resp4, err := http.Get(fmt.Sprintf("%s/api/credits/%d", testServer.URL, createdCredit.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp4.StatusCode)
	}

	// Test Get Credits
	resp5, err := http.Get(testServer.URL + "/api/credits")
	if err != nil {
		t.Fatal(err)
	}
	defer resp5.Body.Close()

	if resp5.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp5.StatusCode)
	}

	var credits []models.Credit
	json.NewDecoder(resp5.Body).Decode(&credits)
	
	if len(credits) != 1 {
		t.Errorf("Expected 1 credit, got %d", len(credits))
	}

	// Test Get Credits by Client
	resp6, err := http.Get(fmt.Sprintf("%s/api/clients/%d/credits", testServer.URL, createdClient.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp6.Body.Close()

	if resp6.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp6.StatusCode)
	}

	// Test Get Credits by Bank
	resp7, err := http.Get(fmt.Sprintf("%s/api/banks/%d/credits", testServer.URL, createdBank.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp7.Body.Close()

	if resp7.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp7.StatusCode)
	}

	// Test Update Credit
	updatedCredit := createdCredit
	updatedCredit.Status = models.CreditStatusApproved
	updatedCredit.MaxPayment = 2000.0
	
	jsonData4, _ := json.Marshal(updatedCredit)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/credits/%d", testServer.URL, createdCredit.ID), bytes.NewBuffer(jsonData4))
	req.Header.Set("Content-Type", "application/json")
	
	httpClient := &http.Client{}
	resp8, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp8.Body.Close()

	if resp8.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp8.StatusCode)
	}

	// Test Delete Credit
	req2, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/credits/%d", testServer.URL, createdCredit.ID), nil)
	resp9, err := httpClient.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp9.Body.Close()

	if resp9.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp9.StatusCode)
	}
}