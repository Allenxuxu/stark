package selector

import (
	"context"

	"github.com/Allenxuxu/stark/registry"
)

type Options struct {
	Strategy Strategy
	Filters  []registry.Filter

	Context context.Context
}

type Option func(*Options)

func Filter(fn ...registry.Filter) Option {
	return func(o *Options) {
		o.Filters = append(o.Filters, fn...)
	}
}

func BalanceStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}
