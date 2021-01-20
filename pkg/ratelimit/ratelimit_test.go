// +build !race

package ratelimit

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/atomic"
)

func ExampleRateLimitTake() {
	rl := New(100) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev))
		}
		prev = now
	}

	// Output:
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}

func ExampleRateLimitAllow() {
	rl := New(100) // per second

	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	time.Sleep(1 * time.Millisecond)
	fmt.Println(rl.Allow())
	time.Sleep(1 * time.Millisecond)
	fmt.Println(rl.Allow())
	time.Sleep(8 * time.Millisecond)
	fmt.Println(rl.Allow())

	time.Sleep(10 * time.Millisecond)
	fmt.Println(rl.Allow())

	time.Sleep(2 * time.Millisecond)
	fmt.Println(rl.Allow())

	// Output:
	// true
	// false
	// false
	// false
	// true
	// true
	// false

}

func TestRateLimitTake(t *testing.T) {
	testRateLimitTake(t, 10, time.Second)
	testRateLimitTake(t, 100, time.Second)
	testRateLimitTake(t, 1000, time.Second)
	testRateLimitTake(t, 10000, time.Second)
	testRateLimitTake(t, 100000, time.Second)
	testRateLimitTake(t, 100, time.Second*5)
}

func testRateLimitTake(t *testing.T, rate int, per time.Duration) {
	rl := New(rate, Per(per))
	exit := make(chan struct{})
	var count atomic.Uint64

	time.AfterFunc(per, func() {
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

				}
			}
		}()
	}

	<-exit
	c := int(count.Load())
	if c > (rate + 2) {
		t.Fatal("rate: ", rate, "per: ", per, "taken: ", c)
	}

	t.Log("rate: ", rate, "per: ", per, "taken: ", c)
}

func TestRateLimitAllow(t *testing.T) {
	testRateLimitAllow(t, 1, time.Second)
	testRateLimitAllow(t, 10, time.Second)
	testRateLimitAllow(t, 100, time.Second)
	testRateLimitAllow(t, 1000, time.Second)
	testRateLimitAllow(t, 10000, time.Second)
	testRateLimitAllow(t, 100000, time.Second)
	testRateLimitAllow(t, 1000000, time.Second)

	testRateLimitAllow(t, 1, time.Second*5)
	testRateLimitAllow(t, 100, time.Second*5)
}

func testRateLimitAllow(t *testing.T, rate int, per time.Duration) {
	rl := New(rate, Per(per))
	exit := make(chan struct{})
	var count atomic.Uint64

	time.AfterFunc(per, func() {
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
				}
			}
		}()
	}

	<-exit
	c := int(count.Load())
	if c > (rate + 2) {
		t.Fatal("rate: ", rate, "per: ", per, "allow: ", c)
	}
	t.Log("rate: ", rate, "per: ", per, "allow: ", c)
}
