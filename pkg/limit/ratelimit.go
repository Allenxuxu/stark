package limit

import "time"

type RateLimit interface {
	Allow() bool
	Take() time.Time
}
