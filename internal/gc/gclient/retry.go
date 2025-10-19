package gclient

import (
	"context"
	"github.com/cenkalti/backoff/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRetryInterceptor() grpc.UnaryClientInterceptor {
	times := [3]int{1, 3, 5}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 0
		_, err := backoff.Retry(ctx, func() (bool, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return true, nil
			}
			if shouldRetry(err) && len(times) > attempt {
				at := attempt
				attempt++
				return false, backoff.RetryAfter(times[at])
			}
			return false, backoff.Permanent(err)
		})
		return err
	}
}

func shouldRetry(err error) bool {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unavailable:
			return true
		case codes.DeadlineExceeded:
			return true
		case codes.Aborted:
			return true
		case codes.Internal:
			return true
		}
	}
	return false
}
