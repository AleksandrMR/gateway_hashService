package main

import (
	"context"
	"fmt"
	"github.com/AleksandrMR/gateway_hashService/internal/config"
	desc "github.com/AleksandrMR/proto_hashService/gen/hashService_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	conf := config.MustLoad()
	log := setupLogger(conf.Env)
	log.Info("starting application", slog.Any("conf", conf))

	ctx := context.Background()
	grpcAddress := getServerAddress(conf.GRPC.Address, conf.GRPC.Port)
	httpAddress := getServerAddress(conf.HTTP.Address, conf.HTTP.Port)
	if err := startHttpServer(ctx, grpcAddress, httpAddress); err != nil {
		panic(err)
	}
	log.Info("http server listening at %v\n", slog.Any("httpAddress", httpAddress))
}

func startHttpServer(
	ctx context.Context,
	grpcAddress string,
	httpAddress string,
) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := desc.RegisterHashServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(httpAddress, mux)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func getServerAddress(address string, port int) string {
	var serverAddress string
	serverAddress = address + fmt.Sprintf(":%d", port)
	return serverAddress
}
