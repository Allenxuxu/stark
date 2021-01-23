package limit

import "time"

type Options struct {
	Per time.Duration
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
