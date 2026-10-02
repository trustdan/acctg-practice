package tutor

import (
	"time"
)

// FactoryOptions configures tutor creation.
type FactoryOptions struct {
	Provider       string        // "offline", "chatgpt_plan", "anthropic", "google", "openai"
	Model          string        // Model override (e.g. "gpt-4o", "claude-3-5-sonnet-20241022")
	Timeout        time.Duration // Default 5s
	Budget         *Budget       // Optional session limiter
	AuthStore      *AuthStore    // Credentials store
	AnthropicURL   string        // Optional override
	GeminiURL      string        // Optional override
	OpenAIURL      string        // Optional override
	ChatGPTAuthURL string        // Optional override
	ChatGPTAPIURL  string        // Optional override
}

// BuildTutor creates the configured Tutor, wrapping any network provider in FallbackTutor.
func BuildTutor(opts FactoryOptions) Tutor {
	offline := NewOfflineTutor()
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	provider := opts.Provider
	if provider == "" && opts.AuthStore != nil {
		provider = opts.AuthStore.GetConfig().ActiveProvider
	}
	if provider == "" {
		provider = ProviderOffline
	}

	model := opts.Model
	if model == "" && opts.AuthStore != nil {
		model = opts.AuthStore.ResolveModel(provider)
	}

	switch provider {
	case ProviderChatGPTPlan:
		chatgpt := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{
			AuthStore: opts.AuthStore,
			AuthURL:   opts.ChatGPTAuthURL,
			APIURL:    opts.ChatGPTAPIURL,
			Model:     model,
			Budget:    opts.Budget,
		})
		return NewFallbackTutor(chatgpt, offline, timeout)

	case ProviderAnthropic:
		anthropic := NewAnthropicTutor(APIKeyConfig{
			AuthStore: opts.AuthStore,
			BaseURL:   opts.AnthropicURL,
			Model:     model,
			Budget:    opts.Budget,
		})
		return NewFallbackTutor(anthropic, offline, timeout)

	case ProviderGoogle:
		gemini := NewGeminiTutor(APIKeyConfig{
			AuthStore: opts.AuthStore,
			BaseURL:   opts.GeminiURL,
			Model:     model,
			Budget:    opts.Budget,
		})
		return NewFallbackTutor(gemini, offline, timeout)

	case ProviderOpenAI:
		openai := NewOpenAIAPITutor(APIKeyConfig{
			AuthStore: opts.AuthStore,
			BaseURL:   opts.OpenAIURL,
			Model:     model,
			Budget:    opts.Budget,
		})
		return NewFallbackTutor(openai, offline, timeout)

	default: // ProviderOffline
		return offline
	}
}
