// Package static provides a static resolver which returns the name/ip passed in without any change
package static

import (
	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/Allenxuxu/stark/pkg/selector"
)

// staticSelector is a static selector
type staticSelector struct {
	opts selector.Options

	nodes []*registry.Node
}

func (s *staticSelector) Options() selector.Options {
	return s.opts
}

func (s *staticSelector) Next(service string, opts ...selector.SelectOption) (*registry.Node, error) {
	rs := &registry.Service{
		Name:      "",
		Version:   "",
		Metadata:  nil,
		Endpoints: nil,
		Nodes:     s.nodes,
	}
	return s.opts.Strategy([]*registry.Service{rs})
}

func (s *staticSelector) Mark(service string, node *registry.Node, err error) {
	return
}

func (s *staticSelector) Reset(service string) {
	return
}

func (s *staticSelector) Close() error {
	return nil
}

func (s *staticSelector) String() string {
	return "static"
}

func NewSelector(nodes []*registry.Node, opts ...selector.Option) selector.Selector {
	options := selector.Options{
		Strategy: selector.RoundRobin(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &staticSelector{
		opts:  options,
		nodes: nodes,
	}
}
