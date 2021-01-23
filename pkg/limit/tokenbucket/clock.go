package tokenbucket

import "time"

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (r *realClock) Now() time.Time {
	return time.Now()
}

// mockClock for unit test
type mockClock struct {
	Time time.Time
}

func (m *mockClock) Now() time.Time {
	return m.Time
}
