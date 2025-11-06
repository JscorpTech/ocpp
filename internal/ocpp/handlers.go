package ocpp

import (
	"context"
	"time"

	"github.com/JscorpTech/ocpp/internal/client"
	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/domain"
	"github.com/JscorpTech/ocpp/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/voltbras/go-ocpp/cs"
	"github.com/voltbras/go-ocpp/messages/v1x/cpreq"
	"github.com/voltbras/go-ocpp/messages/v1x/cpresp"
	"go.uber.org/zap"
)

type Handlers struct {
	Logger            *zap.Logger
	redis             *redis.Client
	ctx               context.Context
	metadata          cs.ChargePointRequestMetadata
	event             services.EventService
	transactionClient client.TransactionClient
}

func NewHandler(ctx context.Context, logger *zap.Logger, rdb *redis.Client, metadata cs.ChargePointRequestMetadata, cfg *config.Config, event services.EventService) *Handlers {
	return &Handlers{
		Logger:            logger,
		redis:             rdb,
		ctx:               ctx,
		metadata:          metadata,
		event:             event,
		transactionClient: client.NewTransactionClient(cfg),
	}
}

func (h *Handlers) MeterValues(req *cpreq.MeterValues) (cpresp.ChargePointResponse, error) {
	event := domain.Event{
		Event: domain.MeterValuesEvent,
		Data: domain.MeterValues{
			Conn:          req.ConnectorId,
			TransactionId: req.TransactionId,
			MeterValue:    req.MeterValue,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.MeterValues{}, nil
}

func (h *Handlers) StartTransaction(req *cpreq.StartTransaction) (cpresp.ChargePointResponse, error) {
	transaction, err := h.transactionClient.GetTransactionFromTag(req.IdTag)
	if err != nil {
		panic(err)
	}
	event := domain.Event{
		Event: domain.StartTransactionEvent,
		Data: domain.StartTransaction{
			Charger:    h.metadata.ChargePointID,
			Conn:       req.ConnectorId,
			Tag:        req.IdTag,
			MeterStart: req.MeterStart,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.StartTransaction{
		IdTagInfo: &cpresp.IdTagInfo{
			Status: "Accepted",
		},
		TransactionId: int32(transaction.Data.Id),
	}, nil
}

func (h *Handlers) StopTransaction(req *cpreq.StopTransaction) (cpresp.ChargePointResponse, error) {
	event := domain.Event{
		Event: domain.StopTransactionEvent,
		Data: domain.StopTransaction{
			Charger:       h.metadata.ChargePointID,
			TransactionId: req.TransactionId,
			Reason:        req.Reason,
			MeterStop:     req.MeterStop,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.StopTransaction{
		IdTagInfo: &cpresp.IdTagInfo{
			Status: "Accepted",
		},
	}, nil
}

func (h *Handlers) Heartbeart(req *cpreq.Heartbeat) (cpresp.ChargePointResponse, error) {
	event := domain.Event{
		Event: domain.HealthEvent,
		Data: domain.Healthcheck{
			Charger: h.metadata.ChargePointID,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.Heartbeat{CurrentTime: time.Now()}, nil
}

func (h *Handlers) StatusNotification(req *cpreq.StatusNotification) (cpresp.ChargePointResponse, error) {
	event := domain.Event{
		Event: domain.ChangeConnectorStatusEvent,
		Data: domain.ChangeConnectorStatus{
			Charger: h.metadata.ChargePointID,
			Conn:    req.ConnectorId,
			Status:  req.Status,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.StatusNotification{}, nil
}

func (h *Handlers) Authorize(req *cpreq.Authorize) (cpresp.ChargePointResponse, error) {
	return &cpresp.Authorize{IdTagInfo: &cpresp.IdTagInfo{
		Status: "Accepted",
	}}, nil
}

func (h *Handlers) BootNotification(req *cpreq.BootNotification) (cpresp.ChargePointResponse, error) {
	return &cpresp.BootNotification{
		Status:      "Accepted",
		CurrentTime: time.Now(),
		Interval:    60,
	}, nil
}

func (h *Handlers) DataTransfer(req *cpreq.DataTransfer) (cpresp.ChargePointResponse, error) {
	event := domain.Event{
		Event: domain.DataTransferEvent,
		Data: &domain.DataTransfer{
			VendorId:  req.VendorId,
			MessageId: req.MessageId,
			Data:      req.Data,
		},
	}
	h.event.SendEvent(h.ctx, h.redis, &event, h.Logger)
	return &cpresp.DataTransfer{
		Status: "Accepted",
	}, nil
}
