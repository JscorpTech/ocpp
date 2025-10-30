package services

import (
	"context"
	"encoding/json"

	"github.com/JscorpTech/ocpp/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type EventService interface {
	SendEvent(context.Context, *redis.Client, domain.Event, *zap.Logger)
}

type eventService struct{}

func NewEventService() *eventService {
	return &eventService{}
}

func (e *eventService) SendEvent(ctx context.Context, rdb *redis.Client, event domain.Event, log *zap.Logger) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Error("Event encode error", zap.Error(err))
	}
	rdb.RPush(ctx, "events", payload)
}
