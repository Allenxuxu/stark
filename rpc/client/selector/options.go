package selector

import (
	"context"

	"github.com/Allenxuxu/stark/registry"
)

type Options struct {
	Balancer string
	Filters  []registry.Filter

	Context context.Context
}

type Option func(*Options)

func BalancerName(b string) Option {
	return func(o *Options) {
		o.Balancer = b
	}
}

func Filter(fn ...registry.Filter) Option {
	return func(o *Options) {
		o.Filters = append(o.Filters, fn...)
	}
}
