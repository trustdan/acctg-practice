package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestModelCacheDefaultsAndPersistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "models_cache_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cacheFile := filepath.Join(tempDir, "models_cache.json")

	// 1. Uninitialized cache returns default models
	cache, err := NewModelCache(cacheFile)
	if err != nil {
		t.Fatalf("failed creating model cache: %v", err)
	}

	antDefaults := cache.GetModels(ProviderAnthropic)
	if len(antDefaults) == 0 {
		t.Fatalf("expected non-empty default models for anthropic")
	}
	if antDefaults[0].ID != "claude-3-5-haiku-20241022" {
		t.Errorf("expected default haiku, got %s", antDefaults[0].ID)
	}

	gemDefaults := cache.GetModels(ProviderGoogle)
	if len(gemDefaults) == 0 || gemDefaults[0].ID != "gemini-1.5-flash" {
		t.Errorf("unexpected gemini defaults: %+v", gemDefaults)
	}

	oaiDefaults := cache.GetModels(ProviderOpenAI)
	if len(oaiDefaults) == 0 || oaiDefaults[0].ID != "gpt-4o-mini" {
		t.Errorf("unexpected openai defaults: %+v", oaiDefaults)
	}

	// 2. Set discovered models and verify persistence
	customModels := []ModelInfo{
		{ID: "claude-custom-1", DisplayName: "Custom Claude 1", Provider: ProviderAnthropic},
		{ID: "claude-custom-2", DisplayName: "Custom Claude 2", Provider: ProviderAnthropic},
	}
	if err := cache.SetModels(ProviderAnthropic, customModels); err != nil {
		t.Fatalf("failed setting models: %v", err)
	}

	retrieved := cache.GetModels(ProviderAnthropic)
	if len(retrieved) != 2 || retrieved[0].ID != "claude-custom-1" {
		t.Fatalf("expected retrieved custom models, got %+v", retrieved)
	}

	// 3. Reload from disk
	reloaded, err := NewModelCache(cacheFile)
	if err != nil {
		t.Fatalf("failed reloading cache: %v", err)
	}
	reloadedModels := reloaded.GetModels(ProviderAnthropic)
	if len(reloadedModels) != 2 || reloadedModels[0].ID != "claude-custom-1" {
		t.Errorf("expected reloaded custom models from disk, got %+v", reloadedModels)
	}

	// Non-updated provider still returns defaults
	if models := reloaded.GetModels(ProviderGoogle); len(models) == 0 {
		t.Errorf("expected default models for non-updated provider")
	}
}

func TestDiscoverOpenAIModelsMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			http.NotFound(w, r)
			return
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-openai-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp := map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "text-embedding-3-small", "created": 1705948997},
				{"id": "gpt-4o", "created": 1715368132},
				{"id": "whisper-1", "created": 1677532384},
				{"id": "gpt-4o-mini", "created": 1721172741},
				{"id": "o3-mini", "created": 1737500000},
				{"id": "dall-e-3", "created": 1698785189},
				{"id": "o1-preview", "created": 1725000000},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	ctx := context.Background()
	models, err := DiscoverOpenAIModels(ctx, "test-openai-key", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("failed discovering openai models: %v", err)
	}

	// Non-chat models should be filtered out
	for _, m := range models {
		if m.ID == "text-embedding-3-small" || m.ID == "whisper-1" || m.ID == "dall-e-3" {
			t.Errorf("expected non-chat model %s to be filtered out", m.ID)
		}
	}

	// Priority sorting: gpt-4o-mini (1), gpt-4o (2), o3-mini (3)
	if len(models) < 3 {
		t.Fatalf("expected at least 3 models, got %d", len(models))
	}
	if models[0].ID != "gpt-4o-mini" {
		t.Errorf("expected first model to be gpt-4o-mini, got %s", models[0].ID)
	}
	if models[1].ID != "gpt-4o" {
		t.Errorf("expected second model to be gpt-4o, got %s", models[1].ID)
	}
	if models[2].ID != "o3-mini" {
		t.Errorf("expected third model to be o3-mini, got %s", models[2].ID)
	}
}

func TestDiscoverAnthropicModelsMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("x-api-key") != "test-anthropic-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			http.Error(w, "Invalid version", http.StatusBadRequest)
			return
		}

		resp := map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "claude-3-5-haiku-20241022", "display_name": "Claude 3.5 Haiku", "type": "model"},
				{"id": "claude-3-5-sonnet-20241022", "display_name": "Claude 3.5 Sonnet", "type": "model"},
				{"id": "claude-3-7-sonnet-20250219", "display_name": "Claude 3.7 Sonnet", "type": "model"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	ctx := context.Background()
	models, err := DiscoverAnthropicModels(ctx, "test-anthropic-key", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("failed discovering anthropic models: %v", err)
	}

	if len(models) != 3 {
		t.Fatalf("expected 3 models, got %d", len(models))
	}
	if models[0].ID != "claude-3-5-haiku-20241022" {
		t.Errorf("expected first model haiku, got %s", models[0].ID)
	}
	if models[1].ID != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected second model sonnet, got %s", models[1].ID)
	}
}

func TestDiscoverGeminiModelsMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-gemini-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp := map[string]interface{}{
			"models": []map[string]interface{}{
				{
					"name":                       "models/gemini-1.5-flash",
					"displayName":                "Gemini 1.5 Flash",
					"supportedGenerationMethods": []string{"generateContent", "countTokens"},
				},
				{
					"name":                       "models/text-embedding-004",
					"displayName":                "Text Embedding",
					"supportedGenerationMethods": []string{"embedContent"},
				},
				{
					"name":                       "models/gemini-1.5-pro",
					"displayName":                "Gemini 1.5 Pro",
					"supportedGenerationMethods": []string{"generateContent"},
				},
				{
					"name":                       "models/gemini-2.0-flash",
					"displayName":                "Gemini 2.0 Flash",
					"supportedGenerationMethods": []string{"generateContent"},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	ctx := context.Background()
	models, err := DiscoverGeminiModels(ctx, "test-gemini-key", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("failed discovering gemini models: %v", err)
	}

	// embedContent-only model should be filtered out
	for _, m := range models {
		if m.ID == "text-embedding-004" {
			t.Errorf("expected embedding model to be filtered out")
		}
		if strings.HasPrefix(m.ID, "models/") {
			t.Errorf("expected models/ prefix to be stripped, got %s", m.ID)
		}
	}

	if len(models) != 3 {
		t.Fatalf("expected 3 models supporting generateContent, got %d", len(models))
	}
	if models[0].ID != "gemini-1.5-flash" {
		t.Errorf("expected first model gemini-1.5-flash, got %s", models[0].ID)
	}
}

func TestModelSelectionPersistenceAndResolution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "model_selection_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")
	store, err := NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed creating auth store: %v", err)
	}

	// 1. Initial resolution returns defaults
	if m := store.ResolveModel(ProviderAnthropic); m != DefaultAnthropicModel {
		t.Errorf("expected default anthropic model %s, got %s", DefaultAnthropicModel, m)
	}
	if m := store.ResolveModel(ProviderGoogle); m != DefaultGeminiModel {
		t.Errorf("expected default gemini model %s, got %s", DefaultGeminiModel, m)
	}
	if m := store.ResolveModel(ProviderOpenAI); m != DefaultOpenAIAPIModel {
		t.Errorf("expected default openai model %s, got %s", DefaultOpenAIAPIModel, m)
	}

	// 2. Select custom models and verify resolution
	if err := store.SetSelectedModel(ProviderAnthropic, "claude-3-7-sonnet-20250219"); err != nil {
		t.Fatalf("failed setting anthropic model: %v", err)
	}
	if err := store.SetSelectedModel(ProviderGoogle, "gemini-2.0-flash"); err != nil {
		t.Fatalf("failed setting gemini model: %v", err)
	}
	if err := store.SetSelectedModel(ProviderOpenAI, "o3-mini"); err != nil {
		t.Fatalf("failed setting openai model: %v", err)
	}

	if m := store.ResolveModel(ProviderAnthropic); m != "claude-3-7-sonnet-20250219" {
		t.Errorf("expected selected anthropic model claude-3-7-sonnet-20250219, got %s", m)
	}
	if m := store.ResolveModel(ProviderGoogle); m != "gemini-2.0-flash" {
		t.Errorf("expected selected gemini model gemini-2.0-flash, got %s", m)
	}
	if m := store.ResolveModel(ProviderOpenAI); m != "o3-mini" {
		t.Errorf("expected selected openai model o3-mini, got %s", m)
	}

	// 3. Verify persistence across store reloads
	reloaded, err := NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed reloading auth store: %v", err)
	}
	if m := reloaded.ResolveModel(ProviderAnthropic); m != "claude-3-7-sonnet-20250219" {
		t.Errorf("expected persisted anthropic model claude-3-7-sonnet-20250219, got %s", m)
	}
	if m := reloaded.GetSelectedModel(ProviderOpenAI); m != "o3-mini" {
		t.Errorf("expected GetSelectedModel o3-mini, got %s", m)
	}

	// 4. Change selected model again (dynamic switching)
	if err := reloaded.SetSelectedModel(ProviderAnthropic, "claude-3-5-sonnet-20241022"); err != nil {
		t.Fatalf("failed changing model: %v", err)
	}
	if m := reloaded.ResolveModel(ProviderAnthropic); m != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected switched model claude-3-5-sonnet-20241022, got %s", m)
	}

	// 5. BuildTutor respects user-selected model
	tut := BuildTutor(FactoryOptions{
		Provider:  ProviderAnthropic,
		AuthStore: reloaded,
	})
	fb, ok := tut.(*FallbackTutor)
	if !ok {
		t.Fatalf("expected FallbackTutor")
	}
	ant, ok := fb.primary.(*AnthropicTutor)
	if !ok {
		t.Fatalf("expected primary to be AnthropicTutor")
	}
	if ant.GetModel() != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected tutor to have active model claude-3-5-sonnet-20241022, got %s", ant.GetModel())
	}
}

