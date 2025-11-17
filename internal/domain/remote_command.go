package domain

import (
	"encoding/json"
)

type RemoteCommand string

var (
	RemoteStartTransaction RemoteCommand = "remote_start_transaction"
	RemoteStopTransaction  RemoteCommand = "remote_stop_transaction"
	GetConfiguration       RemoteCommand = "get_configuration"
	ChangeConfiguration    RemoteCommand = "change_configuration"
)

type RemoteCommandRes struct {
	Detail string `json:"detail"`
	Data   any    `json:"data"`
}

type RemoteCommandReq struct {
	CpID    string          `json:"cp_id"`
	Command RemoteCommand   `json:"command"`
	Data    json.RawMessage `json:"data"`
}

type RemoteStartTransactionReq struct {
	Tag         string `json:"tag"`
	ConnectorID int32  `json:"connector_id"`
}

type RemoteStartTransactionRes struct {
	Status string `json:"status"`
}
type RemoteStopTransactionReq struct {
	TransactionId int32 `json:"transaction_id"`
}
type RemoteStopTransactionRes struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type GetConfigurationReq struct {
	Key []string `json:"key"`
}

type ChangeConfigurationReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
