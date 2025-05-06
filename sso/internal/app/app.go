package app

import (
	"log/slog"
	"time"

	grpc_app "github.com/kiriksik/go_grpc/internal/app/grpc"
	service "github.com/kiriksik/go_grpc/internal/services/auth"
	"github.com/kiriksik/go_grpc/internal/services/storage/sqlite"
)

type App struct {
	GRPCSrv *grpc_app.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTl time.Duration,
) *App {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := service.New(log, storage, tokenTTl)
	grpcApp := grpc_app.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}

// func
