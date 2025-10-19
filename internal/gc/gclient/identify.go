package gclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryIdentifyInterceptor(ip string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.New(map[string]string{"X-Real-IP": ip})
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
	}
}
