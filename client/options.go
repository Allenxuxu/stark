package client

import (
	"time"

	"google.golang.org/grpc"
)

type Options struct {
	Timeout time.Duration

	GrpcOpts []grpc.DialOption
}

type Option func(*Options)

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func UnaryClientInterceptor(u grpc.UnaryClientInterceptor) Option {
	return func(o *Options) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.WithChainUnaryInterceptor(u))
	}
}

func StreamClientInterceptors(u grpc.StreamClientInterceptor) Option {
	return func(o *Options) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.WithChainStreamInterceptor(u))
	}
}

func DialOption(do ...grpc.DialOption) Option {
	return func(o *Options) {
		o.GrpcOpts = append(o.GrpcOpts, do...)
	}
}
