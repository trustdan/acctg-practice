package tutor

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// SimulatedConfig configures the behavior of SimulatedProvider.
type SimulatedConfig struct {
	Delay            time.Duration
	FailWith         error
	FailRate         float64
	MaxResponseBytes int
	Budget           *Budget
	CustomHint       string
	CustomExplain    string
	RNG              *rand.Rand
}

// SimulatedProvider is a testable asynchronous tutor provider that simulates
// latency, network faults, size limits, and token budgeting.
type SimulatedProvider struct {
	mu  sync.Mutex
	cfg SimulatedConfig
}

// NewSimulatedProvider creates a new SimulatedProvider with the given configuration.
func NewSimulatedProvider(cfg SimulatedConfig) *SimulatedProvider {
	if cfg.RNG == nil {
		cfg.RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &SimulatedProvider{
		cfg: cfg,
	}
}

// Name implements Tutor.
func (s *SimulatedProvider) Name() string {
	return "SimulatedProvider"
}

// Hint implements Tutor.
func (s *SimulatedProvider) Hint(ctx context.Context, req Request) (Response, error) {
	return s.process(ctx, req, true)
}

// Explain implements Tutor.
func (s *SimulatedProvider) Explain(ctx context.Context, req Request) (Response, error) {
	return s.process(ctx, req, false)
}

func (s *SimulatedProvider) process(ctx context.Context, req Request, isHint bool) (Response, error) {
	s.mu.Lock()
	budget := s.cfg.Budget
	delay := s.cfg.Delay
	failWith := s.cfg.FailWith
	failRate := s.cfg.FailRate
	maxBytes := s.cfg.MaxResponseBytes
	customHint := s.cfg.CustomHint
	customExplain := s.cfg.CustomExplain
	rng := s.cfg.RNG
	s.mu.Unlock()

	// 1. Budget check
	if budget != nil {
		if err := budget.Check(); err != nil {
			return Response{}, err
		}
	}

	// 2. Simulated latency with strict context cancellation / timeout check
	if delay > 0 {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		case <-timer.C:
		}
	} else if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	// 3. Simulated failure injection
	if failWith != nil {
		return Response{}, failWith
	}
	if failRate > 0 && rng != nil {
		s.mu.Lock()
		roll := rng.Float64()
		s.mu.Unlock()
		if roll < failRate {
			return Response{}, ErrProviderUnavailable
		}
	}

	// 4. Generate text
	var text string
	if isHint {
		if customHint != "" {
			text = customHint
		} else {
			offline := NewOfflineTutor()
			resp, _ := offline.Hint(ctx, req)
			text = "[Simulated] " + resp.Text
		}
	} else {
		if customExplain != "" {
			text = customExplain
		} else {
			offline := NewOfflineTutor()
			resp, _ := offline.Explain(ctx, req)
			text = "[Simulated] " + resp.Text
		}
	}

	// 5. Size limit check
	if maxBytes > 0 && len(text) > maxBytes {
		return Response{}, ErrResponseTooLarge
	}

	// 6. Token budgeting deduction
	tokens := (len(text) + 3) / 4
	if budget != nil {
		if err := budget.RecordUsage(tokens); err != nil {
			return Response{}, err
		}
	}

	return Response{
		Text:        text,
		Provider:    "simulated",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}
