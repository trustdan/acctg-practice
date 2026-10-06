package tutor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// LM Studio local server defaults.
const (
	DefaultLMStudioURL       = "http://localhost:1234/v1"
	DefaultLMStudioMaxTokens = 4096             // Room for reasoning models to finish thinking.
	DefaultLMStudioTimeout   = 90 * time.Second // Covers cold model loads and CPU inference.
)

// ErrNoLMStudioModel means the server answered but reported no usable chat model.
var ErrNoLMStudioModel = errors.New("LM Studio is running but no model is loaded: load a model in LM Studio's Developer tab (or run `lms load <model>`)")

// LMStudioToken returns LM_API_TOKEN, used only when the learner has enabled
// authentication in LM Studio's server settings. It is never persisted.
func LMStudioToken() string {
	return strings.TrimSpace(os.Getenv("LM_API_TOKEN"))
}

// ValidateLMStudioURL accepts http(s) URLs on loopback hosts. Other hosts are
// rejected unless allowRemote is set, so learner context stays on the machine
// by default.
func ValidateLMStudioURL(raw string, allowRemote bool) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return fmt.Errorf("invalid LM Studio URL %q (expected e.g. %s)", raw, DefaultLMStudioURL)
	}
	if allowRemote || isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("LM Studio URL %q is not on this machine; only localhost, 127.0.0.1 and ::1 are allowed unless you opt in with --lmstudio-allow-remote or LMSTUDIO_ALLOW_REMOTE=1", raw)
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// LMStudioConfig configures the local LM Studio tutor.
type LMStudioConfig struct {
	AuthStore   *AuthStore
	BaseURL     string
	Model       string
	Token       string
	MaxTokens   int
	AllowRemote bool
	HTTPClient  *http.Client
	Budget      *Budget
}

// LMStudioTutor calls LM Studio's OpenAI-compatible /v1/chat/completions endpoint.
type LMStudioTutor struct {
	mu  sync.Mutex
	cfg LMStudioConfig
}

// NewLMStudioTutor constructs an LMStudioTutor. Construction makes no network call.
func NewLMStudioTutor(cfg LMStudioConfig) *LMStudioTutor {
	if cfg.AuthStore != nil {
		if cfg.BaseURL == "" {
			cfg.BaseURL = cfg.AuthStore.ResolveLMStudioURL()
		}
		if cfg.Model == "" {
			cfg.Model = cfg.AuthStore.ResolveModel(ProviderLMStudio)
		}
		if cfg.MaxTokens <= 0 {
			cfg.MaxTokens = cfg.AuthStore.ResolveLMStudioMaxTokens()
		}
		cfg.AllowRemote = cfg.AllowRemote || cfg.AuthStore.LMStudioAllowRemote()
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultLMStudioURL
	}
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = DefaultLMStudioMaxTokens
	}
	if cfg.Token == "" {
		cfg.Token = LMStudioToken()
	}
	if cfg.HTTPClient == nil {
		// Request deadlines come from the caller's context (FallbackTutor).
		cfg.HTTPClient = &http.Client{}
	}
	return &LMStudioTutor{cfg: cfg}
}

// SetModel updates the active LM Studio model.
func (l *LMStudioTutor) SetModel(model string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cfg.Model = model
}

// GetModel returns the configured model, or "" when the first loaded model is used.
func (l *LMStudioTutor) GetModel() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.cfg.Model
}

func (l *LMStudioTutor) Name() string {
	return "LMStudioTutor"
}

func (l *LMStudioTutor) Hint(ctx context.Context, req Request) (Response, error) {
	return l.call(ctx, req, true)
}

func (l *LMStudioTutor) Explain(ctx context.Context, req Request) (Response, error) {
	return l.call(ctx, req, false)
}

func (l *LMStudioTutor) call(ctx context.Context, req Request, isHint bool) (Response, error) {
	if err := ValidateLMStudioURL(l.cfg.BaseURL, l.cfg.AllowRemote); err != nil {
		return Response{}, err
	}
	if l.cfg.Budget != nil {
		if err := l.cfg.Budget.Check(); err != nil {
			return Response{}, err
		}
	}

	model := l.GetModel()
	if model == "" {
		models, err := DiscoverLMStudioModels(ctx, l.cfg.BaseURL, l.cfg.Token, l.cfg.HTTPClient)
		if err != nil {
			return Response{}, err
		}
		model = models[0].ID
		l.SetModel(model)
	}

	return callChatCompletions(ctx, chatCompletionsCall{
		BaseURL:    l.cfg.BaseURL,
		Token:      l.cfg.Token,
		Model:      model,
		MaxTokens:  l.cfg.MaxTokens,
		HTTPClient: l.cfg.HTTPClient,
		Budget:     l.cfg.Budget,
		Label:      "lm studio",
		Provider:   ProviderLMStudio,
	}, req, isHint)
}

// DiscoverLMStudioModels lists chat models from LM Studio's GET /v1/models.
// Embedding models are hidden. With just-in-time loading enabled, LM Studio may
// also list downloaded models that load on first use.
func DiscoverLMStudioModels(ctx context.Context, baseURL, token string, client *http.Client) ([]ModelInfo, error) {
	if baseURL == "" {
		baseURL = DefaultLMStudioURL
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(baseURL, "/")+"/models", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("LM Studio server not reachable at %s: start it in LM Studio's Developer tab or run `lms server start` (%w)", baseURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("LM Studio rejected the request (HTTP %d): authentication is enabled in its server settings; set LM_API_TOKEN", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LM Studio models error (HTTP %d)", resp.StatusCode)
	}

	var parsed struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed decoding LM Studio models json: %w", err)
	}

	var models []ModelInfo
	now := time.Now().UTC()
	for _, item := range parsed.Data {
		if item.ID == "" || strings.Contains(strings.ToLower(item.ID), "embed") {
			continue
		}
		models = append(models, ModelInfo{ID: item.ID, DisplayName: item.ID, Provider: ProviderLMStudio, DiscoveredAt: now})
	}
	if len(models) == 0 {
		return nil, ErrNoLMStudioModel
	}
	return models, nil
}
