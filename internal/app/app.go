package app

import (
	"errors"
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

type HttpServer struct {
	log         *slog.Logger
	ctx         context.Context
	http        *http.Server
	grpcAddress string
	httpAddress string
}

func New(
	cnf *config.Config,
	log *slog.Logger,
) *HttpServer {
	grpcAddress := getServerAddress(cnf.GRPC.Address, cnf.GRPC.Port)
	httpAddress := getServerAddress(cnf.HTTP.Address, cnf.HTTP.Port)

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := desc.RegisterHashServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}
	return &HttpServer{
		log:         log,
		ctx:         ctx,
		http:        httpServer,
		grpcAddress: grpcAddress,
		httpAddress: httpAddress,
	}
}

func (srv *HttpServer) MustRun() {
	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (srv *HttpServer) Run() error {
	srv.log.Info("http server listening at", slog.Any("httpAddress", srv.httpAddress))
	return srv.http.ListenAndServe()
}

func (srv *HttpServer) Stop() {
	const op = "httpServer.Stop"
	srv.log.With(slog.String("op", op)).
		Info("stopping HTTP server", slog.String("port", srv.httpAddress))

	if err := srv.http.Shutdown(srv.ctx); err != nil {
		srv.log.Info("Server Shutdown Failed:", slog.Any("error", err))
	}
	srv.log.Info("Graceful shutdown complete.")
}

func getServerAddress(address string, port int) string {
	var serverAddress string
	serverAddress = address + fmt.Sprintf(":%d", port)
	return serverAddress
}
