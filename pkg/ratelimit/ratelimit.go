package ratelimit // fork from "go.uber.org/ratelimit"

import (
	"time"

	"sync/atomic"
	"unsafe"
)

type state struct {
	last     time.Time
	sleepFor time.Duration
}

type Option func(l *atomicLimiter)

// Per allows configuring limits for different time windows.
//
// The default window is one second, so New(100) produces a one hundred per
// second (100 Hz) rate limiter.
//
// New(2, Per(60*time.Second)) creates a 2 per minute rate limiter.
func Per(per time.Duration) Option {
	return func(l *atomicLimiter) {
		l.per = per
	}
}

type atomicLimiter struct {
	state unsafe.Pointer
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	padding [56]byte // cache line size - state pointer size = 64 - 8; created to avoid false sharing.

	per        time.Duration
	perRequest time.Duration
	maxSlack   time.Duration
}

// newAtomicBased returns a new atomic based limiter.
func New(rate int, opts ...Option) *atomicLimiter {
	l := &atomicLimiter{
		per:      time.Second,
		maxSlack: -10 * time.Second / time.Duration(rate),
	}

	for _, o := range opts {
		o(l)
	}

	l.perRequest = l.per / time.Duration(rate)

	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}

func (t *atomicLimiter) Allow() bool {
	newState := state{}
	now := time.Now()

	previousStatePointer := atomic.LoadPointer(&t.state)
	oldState := (*state)(previousStatePointer)

	newState = state{}
	newState.last = now

	// If this is our first request, then we allow it.
	if oldState.last.IsZero() {
		return atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}

	if (t.perRequest - now.Sub(oldState.last)) > 0 {
		return false
	} else {
		newState.last = now
		return atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *atomicLimiter) Take() time.Time {
	newState := state{}
	taken := false
	for !taken {
		now := time.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)

		newState = state{}
		newState.last = now

		// If this is our first request, then we allow it.
		if oldState.last.IsZero() {
			taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}

		// sleepFor calculates how much time we should sleep based on
		// the perRequest budget and how long the last request took.
		// Since the request may take longer than the budget, this number
		// can get negative, and is summed across requests.
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		// We shouldn't allow sleepFor to get too negative, since it would mean that
		// a service that slowed down a lot for a short period of time would get
		// a much higher RPS following that.
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			newState.last = newState.last.Add(newState.sleepFor)
		}
		taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
	time.Sleep(newState.sleepFor)
	return newState.last
}
