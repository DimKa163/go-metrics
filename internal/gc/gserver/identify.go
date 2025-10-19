package gserver

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	metadata2 "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

func UnaryIdentifyInterceptor(ipNet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		metadata, ok := metadata2.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "missing header Real-IP")
		}
		val := metadata.Get("X-Real-IP")
		if len(val) == 0 {
			return nil, status.Errorf(codes.PermissionDenied, "missing header Real-IP")
		}
		ip := net.ParseIP(val[0])

		if !ipNet.Contains(ip) {
			return nil, status.Errorf(codes.PermissionDenied, "not trustable IP")
		}
		return handler(ctx, req)
	}
}
