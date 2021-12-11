package registry

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/registry/cache"
	sr "github.com/Allenxuxu/stark/rpc/client/resolver"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"google.golang.org/grpc/resolver"
)

const scheme = "stark-registry"

var _selector atomic.Value

func registerSelector(s selector.Selector) {
	_selector.Store(s)
}

func init() {
	resolver.Register(sr.NewBuilder(scheme, &_selector))
}

type registrySelector struct {
	opts selector.Options
	rc   cache.Cache
}

func NewSelector(rg registry.Registry, opt ...selector.Option) (selector.Selector, error) {
	var opts selector.Options
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

	// fixme do better
	registerSelector(s)

	return s, nil
}

func (c *registrySelector) Options() selector.Options {
	return c.opts
}

func (c *registrySelector) GetService(service string) ([]*registry.Service, error) {
	services, err := c.rc.GetService(service)
	if err != nil {
		if err == registry.ErrNotFound {
			return nil, selector.ErrNotFound
		}
		return nil, err
	}

	for _, filter := range c.opts.Filters {
		services = filter(services)
	}

	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return services, nil
}

func (c *registrySelector) Watch(service string) (registry.Watcher, error) {
	return c.rc.Watch(registry.WatchService(service))
}

// Close stops the watcher and destroys the cache
func (c *registrySelector) Close() error {
	c.rc.Stop()

	return nil
}

func (c *registrySelector) Address(service string) string {
	return fmt.Sprintf("%s:///%s", scheme, service)
}

func (c *registrySelector) String() string {
	return "registry"
}
