package main

import (
	"fmt"
	"github.com/Pinkman-77/grpc-auth/pkg/app"
	"github.com/Pinkman-77/grpc-auth/pkg/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogs(cfg.Env)

	log.Info(
		"starting the app",
		slog.Any("cfg", cfg.Env),
		slog.Int("Port", cfg.Grpc.Port),
	)

	application := app.NewApp(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTll)

	if err := application.Run(); err != nil {
		log.Error("failed to start gRPC server", slog.Any("error", err))
		os.Exit(1)
	}

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warm message")

	fmt.Println(cfg)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sign := <-quit

	log.Info("signal received", slog.Any("sign", sign))

	application.Stop()

}
func setupLogs(env string) *slog.Logger {

	const (
		envLocal = "local"
		envDev   = "dev"
		envProd  = "prod"
	)

	var log *slog.Logger

	switch env {
	case envLocal:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
