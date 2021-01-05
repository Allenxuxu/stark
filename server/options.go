package server

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Options struct {
	Metadata map[string]string
	Name     string
	Address  string
	Id       string
	Version  string

	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration

	GrpcOpts []grpc.ServerOption
}

type Option func(*Options)

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Id Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Metadata associated with the server
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// RegisterCheck run func before registry service
func RegisterCheck(fn func(context.Context) error) Option {
	return func(o *Options) {
		o.RegisterCheck = fn
	}
}

// RegisterTTL register the service with a TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// RegisterInterval register the service with at interval
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// UnaryServerInterceptor to be used to configure gRPC options
func UnaryServerInterceptor(u grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.ChainUnaryInterceptor(u))
	}
}

// StreamServerInterceptor to be used to configure gRPC options
func StreamServerInterceptor(u grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.ChainStreamInterceptor(u))
	}
}
