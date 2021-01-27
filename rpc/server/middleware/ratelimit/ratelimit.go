package ratelimit

import (
	"context"

	"github.com/Allenxuxu/stark/pkg/limit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(limiter limit.RateLimit) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Allow() {
			return handler(ctx, req)
		}

		return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc rate limit middleware, please retry later.", info.FullMethod)
	}
}

func StreamServerInterceptor(limiter limit.RateLimit) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Allow() {
			return handler(srv, stream)
		}

		return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc rate limit middleware, please retry later.", info.FullMethod)
	}
}
