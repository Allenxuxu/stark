package leakybucket // fork from "go.uber.org/ratelimit"

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/Allenxuxu/stark/pkg/limit"
	uAtomic "go.uber.org/atomic"
)

type state struct {
	last     time.Time
	sleepFor time.Duration
}

type limiter struct {
	state unsafe.Pointer
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	padding [56]byte // cache line size - state pointer size = 64 - 8; created to avoid false sharing.

	perRequest *uAtomic.Int64
	maxSlack   time.Duration

	opts limit.Options
}

func New(rate int64, opts ...limit.Option) limit.RateLimit {
	return newLimit(rate, opts...)
}

func newLimit(rate int64, opts ...limit.Option) *limiter {
	options := limit.Options{
		Per: time.Second,
	}

	for _, o := range opts {
		o(&options)
	}

	l := &limiter{
		maxSlack:   -10 * time.Second / time.Duration(rate),
		opts:       options,
		perRequest: uAtomic.NewInt64(0),
	}
	l.perRequest.Store(int64(l.opts.Per / time.Duration(rate)))

	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))

	if l.opts.DynamicLimitLoop != nil {
		go l.opts.DynamicLimitLoop(l.perRequest, rate)
	}
	return l
}

func (t *limiter) Allow() bool {
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

	if (time.Duration(t.perRequest.Load()) - now.Sub(oldState.last)) > 0 {
		return false
	} else {
		newState.last = now
		return atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *limiter) Take() time.Time {
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
		newState.sleepFor += time.Duration(t.perRequest.Load()) - now.Sub(oldState.last)
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
