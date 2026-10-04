package tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Default models and API URLs
const (
	DefaultAnthropicURL   = "https://api.anthropic.com/v1"
	DefaultAnthropicModel = "claude-3-5-haiku-20241022"
	DefaultGeminiURL      = "https://generativelanguage.googleapis.com/v1beta"
	DefaultGeminiModel    = "gemini-1.5-flash"
	DefaultOpenAIAPIModel = "gpt-4o-mini"
)

// APIKeyConfig holds shared configuration for API key based tutors.
type APIKeyConfig struct {
	AuthStore  *AuthStore
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
	Budget     *Budget
}

// -------------------------------------------------------------------------
// Anthropic Tutor
// -------------------------------------------------------------------------

// AnthropicTutor calls the Anthropic Messages API using ANTHROPIC_API_KEY.
type AnthropicTutor struct {
	cfg APIKeyConfig
}

// NewAnthropicTutor constructs an AnthropicTutor.
func NewAnthropicTutor(cfg APIKeyConfig) *AnthropicTutor {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultAnthropicURL
	}
	if cfg.Model == "" && cfg.AuthStore != nil {
		cfg.Model = cfg.AuthStore.ResolveModel(ProviderAnthropic)
	}
	if cfg.Model == "" {
		cfg.Model = DefaultAnthropicModel
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &AnthropicTutor{cfg: cfg}
}

// SetModel updates the active Anthropic model.
func (a *AnthropicTutor) SetModel(model string) {
	a.cfg.Model = model
}

// GetModel returns the current Anthropic model.
func (a *AnthropicTutor) GetModel() string {
	return a.cfg.Model
}

func (a *AnthropicTutor) Name() string {
	return "AnthropicTutor"
}

func (a *AnthropicTutor) Hint(ctx context.Context, req Request) (Response, error) {
	return a.call(ctx, req, true)
}

func (a *AnthropicTutor) Explain(ctx context.Context, req Request) (Response, error) {
	return a.call(ctx, req, false)
}

func (a *AnthropicTutor) getKey() string {
	if a.cfg.APIKey != "" {
		return a.cfg.APIKey
	}
	if a.cfg.AuthStore != nil {
		return a.cfg.AuthStore.ResolveKey(ProviderAnthropic)
	}
	return ""
}

