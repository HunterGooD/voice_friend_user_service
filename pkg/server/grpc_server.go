package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HunterGooD/voice_friend_user_service/pkg/logger"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	server *grpc.Server
	log    logger.Logger
}

func NewGRPCServer(log logger.Logger, retriesCount int, timeout time.Duration) *GRPCServer {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", map[string]any{"panic": fmt.Sprintf("panic %+v", p)})
			return status.Errorf(codes.Internal, "%s", p)
		}),
	}
	loggingOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(
			// grpclog.StartCall, grpclog.FinishCall,
			grpclog.PayloadReceived, grpclog.PayloadSent,
		),
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		grpclog.UnaryServerInterceptor(interceptorLogger(log), loggingOpts...),
	))
	return &GRPCServer{server, log}
}

func (s *GRPCServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.log.Error("Error init listener: ", map[string]any{
			"error": err,
		})
		return err
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// wait interput signal
		<-quit

		s.log.Info("Initiating graceful shutdown...")
		timer := time.AfterFunc(10*time.Second, func() {
			s.log.Warn("Server couldn't stop gracefully in time. Doing force stop.")
			s.server.Stop()
		})
		defer timer.Stop()

		s.server.GracefulStop()
		s.log.Info("GRPC server stopped")
	}()

	s.log.Info("GRPC server start", map[string]any{
		"grpc_addr": address,
	})
	if err := s.server.Serve(listener); err != nil {
		s.log.Error("Server error: ", map[string]any{
			"error": err,
		})
		return err
	}
	return nil
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

func interceptorLogger(l logger.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, int(lvl), msg, fields...)
	})
}
