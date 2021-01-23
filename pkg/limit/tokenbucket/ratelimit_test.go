// +build !race

package tokenbucket

import (
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/Allenxuxu/stark/pkg/limit"
	"github.com/stretchr/testify/assert"
)

func newMockLimit(rate, cap int64, mc *mockClock, opts ...limit.Option) *limiter {
	rl := newLimit(rate, cap, opts...)
	rl.time = mc

	return rl
}

func TestTokenBucketTokenLimitAllow(t *testing.T) {
	mc := &mockClock{}
	mc.Time = time.Now()

	rl := newMockLimit(100, 100, mc)
	assert.Equal(t, rl.cap, int64(100))
	assert.Equal(t, rl.opts.Per, time.Second)
	assert.Equal(t, rl.perToken, time.Millisecond*10)
	assert.Equal(t, int64(100), rl.availableTokens())

	for i := 0; i < 100; i++ {
		assert.True(t, rl.Allow())
	}
	assert.False(t, rl.Allow())
	assert.Equal(t, int64(0), rl.availableTokens())

	mc.Time = mc.Time.Add(time.Second)
	assert.True(t, rl.Allow())
	assert.Equal(t, int64(99), rl.availableTokens())
	for i := 0; i < 99; i++ {
		assert.True(t, rl.Allow())
	}
	assert.Equal(t, int64(0), rl.availableTokens())

	mc.Time = mc.Time.Add(time.Millisecond * 10)
	assert.Equal(t, int64(0), rl.availableTokens())

	assert.True(t, rl.Allow())
	assert.Equal(t, int64(0), rl.availableTokens())
}

func TestTokenBucketRateLimitTaken(t *testing.T) {
	testRateLimitTaken(t, 20, 100, time.Second)
	testRateLimitTaken(t, 10, 1000, time.Second)
	testRateLimitTaken(t, 100, 10000, time.Second)
	testRateLimitTaken(t, 1000, 100000, time.Second)
	testRateLimitTaken(t, 10000, 1000000, time.Second)
	testRateLimitTaken(t, 100000, 1000000, time.Second)

	testRateLimitTaken(t, 1000, 1, time.Second*5)
}

func testRateLimitTaken(t *testing.T, rate, cap int64, per time.Duration) {
	rl := New(rate, cap, limit.Per(per))
	exit := make(chan struct{})
	var count atomic.Int64

	time.AfterFunc(per*2, func() {
		close(exit)
	})

	for i := 0; i < 1000; i++ {
		go func() {
			for {
				select {
				case <-exit:
					return

				default:
					rl.Take()
					count.Inc()

					time.Sleep(2 * time.Millisecond)
				}
			}
		}()
	}

	<-exit
	c := count.Load()
	should := cap + (rate * 2)
	assert.LessOrEqual(t, c, should, "rate %v,per: %v, allow: %v but should be %v", rate, per, c, should)
}

func TestTokenBucketRateLimitAllow(t *testing.T) {
	testRateLimitAllow(t, 10, 1000, time.Second)
	testRateLimitAllow(t, 100, 10000, time.Second)
	testRateLimitAllow(t, 1000, 100000, time.Second)
	testRateLimitAllow(t, 10000, 1000000, time.Second)
	testRateLimitAllow(t, 100000, 1000000, time.Second)

	testRateLimitAllow(t, 1000, 10, time.Second*5)
}

func testRateLimitAllow(t *testing.T, rate, cap int64, per time.Duration) {
	rl := New(rate, cap, limit.Per(per))
	exit := make(chan struct{})
	var count atomic.Int64

	time.AfterFunc(per*2, func() {
		close(exit)
	})

	for i := 0; i < 1000; i++ {
		go func() {
			for {
				select {
				case <-exit:
					return

				default:
					if rl.Allow() {
						count.Inc()
					}

					time.Sleep(2 * time.Millisecond)
				}
			}
		}()
	}

	<-exit
	c := count.Load()
	should := cap + (rate * 2)
	assert.LessOrEqual(t, c, should, "rate %v,per: %v, cap %v, allow: %v but should be %v", rate, per, cap, c, should)
}
