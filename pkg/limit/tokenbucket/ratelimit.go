package tokenbucket

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/Allenxuxu/stark/pkg/limit"
)

type state struct {
	last            time.Time
	availableTokens int64
}

type limiter struct {
	cap      int64
	perToken time.Duration
	state    unsafe.Pointer

	opts limit.Options

	time clock // for unit test
}

func New(rate, cap int64, opts ...limit.Option) limit.RateLimit {
	return newLimit(rate, cap, opts...)
}

func newLimit(rate, cap int64, opts ...limit.Option) *limiter {
	options := limit.Options{
		Per: time.Second,
	}

	for _, o := range opts {
		o(&options)
	}

	l := &limiter{
		cap:  cap,
		opts: options,
		time: &realClock{},
	}

	l.perToken = l.opts.Per / time.Duration(rate)

	s := state{
		last:            l.time.Now(),
		availableTokens: cap,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&s))

	return l
}

func (l *limiter) Take() time.Time {
	newState := state{}
	taken := false
	for !taken {
		previousStatePointer := atomic.LoadPointer(&l.state)
		oldState := (*state)(previousStatePointer)

		last := oldState.last
		now := l.time.Now()

		newState.last = now
		newState.availableTokens = min(l.cap, oldState.availableTokens+int64(now.Sub(last)/l.perToken))

		if newState.availableTokens > 0 {
			newState.availableTokens--
			taken = atomic.CompareAndSwapPointer(&l.state, previousStatePointer, unsafe.Pointer(&newState))
		}
		if !taken {
			time.Sleep(l.perToken)
		}
	}

	return newState.last
}

func (l *limiter) Allow() bool {
	previousStatePointer := atomic.LoadPointer(&l.state)
	oldState := (*state)(previousStatePointer)

	last := oldState.last
	now := l.time.Now()

	newState := state{
		last:            now,
		availableTokens: min(l.cap, oldState.availableTokens+int64(now.Sub(last)/l.perToken)),
	}

	if newState.availableTokens > 0 {
		newState.availableTokens--
		return atomic.CompareAndSwapPointer(&l.state, previousStatePointer, unsafe.Pointer(&newState))
	}

	return false
}

func (l *limiter) availableTokens() int64 {
	previousStatePointer := atomic.LoadPointer(&l.state)
	oldState := (*state)(previousStatePointer)

	return oldState.availableTokens
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}
