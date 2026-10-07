package tutor_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/tutor"
)

// Stage 30b: Anthropic, Gemini and ChatGPT plan streaming over fake SSE servers.

func anthropicEvent(name, data string) string {
	return "event: " + name + "\ndata: " + data + "\n\n"
}

func anthropicText(text string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"type": "content_block_delta", "index": 1,
		"delta": map[string]string{"type": "text_delta", "text": text},
	})
	return anthropicEvent("content_block_delta", string(b))
}

const (
	anthropicStart    = "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"usage\":{\"input_tokens\":25,\"output_tokens\":1}}}\n\n"
	anthropicThinking = "event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"thinking_delta\",\"thinking\":\"secret reasoning\"}}\n\n"
	anthropicPing     = "event: ping\ndata: {\"type\": \"ping\"}\n\n"
	anthropicUsage    = "event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":15}}\n\n"
	anthropicStop     = "event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
)

func anthropicStream(url string, budget *tutor.Budget) *tutor.AnthropicTutor {
	return tutor.NewAnthropicTutor(tutor.APIKeyConfig{APIKey: "test-anthropic-key", BaseURL: url, Model: "m", Budget: budget})
}

func TestAnthropicStreamTextThinkingAndUsage(t *testing.T) {
	var gotKey, gotPath string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey, gotPath = r.Header.Get("x-api-key"), r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "text/event-stream")
		full := anthropicStart + anthropicThinking + anthropicPing + anthropicText("Debits are ") + anthropicText("on the left.") + anthropicUsage + anthropicStop
		// Split mid-line to mimic arbitrary network reads.
		for i := 0; i < len(full); i += 23 {
			end := min(i+23, len(full))
			_, _ = w.Write([]byte(full[i:end]))
			w.(http.Flusher).Flush()
		}
	}))
	defer srv.Close()

	var updates []tutor.StreamUpdate
	resp, err := anthropicStream(srv.URL, nil).ExplainStream(context.Background(), tutor.Request{}, collect(&updates))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotKey != "test-anthropic-key" || gotPath != "/messages" || body["stream"] != true {
		t.Fatalf("request: key=%q path=%q stream=%v", gotKey, gotPath, body["stream"])
	}
	if resp.Text != "Debits are on the left." || resp.Provider != "anthropic" || resp.TokensUsed != 15 || resp.Incomplete {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(updates) != 3 || !updates[0].Thinking || updates[0].Text != "" {
		t.Fatalf("expected a thinking update then two text updates, got %+v", updates)
	}
	for _, u := range updates {
		if strings.Contains(u.Text, "secret") {
			t.Fatalf("thinking text leaked: %+v", u)
		}
	}
	if updates[1].Text != "Debits are" || updates[1].Thinking {
		t.Fatalf("unexpected first text update: %+v", updates[1])
	}
}

func TestAnthropicStreamErrorEventKeepsPartialReply(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: anthropicStart + anthropicText("Revenue records earning") +
		anthropicEvent("error", `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`)}})
	ft := tutor.NewFallbackTutor(anthropicStream(url, nil), nil, 2*time.Second)
	resp, err := ft.ExplainStream(context.Background(), tutor.Request{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Incomplete || resp.Text != "Revenue records earning" || resp.Fallback {
		t.Fatalf("expected an incomplete partial reply, got %+v", resp)
	}
	if !strings.Contains(resp.FallbackReason, "overloaded_error") || !strings.Contains(resp.FallbackReason, "not saved") {
		t.Fatalf("unexpected reason: %q", resp.FallbackReason)
	}
}

func TestAnthropicStreamEndWithoutMessageStopIsIncomplete(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: anthropicStart + anthropicText("Cash is an asset")}})
	budget := tutor.NewBudget(10, 0)
	_, err := anthropicStream(url, budget).HintStream(context.Background(), tutor.Request{}, func(tutor.StreamUpdate) {})
	if err == nil || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected an interrupted stream, got %v", err)
	}
	if reqs, _ := budget.Consumed(); reqs != 1 {
		t.Fatalf("expected one recorded request on abort, got %d", reqs)
	}
}

func TestAnthropicStreamErrorBeforeTextFallsBack(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: anthropicStart +
		anthropicEvent("error", `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`)}})
	ft := tutor.NewFallbackTutor(anthropicStream(url, nil), nil, 2*time.Second)
	resp, err := ft.HintStream(context.Background(), tutor.Request{}, nil)
	if err != nil || !resp.Fallback || resp.Incomplete || !strings.Contains(resp.FallbackReason, "Overloaded") {
		t.Fatalf("expected offline fallback, got %+v, %v", resp, err)
	}
}

