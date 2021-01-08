package registry

import (
	"context"
	"time"

	"github.com/Allenxuxu/stark/pkg/selector"
)

// Set the registry cache ttl
func TTL(t time.Duration) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "selector_ttl", t)
	}
}
