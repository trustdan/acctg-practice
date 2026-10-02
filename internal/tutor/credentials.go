package tutor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Supported provider identifiers
const (
	ProviderOffline     = "offline"
	ProviderChatGPTPlan = "chatgpt_plan"
	ProviderAnthropic   = "anthropic"
	ProviderGoogle      = "google"
	ProviderOpenAI      = "openai"
)

// OAuthToken holds OAuth credentials for ChatGPT Plus subscription plan access.
type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	ClientID     string    `json:"client_id,omitempty"`
	PlanType     string    `json:"plan_type,omitempty"`
	Scope        string    `json:"scope,omitempty"`
}

// IsExpired checks if the access token has expired (with 60-second grace window).
func (t *OAuthToken) IsExpired() bool {
	if t == nil || t.AccessToken == "" {
		return true
	}
	return time.Now().UTC().Add(60 * time.Second).After(t.ExpiresAt)
}

// AuthConfig stores learner tutor settings and credentials in local user storage.
//
// Invariant (AGENT-CONTRACT):
// Never include secrets in commits or plain tracked files. Credentials are stored
// strictly in the OS user data directory with restricted file permissions (0600).
type AuthConfig struct {
	ActiveProvider   string      `json:"active_provider"`
	AnthropicKey     string      `json:"anthropic_key,omitempty"`
	GoogleKey        string      `json:"google_key,omitempty"`
	OpenAIKey        string      `json:"openai_key,omitempty"`
	ChatGPTPlanToken *OAuthToken `json:"chatgpt_plan_token,omitempty"`
}

// AuthStore provides thread-safe access to persistent tutor credentials.
type AuthStore struct {
	mu       sync.Mutex
	filePath string
	config   AuthConfig
}

// DefaultAuthPath returns the default path to credentials file in user data directory.
func DefaultAuthPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, hErr := os.UserHomeDir()
		if hErr != nil {
			return "", fmt.Errorf("could not determine user data dir: %w", err)
		}
		configDir = filepath.Join(home, ".config")
	}
	appDir := filepath.Join(configDir, "acctg-practice")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", fmt.Errorf("could not create application data directory %s: %w", appDir, err)
	}
	return filepath.Join(appDir, "tutor_auth.json"), nil
}

// NewAuthStore loads or initializes an AuthStore at the given path.
func NewAuthStore(filePath string) (*AuthStore, error) {
	store := &AuthStore{
		filePath: filePath,
		config: AuthConfig{
			ActiveProvider: ProviderOffline,
		},
	}

	if filePath == "" {
		return store, nil
	}

	data, err := os.ReadFile(filePath)
	if err == nil {
		var loaded AuthConfig
		if err := json.Unmarshal(data, &loaded); err == nil {
			store.config = loaded
			if store.config.ActiveProvider == "" {
				store.config.ActiveProvider = ProviderOffline
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed reading auth config: %w", err)
	}

	return store, nil
}

// GetConfig returns a copy of current auth config.
func (s *AuthStore) GetConfig() AuthConfig {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config
}

// SetActiveProvider updates the active provider and persists changes.
func (s *AuthStore) SetActiveProvider(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.ActiveProvider = provider
	return s.saveLocked()
}

// SetAPIKey saves an API key for a specific provider and persists changes.
func (s *AuthStore) SetAPIKey(provider string, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trimmed := strings.TrimSpace(key)
	switch provider {
	case ProviderAnthropic:
		s.config.AnthropicKey = trimmed
	case ProviderGoogle:
		s.config.GoogleKey = trimmed
	case ProviderOpenAI:
		s.config.OpenAIKey = trimmed
	default:
		return fmt.Errorf("unknown api key provider: %s", provider)
	}

	return s.saveLocked()
}

// SetChatGPTPlanToken saves OAuth tokens for ChatGPT Plus and persists changes.
func (s *AuthStore) SetChatGPTPlanToken(token *OAuthToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.ChatGPTPlanToken = token
	return s.saveLocked()
}

// ClearCredentials clears stored credentials for a provider.
func (s *AuthStore) ClearCredentials(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch provider {
	case ProviderChatGPTPlan:
		s.config.ChatGPTPlanToken = nil
	case ProviderAnthropic:
		s.config.AnthropicKey = ""
	case ProviderGoogle:
		s.config.GoogleKey = ""
	case ProviderOpenAI:
		s.config.OpenAIKey = ""
	}

	return s.saveLocked()
}

// ResolveKey returns the active key for a provider, checking environment variables first.
func (s *AuthStore) ResolveKey(provider string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch provider {
	case ProviderAnthropic:
		if env := os.Getenv("ANTHROPIC_API_KEY"); env != "" {
			return strings.TrimSpace(env)
		}
		return s.config.AnthropicKey

	case ProviderGoogle:
		if env := os.Getenv("GEMINI_API_KEY"); env != "" {
			return strings.TrimSpace(env)
		}
		return s.config.GoogleKey

	case ProviderOpenAI:
		if env := os.Getenv("OPENAI_API_KEY"); env != "" {
			return strings.TrimSpace(env)
		}
		return s.config.OpenAIKey

	default:
		return ""
	}
}

// IsConfigured returns true if the specified provider has valid credentials available.
func (s *AuthStore) IsConfigured(provider string) bool {
	switch provider {
	case ProviderOffline:
		return true
	case ProviderChatGPTPlan:
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.config.ChatGPTPlanToken != nil && s.config.ChatGPTPlanToken.AccessToken != ""
	case ProviderAnthropic, ProviderGoogle, ProviderOpenAI:
		return s.ResolveKey(provider) != ""
	default:
		return false
	}
}

func (s *AuthStore) saveLocked() error {
	if s.filePath == "" {
		return nil
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed creating credentials dir: %w", err)
	}

	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed encoding credentials: %w", err)
	}

	// Invariant (AGENT-CONTRACT): 0600 permissions strictly enforced
	if err := os.WriteFile(s.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed writing credentials file: %w", err)
	}

	return nil
}
