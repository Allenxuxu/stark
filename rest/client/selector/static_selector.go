package static

import (
	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/rest/client/selector"
)

type staticSelector struct {
	opts selector.Options

	service []*registry.Service
}

func New(service []*registry.Service, opts ...selector.Option) *staticSelector {
	var options selector.Options
	for _, o := range opts {
		o(&options)
	}

	if options.Strategy == nil {
		options.Strategy = selector.RoundRobin()
	}

	s := &staticSelector{
		service: service,
		opts:    options,
	}

	for _, filter := range s.opts.Filters {
		s.service = filter(s.service)
	}

	return s
}

func (s staticSelector) Options() selector.Options {
	return s.opts
}

func (s staticSelector) Next(service string) (*registry.Node, error) {
	if len(s.service) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return s.opts.Strategy(s.service)
}

func (s staticSelector) Close() error {
	return nil
}

func (s staticSelector) String() string {
	return "static"
}
