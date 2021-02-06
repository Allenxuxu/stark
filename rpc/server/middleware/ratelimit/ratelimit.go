package ratelimit

import (
	"context"

	"github.com/Allenxuxu/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(limiter ratelimit.RateLimit) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Allow() {
			return handler(ctx, req)
		}

		return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc rate limit middleware, please retry later.", info.FullMethod)
	}
}

func StreamServerInterceptor(limiter ratelimit.RateLimit) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Allow() {
			return handler(srv, stream)
		}

		return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc rate limit middleware, please retry later.", info.FullMethod)
	}
}
