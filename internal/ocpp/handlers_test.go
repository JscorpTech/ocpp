package ocpp

import (
	"context"
	"testing"
	"time"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/voltbras/go-ocpp/cs"
	"github.com/voltbras/go-ocpp/messages/v1x/cpreq"
	"github.com/voltbras/go-ocpp/messages/v1x/cpresp"
	"go.uber.org/zap"
)

func setupTestHandler() *Handlers {
	logger, _ := zap.NewDevelopment()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	metadata := cs.ChargePointRequestMetadata{
		ChargePointID: "test-charger-001",
	}
	cfg := &config.Config{
		BaseUrl: "http://localhost:8000",
		Addr:    ":8080",
	}

	event := services.NewEventService()
	return NewHandler(ctx, logger, rdb, metadata, cfg, event)
}

func TestNewHandler(t *testing.T) {
	handler := setupTestHandler()
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}

	if handler.Logger == nil {
		t.Error("Logger should not be nil")
	}

	if handler.metadata.ChargePointID != "test-charger-001" {
		t.Errorf("ChargePointID = %v, want test-charger-001", handler.metadata.ChargePointID)
	}
}

func TestHandlers_BootNotification(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.BootNotification{
		ChargePointVendor:       "TestVendor",
		ChargePointModel:        "TestModel",
		ChargePointSerialNumber: "SN12345",
	}

	resp, err := handler.BootNotification(req)
	if err != nil {
		t.Fatalf("BootNotification() error = %v", err)
	}

	bootResp, ok := resp.(*cpresp.BootNotification)
	if !ok {
		t.Fatal("Response is not *cpresp.BootNotification")
	}

	if bootResp.Status != "Accepted" {
		t.Errorf("Status = %v, want Accepted", bootResp.Status)
	}

	if bootResp.Interval != 60 {
		t.Errorf("Interval = %v, want 60", bootResp.Interval)
	}

	if bootResp.CurrentTime.IsZero() {
		t.Error("CurrentTime should not be zero")
	}
}

func TestHandlers_Authorize(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.Authorize{
		IdTag: "RFID-12345",
	}

	resp, err := handler.Authorize(req)
	if err != nil {
		t.Fatalf("Authorize() error = %v", err)
	}

	authResp, ok := resp.(*cpresp.Authorize)
	if !ok {
		t.Fatal("Response is not *cpresp.Authorize")
	}

	if authResp.IdTagInfo == nil {
		t.Fatal("IdTagInfo should not be nil")
	}

	if authResp.IdTagInfo.Status != "Accepted" {
		t.Errorf("Status = %v, want Accepted", authResp.IdTagInfo.Status)
	}
}

func TestHandlers_Heartbeat(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.Heartbeat{}

	resp, err := handler.Heartbeart(req)
	if err != nil {
		t.Fatalf("Heartbeat() error = %v", err)
	}

	heartbeatResp, ok := resp.(*cpresp.Heartbeat)
	if !ok {
		t.Fatal("Response is not *cpresp.Heartbeat")
	}

	if heartbeatResp.CurrentTime.IsZero() {
		t.Error("CurrentTime should not be zero")
	}

	if time.Since(heartbeatResp.CurrentTime) > time.Second {
		t.Error("CurrentTime should be close to now")
	}
}

func TestHandlers_StatusNotification(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.StatusNotification{
		ConnectorId: 1,
		Status:      "Available",
		ErrorCode:   "NoError",
	}

	resp, err := handler.StatusNotification(req)
	if err != nil {
		t.Fatalf("StatusNotification() error = %v", err)
	}

	_, ok := resp.(*cpresp.StatusNotification)
	if !ok {
		t.Fatal("Response is not *cpresp.StatusNotification")
	}
}

func TestHandlers_MeterValues(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.MeterValues{
		ConnectorId:   1,
		TransactionId: 123,
	}

	resp, err := handler.MeterValues(req)
	if err != nil {
		t.Fatalf("MeterValues() error = %v", err)
	}

	_, ok := resp.(*cpresp.MeterValues)
	if !ok {
		t.Fatal("Response is not *cpresp.MeterValues")
	}
}

func TestHandlers_StopTransaction(t *testing.T) {
	handler := setupTestHandler()

	req := &cpreq.StopTransaction{
		TransactionId: 123,
		Timestamp:     time.Now(),
		MeterStop:     5000,
		Reason:        "Local",
	}

	resp, err := handler.StopTransaction(req)
	if err != nil {
		t.Fatalf("StopTransaction() error = %v", err)
	}

	stopResp, ok := resp.(*cpresp.StopTransaction)
	if !ok {
		t.Fatal("Response is not *cpresp.StopTransaction")
	}

	if stopResp.IdTagInfo == nil {
		t.Fatal("IdTagInfo should not be nil")
	}

	if stopResp.IdTagInfo.Status != "Accepted" {
		t.Errorf("Status = %v, want Accepted", stopResp.IdTagInfo.Status)
	}
}