func (a *AnthropicTutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	if a.cfg.Budget != nil {
		if err := a.cfg.Budget.Check(); err != nil {
			return Response{}, err
		}
	}

	key := a.getKey()
	if key == "" {
		return Response{}, errors.New("anthropic api key not configured (set ANTHROPIC_API_KEY or configure via 't')")
	}

	endpoint := strings.TrimSuffix(a.cfg.BaseURL, "/") + "/messages"

	payload := map[string]interface{}{
		"model":      a.cfg.Model,
		"max_tokens": 1024,
		"system":     SystemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": FormatUserPrompt(req, isHint)},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("x-api-key", key)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("content-type", "application/json")

	httpResp, err := a.cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("anthropic api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("anthropic api error (%d): %s", httpResp.StatusCode, string(bodyBytes))
	}

	var anthropicResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(bodyBytes, &anthropicResp); err != nil {
		return Response{}, fmt.Errorf("failed decoding anthropic response: %w", err)
	}

	var sb strings.Builder
	for _, c := range anthropicResp.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	text := strings.TrimSpace(CleanLaTeXMath(sb.String()))
	if text == "" {
		return Response{}, errors.New("empty text response from anthropic")
	}

	tokens := anthropicResp.Usage.OutputTokens
	if tokens <= 0 {
		tokens = (len(text) + 3) / 4
	}
	if a.cfg.Budget != nil {
		_ = a.cfg.Budget.RecordUsage(tokens)
	}

	return Response{
		Text:        text,
		Provider:    "anthropic",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// -------------------------------------------------------------------------
// Google Gemini Tutor
// -------------------------------------------------------------------------

// GeminiTutor calls the Google Gemini generateContent API using GEMINI_API_KEY.
type GeminiTutor struct {
	cfg APIKeyConfig
}

// NewGeminiTutor constructs a GeminiTutor.
func NewGeminiTutor(cfg APIKeyConfig) *GeminiTutor {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultGeminiURL
	}
	if cfg.Model == "" && cfg.AuthStore != nil {
		cfg.Model = cfg.AuthStore.ResolveModel(ProviderGoogle)
	}
	if cfg.Model == "" {
		cfg.Model = DefaultGeminiModel
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &GeminiTutor{cfg: cfg}
}

// SetModel updates the active Gemini model.
func (g *GeminiTutor) SetModel(model string) {
	g.cfg.Model = model
}

// GetModel returns the current Gemini model.
func (g *GeminiTutor) GetModel() string {
	return g.cfg.Model
}

func (g *GeminiTutor) Name() string {
	return "GeminiTutor"
}

func (g *GeminiTutor) Hint(ctx context.Context, req Request) (Response, error) {
	return g.call(ctx, req, true)
}

func (g *GeminiTutor) Explain(ctx context.Context, req Request) (Response, error) {
	return g.call(ctx, req, false)
}

func (g *GeminiTutor) getKey() string {
	if g.cfg.APIKey != "" {
		return g.cfg.APIKey
	}
	if g.cfg.AuthStore != nil {
		return g.cfg.AuthStore.ResolveKey(ProviderGoogle)
	}
	return ""
}

func (g *GeminiTutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	if g.cfg.Budget != nil {
		if err := g.cfg.Budget.Check(); err != nil {
			return Response{}, err
		}
	}

	key := g.getKey()
	if key == "" {
		return Response{}, errors.New("gemini api key not configured (set GEMINI_API_KEY or configure via 't')")
	}

	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s",
		strings.TrimSuffix(g.cfg.BaseURL, "/"), g.cfg.Model, key)

	payload := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{
				{"text": SystemPrompt},
			},
		},
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": FormatUserPrompt(req, isHint)},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("content-type", "application/json")

	httpResp, err := g.cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("gemini api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("gemini api error (%d): %s", httpResp.StatusCode, string(bodyBytes))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return Response{}, fmt.Errorf("failed decoding gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return Response{}, errors.New("empty response received from gemini")
	}

	text := strings.TrimSpace(CleanLaTeXMath(geminiResp.Candidates[0].Content.Parts[0].Text))
	tokens := geminiResp.UsageMetadata.CandidatesTokenCount
	if tokens <= 0 {
		tokens = (len(text) + 3) / 4
	}
	if g.cfg.Budget != nil {
		_ = g.cfg.Budget.RecordUsage(tokens)
	}

	return Response{
		Text:        text,
		Provider:    "google",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// -------------------------------------------------------------------------
// OpenAI API Tutor (Commercial Key)
// -------------------------------------------------------------------------

// OpenAIAPITutor calls OpenAI Chat Completions API using OPENAI_API_KEY.
type OpenAIAPITutor struct {
	cfg APIKeyConfig
}

// NewOpenAIAPITutor constructs an OpenAIAPITutor.
func NewOpenAIAPITutor(cfg APIKeyConfig) *OpenAIAPITutor {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultOpenAIAPIURL
	}
	if cfg.Model == "" && cfg.AuthStore != nil {
		cfg.Model = cfg.AuthStore.ResolveModel(ProviderOpenAI)
	}
	if cfg.Model == "" {
		cfg.Model = DefaultOpenAIAPIModel
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &OpenAIAPITutor{cfg: cfg}
}

// SetModel updates the active OpenAI model.
func (o *OpenAIAPITutor) SetModel(model string) {
	o.cfg.Model = model
}

// GetModel returns the current OpenAI model.
func (o *OpenAIAPITutor) GetModel() string {
	return o.cfg.Model
}

func (o *OpenAIAPITutor) Name() string {
	return "OpenAIAPITutor"
}

func (o *OpenAIAPITutor) Hint(ctx context.Context, req Request) (Response, error) {
	return o.call(ctx, req, true)
}

func (o *OpenAIAPITutor) Explain(ctx context.Context, req Request) (Response, error) {
	return o.call(ctx, req, false)
}

func (o *OpenAIAPITutor) getKey() string {
	if o.cfg.APIKey != "" {
		return o.cfg.APIKey
	}
	if o.cfg.AuthStore != nil {
		return o.cfg.AuthStore.ResolveKey(ProviderOpenAI)
	}
	return ""
}

func (o *OpenAIAPITutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	if o.cfg.Budget != nil {
		if err := o.cfg.Budget.Check(); err != nil {
			return Response{}, err
		}
	}

	key := o.getKey()
	if key == "" {
		return Response{}, errors.New("openai api key not configured (set OPENAI_API_KEY or configure via 't')")
	}

	endpoint := strings.TrimSuffix(o.cfg.BaseURL, "/") + "/chat/completions"

	payload := map[string]interface{}{
		"model":      o.cfg.Model,
		"max_tokens": 800,
		"messages": []map[string]string{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": FormatUserPrompt(req, isHint)},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+key)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := o.cfg.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("openai api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("openai api error (%d): %s", httpResp.StatusCode, string(bodyBytes))
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return Response{}, fmt.Errorf("failed decoding openai response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return Response{}, errors.New("empty choices received from openai api")
	}

	text := strings.TrimSpace(CleanLaTeXMath(chatResp.Choices[0].Message.Content))
	tokens := chatResp.Usage.CompletionTokens
	if tokens <= 0 {
		tokens = (len(text) + 3) / 4
	}
	if o.cfg.Budget != nil {
		_ = o.cfg.Budget.RecordUsage(tokens)
	}

	return Response{
		Text:        text,
		Provider:    "openai",
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}
