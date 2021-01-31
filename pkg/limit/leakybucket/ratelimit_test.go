// +build !race

package leakybucket

import (
	"runtime"
	"testing"
	"time"

	"github.com/Allenxuxu/stark/pkg/limit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestLeakyBucketRateLimitTake(t *testing.T) {
	testRateLimitTake(t, 100, time.Second)
	testRateLimitTake(t, 1000, time.Second)
	testRateLimitTake(t, 10000, time.Second)
	testRateLimitTake(t, 100000, time.Second)

	testRateLimitTake(t, 1000, time.Second*5)

}

func testRateLimitTake(t *testing.T, rate int64, per time.Duration) {
	rl := newLimit(rate, limit.Per(per))
	exit := make(chan struct{})
	var count atomic.Uint64

	time.AfterFunc(per, func() {
		close(exit)
	})

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-exit:
					return

				default:
					rl.Take()
					count.Inc()
				}
			}
		}()
	}

	<-exit
	elapsed := per
	ideal := int64((time.Duration(rate) * elapsed) / per)

	c := int64(count.Load())
	should := int64(float64(ideal) * 1.01)
	assert.LessOrEqual(t, c, should, "rate %v,per: %v, allow: %v but should less or equal %v", rate, per, c, should)
}

func TestLeakyBucketRateLimitAllow(t *testing.T) {
	testRateLimitAllow(t, 100, time.Second)
	testRateLimitAllow(t, 1000, time.Second)
	testRateLimitAllow(t, 10000, time.Second)
	testRateLimitTake(t, 100000, time.Second)

	testRateLimitTake(t, 1000, time.Second*5)

}

func testRateLimitAllow(t *testing.T, rate int64, per time.Duration) {
	rl := newLimit(rate, limit.Per(per))
	exit := make(chan struct{})
	var count atomic.Uint64

	time.AfterFunc(per, func() {
		close(exit)
	})

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-exit:
					return

				default:
					if rl.Allow() {
						count.Inc()
					}
					time.Sleep(time.Microsecond)
				}
			}
		}()
	}

	<-exit
	elapsed := per
	ideal := int64((time.Duration(rate) * elapsed) / per)

	c := int64(count.Load())
	should := int64(float64(ideal) * 1.01)
	assert.LessOrEqual(t, c, should, "rate %v,per: %v, allow: %v but should less or equal %v", rate, per, c, should)
}
