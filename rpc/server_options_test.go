package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryServerInterceptor(t *testing.T) {
	opts := ServerOptions{}
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return nil, nil
	}

	n := 10
	for i := 0; i < n; i++ {
		opt := UnaryServerInterceptor(f)
		opt(&opts)
	}

	assert.Equal(t, len(opts.GrpcOpts), n)
}

func TestStreamServerInterceptor(t *testing.T) {
	opts := ServerOptions{}
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}

	n := 10
	for i := 0; i < n; i++ {
		opt := StreamServerInterceptor(f)
		opt(&opts)
	}

	assert.Equal(t, len(opts.GrpcOpts), n)
}

func TestGrpcOptions(t *testing.T) {
	opts := ServerOptions{}
	n := 10

	f1 := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
	for i := 0; i < n; i++ {
		opt := StreamServerInterceptor(f1)
		opt(&opts)
	}

	f2 := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return nil, nil
	}
	for i := 0; i < n; i++ {
		opt := UnaryServerInterceptor(f2)
		opt(&opts)
	}

	assert.Equal(t, len(opts.GrpcOpts), n*2)
}
