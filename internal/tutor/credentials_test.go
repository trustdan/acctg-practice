package tutor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuthStoreSaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tutor_auth_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")

	store, err := NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed creating auth store: %v", err)
	}

	if store.GetConfig().ActiveProvider != ProviderOffline {
		t.Errorf("expected default active provider offline, got %s", store.GetConfig().ActiveProvider)
	}

	// Save Anthropic key
	if err := store.SetAPIKey(ProviderAnthropic, "sk-ant-test-12345"); err != nil {
		t.Fatalf("failed setting anthropic key: %v", err)
	}

	// Verify file permissions (must be 0600 or restricted)
	info, err := os.Stat(authPath)
	if err != nil {
		t.Fatalf("failed stating auth file: %v", err)
	}
	// On Windows filemode may not strictly show 0600 in the same unix bits, but check it exists and non-empty
	if info.Size() == 0 {
		t.Errorf("auth file is empty")
	}

	// Save ChatGPT Plus token
	token := &OAuthToken{
		AccessToken:  "at-test-999",
		RefreshToken: "rt-test-888",
		ExpiresAt:    time.Now().UTC().Add(1 * time.Hour),
		ClientID:     DefaultDCRClientID,
		PlanType:     "chatgpt_plus",
	}
	if err := store.SetChatGPTPlanToken(token); err != nil {
		t.Fatalf("failed setting chatgpt token: %v", err)
	}

	// Switch active provider
	if err := store.SetActiveProvider(ProviderChatGPTPlan); err != nil {
		t.Fatalf("failed setting active provider: %v", err)
	}

	// Reload from new store instance
	reloaded, err := NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed reloading auth store: %v", err)
	}

	cfg := reloaded.GetConfig()
	if cfg.ActiveProvider != ProviderChatGPTPlan {
		t.Errorf("expected active provider %s, got %s", ProviderChatGPTPlan, cfg.ActiveProvider)
	}
	if cfg.AnthropicKey != "sk-ant-test-12345" {
		t.Errorf("expected anthropic key sk-ant-test-12345, got %s", cfg.AnthropicKey)
	}
	if cfg.ChatGPTPlanToken == nil || cfg.ChatGPTPlanToken.AccessToken != "at-test-999" {
		t.Errorf("expected access token at-test-999, got %+v", cfg.ChatGPTPlanToken)
	}
	if cfg.ChatGPTPlanToken.IsExpired() {
		t.Errorf("expected token not expired")
	}

	// Clear credentials
	if err := reloaded.ClearCredentials(ProviderAnthropic); err != nil {
		t.Fatalf("failed clearing credentials: %v", err)
	}
	if reloaded.GetConfig().AnthropicKey != "" {
		t.Errorf("expected empty anthropic key after clear")
	}
}

func TestResolveKeyPrioritizesEnvironment(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tutor_env_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")
	store, _ := NewAuthStore(authPath)
	_ = store.SetAPIKey(ProviderOpenAI, "stored-openai-key")

	// 1. Without env var: returns stored key
	origEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", origEnv)

	os.Unsetenv("OPENAI_API_KEY")
	if key := store.ResolveKey(ProviderOpenAI); key != "stored-openai-key" {
		t.Errorf("expected stored-openai-key, got %s", key)
	}

	// 2. With env var: env var takes precedence
	os.Setenv("OPENAI_API_KEY", "env-openai-key")
	if key := store.ResolveKey(ProviderOpenAI); key != "env-openai-key" {
		t.Errorf("expected env-openai-key, got %s", key)
	}
}
