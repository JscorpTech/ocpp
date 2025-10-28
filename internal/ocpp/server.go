package ocpp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/voltbras/go-ocpp"
	"github.com/voltbras/go-ocpp/cs"
	"github.com/voltbras/go-ocpp/messages/v1x/cpreq"
	"github.com/voltbras/go-ocpp/messages/v1x/cpresp"
	"go.uber.org/zap"
)

type Server struct {
	Cfg    *config.Config
	Ctx    context.Context
	Logger *zap.Logger
	Redis  *redis.Client
}

type RemoteCommand struct {
	CpID string         `json:"CpID"`
	Data map[string]any `json:"data"`
}

func NewServer(ctx context.Context, cfg *config.Config, logger *zap.Logger, rdb *redis.Client) *Server {
	return &Server{
		Cfg:    cfg,
		Ctx:    ctx,
		Logger: logger,
		Redis:  rdb,
	}
}

func (s *Server) Run() {
	ocpp.SetDebugLogger(log.New(os.Stdout, "DEBUG:", log.Ltime))
	ocpp.SetErrorLogger(log.New(os.Stderr, "ERROR:", log.Ltime))
	csys := cs.New()
	go s.RemoteCommandWorker(csys)
	go csys.Run(s.Cfg.Addr, func(req cpreq.ChargePointRequest, metadata cs.ChargePointRequestMetadata) (cpresp.ChargePointResponse, error) {
		fmt.Printf("EXAMPLE(MAIN): Request from %s\n", metadata.ChargePointID)
		handler := NewHandler(s.Logger)
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
		default:
			fmt.Printf("EXAMPLE(MAIN): action not supported: %s\n", req.Action())
			return nil, errors.New("Response not supported")
		}
	})
}

func (s *Server) RemoteCommandWorker(csys cs.CentralSystem) {
	pubsub := s.Redis.Subscribe(s.Ctx, "commands")
	fmt.Println("Remote command worker ishga tushdi")
	for {
		message, err := pubsub.ReceiveMessage(s.Ctx)
		if err != nil {
			s.Logger.Error("remote command message receive error", zap.Error(err))
		}
		var data RemoteCommand
		if err := json.Unmarshal([]byte(message.Payload), &data); err != nil {
			s.Logger.Error("Remote command message decode error", zap.Error(err))
		}
		csys.WriteJson(data.CpID, data.Data)
	}
}
