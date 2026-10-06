package tutor_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

var lmReq = tutor.Request{
	ProblemPrompt: "The company pays $1,000 cash for rent today.",
	FamilyID:      "pay_rent_cash",
	Stage:         domain.StageIdentifyAccount,
	ConceptID:     "rent_expense_vs_cash",
	CausalHint:    "Rent is consumed in the period.",
}

// fakeLMStudio serves /v1/models and /v1/chat/completions like LM Studio.
type fakeLMStudio struct {
	models     []string
	content    string
	delay      time.Duration
	authHeader atomic.Value
	model      atomic.Value
	maxTokens  atomic.Int64
	calls      atomic.Int64
}

func (f *fakeLMStudio) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		f.calls.Add(1)
		f.authHeader.Store(r.Header.Get("Authorization"))
		data := []map[string]string{}
		for _, id := range f.models {
			data = append(data, map[string]string{"id": id, "object": "model"})
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"object": "list", "data": data})
	})
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		f.calls.Add(1)
		f.authHeader.Store(r.Header.Get("Authorization"))
		var body struct {
			Model     string `json:"model"`
			MaxTokens int64  `json:"max_tokens"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.model.Store(body.Model)
		f.maxTokens.Store(body.MaxTokens)
		if f.delay > 0 {
			select {
			case <-time.After(f.delay):
			case <-r.Context().Done():
				return
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{"message": map[string]string{
				"role": "assistant", "content": f.content, "reasoning_content": "secret chain of thought",
			}}},
			"usage": map[string]int{"completion_tokens": 12},
		})
	})
	return mux
}

func startFake(t *testing.T, f *fakeLMStudio) string {
	t.Helper()
	srv := httptest.NewServer(f.handler())
	t.Cleanup(srv.Close)
	return srv.URL + "/v1"
}

func closedServerURL() string {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL + "/v1"
	srv.Close()
	return url
}

func TestLMStudioTutorSuccessUsesFirstModelWithoutKey(t *testing.T) {
	t.Setenv("LM_API_TOKEN", "")
	f := &fakeLMStudio{models: []string{"text-embedding-nomic", "qwen3-8b", "llama-3.2-3b"}, content: "<think>debits are...</think>What did the company receive?"}
	tut := tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: startFake(t, f)})

	resp, err := tut.Hint(context.Background(), lmReq)
	if err != nil {
		t.Fatalf("hint failed: %v", err)
	}
	if resp.Text != "What did the company receive?" {
		t.Errorf("think block not stripped or reasoning_content leaked: %q", resp.Text)
	}
	if resp.Provider != tutor.ProviderLMStudio || resp.TokensUsed != 12 {
		t.Errorf("unexpected provenance: %+v", resp)
	}
	if got := f.model.Load(); got != "qwen3-8b" {
		t.Errorf("expected first non-embedding model, got %v", got)
	}
	if got := f.maxTokens.Load(); got != tutor.DefaultLMStudioMaxTokens {
		t.Errorf("expected local default max_tokens %d, got %d", tutor.DefaultLMStudioMaxTokens, got)
	}
	if got := f.authHeader.Load(); got != "" {
		t.Errorf("expected no Authorization header without LM_API_TOKEN, got %q", got)
	}
}

func TestLMStudioTutorSendsOptionalToken(t *testing.T) {
	t.Setenv("LM_API_TOKEN", "lm-secret")
	f := &fakeLMStudio{models: []string{"m"}, content: "Hint."}
	tut := tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: startFake(t, f), Model: "m"})
	if _, err := tut.Explain(context.Background(), lmReq); err != nil {
		t.Fatalf("explain failed: %v", err)
	}
	if got := f.authHeader.Load(); got != "Bearer lm-secret" {
		t.Errorf("expected bearer token, got %q", got)
	}
}

func TestStripThinkBlocks(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"<think>a\nb</think>answer", "answer"},
		{"<THINK>x</THINK>one <think>y</think>two", "one two"},
		{"answer<think>truncated mid-thought", "answer"},
		{"<think>only thinking", ""},
	}
	for _, c := range cases {
		if got := tutor.StripThinkBlocks(c.in); got != c.want {
			t.Errorf("StripThinkBlocks(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLMStudioOnlyThinkingFallsBackToOffline(t *testing.T) {
	f := &fakeLMStudio{models: []string{"m"}, content: "<think>ran out of tokens"}
	tut := tutor.NewFallbackTutor(tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: startFake(t, f), Model: "m"}), nil, 5*time.Second)
	resp, err := tut.Hint(context.Background(), lmReq)
	if err != nil || !resp.Fallback || !strings.Contains(resp.FallbackReason, "max_tokens") {
		t.Fatalf("expected offline fallback explaining max_tokens, got %+v err=%v", resp, err)
	}
}

func TestLMStudioConnectionRefusedFallsBack(t *testing.T) {
	tut := tutor.NewFallbackTutor(tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: closedServerURL()}), nil, 5*time.Second)
	resp, err := tut.Hint(context.Background(), lmReq)
	if err != nil {
		t.Fatalf("expected fallback, got error %v", err)
	}
	if !resp.Fallback || !strings.Contains(resp.FallbackReason, "lms server start") {
		t.Errorf("expected server-not-running guidance, got %+v", resp)
	}
}

func TestLMStudioTimeoutFallsBack(t *testing.T) {
	f := &fakeLMStudio{models: []string{"m"}, content: "late", delay: 2 * time.Second}
	tut := tutor.NewFallbackTutor(tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: startFake(t, f), Model: "m"}), nil, 100*time.Millisecond)
	start := time.Now()
	resp, err := tut.Hint(context.Background(), lmReq)
	if err != nil || !resp.Fallback {
		t.Fatalf("expected fallback on timeout, got %+v err=%v", resp, err)
	}
	if time.Since(start) > time.Second {
		t.Errorf("timeout did not cut the request short")
	}
}

func TestLMStudioCancellationReturnsCallerError(t *testing.T) {
	f := &fakeLMStudio{models: []string{"m"}, content: "late", delay: 2 * time.Second}
	tut := tutor.NewFallbackTutor(tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: startFake(t, f), Model: "m"}), nil, 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(50*time.Millisecond, cancel)
	if _, err := tut.Hint(ctx, lmReq); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected caller cancellation, got %v", err)
	}
}

func TestDiscoverLMStudioModels(t *testing.T) {
	f := &fakeLMStudio{models: []string{"text-embedding-nomic-embed-text-v1.5", "qwen3-8b"}}
	models, err := tutor.DiscoverLMStudioModels(context.Background(), startFake(t, f), "", nil)
	if err != nil || len(models) != 1 || models[0].ID != "qwen3-8b" || models[0].Provider != tutor.ProviderLMStudio {
		t.Fatalf("unexpected models %+v err=%v", models, err)
	}

	empty := &fakeLMStudio{models: []string{"text-embedding-only"}}
	if _, err := tutor.DiscoverLMStudioModels(context.Background(), startFake(t, empty), "", nil); !errors.Is(err, tutor.ErrNoLMStudioModel) {
		t.Fatalf("expected ErrNoLMStudioModel, got %v", err)
	}

	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer authSrv.Close()
	if _, err := tutor.DiscoverLMStudioModels(context.Background(), authSrv.URL+"/v1", "", nil); err == nil || !strings.Contains(err.Error(), "LM_API_TOKEN") {
		t.Fatalf("expected LM_API_TOKEN guidance, got %v", err)
	}
}

func TestValidateLMStudioURL(t *testing.T) {
	cases := []struct {
		url         string
		allowRemote bool
		ok          bool
	}{
		{"http://localhost:1234/v1", false, true},
		{"http://127.0.0.1:1234/v1", false, true},
		{"http://[::1]:1234/v1", false, true},
		{"http://192.168.1.20:1234/v1", false, false},
		{"https://lm.example.com/v1", false, false},
		{"http://192.168.1.20:1234/v1", true, true},
		{"localhost:1234", false, false},
		{"ftp://localhost/v1", false, false},
	}
	for _, c := range cases {
		err := tutor.ValidateLMStudioURL(c.url, c.allowRemote)
		if (err == nil) != c.ok {
			t.Errorf("ValidateLMStudioURL(%q, %v) err=%v, want ok=%v", c.url, c.allowRemote, err, c.ok)
		}
	}
}

func TestLMStudioRejectsRemoteHostWithoutRequest(t *testing.T) {
	tut := tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: "http://10.0.0.5:1234/v1", Model: "m"})
	if _, err := tut.Hint(context.Background(), lmReq); err == nil || !strings.Contains(err.Error(), "not on this machine") {
		t.Fatalf("expected loopback rejection, got %v", err)
	}
}

func TestBuildTutorLMStudioMakesNoNetworkCall(t *testing.T) {
	f := &fakeLMStudio{models: []string{"m"}, content: "ok"}
	baseURL := startFake(t, f)
	store, _ := tutor.NewAuthStore("")
	if err := store.SetLMStudioURL(baseURL); err != nil {
		t.Fatal(err)
	}
	if !store.IsConfigured(tutor.ProviderLMStudio) {
		t.Fatal("LM Studio should count as configured with no key")
	}
	built := tutor.BuildTutor(tutor.FactoryOptions{Provider: tutor.ProviderLMStudio, AuthStore: store, Timeout: 5 * time.Second})
	if f.calls.Load() != 0 {
		t.Fatalf("building the tutor contacted the server %d times", f.calls.Load())
	}
	fb, ok := built.(*tutor.FallbackTutor)
	if !ok || fb.Timeout() < tutor.DefaultLMStudioTimeout {
		t.Fatalf("expected FallbackTutor with >= %s timeout, got %T", tutor.DefaultLMStudioTimeout, built)
	}
	resp, err := built.Hint(context.Background(), lmReq)
	if err != nil || resp.Fallback || resp.Text != "ok" {
		t.Fatalf("expected live LM Studio reply, got %+v err=%v", resp, err)
	}
}

func TestLMStudioURLResolution(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", "")
	t.Setenv("LMSTUDIO_ALLOW_REMOTE", "")
	store, _ := tutor.NewAuthStore("")
	if got := store.ResolveLMStudioURL(); got != tutor.DefaultLMStudioURL {
		t.Errorf("default URL = %q", got)
	}
	_ = store.SetLMStudioURL("http://127.0.0.1:4321/v1")
	if got := store.ResolveLMStudioURL(); got != "http://127.0.0.1:4321/v1" {
		t.Errorf("stored URL = %q", got)
	}
	t.Setenv("LMSTUDIO_BASE_URL", "http://localhost:9999/v1")
	if got := store.ResolveLMStudioURL(); got != "http://localhost:9999/v1" {
		t.Errorf("env URL = %q", got)
	}
	if store.LMStudioAllowRemote() {
		t.Error("remote hosts must be off by default")
	}
	t.Setenv("LMSTUDIO_ALLOW_REMOTE", "1")
	if !store.LMStudioAllowRemote() {
		t.Error("LMSTUDIO_ALLOW_REMOTE=1 should opt in")
	}
	if store.ResolveLMStudioMaxTokens() != tutor.DefaultLMStudioMaxTokens {
		t.Error("unexpected default max tokens")
	}
}
