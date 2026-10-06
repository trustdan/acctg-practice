package tutor_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/tutor"
)

// sseStep is one write to a fake SSE stream: raw bytes, then an optional pause.
type sseStep struct {
	raw   string
	pause time.Duration
}

func sseChunk(content string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"choices": []map[string]interface{}{{"delta": map[string]string{"content": content}}},
	})
	return "data: " + string(b) + "\n\n"
}

func sseReasoning(text string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"choices": []map[string]interface{}{{"delta": map[string]string{"reasoning_content": text}}},
	})
	return "data: " + string(b) + "\n\n"
}

const sseUsage = `data: {"choices":[],"usage":{"completion_tokens":7}}` + "\n\n"
const sseDone = "data: [DONE]\n\n"

// startSSE serves /v1/chat/completions as the given steps. It records the
// request body and reports when the client disconnects.
func startSSE(t *testing.T, steps []sseStep) (url string, body *map[string]interface{}, gone chan struct{}) {
	t.Helper()
	body = &map[string]interface{}{}
	gone = make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(body)
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		defer close(gone)
		for _, s := range steps {
			if s.raw != "" {
				_, _ = w.Write([]byte(s.raw))
				flusher.Flush()
			}
			if s.pause > 0 {
				select {
				case <-time.After(s.pause):
				case <-r.Context().Done():
					return
				}
			}
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL + "/v1", body, gone
}

func lmStream(url string) *tutor.LMStudioTutor {
	return tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: url, Model: "m"})
}

func collect(updates *[]tutor.StreamUpdate) func(tutor.StreamUpdate) {
	return func(u tutor.StreamUpdate) { *updates = append(*updates, u) }
}

func TestStreamParsesChunksAcrossSplitLines(t *testing.T) {
	full := sseChunk("Debits are ") + sseChunk("on the left.") + sseUsage + sseDone + sseChunk("ignored after DONE")
	// Split mid-line and mid-JSON to mimic arbitrary network reads.
	var steps []sseStep
	for i := 0; i < len(full); i += 13 {
		end := i + 13
		if end > len(full) {
			end = len(full)
		}
		steps = append(steps, sseStep{raw: full[i:end]})
	}
	url, body, _ := startSSE(t, steps)

	var updates []tutor.StreamUpdate
	resp, err := lmStream(url).HintStream(context.Background(), lmReq, collect(&updates))
	if err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	if resp.Text != "Debits are on the left." || resp.TokensUsed != 7 || resp.Incomplete {
		t.Errorf("unexpected response: %+v", resp)
	}
	if len(updates) != 2 || updates[0].Text != "Debits are" || updates[1].Text != "Debits are on the left." || updates[1].Provider != tutor.ProviderLMStudio {
		t.Errorf("unexpected updates: %+v", updates)
	}
	if (*body)["stream"] != true {
		t.Errorf("request did not ask for a stream: %v", *body)
	}
}

func TestStreamHidesThinkBlocksSplitAcrossChunks(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: sseChunk("<thi") + sseChunk("nk>secret plan") + sseChunk(" more</th") +
		sseChunk("ink>What did the company ") + sseChunk("receive?") + sseDone}})

	var updates []tutor.StreamUpdate
	resp, err := lmStream(url).HintStream(context.Background(), lmReq, collect(&updates))
	if err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	sawThinking := false
	for _, u := range updates {
		if strings.Contains(u.Text, "secret") || strings.Contains(u.Text, "<") {
			t.Errorf("reasoning leaked into visible text: %q", u.Text)
		}
		sawThinking = sawThinking || u.Thinking
	}
	if !sawThinking {
		t.Error("expected a Thinking update while the think block was open")
	}
	if resp.Text != "What did the company receive?" {
		t.Errorf("unexpected final text: %q", resp.Text)
	}
}

func TestStreamReasoningContentIsThinkingOnly(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: sseReasoning("hidden chain") + sseChunk("Visible.") + sseDone}})

	var updates []tutor.StreamUpdate
	resp, err := lmStream(url).HintStream(context.Background(), lmReq, collect(&updates))
	if err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	if len(updates) != 2 || !updates[0].Thinking || updates[0].Text != "" || updates[1].Thinking {
		t.Errorf("unexpected updates: %+v", updates)
	}
	if strings.Contains(resp.Text, "hidden") {
		t.Errorf("reasoning_content leaked: %q", resp.Text)
	}
}

func TestStreamThinkingOnlyFallsBackOffline(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: sseChunk("<think>never finished") + sseDone}})
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 2*time.Second)
	resp, err := ft.HintStream(context.Background(), lmReq, nil)
	if err != nil {
		t.Fatalf("expected offline fallback, got %v", err)
	}
	if !resp.Fallback || !strings.Contains(resp.FallbackReason, "max_tokens") {
		t.Errorf("expected thinking-only fallback, got %+v", resp)
	}
}

func TestStreamErrorBeforeFirstChunkFallsBack(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model crashed", http.StatusInternalServerError)
	}))
	defer srv.Close()
	ft := tutor.NewFallbackTutor(lmStream(srv.URL+"/v1"), nil, 2*time.Second)
	resp, err := ft.ExplainStream(context.Background(), lmReq, nil)
	if err != nil || !resp.Fallback || resp.Incomplete {
		t.Fatalf("expected offline fallback, got %+v, %v", resp, err)
	}
}

