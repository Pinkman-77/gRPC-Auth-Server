package grpcapp

import (
	"fmt"
	"github.com/Pinkman-77/grpc-auth/pkg/service/auth"
	"log/slog"
	"net"

	authgRPC "github.com/Pinkman-77/grpc-auth/pkg/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService *auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgRPC.Register(gRPCServer, authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const operation = "grpcapp.Run"

	log := a.log.With(
		slog.String("operation", operation),
		slog.Int("port", a.port),
	)

	log.Info("starting grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func (a *App) Stop() {
	const operation = "grpcapp.Stop"

	a.log.With(slog.String("operation", operation)).Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
