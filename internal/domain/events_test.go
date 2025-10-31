package domain

import (
	"encoding/json"
	"testing"
)

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventTypes
		want      string
	}{
		{"change connector status", ChangeConnectorStatusEvent, "change_connector_status"},
		{"start transaction", StartTransactionEvent, "start_transaction"},
		{"stop transaction", StopTransactionEvent, "stop_transaction"},
		{"meter values", MeterValuesEvent, "meter_value"},
		{"health", HealthEvent, "health"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.eventType) != tt.want {
				t.Errorf("EventType = %v, want %v", tt.eventType, tt.want)
			}
		})
	}
}

func TestEvent_JSONMarshaling(t *testing.T) {
	event := Event{
		Event: HealthEvent,
		Data: Healthcheck{
			Charger: "test-charger",
		},
	}

	// Marshal
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Unmarshal
	var unmarshaled Event
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if unmarshaled.Event != event.Event {
		t.Errorf("Event type = %v, want %v", unmarshaled.Event, event.Event)
	}
}

func TestChangeConnectorStatus(t *testing.T) {
	status := ChangeConnectorStatus{
		Charger: "charger-001",
		Conn:    1,
		Status:  "Available",
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled ChangeConnectorStatus
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Charger != status.Charger {
		t.Errorf("Charger = %v, want %v", unmarshaled.Charger, status.Charger)
	}
	if unmarshaled.Conn != status.Conn {
		t.Errorf("Conn = %v, want %v", unmarshaled.Conn, status.Conn)
	}
	if unmarshaled.Status != status.Status {
		t.Errorf("Status = %v, want %v", unmarshaled.Status, status.Status)
	}
}

func TestStartTransaction(t *testing.T) {
	tx := StartTransaction{
		Charger:    "charger-001",
		Conn:       1,
		Tag:        "RFID-12345",
		MeterStart: 0,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled StartTransaction
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Tag != tx.Tag {
		t.Errorf("Tag = %v, want %v", unmarshaled.Tag, tx.Tag)
	}
	if unmarshaled.MeterStart != tx.MeterStart {
		t.Errorf("MeterStart = %v, want %v", unmarshaled.MeterStart, tx.MeterStart)
	}
}

func TestStopTransaction(t *testing.T) {
	tx := StopTransaction{
		Charger:       "charger-001",
		TransactionId: 123,
		Reason:        "Local",
		MeterStop:     5000,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled StopTransaction
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.TransactionId != tx.TransactionId {
		t.Errorf("TransactionId = %v, want %v", unmarshaled.TransactionId, tx.TransactionId)
	}
	if unmarshaled.MeterStop != tx.MeterStop {
		t.Errorf("MeterStop = %v, want %v", unmarshaled.MeterStop, tx.MeterStop)
	}
}

func TestMeterValues(t *testing.T) {
	meter := MeterValues{
		Conn:          1,
		TransactionId: 123,
		MeterValue:    map[string]interface{}{"value": 1000},
	}

	data, err := json.Marshal(meter)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled MeterValues
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Conn != meter.Conn {
		t.Errorf("Conn = %v, want %v", unmarshaled.Conn, meter.Conn)
	}
	if unmarshaled.TransactionId != meter.TransactionId {
		t.Errorf("TransactionId = %v, want %v", unmarshaled.TransactionId, meter.TransactionId)
	}
}

func TestHealthcheck(t *testing.T) {
	health := Healthcheck{
		Charger: "charger-001",
	}

	data, err := json.Marshal(health)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled Healthcheck
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Charger != health.Charger {
		t.Errorf("Charger = %v, want %v", unmarshaled.Charger, health.Charger)
	}
}
