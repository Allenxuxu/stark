package registry

import (
	"context"
	"time"

	"github.com/Allenxuxu/stark/rpc/client/selector"
)

const ttlKey = "selector_ttl"

// Set the registry cache ttl
func TTL(t time.Duration) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ttlKey, t)
	}
}