func TestAnthropicStreamSizeLimit(t *testing.T) {
	big := strings.Repeat("x", 64*1024)
	steps := []sseStep{{raw: anthropicStart}}
	for i := 0; i < 20; i++ {
		steps = append(steps, sseStep{raw: anthropicText(big)})
	}
	url, _, _ := startSSE(t, append(steps, sseStep{raw: anthropicStop}))
	_, err := anthropicStream(url, nil).ExplainStream(context.Background(), tutor.Request{}, func(tutor.StreamUpdate) {})
	if !errors.Is(err, tutor.ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got %v", err)
	}
}

func geminiChunk(parts []map[string]interface{}, finish string, tokens int) string {
	cand := map[string]interface{}{"content": map[string]interface{}{"role": "model", "parts": parts}}
	if finish != "" {
		cand["finishReason"] = finish
	}
	msg := map[string]interface{}{"candidates": []interface{}{cand}}
	if tokens > 0 {
		msg["usageMetadata"] = map[string]int{"promptTokenCount": 40, "candidatesTokenCount": tokens}
	}
	b, _ := json.Marshal(msg)
	return "data: " + string(b) + "\r\n\r\n"
}

func geminiText(text string) []map[string]interface{} {
	return []map[string]interface{}{{"text": text}}
}

func geminiStream(url string) *tutor.GeminiTutor {
	return tutor.NewGeminiTutor(tutor.APIKeyConfig{APIKey: "test-gemini-key", BaseURL: url, Model: "gemini-test"})
}

func TestGeminiStreamTextThoughtsAndUsage(t *testing.T) {
	var gotKey, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey, gotPath, gotQuery = r.Header.Get("x-goog-api-key"), r.URL.Path, r.URL.RawQuery
		w.Header().Set("Content-Type", "text/event-stream")
		thought := []map[string]interface{}{{"text": "secret reasoning", "thought": true}}
		for _, c := range []string{
			geminiChunk(thought, "", 0),
			geminiChunk(geminiText("Unearned revenue "), "", 3),
			geminiChunk(geminiText("is a liability."), "STOP", 9),
		} {
			_, _ = w.Write([]byte(c))
			w.(http.Flusher).Flush()
		}
	}))
	defer srv.Close()

	var updates []tutor.StreamUpdate
	resp, err := geminiStream(srv.URL).HintStream(context.Background(), tutor.Request{}, collect(&updates))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/models/gemini-test:streamGenerateContent" || gotQuery != "alt=sse" || gotKey != "test-gemini-key" {
		t.Fatalf("request: path=%q query=%q key=%q (the key must not be in the URL)", gotPath, gotQuery, gotKey)
	}
	if resp.Text != "Unearned revenue is a liability." || resp.Provider != "google" || resp.TokensUsed != 9 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(updates) != 3 || !updates[0].Thinking || updates[0].Text != "" || updates[1].Text != "Unearned revenue" {
		t.Fatalf("unexpected updates: %+v", updates)
	}
}

func TestGeminiStreamEndWithoutFinishReasonIsIncomplete(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: geminiChunk(geminiText("Prepaid rent is"), "", 0)}})
	ft := tutor.NewFallbackTutor(geminiStream(url), nil, 2*time.Second)
	resp, err := ft.ExplainStream(context.Background(), tutor.Request{}, nil)
	if err != nil || !resp.Incomplete || resp.Text != "Prepaid rent is" {
		t.Fatalf("expected an incomplete partial reply, got %+v, %v", resp, err)
	}
}

func TestGeminiStreamErrorChunkBeforeTextFallsBack(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: `data: {"error":{"code":503,"message":"The model is overloaded.","status":"UNAVAILABLE"}}` + "\n\n"}})
	ft := tutor.NewFallbackTutor(geminiStream(url), nil, 2*time.Second)
	resp, err := ft.HintStream(context.Background(), tutor.Request{}, nil)
	if err != nil || !resp.Fallback || !strings.Contains(resp.FallbackReason, "UNAVAILABLE") {
		t.Fatalf("expected offline fallback, got %+v, %v", resp, err)
	}
}

func TestGeminiStreamCancellationClosesBody(t *testing.T) {
	url, _, gone := startSSE(t, []sseStep{{raw: geminiChunk(geminiText("First words"), "", 0), pause: 10 * time.Second}})
	ft := tutor.NewFallbackTutor(geminiStream(url), nil, 20*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	_, err := ft.ExplainStream(ctx, tutor.Request{}, func(tutor.StreamUpdate) { cancel() })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected caller cancellation, got %v", err)
	}
	waitClosed(t, gone)
}

