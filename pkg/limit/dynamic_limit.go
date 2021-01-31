package limit

import (
	"math"
)

func updateEstimatedLimit(expect, current float64, limit, minLimit, MaxLimit int64) int64 {
	diff := float64(limit) * (1 - expect/current)
	alpha := 3 * math.Log10(float64(limit))
	beta := 6 * math.Log10(float64(limit))

	if diff < alpha {
		limit += int64(math.Log10(float64(limit)))
	} else if diff > beta {
		limit -= int64(math.Log10(float64(limit)))
	}
	// else : do nothing

	if limit > MaxLimit {
		limit = MaxLimit
	}
	if limit < minLimit {
		limit = minLimit
	}

	return limit
}
