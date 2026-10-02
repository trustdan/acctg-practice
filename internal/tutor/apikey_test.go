package tutor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAnthropicTutorSuccessAndHeaders(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			http.NotFound(w, r)
			return
		}

		// Verify Anthropic-specific headers
		if r.Header.Get("x-api-key") != "test-anthropic-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			http.Error(w, "bad version", http.StatusBadRequest)
			return
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		_ = json.Unmarshal(bodyBytes, &body)
		if body["system"] == "" {
			http.Error(w, "missing system prompt", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": "Anthropic Socratic hint: Cash is an asset that increased on debit."},
			},
			"usage": map[string]int{
				"output_tokens": 15,
			},
		})
	}))
	defer mockServer.Close()

	tutor := NewAnthropicTutor(APIKeyConfig{
		APIKey:     "test-anthropic-key",
		BaseURL:    mockServer.URL,
		HTTPClient: mockServer.Client(),
	})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Acme receives cash."}

	resp, err := tutor.Hint(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Anthropic Socratic hint") {
		t.Errorf("unexpected text: %s", resp.Text)
	}
	if resp.Provider != "anthropic" {
		t.Errorf("expected provider 'anthropic', got %s", resp.Provider)
	}
	if resp.TokensUsed != 15 {
		t.Errorf("expected 15 tokens, got %d", resp.TokensUsed)
	}
}

func TestGeminiTutorSuccessAndFormat(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "generateContent") {
			http.NotFound(w, r)
			return
		}

		// Verify Gemini API key in query parameter
		if r.URL.Query().Get("key") != "test-gemini-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]string{
							{"text": "Gemini Socratic hint: Did we perform work or just take a deposit?"},
						},
					},
				},
			},
			"usageMetadata": map[string]int{
				"candidatesTokenCount": 18,
			},
		})
	}))
	defer mockServer.Close()

	tutor := NewGeminiTutor(APIKeyConfig{
		APIKey:     "test-gemini-key",
		BaseURL:    mockServer.URL,
		HTTPClient: mockServer.Client(),
	})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Customer advance payment."}

	resp, err := tutor.Hint(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "Gemini Socratic hint") {
		t.Errorf("unexpected text: %s", resp.Text)
	}
	if resp.Provider != "google" {
		t.Errorf("expected provider 'google', got %s", resp.Provider)
	}
	if resp.TokensUsed != 18 {
		t.Errorf("expected 18 tokens, got %d", resp.TokensUsed)
	}
}

func TestOpenAIAPITutorSuccessAndBearerAuth(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}

		if r.Header.Get("Authorization") != "Bearer test-openai-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "OpenAI Socratic hint: Verify debit equals credit.",
					},
				},
			},
			"usage": map[string]int{
				"completion_tokens": 12,
			},
		})
	}))
	defer mockServer.Close()

	tutor := NewOpenAIAPITutor(APIKeyConfig{
		APIKey:     "test-openai-key",
		BaseURL:    mockServer.URL,
		HTTPClient: mockServer.Client(),
	})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Balanced journal entry."}

	resp, err := tutor.Hint(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Text, "OpenAI Socratic hint") {
		t.Errorf("unexpected text: %s", resp.Text)
	}
	if resp.Provider != "openai" {
		t.Errorf("expected provider 'openai', got %s", resp.Provider)
	}
}

func TestMissingAPIKeysReturnClearError(t *testing.T) {
	anthropic := NewAnthropicTutor(APIKeyConfig{})
	gemini := NewGeminiTutor(APIKeyConfig{})
	openai := NewOpenAIAPITutor(APIKeyConfig{})

	ctx := context.Background()
	req := Request{ProblemPrompt: "Test"}

	if _, err := anthropic.Hint(ctx, req); err == nil || !strings.Contains(err.Error(), "anthropic api key not configured") {
		t.Errorf("expected missing anthropic key error, got: %v", err)
	}
	if _, err := gemini.Hint(ctx, req); err == nil || !strings.Contains(err.Error(), "gemini api key not configured") {
		t.Errorf("expected missing gemini key error, got: %v", err)
	}
	if _, err := openai.Hint(ctx, req); err == nil || !strings.Contains(err.Error(), "openai api key not configured") {
		t.Errorf("expected missing openai key error, got: %v", err)
	}
}

func TestNetworkTutorFallbackToOfflineOnTimeout(t *testing.T) {
	// Hang server
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer slowServer.Close()

	primary := NewAnthropicTutor(APIKeyConfig{
		APIKey:     "key",
		BaseURL:    slowServer.URL,
		HTTPClient: slowServer.Client(),
	})
	offline := NewOfflineTutor()
	wrapper := NewFallbackTutor(primary, offline, 50*time.Millisecond)

	req := Request{
		ProblemPrompt: "Acme pays rent.",
		CausalHint:    "Rent is an expense.",
	}

	resp, err := wrapper.Hint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Fallback {
		t.Errorf("expected Fallback=true on timeout")
	}
	if resp.Provider != "offline" {
		t.Errorf("expected offline provider from fallback, got %s", resp.Provider)
	}
}