func responsesEvent(v map[string]interface{}) string {
	b, _ := json.Marshal(v)
	return "data: " + string(b) + "\n\n"
}

func chatGPTStream(t *testing.T, url string) *tutor.OpenAIChatGPTPlanTutor {
	t.Helper()
	store, err := tutor.NewAuthStore(t.TempDir() + "/auth.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = store.SetChatGPTPlanToken(&tutor.OAuthToken{
		AccessToken: "test-valid-jwt", Subject: "user", Scope: tutor.DefaultOAuthScope, ClientID: "oaiapp_test",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	})
	return tutor.NewOpenAIChatGPTPlanTutor(tutor.ChatGPTPlanConfig{AuthStore: store, APIURL: url, Model: "m"})
}

func TestChatGPTPlanStreamSurfacesDeltas(t *testing.T) {
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "text/event-stream")
		for _, e := range []string{
			"event: response.created\n" + responsesEvent(map[string]interface{}{"type": "response.created"}),
			responsesEvent(map[string]interface{}{"type": "response.reasoning_summary_text.delta", "delta": "secret reasoning"}),
			responsesEvent(map[string]interface{}{"type": "response.output_text.delta", "delta": "Consider the economic "}),
			responsesEvent(map[string]interface{}{"type": "response.output_text.delta", "delta": "reality."}),
			responsesEvent(map[string]interface{}{"type": "response.completed", "response": map[string]interface{}{"usage": map[string]int{"output_tokens": 6}}}),
		} {
			_, _ = w.Write([]byte(e))
			w.(http.Flusher).Flush()
		}
	}))
	defer srv.Close()

	var updates []tutor.StreamUpdate
	ft := tutor.NewFallbackTutor(chatGPTStream(t, srv.URL), nil, 2*time.Second)
	if !ft.CanStream() {
		t.Fatal("ChatGPT plan should stream")
	}
	resp, err := ft.HintStream(context.Background(), tutor.Request{}, collect(&updates))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body["stream"] != true || body["store"] != false {
		t.Fatalf("subscription invariants not sent: %v", body)
	}
	if resp.Text != "Consider the economic reality." || resp.Provider != "chatgpt_plan" || resp.TokensUsed != 6 || resp.Fallback {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(updates) != 3 || !updates[0].Thinking || updates[0].Text != "" || updates[1].Text != "Consider the economic" {
		t.Fatalf("unexpected updates: %+v", updates)
	}
}

func TestChatGPTPlanStreamEndBeforeCompletedIsIncomplete(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: responsesEvent(map[string]interface{}{"type": "response.output_text.delta", "delta": "Debit means left"}) + "data: [DONE]\n\n"}})
	ft := tutor.NewFallbackTutor(chatGPTStream(t, url), nil, 2*time.Second)
	resp, err := ft.ExplainStream(context.Background(), tutor.Request{}, nil)
	if err != nil || !resp.Incomplete || resp.Text != "Debit means left" || !strings.Contains(resp.FallbackReason, "response.completed") {
		t.Fatalf("expected an incomplete partial reply, got %+v, %v", resp, err)
	}
}

func TestChatGPTPlanFailedEventBeforeTextFallsBack(t *testing.T) {
	url, _, _ := startSSE(t, []sseStep{{raw: responsesEvent(map[string]interface{}{"type": "response.failed", "response": map[string]interface{}{"error": map[string]string{"code": "server_error"}}})}})
	ft := tutor.NewFallbackTutor(chatGPTStream(t, url), nil, 2*time.Second)
	resp, err := ft.HintStream(context.Background(), tutor.Request{}, nil)
	if err != nil || !resp.Fallback || !strings.Contains(resp.FallbackReason, "server_error") {
		t.Fatalf("expected offline fallback, got %+v, %v", resp, err)
	}
}

func TestAllNetworkProvidersStream(t *testing.T) {
	for _, provider := range []string{tutor.ProviderAnthropic, tutor.ProviderGoogle, tutor.ProviderChatGPTPlan, tutor.ProviderOpenAI, tutor.ProviderLMStudio} {
		ft, ok := tutor.BuildTutor(tutor.FactoryOptions{Provider: provider}).(*tutor.FallbackTutor)
		if !ok || !ft.CanStream() {
			t.Errorf("%s: expected a streaming FallbackTutor", provider)
		}
	}
}
