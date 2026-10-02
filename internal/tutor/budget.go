package tutor

import (
	"sync"
)

// Budget manages request and token allowances within a practice session.
// Invariant (AGENT-CONTRACT): Provider requests must be bounded to prevent
// runaway loops or unexpected usage.
type Budget struct {
	mu               sync.Mutex
	maxRequests      int
	maxTokens        int
	requestsConsumed int
	tokensConsumed   int
}

// NewBudget creates a new session budget limiter.
// A value of <= 0 indicates unlimited for that metric.
func NewBudget(maxRequests, maxTokens int) *Budget {
	return &Budget{
		maxRequests: maxRequests,
		maxTokens:   maxTokens,
	}
}

// Check returns ErrBudgetExceeded if either limit has been reached.
func (b *Budget) Check() error {
	if b == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.maxRequests > 0 && b.requestsConsumed >= b.maxRequests {
		return ErrBudgetExceeded
	}
	if b.maxTokens > 0 && b.tokensConsumed >= b.maxTokens {
		return ErrBudgetExceeded
	}
	return nil
}

// RecordUsage records a completed request and its consumed tokens.
// Returns ErrBudgetExceeded if limits were already passed.
func (b *Budget) RecordUsage(tokens int) error {
	if b == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.maxRequests > 0 && b.requestsConsumed >= b.maxRequests {
		return ErrBudgetExceeded
	}
	if b.maxTokens > 0 && b.tokensConsumed >= b.maxTokens {
		return ErrBudgetExceeded
	}

	b.requestsConsumed++
	b.tokensConsumed += tokens
	return nil
}

// RemainingRequests returns the number of requests left before hitting the limit,
// or -1 if unlimited.
func (b *Budget) RemainingRequests() int {
	if b == nil || b.maxRequests <= 0 {
		return -1
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	rem := b.maxRequests - b.requestsConsumed
	if rem < 0 {
		return 0
	}
	return rem
}

// RemainingTokens returns the number of tokens left before hitting the limit,
// or -1 if unlimited.
func (b *Budget) RemainingTokens() int {
	if b == nil || b.maxTokens <= 0 {
		return -1
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	rem := b.maxTokens - b.tokensConsumed
	if rem < 0 {
		return 0
	}
	return rem
}

// Consumed returns total requests and tokens used.
func (b *Budget) Consumed() (int, int) {
	if b == nil {
		return 0, 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.requestsConsumed, b.tokensConsumed
}

// Reset resets consumed counts to zero.
func (b *Budget) Reset() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requestsConsumed = 0
	b.tokensConsumed = 0
}
