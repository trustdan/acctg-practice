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

func (a *AnthropicTutor) HintStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return a.stream(ctx, req, true, onUpdate)
}

func (a *AnthropicTutor) ExplainStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return a.stream(ctx, req, false, onUpdate)
}

// newRequest checks the budget and key and builds the POST /messages request.
func (a *AnthropicTutor) newRequest(ctx context.Context, req Request, isHint, stream bool) (*http.Request, error) {
	if a.cfg.Budget != nil {
		if err := a.cfg.Budget.Check(); err != nil {
			return nil, err
		}
	}

	key := a.getKey()
	if key == "" {
		return nil, errors.New("anthropic api key not configured (set ANTHROPIC_API_KEY or configure via 't')")
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
	if stream {
		payload["stream"] = true
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("x-api-key", key)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("content-type", "application/json")
	if stream {
		httpReq.Header.Set("accept", "text/event-stream")
	}
	return httpReq, nil
}

// stream sends a streamed Messages request. text_delta events are shown;
// thinking_delta events only mark the reply as thinking.
func (a *AnthropicTutor) stream(ctx context.Context, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
	httpReq, err := a.newRequest(ctx, req, isHint, true)
	if err != nil {
		return Response{}, err
	}

	httpResp, err := streamingClient(a.cfg.HTTPClient).Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("anthropic api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("anthropic api error (%d): %s", httpResp.StatusCode, httpErrorBody(httpResp))
	}

	reply := &streamReply{label: "anthropic", provider: "anthropic", budget: a.cfg.Budget, onUpdate: onUpdate}
	defer reply.recordBudget()
	finished := false
	err = readSSE(httpResp.Body, "anthropic", &finished, func(data string) (bool, error) {
		var event struct {
			Type  string `json:"type"`
			Delta struct {
				Type     string `json:"type"`
				Text     string `json:"text"`
				Thinking string `json:"thinking"`
			} `json:"delta"`
			Usage *struct {
				OutputTokens int `json:"output_tokens"`
			} `json:"usage"`
			Error struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return false, fmt.Errorf("failed decoding anthropic stream: %w", err)
		}
		switch event.Type {
		case "content_block_delta":
			switch event.Delta.Type {
			case "text_delta":
				if event.Delta.Text != "" {
					reply.add(event.Delta.Text, false)
				}
			case "thinking_delta":
				reply.add("", true)
			}
		case "message_delta":
			// message_delta usage is cumulative.
			if event.Usage != nil {
				reply.usage = event.Usage.OutputTokens
			}
		case "message_stop":
			finished = true
			return true, nil
		case "error":
			return false, fmt.Errorf("anthropic stream error: %s: %s", event.Error.Type, event.Error.Message)
		}
		return false, nil
	})
	if err != nil {
		return reply.partial(err)
	}
	return reply.result()
}

func (a *AnthropicTutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	httpReq, err := a.newRequest(ctx, req, isHint, false)
	if err != nil {
		return Response{}, err
	}

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

func (g *GeminiTutor) HintStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return g.stream(ctx, req, true, onUpdate)
}

func (g *GeminiTutor) ExplainStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return g.stream(ctx, req, false, onUpdate)
}

// newRequest checks the budget and key and builds a generateContent request,
// or a streamGenerateContent?alt=sse request when stream is set.
func (g *GeminiTutor) newRequest(ctx context.Context, req Request, isHint, stream bool) (*http.Request, error) {
	if g.cfg.Budget != nil {
		if err := g.cfg.Budget.Check(); err != nil {
			return nil, err
		}
	}

	key := g.getKey()
	if key == "" {
		return nil, errors.New("gemini api key not configured (set GEMINI_API_KEY or configure via 't')")
	}

	base := strings.TrimSuffix(g.cfg.BaseURL, "/")
	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s", base, g.cfg.Model, key)
	if stream {
		// The key goes in a header so it cannot appear in a transport error,
		// which may be shown as a fallback reason.
		endpoint = fmt.Sprintf("%s/models/%s:streamGenerateContent?alt=sse", base, g.cfg.Model)
	}

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
		return nil, fmt.Errorf("failed encoding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("content-type", "application/json")
	if stream {
		httpReq.Header.Set("x-goog-api-key", key)
		httpReq.Header.Set("accept", "text/event-stream")
	}
	return httpReq, nil
}

// stream sends a streamGenerateContent request. Each SSE chunk is a full
// GenerateContentResponse; parts marked thought only mark the reply as
// thinking, and a finishReason marks the stream complete.
func (g *GeminiTutor) stream(ctx context.Context, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
	httpReq, err := g.newRequest(ctx, req, isHint, true)
	if err != nil {
		return Response{}, err
	}

	httpResp, err := streamingClient(g.cfg.HTTPClient).Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("gemini api request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("gemini api error (%d): %s", httpResp.StatusCode, httpErrorBody(httpResp))
	}

	reply := &streamReply{label: "gemini", provider: "google", budget: g.cfg.Budget, onUpdate: onUpdate}
	defer reply.recordBudget()
	finished := false
	err = readSSE(httpResp.Body, "gemini", &finished, func(data string) (bool, error) {
		var chunk struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text    string `json:"text"`
						Thought bool   `json:"thought"`
					} `json:"parts"`
				} `json:"content"`
				FinishReason string `json:"finishReason"`
			} `json:"candidates"`
			UsageMetadata *struct {
				CandidatesTokenCount int `json:"candidatesTokenCount"`
			} `json:"usageMetadata"`
			Error *struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return false, fmt.Errorf("failed decoding gemini stream: %w", err)
		}
		if chunk.Error != nil {
			return false, fmt.Errorf("gemini stream error: %s: %s", chunk.Error.Status, chunk.Error.Message)
		}
		// candidatesTokenCount is cumulative across chunks.
		if chunk.UsageMetadata != nil && chunk.UsageMetadata.CandidatesTokenCount > 0 {
			reply.usage = chunk.UsageMetadata.CandidatesTokenCount
		}
		if len(chunk.Candidates) > 0 {
			cand := chunk.Candidates[0]
			for _, part := range cand.Content.Parts {
				if part.Thought {
					reply.add("", true)
				} else if part.Text != "" {
					reply.add(part.Text, false)
				}
			}
			if cand.FinishReason != "" {
				finished = true
			}
		}
		return false, nil
	})
	if err != nil {
		return reply.partial(err)
	}
	return reply.result()
}

