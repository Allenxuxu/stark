// +build !race

package leakybucket

import (
	"fmt"
	"time"
)

func ExampleLeakyBucketRateLimitTake() {
	rl := New(100) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev).Milliseconds())
		}
		prev = now
	}

	// Output:
	// 1 10
	// 2 10
	// 3 10
	// 4 10
	// 5 10
	// 6 10
	// 7 10
	// 8 10
	// 9 10
}

func ExampleLeakyBucketRateLimitAllow() {
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
