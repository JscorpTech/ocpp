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
	Cfg    *config.Config
	Ctx    context.Context
	Logger *zap.Logger
	Redis  *redis.Client
}

func NewServer(ctx context.Context, cfg *config.Config, logger *zap.Logger, rdb *redis.Client) *Server {
	return &Server{
		Cfg:    cfg,
		Ctx:    ctx,
		Logger: logger,
		Redis:  rdb,
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

func (s *Server) Run() {
	ocpp.SetDebugLogger(log.New(os.Stdout, "DEBUG:", log.Ltime))
	ocpp.SetErrorLogger(log.New(os.Stderr, "ERROR:", log.Ltime))
	csys := cs.New()

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
			s.Logger.Warn("Invalid request body", zap.Any("req", string(dataByte)))
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
			s.Logger.Error("Station error", zap.Error(err))
			return
		}
		switch req.Command {
		case domain.RemoteStartTransaction:
			var data domain.RemoteStartTransactionReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.RemoteStartTransaction{
				IdTag: data.Tag,
			})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Internal server error"}, http.StatusBadRequest)
				s.Logger.Error("remote start transaction error", zap.Error(err))
			}
			res := resp.(*csresp.RemoteStartTransaction)
			writeJson(w, domain.RemoteStartTransactionRes{Status: res.Status}, http.StatusOK)
		case domain.RemoteStopTransaction:
			var data domain.RemoteStopTransactionReq
			json.Unmarshal(req.Data, &data)
			resp, err := station.Send(&csreq.RemoteStopTransaction{TransactionId: data.TransactionId})
			if err != nil {
				writeJson(w, domain.ErrorResponse{Detail: "Internal server error"}, http.StatusBadRequest)
				s.Logger.Error("remote stop transaction error", zap.Error(err))
			}
			res := resp.(*csresp.RemoteStopTransaction)
			writeJson(w, domain.RemoteStopTransactionRes{Status: res.Status}, http.StatusOK)
		default:
			writeJson(w, domain.ErrorResponse{Detail: "Invalid command"}, http.StatusBadRequest)
			s.Logger.Info("Invalid command")
		}
	})

	go csys.Run(s.Cfg.Addr, func(req cpreq.ChargePointRequest, metadata cs.ChargePointRequestMetadata) (cpresp.ChargePointResponse, error) {
		handler := NewHandler(s.Ctx, s.Logger, s.Redis, metadata, s.Cfg)
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
