package gserver

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logging.Log.Info(
			"got incoming grpc request",
			zap.String("method", info.FullMethod),
			zap.Any("req", req),
		)
		startTime := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(startTime)
		if err != nil {
			logging.Log.Warn("Processed with error", zap.String("error", err.Error()))
		}
		logging.Log.Info("grpc request processed", zap.Duration("elapsed", elapsed))
		return resp, err
	}
}
