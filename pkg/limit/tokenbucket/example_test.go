// +build !race

package tokenbucket

import (
	"fmt"
	"time"
)

func ExampleTokenBucketRateLimitTake() {
	rl := New(100, 100) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev).Milliseconds())
		}
		prev = now
	}

	// Output:
	// 1 0
	// 2 0
	// 3 0
	// 4 0
	// 5 0
	// 6 0
	// 7 0
	// 8 0
	// 9 0
}

func ExampleTokenBucketRateLimitAllow() {
	rl := New(100, 2) // per second

	for i := 0; i < 2; i++ {
		fmt.Println(rl.Allow())
	}
	fmt.Println(rl.Allow())
	time.Sleep(1 * time.Millisecond)
	fmt.Println(rl.Allow())
	time.Sleep(10 * time.Millisecond)
	fmt.Println(rl.Allow())

	// Output:
	// true
	// true
	// false
	// false
	// true
}
