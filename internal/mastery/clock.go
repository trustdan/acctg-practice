package mastery

import (
	"time"
)

// Clock provides an interface for time injection, enabling deterministic decay and schedule testing.
type Clock interface {
	Now() time.Time
}

// RealClock uses the system UTC clock.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

// MockClock allows advancing and controlling time in deterministic unit tests.
type MockClock struct {
	current time.Time
}

func NewMockClock(start time.Time) *MockClock {
	return &MockClock{current: start.UTC()}
}

func (m *MockClock) Now() time.Time {
	return m.current
}

func (m *MockClock) Advance(d time.Duration) {
	m.current = m.current.Add(d)
}

func (m *MockClock) Set(t time.Time) {
	m.current = t.UTC()
}
