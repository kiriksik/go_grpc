package grpc_app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/kiriksik/go_grpc/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService auth.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	auth.RegisterAuthServer(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic("cant start gRPC server:" + err.Error())
	}
}
func (a *App) Run() error {
	const op = "grpc_app.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:, %w", op, err)
	}

	log.Info("gRPC server started", slog.String("addres", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s:, %w", op, err)
	}

	return nil
}

func (a *App) Stop() error {
	const op = "grpc_app.Stop"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	a.gRPCServer.GracefulStop()
	log.Info("stopping gRPC server", slog.Int("port", a.port))

	return nil
}