func TestChatGPTClientIDConfigurationAndResolution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "client_id_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")
	store, _ := NewAuthStore(authPath)

	// Default
	if id := store.ResolveChatGPTClientID(); id != DefaultDCRClientID {
		t.Errorf("expected default client ID %s, got %s", DefaultDCRClientID, id)
	}

	// Configured in store
	if err := store.SetChatGPTClientID("custom-org-client-id"); err != nil {
		t.Fatalf("failed setting client id: %v", err)
	}
	if id := store.ResolveChatGPTClientID(); id != DefaultDCRClientID {
		t.Errorf("expected custom-org-client-id, got %s", id)
	}

	// Environment variable takes precedence
	origEnv := os.Getenv("OPENAI_OAUTH_CLIENT_ID")
	defer os.Setenv("OPENAI_OAUTH_CLIENT_ID", origEnv)

	os.Setenv("OPENAI_OAUTH_CLIENT_ID", "env-org-client-id")
	if id := store.ResolveChatGPTClientID(); id != DefaultDCRClientID {
		t.Errorf("expected env-org-client-id, got %s", id)
	}
}

func TestChatGPTAccountCatalogSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer account-token" {
			t.Error("missing OAuth bearer")
		}
		fmt.Fprint(w, `{"models":[{"slug":"gpt-new","display_name":"Current Model","visibility":"list"},{"slug":"hidden","visibility":"hidden"},{"slug":"gpt-other","display_name":"Other Model","visibility":"list"}]}`)
	}))
	defer server.Close()
	models, err := DiscoverChatGPTModels(context.Background(), "account-token", server.URL, server.Client())
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 2 || models[0].ID != "gpt-new" || models[0].DisplayName != "Current Model" || models[1].ID != "gpt-other" {
		t.Fatalf("wrong account catalog: %+v", models)
	}
}
func TestChatGPTDiscoveryDoesNotReplaceFailureWithDefaults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "denied", 403) }))
	defer server.Close()
	store, _ := NewAuthStore("")
	_ = store.SetChatGPTPlanToken(&OAuthToken{AccessToken: "account-token", Subject: "user", Scope: DefaultOAuthScope, ClientID: "oaiapp_test", ExpiresAt: time.Now().Add(time.Hour)})
	models, err := DiscoverProviderModels(context.Background(), ProviderChatGPTPlan, store, server.URL, server.Client())
	if err == nil || len(models) != 0 {
		t.Fatal("discovery failure silently became defaults")
	}
}
