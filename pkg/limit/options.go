package limit

import (
	"context"
	"time"

	"go.uber.org/atomic"
)

type Options struct {
	Per              time.Duration
	DynamicLimitLoop func(perTime *atomic.Int64, rate int64)
}

type Option func(l *Options)

// Per allows configuring limits for different time windows.
//
// The default window is one second, so New(100) produces a one hundred per
// second (100 Hz) rate limiter.
//
// New(2, Per(60*time.Second)) creates a 2 per minute rate limiter.
func Per(per time.Duration) Option {
	return func(o *Options) {
		o.Per = per
	}
}

func DynamicLimit(ctx context.Context, except float64, current chan float64, minLimit, maxLimit int64) Option {
	return func(o *Options) {
		o.DynamicLimitLoop = func(perTime *atomic.Int64, rate int64) {
			limit := rate

		LOOP:
			for {
				select {
				case c := <-current:
					limit = updateEstimatedLimit(except, c, limit, minLimit, maxLimit)
					perTime.Store(int64(o.Per / time.Duration(limit)))
				case <-ctx.Done():
					break LOOP
				}
			}
		}
	}
}
