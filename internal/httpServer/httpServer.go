package httpServer

import (
	"fmt"
	"github.com/AleksandrMR/gateway_hashService/internal/config"
	desc "github.com/AleksandrMR/proto_hashService/gen/hashService_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
)

func Start(
	ctx context.Context,
	cnf *config.Config,
	log *slog.Logger,
) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	grpcAddress := getServerAddress(cnf.GRPC.Address, cnf.GRPC.Port)
	httpAddress := getServerAddress(cnf.HTTP.Address, cnf.HTTP.Port)
	err := desc.RegisterHashServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}
	log.Info("http server listening at", slog.Any("httpAddress", httpAddress))
	return http.ListenAndServe(httpAddress, mux)
}

func Stop() {

}

func getServerAddress(address string, port int) string {
	var serverAddress string
	serverAddress = address + fmt.Sprintf(":%d", port)
	return serverAddress
}
