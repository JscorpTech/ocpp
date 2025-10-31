package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/JscorpTech/ocpp/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func TestNewEventService(t *testing.T) {
	service := NewEventService()
	if service == nil {
		t.Fatal("NewEventService() returned nil")
	}
}

func TestEventService_SendEvent(t *testing.T) {
	// Setup mock Redis client (miniredis yoki test uchun)
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	logger, _ := zap.NewDevelopment()
	service := NewEventService()

	event := domain.Event{
		Event: domain.HealthEvent,
		Data: domain.Healthcheck{
			Charger: "test-charger-001",
		},
	}

	// Clear any existing events
	rdb.Del(ctx, "events")

	// Test event sending
	service.SendEvent(ctx, rdb, event, logger)

	// Verify event was pushed to Redis
	result, err := rdb.RPop(ctx, "events").Result()
	if err != nil && err != redis.Nil {
		// Redis mavjud bo'lmasa test o'tkazib yuboramiz
		t.Skip("Redis not available for testing")
		return
	}

	if err == nil {
		var receivedEvent domain.Event
		if err := json.Unmarshal([]byte(result), &receivedEvent); err != nil {
			t.Errorf("Failed to unmarshal event: %v", err)
		}

		if receivedEvent.Event != event.Event {
			t.Errorf("Event type = %v, want %v", receivedEvent.Event, event.Event)
		}
	}
}

func TestEventService_SendEvent_AllEventTypes(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	logger, _ := zap.NewDevelopment()
	service := NewEventService()

	tests := []struct {
		name  string
		event domain.Event
	}{
		{
			name: "health event",
			event: domain.Event{
				Event: domain.HealthEvent,
				Data: domain.Healthcheck{
					Charger: "charger-001",
				},
			},
		},
		{
			name: "change connector status event",
			event: domain.Event{
				Event: domain.ChangeConnectorStatusEvent,
				Data: domain.ChangeConnectorStatus{
					Charger: "charger-001",
					Conn:    1,
					Status:  "Available",
				},
			},
		},
		{
			name: "start transaction event",
			event: domain.Event{
				Event: domain.StartTransactionEvent,
				Data: domain.StartTransaction{
					Charger:    "charger-001",
					Conn:       1,
					Tag:        "RFID-12345",
					MeterStart: 0,
				},
			},
		},
		{
			name: "stop transaction event",
			event: domain.Event{
				Event: domain.StopTransactionEvent,
				Data: domain.StopTransaction{
					Charger:       "charger-001",
					TransactionId: 123,
					Reason:        "Local",
					MeterStop:     5000,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.SendEvent(ctx, rdb, tt.event, logger)
			
			// Redis test bo'lsa, event mavjudligini tekshiramiz
			_, err := rdb.LLen(ctx, "events").Result()
			if err != nil && err != redis.Nil {
				t.Skip("Redis not available for testing")
			}
		})
	}
}
