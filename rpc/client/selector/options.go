package selector

import (
	"context"
	"time"
)

type Options struct {
	TTL      time.Duration
	Balancer string
	Filters  []Filter

	Context context.Context
}

type Option func(*Options)

func WithBalancerName(b string) Option {
	return func(o *Options) {
		o.Balancer = b
	}
}

func WithFilter(fn ...Filter) Option {
	return func(o *Options) {
		o.Filters = append(o.Filters, fn...)
	}
}
