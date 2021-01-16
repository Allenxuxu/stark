package rest

import (
	"time"
)

type ClientOptions struct {
	Timeout time.Duration
	Scheme  string
}

type ClientOption func(*ClientOptions)

func Timeout(t time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = t
	}
}

func Scheme(s string) ClientOption {
	return func(o *ClientOptions) {
		o.Scheme = s
	}
}
