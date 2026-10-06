package tutor

import (
	"context"
	"fmt"
	"sync/atomic"
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
	idle     time.Duration
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
		idle:     StreamIdleTimeout,
	}
}

// SetStreamIdleTimeout overrides StreamIdleTimeout for this tutor.
func (f *FallbackTutor) SetStreamIdleTimeout(d time.Duration) {
	if d > 0 {
		f.idle = d
	}
}

// Name implements Tutor.
func (f *FallbackTutor) Name() string {
	if f.primary != nil {
		return fmt.Sprintf("%s-with-fallback", f.primary.Name())
	}
	return f.fallback.Name()
}

// Timeout returns the per-request deadline given to the primary provider.
// Callers should allow at least this long so a primary timeout can fall back.
func (f *FallbackTutor) Timeout() time.Duration {
	return f.timeout
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

// StreamIdleTimeout is how long a streamed reply may go without a chunk before
// it is stopped. The first chunk may take up to the provider timeout instead.
const StreamIdleTimeout = 30 * time.Second

// CanStream reports whether the primary provider streams replies.
func (f *FallbackTutor) CanStream() bool {
	_, ok := f.primary.(StreamingTutor)
	return ok
}

// HintStream implements StreamingTutor.
func (f *FallbackTutor) HintStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return f.stream(ctx, req, true, onUpdate)
}

// ExplainStream implements StreamingTutor.
func (f *FallbackTutor) ExplainStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return f.stream(ctx, req, false, onUpdate)
}

// stream replaces the total request deadline with a first-chunk deadline (the
// provider timeout) and an idle deadline between chunks. An error before any
// visible text falls back offline; an error after it keeps the partial reply,
// marked Incomplete. Caller cancellation returns the caller's error.
func (f *FallbackTutor) stream(ctx context.Context, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
	st, ok := f.primary.(StreamingTutor)
	if !ok {
		if isHint {
			return f.Hint(ctx, req)
		}
		return f.Explain(ctx, req)
	}
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	streamCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	var started atomic.Bool
	timer := time.AfterFunc(f.timeout, func() {
		if started.Load() {
			cancel(fmt.Errorf("%w: reply paused for more than %s", ErrTimeout, f.idle))
		} else {
			cancel(fmt.Errorf("%w: no reply within %s", ErrTimeout, f.timeout))
		}
	})
	defer timer.Stop()
	watch := func(u StreamUpdate) {
		started.Store(true)
		timer.Reset(f.idle)
		if onUpdate != nil {
			onUpdate(u)
		}
	}

	var resp Response
	var err error
	if isHint {
		resp, err = st.HintStream(streamCtx, req, watch)
	} else {
		resp, err = st.ExplainStream(streamCtx, req, watch)
	}
	if err == nil {
		return resp, nil
	}
	if ctx.Err() != nil {
		return Response{}, ctx.Err()
	}
	if cause := context.Cause(streamCtx); cause != nil {
		err = cause
	}
	if resp.Text != "" {
		resp.Incomplete = true
		resp.FallbackReason = "Reply stopped early (" + err.Error() + "); this partial reply is not saved."
		resp.GeneratedAt = time.Now().UTC()
		return resp, nil
	}

	var fallbackResp Response
	var ferr error
	if isHint {
		fallbackResp, ferr = f.fallback.Hint(ctx, req)
	} else {
		fallbackResp, ferr = f.fallback.Explain(ctx, req)
	}
	if ferr != nil {
		return Response{}, ferr
	}
	fallbackResp.Fallback = true
	fallbackResp.FallbackReason = err.Error()
	return fallbackResp, nil
}
