package registry

import (
	"time"

	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/Allenxuxu/stark/pkg/registry/cache"
	"github.com/Allenxuxu/stark/pkg/registry/mdns"
	"github.com/Allenxuxu/stark/pkg/selector"
)

type registrySelector struct {
	opts selector.Options
	rc   cache.Cache
}

func (c *registrySelector) newCache() cache.Cache {
	opts := make([]cache.Option, 0, 1)
	if c.opts.Context != nil {
		if t, ok := c.opts.Context.Value("selector_ttl").(time.Duration); ok {
			opts = append(opts, cache.WithTTL(t))
		}
	}
	return cache.New(c.opts.Registry, opts...)
}

func (c *registrySelector) Options() selector.Options {
	return c.opts
}

func (c *registrySelector) Next(service string, opts ...selector.SelectOption) (*registry.Node, error) {
	sopts := selector.SelectOptions{
		Strategy: c.opts.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := c.rc.GetService(service)
	if err != nil {
		if err == registry.ErrNotFound {
			return nil, selector.ErrNotFound
		}
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services)
}

func (c *registrySelector) Mark(service string, node *registry.Node, err error) {
}

func (c *registrySelector) Reset(service string) {
}

// Close stops the watcher and destroys the cache
func (c *registrySelector) Close() error {
	c.rc.Stop()

	return nil
}

func (c *registrySelector) String() string {
	return "registry"
}

func NewSelector(opts ...selector.Option) (selector.Selector, error) {
	var err error
	sopts := selector.Options{
		Strategy: selector.Random(),
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry, err = mdns.NewRegistry()
		if err != nil {
			return nil, err
		}
	}

	s := &registrySelector{
		opts: sopts,
	}
	s.rc = s.newCache()

	return s, nil
}
