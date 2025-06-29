package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	grpc_auth "grpc_auth"

	"github.com/kiriksik/go_grpc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient grpc_auth.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: grpc_auth.NewAuthClient(cc),
	}

}
