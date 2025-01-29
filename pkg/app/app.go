package app

import (
	"github.com/Pinkman-77/grpc-auth/pkg/service/auth"
	"github.com/Pinkman-77/grpc-auth/pkg/storage/sqlite"
	"log/slog"
	"time"

	grpcapp "github.com/Pinkman-77/grpc-auth/pkg/app/grpc"
)

type gRPCApp struct {
	gRPCServer *grpcapp.App
}

func NewApp(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTLL time.Duration,
) *gRPCApp {
	storage, err := sqlite.New(storagePath)

	if err != nil {
		panic(err)
	}

	authService := auth.NewAuth(log, storage, storage, storage, tokenTLL)

	gRPCServer := grpcapp.New(log, authService, port)

	return &gRPCApp{
		gRPCServer: gRPCServer,
	}
}

func (a *gRPCApp) Run() error {
	return a.gRPCServer.Run()
}

func (a *gRPCApp) Stop() {
	a.gRPCServer.Stop()
}