func (g *GeminiTutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	httpReq, err := g.newRequest(ctx, req, isHint, false)
	if err != nil {
		return Response{}, err
	}

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
		// Deadlines come from the caller's context (FallbackTutor); a fixed client
		// timeout would also cut off a streamed reply that is still arriving.
		cfg.HTTPClient = &http.Client{}
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

func (o *OpenAIAPITutor) HintStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return o.stream(ctx, req, true, onUpdate)
}

func (o *OpenAIAPITutor) ExplainStream(ctx context.Context, req Request, onUpdate func(StreamUpdate)) (Response, error) {
	return o.stream(ctx, req, false, onUpdate)
}

func (o *OpenAIAPITutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	c, err := o.prepare()
	if err != nil {
		return Response{}, err
	}
	return callChatCompletions(ctx, c, req, isHint)
}

func (o *OpenAIAPITutor) stream(ctx context.Context, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
	c, err := o.prepare()
	if err != nil {
		return Response{}, err
	}
	return streamChatCompletions(ctx, c, req, isHint, onUpdate)
}

func (o *OpenAIAPITutor) prepare() (chatCompletionsCall, error) {
	if o.cfg.Budget != nil {
		if err := o.cfg.Budget.Check(); err != nil {
			return chatCompletionsCall{}, err
		}
	}

	key := o.getKey()
	if key == "" {
		return chatCompletionsCall{}, errors.New("openai api key not configured (set OPENAI_API_KEY or configure via 't')")
	}

	return chatCompletionsCall{
		BaseURL:    o.cfg.BaseURL,
		Token:      key,
		Model:      o.cfg.Model,
		MaxTokens:  800,
		HTTPClient: o.cfg.HTTPClient,
		Budget:     o.cfg.Budget,
		Label:      "openai api",
		Provider:   "openai",
	}, nil
}
