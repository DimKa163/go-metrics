package keeper

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	services *ServiceContainer
	listener net.Listener
	*grpc.Server
}

func NewGRPCServer(listener net.Listener, server *grpc.Server, services *ServiceContainer) *GRPCServer {
	return &GRPCServer{
		Server:   server,
		listener: listener,
		services: services,
	}
}
func (gs *GRPCServer) ListenAndServe() error {
	logging.Log.Info("Starting gRPC Server", zap.String("address", gs.listener.Addr().String()))
	return gs.Serve(gs.listener)
}

func (gs *GRPCServer) Map() {
	gs.services.GrpcService.Register(gs.Server)
}

func (gs *GRPCServer) Shutdown(_ context.Context) error {
	gs.GracefulStop()
	return nil
}
