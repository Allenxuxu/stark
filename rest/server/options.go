package server

import (
	"time"
)

type Options struct {
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

type Option func(o *Options)

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Unique server id
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

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

func TLS(certFile, keyFile string) Option {
	return func(o *Options) {
		o.CertFile = certFile
		o.KeyFile = keyFile
	}
}
