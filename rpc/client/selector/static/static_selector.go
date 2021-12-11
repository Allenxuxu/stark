// Package static provides a static resolver which returns the name/ip passed in without any change
package static

import (
	"fmt"
	"sync/atomic"

	"github.com/Allenxuxu/stark/registry"
	sr "github.com/Allenxuxu/stark/rpc/client/resolver"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"google.golang.org/grpc/resolver"
)

const scheme = "stark-static"

var _selector atomic.Value

func registerSelector(s selector.Selector) {
	_selector.Store(s)
}

func init() {
	resolver.Register(sr.NewBuilder(scheme, &_selector))
}

// staticSelector is a static selector
type staticSelector struct {
	opts selector.Options

	service []*registry.Service
}

func NewSelector(service []*registry.Service, opts ...selector.Option) selector.Selector {
	var options selector.Options
	for _, o := range opts {
		o(&options)
	}

	s := &staticSelector{
		opts:    options,
		service: service,
	}

	// fixme do better
	registerSelector(s)

	return s
}

func (s *staticSelector) Options() selector.Options {
	return s.opts
}

func (s *staticSelector) GetService(service string) ([]*registry.Service, error) {

	for _, filter := range s.opts.Filters {
		s.service = filter(s.service)
	}

	if len(s.service) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return s.service, nil
}

func (s *staticSelector) Watch(service string) (registry.Watcher, error) {
	return nil, nil
}

func (s *staticSelector) Close() error {
	return nil
}

func (s *staticSelector) Address(service string) string {
	return fmt.Sprintf("%s:///%s", scheme, service)
}

func (s *staticSelector) String() string {
	return "static"
}
