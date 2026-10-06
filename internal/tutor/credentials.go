package tutor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	ProviderLMStudio    = "lmstudio"
)

// OAuthToken holds OAuth credentials for ChatGPT Plus subscription plan access.
type OAuthToken struct {
	IDToken      string    `json:"id_token,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Email        string    `json:"email,omitempty"`
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
	HostID           string            `json:"ext_agent_host_id,omitempty"`
	ActiveProvider   string            `json:"active_provider"`
	AnthropicKey     string            `json:"anthropic_key,omitempty"`
	GoogleKey        string            `json:"google_key,omitempty"`
	OpenAIKey        string            `json:"openai_key,omitempty"`
	ChatGPTPlanToken *OAuthToken       `json:"chatgpt_plan_token,omitempty"`
	ChatGPTClientID  string            `json:"chatgpt_client_id,omitempty"`
	SelectedModels   map[string]string `json:"selected_models,omitempty"`

	// LM Studio local server settings. No key is stored; LM_API_TOKEN is read from the environment only.
	LMStudioURL         string `json:"lmstudio_url,omitempty"`
	LMStudioAllowRemote bool   `json:"lmstudio_allow_remote,omitempty"`
	LMStudioMaxTokens   int    `json:"lmstudio_max_tokens,omitempty"`
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
	copyConfig := s.config
	if s.config.ChatGPTPlanToken != nil {
		token := *s.config.ChatGPTPlanToken
		copyConfig.ChatGPTPlanToken = &token
	}
	if s.config.SelectedModels != nil {
		copyConfig.SelectedModels = make(map[string]string)
		for k, v := range s.config.SelectedModels {
			copyConfig.SelectedModels[k] = v
		}
	}
	return copyConfig
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
	if token != nil {
		copyToken := *token
		token = &copyToken
	}
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
	case ProviderLMStudio:
		s.config.LMStudioURL = ""
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
		return s.config.ChatGPTPlanToken != nil && s.config.ChatGPTPlanToken.AccessToken != "" && s.config.ChatGPTPlanToken.Subject != "" && hasPlanScope(s.config.ChatGPTPlanToken.Scope)
	case ProviderAnthropic, ProviderGoogle, ProviderOpenAI:
		return s.ResolveKey(provider) != ""
	case ProviderLMStudio:
		// No network probe: a server URL (default or override) is all that is required.
		return s.ResolveLMStudioURL() != ""
	default:
		return false
	}
}

// GetSelectedModel returns the user-selected model for a provider if set.
func (s *AuthStore) GetSelectedModel(provider string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.SelectedModels == nil {
		return ""
	}
	return s.config.SelectedModels[provider]
}

// SetSelectedModel updates and persists the chosen model for a provider.
func (s *AuthStore) SetSelectedModel(provider string, model string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.SelectedModels == nil {
		s.config.SelectedModels = make(map[string]string)
	}
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		delete(s.config.SelectedModels, provider)
	} else {
		s.config.SelectedModels[provider] = trimmed
	}
	return s.saveLocked()
}

// ResolveModel returns the active model for a provider, honoring explicit user selection
// first, then falling back to vetted provider defaults.
func (s *AuthStore) ResolveModel(provider string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.SelectedModels != nil {
		if sel := strings.TrimSpace(s.config.SelectedModels[provider]); sel != "" {
			return sel
		}
	}
	switch provider {
	case ProviderAnthropic:
		return DefaultAnthropicModel
	case ProviderGoogle:
		return DefaultGeminiModel
	case ProviderOpenAI:
		return DefaultOpenAIAPIModel
	case ProviderChatGPTPlan:
		return DefaultChatGPTModel
	case ProviderLMStudio:
		return "" // First model reported by the local server.
	default:
		return "offline"
	}
}

// GetChatGPTClientID returns custom client ID if configured.
func (s *AuthStore) GetChatGPTClientID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config.ChatGPTClientID
}

// SetChatGPTClientID saves custom client ID for ChatGPT Plus OAuth.
func (s *AuthStore) SetChatGPTClientID(clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.ChatGPTClientID = strings.TrimSpace(clientID)
	return s.saveLocked()
}

// ResolveChatGPTClientID returns active client ID, checking env var first, then stored config, then default.
func (s *AuthStore) ResolveChatGPTClientID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if token := s.config.ChatGPTPlanToken; token != nil && token.Subject != "" && strings.HasPrefix(token.ClientID, "oaiapp_") {
		return token.ClientID
	}
	return DefaultDCRClientID
}

// SetLMStudioURL saves the LM Studio server base URL; empty restores the default.
func (s *AuthStore) SetLMStudioURL(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.LMStudioURL = strings.TrimSpace(url)
	return s.saveLocked()
}

// SetLMStudioAllowRemote records the explicit opt-in to a non-loopback LM Studio host.
func (s *AuthStore) SetLMStudioAllowRemote(allow bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.LMStudioAllowRemote = allow
	return s.saveLocked()
}

// SetLMStudioMaxTokens saves the completion token limit for LM Studio; 0 restores the default.
func (s *AuthStore) SetLMStudioMaxTokens(n int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n < 0 {
		n = 0
	}
	s.config.LMStudioMaxTokens = n
	return s.saveLocked()
}

// ResolveLMStudioURL returns LMSTUDIO_BASE_URL, then the stored URL, then the default.
func (s *AuthStore) ResolveLMStudioURL() string {
	if env := strings.TrimSpace(os.Getenv("LMSTUDIO_BASE_URL")); env != "" {
		return env
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.LMStudioURL != "" {
		return s.config.LMStudioURL
	}
	return DefaultLMStudioURL
}

// LMStudioAllowRemote reports whether a non-loopback LM Studio host is permitted
// (stored opt-in or LMSTUDIO_ALLOW_REMOTE=1).
func (s *AuthStore) LMStudioAllowRemote() bool {
	if v := strings.TrimSpace(os.Getenv("LMSTUDIO_ALLOW_REMOTE")); v == "1" || strings.EqualFold(v, "true") {
		return true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config.LMStudioAllowRemote
}

// ResolveLMStudioMaxTokens returns the stored LM Studio max_tokens or the local default.
func (s *AuthStore) ResolveLMStudioMaxTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.LMStudioMaxTokens > 0 {
		return s.config.LMStudioMaxTokens
	}
	return DefaultLMStudioMaxTokens
}

func (s *AuthStore) EnsureHostID() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.HostID == "" {
		s.config.HostID = "urn:uuid:" + uuid.NewString()
		if err := s.saveLocked(); err != nil {
			s.config.HostID = ""
			return "", err
		}
	}
	return s.config.HostID, nil
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