func TestStreamFirstTokenTimeoutFallsBack(t *testing.T) {
	url, _, gone := startSSE(t, []sseStep{{pause: 5 * time.Second}, {raw: sseChunk("too late") + sseDone}})
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 150*time.Millisecond)

	start := time.Now()
	resp, err := ft.HintStream(context.Background(), lmReq, nil)
	if err != nil || !resp.Fallback {
		t.Fatalf("expected offline fallback, got %+v, %v", resp, err)
	}
	if !strings.Contains(resp.FallbackReason, "no reply within") {
		t.Errorf("expected first-token timeout reason, got %q", resp.FallbackReason)
	}
	if time.Since(start) > 2*time.Second {
		t.Errorf("first-token timeout took too long: %s", time.Since(start))
	}
	waitClosed(t, gone)
}

func TestStreamIdleTimeoutKeepsPartialReply(t *testing.T) {
	url, _, gone := startSSE(t, []sseStep{{raw: sseChunk("Cash is an asset"), pause: 5 * time.Second}, {raw: sseDone}})
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 2*time.Second)
	ft.SetStreamIdleTimeout(150 * time.Millisecond)

	resp, err := ft.ExplainStream(context.Background(), lmReq, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Incomplete || resp.Fallback || resp.Text != "Cash is an asset" || resp.Provider != tutor.ProviderLMStudio {
		t.Errorf("expected incomplete partial reply, got %+v", resp)
	}
	if !strings.Contains(resp.FallbackReason, "paused") || !strings.Contains(resp.FallbackReason, "not saved") {
		t.Errorf("expected idle reason, got %q", resp.FallbackReason)
	}
	waitClosed(t, gone)
}

func TestStreamSlowButSteadyIsNotIdle(t *testing.T) {
	var steps []sseStep
	for i := 0; i < 6; i++ {
		steps = append(steps, sseStep{raw: sseChunk(fmt.Sprintf("w%d ", i)), pause: 60 * time.Millisecond})
	}
	steps = append(steps, sseStep{raw: sseDone})
	url, _, _ := startSSE(t, steps)
	// Total time exceeds the first-token timeout, but no gap exceeds the idle timeout.
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 200*time.Millisecond)
	ft.SetStreamIdleTimeout(200 * time.Millisecond)

	resp, err := ft.HintStream(context.Background(), lmReq, nil)
	if err != nil || resp.Incomplete || resp.Fallback || resp.Text != "w0 w1 w2 w3 w4 w5" {
		t.Fatalf("expected complete reply, got %+v, %v", resp, err)
	}
}

func TestStreamMidStreamDisconnectIsIncomplete(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: sseChunk("Revenue records ")}}) // Ends without [DONE].
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 2*time.Second)
	resp, err := ft.ExplainStream(context.Background(), lmReq, nil)
	if err != nil || !resp.Incomplete || resp.Text != "Revenue records" {
		t.Fatalf("expected incomplete partial reply, got %+v, %v", resp, err)
	}
}

func TestStreamSizeLimit(t *testing.T) {
	big := strings.Repeat("x", 64*1024)
	var steps []sseStep
	for i := 0; i < 20; i++ {
		steps = append(steps, sseStep{raw: sseChunk(big)})
	}
	url, _, _ := startSSE(t, steps)
	resp, err := lmStream(url).ExplainStream(context.Background(), lmReq, func(tutor.StreamUpdate) {})
	if !errors.Is(err, tutor.ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got %v", err)
	}
	if !resp.Incomplete || resp.Text == "" {
		t.Errorf("expected partial text with the size error, got %d bytes", len(resp.Text))
	}
}

func TestStreamCancellationClosesBodyPromptly(t *testing.T) {
	url, _, gone := startSSE(t, []sseStep{{raw: sseChunk("First words"), pause: 10 * time.Second}, {raw: sseDone}})
	ft := tutor.NewFallbackTutor(lmStream(url), nil, 20*time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()
	resp, err := ft.HintStream(ctx, lmReq, func(u tutor.StreamUpdate) {
		if u.Text != "" {
			cancel()
		}
	})
	if !errors.Is(err, context.Canceled) || resp.Fallback {
		t.Fatalf("expected caller cancellation without fallback, got %+v, %v", resp, err)
	}
	if time.Since(start) > 2*time.Second {
		t.Errorf("cancellation took too long: %s", time.Since(start))
	}
	waitClosed(t, gone)
}

func TestStreamBudgetRecordedOncePerRequest(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: sseChunk("Partial")}}) // Aborted: no [DONE].
	budget := tutor.NewBudget(10, 0)
	lm := tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: url, Model: "m", Budget: budget})
	_, _ = lm.HintStream(context.Background(), lmReq, func(tutor.StreamUpdate) {})
	if req, _ := budget.Consumed(); req != 1 {
		t.Errorf("expected one recorded request for an aborted stream, got %d", req)
	}
}

func TestFallbackStreamWithoutStreamingPrimaryUsesNormalPath(t *testing.T) {
	ft := tutor.NewFallbackTutor(nil, nil, time.Second)
	if ft.CanStream() {
		t.Error("offline-only tutor should not report streaming")
	}
	resp, err := ft.HintStream(context.Background(), lmReq, nil)
	if err != nil || resp.Text == "" {
		t.Fatalf("expected offline reply, got %+v, %v", resp, err)
	}
}

func waitClosed(t *testing.T, gone chan struct{}) {
	t.Helper()
	select {
	case <-gone:
	case <-time.After(2 * time.Second):
		t.Error("server handler still running: response body was not closed")
	}
}
