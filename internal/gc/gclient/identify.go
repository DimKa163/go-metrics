package gclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

var ipAddr string

func init() {
	ipAddr, _ = getLocalIP()
}
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
func UnaryIdentifyInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.New(map[string]string{"X-Real-IP": ipAddr})
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
	}
}
