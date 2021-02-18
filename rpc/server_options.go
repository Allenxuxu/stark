package rpc

import (
	"time"

	"google.golang.org/grpc"
)

const (
	metaDataMetricsAddressKey = "metrics_address"
)

type ServerOptions struct {
	Metadata map[string]string
	Name     string
	Address  string
	Id       string
	Version  string

	MetricsPath      string
	RegisterTTL      time.Duration
	RegisterInterval time.Duration
	GrpcOpts         []grpc.ServerOption
}

type ServerOption func(*ServerOptions)

// Server name
func Name(n string) ServerOption {
	return func(o *ServerOptions) {
		o.Name = n
	}
}

// Id Unique server id
func Id(id string) ServerOption {
	return func(o *ServerOptions) {
		o.Id = id
	}
}

// Version of the service
func Version(v string) ServerOption {
	return func(o *ServerOptions) {
		o.Version = v
	}
}

// Address to bind to - host:port
func Address(a string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = a
	}
}

// Metadata associated with the server
func Metadata(md map[string]string) ServerOption {
	return func(o *ServerOptions) {
		o.Metadata = md
	}
}

func MetricsAddress(a string) ServerOption {
	return func(o *ServerOptions) {
		if o.Metadata == nil {
			o.Metadata = make(map[string]string)
		}

		o.Metadata[metaDataMetricsAddressKey] = a
	}
}

func MetricsPath(p string) ServerOption {
	return func(o *ServerOptions) {
		o.MetricsPath = p
	}
}

// RegisterTTL register the service with a TTL
func RegisterTTL(t time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.RegisterTTL = t
	}
}

// RegisterInterval register the service with at interval
func RegisterInterval(t time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.RegisterInterval = t
	}
}

// UnaryServerInterceptor to be used to configure gRPC options
func UnaryServerInterceptor(u ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.ChainUnaryInterceptor(u...))
	}
}

// StreamServerInterceptor to be used to configure gRPC options
func StreamServerInterceptor(u ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.GrpcOpts = append(o.GrpcOpts, grpc.ChainStreamInterceptor(u...))
	}
}
