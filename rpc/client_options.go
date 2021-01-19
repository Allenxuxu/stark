package rpc

import (
	"time"

	"google.golang.org/grpc"
)

type ClientOptions struct {
	Timeout time.Duration

	GrpcOpts []grpc.DialOption
}

type ClientOption func(*ClientOptions)

func Timeout(t time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = t
	}
}

func UnaryClientInterceptor(u ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.WithChainUnaryInterceptor(u...))
	}
}

func StreamClientInterceptors(u ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.WithChainStreamInterceptor(u...))
	}
}

func GrpcDialOption(do ...grpc.DialOption) ClientOption {
	return func(o *ClientOptions) {
		o.GrpcOpts = append(o.GrpcOpts, do...)
	}
}
