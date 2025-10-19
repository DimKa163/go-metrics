package gclient

import (
	"context"
	"github.com/DimKa163/go-metrics/internal/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func UnaryLoggingInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			logging.Log.Error("invoker error", zap.Error(err))
		} else {
			logging.Log.Debug("invoker took", zap.Duration("duration", time.Since(start)))
		}
		return err
	}
}
