package domain

type EventTypes string

const (
	ChangeConnectorStatusEvent EventTypes = "change_connector_status"
	StartTransactionEvent      EventTypes = "start_transaction"
	StopTransactionEvent       EventTypes = "stop_transaction"
	MeterValuesEvent           EventTypes = "meter_value"
	HealthEvent                EventTypes = "health"
	DataTransferEvent          EventTypes = "data_transfer"
)

type Event struct {
	Event EventTypes `json:"event"`
	Data  any        `json:"data"`
}

type ChangeConnectorStatus struct {
	Charger string `json:"charger"`
	Conn    int    `json:"conn"`
	Status  string `json:"status"`
}

type StartTransaction struct {
	Charger    string `json:"charger"`
	Conn       int    `json:"conn"`
	Tag        string `json:"tag"`
	MeterStart int    `json:"meter_start"`
}

type StopTransaction struct {
	Charger       string `json:"charger"`
	TransactionId int    `json:"transaction_id"`
	Reason        string `json:"reason"`
	MeterStop     int    `json:"meter_stop"`
}

type MeterValues struct {
	Conn          int   `json:"conn"`
	TransactionId int32 `json:"transaction_id"`
	MeterValue    any   `json:"meter_value"`
}

type Healthcheck struct {
	Charger string `json:"charger"`
}

type DataTransfer struct {
	VendorId  string `json:"vendorId"`
	MessageId string `json:"messageId"`
	Data      string `json:"data"`
}
