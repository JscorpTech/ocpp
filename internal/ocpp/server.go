package ocpp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/domain"
	"github.com/JscorpTech/ocpp/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/voltbras/go-ocpp"
	"github.com/voltbras/go-ocpp/cs"
	"github.com/voltbras/go-ocpp/messages/v1x/cpreq"
	"github.com/voltbras/go-ocpp/messages/v1x/cpresp"
	"github.com/voltbras/go-ocpp/messages/v1x/csreq"
	"github.com/voltbras/go-ocpp/messages/v1x/csresp"
	"go.uber.org/zap"
)

type Server struct {
	cfg   *config.Config
	ctx   context.Context
	log   *zap.Logger
	redis *redis.Client
	event services.EventService
}

func NewServer(ctx context.Context, cfg *config.Config, logger *zap.Logger, rdb *redis.Client) *Server {
	return &Server{
		cfg:   cfg,
		ctx:   ctx,
		log:   logger,
		redis: rdb,
		event: services.NewEventService(),
	}
}

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) Validate(w http.ResponseWriter, req domain.RemoteCommandReq) string {
	if req.Command == "" {
		return "Command required"
	}
	if req.CpID == "" {
		return "CpId required"
	}
	if len(req.Data) == 0 {
		return "Data required"
	}
	return ""
}

func (s *Server) Run() error {
	ocpp.SetDebugLogger(log.New(os.Stdout, "DEBUG:", log.Ltime))
	ocpp.SetErrorLogger(log.New(os.Stderr, "ERROR:", log.Ltime))
	csys := cs.New()

	csys.SetChargePointDisconnectionListener(func(CpID string, host string) {
		event := domain.Event{
			Domain: host,
			Event:  domain.DisconnectChargerEvent,
			Data: domain.DisconnectCharger{
				Charger: CpID,
			},
		}
		s.event.SendEvent(s.ctx, s.redis, &event, s.log)
	})

	csys.SetChargePointConnectionListener(func(CpID string, host string) {
		event := domain.Event{
			Domain: host,
			Event:  domain.ConnectChargerEvent,
			Data: domain.ConnectCharger{
				Charger: CpID,
			},
		}
		s.event.SendEvent(s.ctx, s.redis, &event, s.log)
	})

	http.HandleFunc("/command/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeJson(w, domain.ErrorResponse{Detail: "Invalid Method " + r.Method}, http.StatusBadRequest)
			return
		}
		var req domain.RemoteCommandReq
		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			writeJson(w, domain.ErrorResponse{Detail: "Invalid request"}, http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(dataByte, &req); err != nil {
			s.log.Warn("Invalid request body", zap.Any("req", string(dataByte)))
			writeJson(w, domain.ErrorResponse{Detail: "Invalid request body " + err.Error()}, http.StatusBadRequest)
			return
		}
		if res := s.Validate(w, req); res != "" {
			writeJson(w, domain.ErrorResponse{Detail: res}, http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		station, err := csys.GetServiceOf(req.CpID, ocpp.V16, "")
		if err != nil {
			writeJson(w, domain.ErrorResponse{Detail: "Charger not connected"}, http.StatusBadRequest)
			s.log.Error("Station error", zap.Error(err))
			return
		}
		switch req.Command {
		case domain.RemoteStartTransaction:
			var data domain.RemoteStartTransactionReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.RemoteStartTransaction{
				IdTag:       data.Tag,
				ConnectorId: data.ConnectorID,
			})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Internal server error"}, http.StatusBadRequest)
				s.log.Error("remote start transaction error", zap.Error(err))
			}
			res := resp.(*csresp.RemoteStartTransaction)
			writeJson(w, domain.RemoteStartTransactionRes{Status: res.Status}, http.StatusOK)
		case domain.RemoteStopTransaction:
			var data domain.RemoteStopTransactionReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.RemoteStopTransaction{TransactionId: data.TransactionId})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Internal server error"}, http.StatusBadRequest)
				s.log.Error("remote stop transaction error", zap.Error(err))
			}
			res := resp.(*csresp.RemoteStopTransaction)
			writeJson(w, domain.RemoteStopTransactionRes{Status: res.Status}, http.StatusOK)
		case domain.GetConfiguration:
			var data domain.GetConfigurationReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.GetConfiguration{Key: data.Key})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Error"}, http.StatusBadRequest)
			}
			res := resp.(*csresp.GetConfiguration)
			writeJson(w, res, http.StatusOK)
		case domain.ChangeConfiguration:
			var data domain.ChangeConfigurationReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.ChangeConfiguration{Key: data.Key, Value: data.Value})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Error"}, http.StatusBadRequest)
			}
			res := resp.(*csresp.ChangeConfiguration)
			writeJson(w, res, http.StatusOK)
		default:
			writeJson(w, domain.ErrorResponse{Detail: "Invalid command"}, http.StatusBadRequest)
			s.log.Info("Invalid command")
		}
	})

	return csys.Run(s.cfg.Addr, func(req cpreq.ChargePointRequest, metadata cs.ChargePointRequestMetadata) (cpresp.ChargePointResponse, error) {
		handler := NewHandler(s.ctx, s.log, s.redis, metadata, s.cfg, s.event)
		switch req := req.(type) {
		case *cpreq.BootNotification:
			return handler.BootNotification(req)
		case *cpreq.StatusNotification:
			return handler.StatusNotification(req)
		case *cpreq.Authorize:
			return handler.Authorize(req)
		case *cpreq.Heartbeat:
			return handler.Heartbeart(req)
		case *cpreq.MeterValues:
			return handler.MeterValues(req)
		case *cpreq.StartTransaction:
			return handler.StartTransaction(req)
		case *cpreq.StopTransaction:
			return handler.StopTransaction(req)
		case *cpreq.DataTransfer:
			return handler.DataTransfer(req)
		default:
			fmt.Printf("EXAMPLE(MAIN): action not supported: %s\n", req.Action())
			return nil, errors.New("Response not supported")
		}
	})
}
