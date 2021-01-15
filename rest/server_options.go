package rest

import (
	"time"
)

type ServerOptions struct {
	Name     string
	Version  string
	Id       string
	Metadata map[string]string
	Address  string
	CertFile string
	KeyFile  string

	RegisterTTL      time.Duration
	RegisterInterval time.Duration
}

type ServerOption func(o *ServerOptions)

// Server name
func Name(n string) ServerOption {
	return func(o *ServerOptions) {
		o.Name = n
	}
}

// Unique server id
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

// Metadata associated with the service
func Metadata(md map[string]string) ServerOption {
	return func(o *ServerOptions) {
		o.Metadata = md
	}
}

// Address to bind to - host:port
func Address(a string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = a
	}
}

func RegisterTTL(t time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.RegisterTTL = t
	}
}

func RegisterInterval(t time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.RegisterInterval = t
	}
}

func TLS(certFile, keyFile string) ServerOption {
	return func(o *ServerOptions) {
		o.CertFile = certFile
		o.KeyFile = keyFile
	}
}
