package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/JscorpTech/ocpp/internal/config"
)

type Transaction struct {
	Status bool `json:"status"`
	Data   struct {
		Id            int     `json:"id"`
		Conn          int     `json:"conn"`
		Status        string  `json:"status"`
		Limit         string  `json:"limit"`
		Amount        string  `json:"amount"`
		Tag           string  `json:"tag"`
		MeterStart    int     `json:"meter_start"`
		MeterStop     int     `json:"meter_stop"`
		MeterConsumed int     `json:"meter_consumed"`
		Soc           int     `json:"soc"`
		StartDate     string  `json:"start_date"`
		EndDate       *string `json:"end_date"`
	} `json:"data"`
}

type TransactionClient interface {
	GetTransactionFromTag(string) (*Transaction, error)
}

type transactionClient struct {
	Client *http.Client
	Config *config.Config
}

func NewTransactionClient(cfg *config.Config) TransactionClient {
	return &transactionClient{
		Client: &http.Client{},
		Config: cfg,
	}
}

func (t *transactionClient) GetTransactionFromTag(tag string) (*Transaction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/transaction/tag/%s/", t.Config.BaseUrl, tag), nil)
	if err != nil {
		return nil, err
	}
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var transaction Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		fmt.Println(string(body))
		return nil, err
	}
	return &transaction, nil
}
