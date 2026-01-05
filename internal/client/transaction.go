package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JscorpTech/ocpp/internal/config"
)

type Transaction struct {
	Status bool `json:"status"`
	Data   struct {
		Id     int    `json:"id"`
		Status string `json:"status"`
		Tag    string `json:"tag"`
	} `json:"data"`
}

type TransactionClient interface {
	GetTransactionFromTag(string, string) (*Transaction, error)
}

type transactionClient struct {
	Client *http.Client
	Config *config.Config
}

func NewTransactionClient(cfg *config.Config) TransactionClient {
	return &transactionClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		Config: cfg,
	}
}

func (t *transactionClient) GetTransactionFromTag(tag string, host string) (*Transaction, error) {
	if t.Config.BaseUrl == "host" {
		host = "https://" + host
	} else {
		host = t.Config.BaseUrl
	}
	
	// Create request with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/transaction/tag/%s/", host, tag), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()
	
	// Check for non-200 status codes
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(io.LimitReader(res.Body, 1024)) // Limit error body to 1KB
		if err != nil {
			return nil, fmt.Errorf("backend API returned status %d (failed to read error body: %w)", res.StatusCode, err)
		}
		return nil, fmt.Errorf("backend API returned status %d: %s", res.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var transaction Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w (body: %s)", err, string(body))
	}
	
	return &transaction, nil
}
