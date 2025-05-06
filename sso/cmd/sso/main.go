package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiriksik/go_grpc/internal/app"
	"github.com/kiriksik/go_grpc/internal/config"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	// fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.Any("config", cfg))

	ssoApp := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go ssoApp.GRPCSrv.MustRun() //gRPC

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	ssoApp.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		{
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		}
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
