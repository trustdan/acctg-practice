package tutor

import (
	"context"
	"fmt"
	"time"
)

// FallbackTutor wraps a primary tutor provider and falls back cleanly to an
// offline tutor when the primary fails, times out, or exceeds budget.
//
// Invariant (AGENT-CONTRACT):
// A working offline drill takes priority over provider integrations. Any provider
// failure or timeout must fallback cleanly without crashing or breaking the learner session.
type FallbackTutor struct {
	primary  Tutor
	fallback Tutor
	timeout  time.Duration
}

// NewFallbackTutor constructs a FallbackTutor with primary, fallback, and timeout.
// If fallback is nil, OfflineTutor is used by default.
// If timeout <= 0, a default of 3 seconds is applied.
func NewFallbackTutor(primary Tutor, fallback Tutor, timeout time.Duration) *FallbackTutor {
	if fallback == nil {
		fallback = NewOfflineTutor()
	}
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &FallbackTutor{
		primary:  primary,
		fallback: fallback,
		timeout:  timeout,
	}
}

// Name implements Tutor.
func (f *FallbackTutor) Name() string {
	if f.primary != nil {
		return fmt.Sprintf("%s-with-fallback", f.primary.Name())
	}
	return f.fallback.Name()
}

// Hint implements Tutor.
func (f *FallbackTutor) Hint(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	var primaryErr error
	if f.primary != nil {
		childCtx, cancel := context.WithTimeout(ctx, f.timeout)
		resp, err := f.primary.Hint(childCtx, req)
		cancel()
		primaryErr = err

		if err == nil {
			return resp, nil
		}

		// If caller's parent context was cancelled, do not fallback; return caller cancellation
		if ctx.Err() != nil {
			return Response{}, ctx.Err()
		}
	}

	// Primary failed or timed out: fall back cleanly
	fallbackResp, err := f.fallback.Hint(ctx, req)
	if err != nil {
		return Response{}, err
	}
	fallbackResp.Fallback = true
	if primaryErr != nil {
		fallbackResp.FallbackReason = primaryErr.Error()
	}
	return fallbackResp, nil
}

// Explain implements Tutor.
func (f *FallbackTutor) Explain(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	var primaryErr error
	if f.primary != nil {
		childCtx, cancel := context.WithTimeout(ctx, f.timeout)
		resp, err := f.primary.Explain(childCtx, req)
		cancel()
		primaryErr = err

		if err == nil {
			return resp, nil
		}

		// If caller's parent context was cancelled, do not fallback; return caller cancellation
		if ctx.Err() != nil {
			return Response{}, ctx.Err()
		}
	}

	// Primary failed or timed out: fall back cleanly
	fallbackResp, err := f.fallback.Explain(ctx, req)
	if err != nil {
		return Response{}, err
	}
	fallbackResp.Fallback = true
	if primaryErr != nil {
		fallbackResp.FallbackReason = primaryErr.Error()
	}
	return fallbackResp, nil
}
