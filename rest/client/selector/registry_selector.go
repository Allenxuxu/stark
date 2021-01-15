package registry

import (
	"time"

	"github.com/Allenxuxu/stark/log"

	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/registry/cache"
	"github.com/Allenxuxu/stark/rest/client/selector"
)

type registrySelector struct {
	opts selector.Options
	rc   cache.Cache
}

func NewSelector(rg registry.Registry, opt ...selector.Option) (selector.Selector, error) {
	opts := selector.Options{
		Strategy: selector.RoundRobin(),
	}

	for _, opt := range opt {
		opt(&opts)
	}

	cacheOpts := make([]cache.Option, 0, 1)
	if opts.Context != nil {
		if t, ok := opts.Context.Value(ttlKey).(time.Duration); ok {
			cacheOpts = append(cacheOpts, cache.WithTTL(t))
		}
	}

	s := &registrySelector{
		opts: opts,
		rc:   cache.New(rg, cacheOpts...),
	}

	return s, nil
}

func (c *registrySelector) Options() selector.Options {
	return c.opts
}

func (c *registrySelector) Next(service string) (*registry.Node, error) {
	services, err := c.rc.GetService(service)
	if err != nil {
		if err == registry.ErrNotFound {
			return nil, selector.ErrNotFound
		}
		return nil, err
	}

	for _, s := range services {
		for _, node := range s.Nodes {
			log.Info(" registrySelector ", node.Address)
		}
	}

	for _, filter := range c.opts.Filters {
		services = filter(services)
	}

	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return c.opts.Strategy(services)
}

func (c *registrySelector) Close() error {
	c.rc.Stop()

	return nil
}

func (c *registrySelector) String() string {
	return "registry"
}
