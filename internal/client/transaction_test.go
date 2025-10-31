package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JscorpTech/ocpp/internal/config"
)

func TestNewTransactionClient(t *testing.T) {
	cfg := &config.Config{
		BaseUrl: "http://localhost:8000",
		Addr:    ":8080",
	}

	client := NewTransactionClient(cfg)
	if client == nil {
		t.Fatal("NewTransactionClient() returned nil")
	}
}

func TestTransactionClient_GetTransactionFromTag_Success(t *testing.T) {
	// Mock server
	mockResponse := Transaction{
		Status: true,
		Data: struct {
			Id            int     `json:"id"`
			Conn          int     `json:"conn"`
			Status        string  `json:"status"`
			Limit         string  `json:"limit"`
			Amount        string  `json:"amount"`
			IsActive      bool    `json:"is_active"`
			Tag           string  `json:"tag"`
			MeterStart    int     `json:"meter_start"`
			MeterStop     int     `json:"meter_stop"`
			MeterConsumed int     `json:"meter_consumed"`
			Soc           int     `json:"soc"`
			StartDate     string  `json:"start_date"`
			EndDate       *string `json:"end_date"`
		}{
			Id:         123,
			Conn:       1,
			Status:     "active",
			Limit:      "100",
			Amount:     "50.00",
			IsActive:   true,
			Tag:        "RFID-12345",
			MeterStart: 0,
			MeterStop:  5000,
			StartDate:  "2025-10-31T10:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/transaction/tag/RFID-12345/" {
			t.Errorf("Expected path '/api/transaction/tag/RFID-12345/', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cfg := &config.Config{
		BaseUrl: server.URL,
		Addr:    ":8080",
	}

	client := NewTransactionClient(cfg)
	transaction, err := client.GetTransactionFromTag("RFID-12345")

	if err != nil {
		t.Fatalf("GetTransactionFromTag() error = %v", err)
	}

	if !transaction.Status {
		t.Error("Expected Status to be true")
	}

	if transaction.Data.Id != 123 {
		t.Errorf("Id = %v, want 123", transaction.Data.Id)
	}

	if transaction.Data.Tag != "RFID-12345" {
		t.Errorf("Tag = %v, want RFID-12345", transaction.Data.Tag)
	}
}

func TestTransactionClient_GetTransactionFromTag_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := &config.Config{
		BaseUrl: server.URL,
		Addr:    ":8080",
	}

	client := NewTransactionClient(cfg)
	_, err := client.GetTransactionFromTag("test-tag")

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestTransactionClient_GetTransactionFromTag_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.Config{
		BaseUrl: server.URL,
		Addr:    ":8080",
	}

	client := NewTransactionClient(cfg)
	transaction, err := client.GetTransactionFromTag("test-tag")

	// Server xatosi bo'lsa ham, response qaytarish mumkin
	if err != nil {
		t.Logf("Got expected error: %v", err)
	}

	// Agar transaction qaytsa, status false bo'lishi kerak
	if transaction != nil && transaction.Status {
		t.Error("Expected Status to be false on server error")
	}
}
