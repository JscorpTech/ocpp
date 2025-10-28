package ocpp

import (
	"time"

	"github.com/voltbras/go-ocpp/messages/v1x/cpreq"
	"github.com/voltbras/go-ocpp/messages/v1x/cpresp"
	"go.uber.org/zap"
)

type Handlers struct {
	Logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handlers {
	return &Handlers{
		Logger: logger,
	}
}

func (h *Handlers) MeterValues(req *cpreq.MeterValues) (cpresp.ChargePointResponse, error) {
	return &cpresp.MeterValues{}, nil
}

func (h *Handlers) StartTransaction(req *cpreq.StartTransaction) (cpresp.ChargePointResponse, error) {
	return &cpresp.StartTransaction{
		IdTagInfo: &cpresp.IdTagInfo{
			Status: "Accepted",
		},
	}, nil
}

func (h *Handlers) StopTransaction(req *cpreq.StopTransaction) (cpresp.ChargePointResponse, error) {
	return &cpresp.StartTransaction{
		IdTagInfo: &cpresp.IdTagInfo{
			Status: "Accepted",
		},
	}, nil
}
func (h *Handlers) Heartbeart(req *cpreq.Heartbeat) (cpresp.ChargePointResponse, error) {
	return &cpresp.Heartbeat{CurrentTime: time.Now()}, nil
}

func (h *Handlers) StatusNotification(req *cpreq.StatusNotification) (cpresp.ChargePointResponse, error) {
	h.Logger.Info("Keldi bratishka")
	return &cpresp.StatusNotification{}, nil
}

func (h *Handlers) Authorize(req *cpreq.Authorize) (cpresp.ChargePointResponse, error) {
	h.Logger.Info("salom")
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
