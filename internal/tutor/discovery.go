package tutor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// ModelInfo describes an available language model from a provider.
type ModelInfo struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"display_name"`
	Provider     string    `json:"provider"`
	Description  string    `json:"description,omitempty"`
	DiscoveredAt time.Time `json:"discovered_at,omitempty"`
}

// ProviderCatalog stores the discovered model list and fetch timestamp for a provider.
type ProviderCatalog struct {
	Provider    string      `json:"provider"`
	Models      []ModelInfo `json:"models"`
	LastFetched time.Time   `json:"last_fetched"`
}

// ModelCache provides thread-safe in-memory and persistent disk caching for discovered models.
type ModelCache struct {
	mu       sync.RWMutex
	filePath string
	catalogs map[string]ProviderCatalog
}

// DefaultModels returns vetted, canonical default models when offline or before remote discovery.
func DefaultModels(provider string) []ModelInfo {
	now := time.Now().UTC()
	switch provider {
	case ProviderAnthropic:
		return []ModelInfo{
			{ID: "claude-3-5-haiku-20241022", DisplayName: "Claude 3.5 Haiku (Fast & Efficient)", Provider: ProviderAnthropic, DiscoveredAt: now},
			{ID: "claude-3-5-sonnet-20241022", DisplayName: "Claude 3.5 Sonnet (Balanced & Highly Capable)", Provider: ProviderAnthropic, DiscoveredAt: now},
			{ID: "claude-3-7-sonnet-20250219", DisplayName: "Claude 3.7 Sonnet (Hybrid Reasoning)", Provider: ProviderAnthropic, DiscoveredAt: now},
			{ID: "claude-3-opus-20240229", DisplayName: "Claude 3 Opus (Comprehensive)", Provider: ProviderAnthropic, DiscoveredAt: now},
		}

	case ProviderGoogle:
		return []ModelInfo{
			{ID: "gemini-1.5-flash", DisplayName: "Gemini 1.5 Flash (Fast & Lightweight)", Provider: ProviderGoogle, DiscoveredAt: now},
			{ID: "gemini-1.5-pro", DisplayName: "Gemini 1.5 Pro (Advanced Reasoning)", Provider: ProviderGoogle, DiscoveredAt: now},
			{ID: "gemini-2.0-flash", DisplayName: "Gemini 2.0 Flash (Next-Gen Fast)", Provider: ProviderGoogle, DiscoveredAt: now},
			{ID: "gemini-2.0-flash-lite", DisplayName: "Gemini 2.0 Flash Lite (Cost-Efficient)", Provider: ProviderGoogle, DiscoveredAt: now},
		}

	case ProviderOpenAI:
		return []ModelInfo{
			{ID: "gpt-4o-mini", DisplayName: "GPT-4o Mini (Fast & Affordable)", Provider: ProviderOpenAI, DiscoveredAt: now},
			{ID: "gpt-4o", DisplayName: "GPT-4o (Flagship Multimodal)", Provider: ProviderOpenAI, DiscoveredAt: now},
			{ID: "o3-mini", DisplayName: "o3-mini (High-Speed STEM Reasoning)", Provider: ProviderOpenAI, DiscoveredAt: now},
			{ID: "o1-mini", DisplayName: "o1-mini (STEM Reasoning)", Provider: ProviderOpenAI, DiscoveredAt: now},
			{ID: "o1", DisplayName: "o1 (Deep Reasoning)", Provider: ProviderOpenAI, DiscoveredAt: now},
		}

	case ProviderChatGPTPlan:
		return []ModelInfo{
			{ID: "gpt-4o-mini", DisplayName: "GPT-4o Mini (ChatGPT Plus Default)", Provider: ProviderChatGPTPlan, DiscoveredAt: now},
			{ID: "gpt-4o", DisplayName: "GPT-4o (ChatGPT Plus Flagship)", Provider: ProviderChatGPTPlan, DiscoveredAt: now},
			{ID: "o3-mini", DisplayName: "o3-mini (ChatGPT Plus Reasoning)", Provider: ProviderChatGPTPlan, DiscoveredAt: now},
		}

	default:
		return []ModelInfo{
			{ID: "offline", DisplayName: "Offline Deterministic Engine", Provider: ProviderOffline, DiscoveredAt: now},
		}
	}
}

// DefaultModelCachePath returns standard cache file location in user data directory.
func DefaultModelCachePath() (string, error) {
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
	return filepath.Join(appDir, "models_cache.json"), nil
}

// NewModelCache initializes a ModelCache loading from disk if present.
func NewModelCache(filePath string) (*ModelCache, error) {
	cache := &ModelCache{
		filePath: filePath,
		catalogs: make(map[string]ProviderCatalog),
	}

	if filePath == "" {
		return cache, nil
	}

	data, err := os.ReadFile(filePath)
	if err == nil {
		var loaded map[string]ProviderCatalog
		if err := json.Unmarshal(data, &loaded); err == nil && loaded != nil {
			cache.catalogs = loaded
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed reading models cache: %w", err)
	}

	return cache, nil
}

// GetModels returns cached models for the specified provider, falling back to static defaults if unpopulated.
func (c *ModelCache) GetModels(provider string) []ModelInfo {
	if c == nil {
		return DefaultModels(provider)
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	cat, ok := c.catalogs[provider]
	if ok && len(cat.Models) > 0 {
		result := make([]ModelInfo, len(cat.Models))
		copy(result, cat.Models)
		return result
	}

	return DefaultModels(provider)
}

// GetLastFetched returns when the provider catalog was last refreshed, or zero time if never.
func (c *ModelCache) GetLastFetched(provider string) time.Time {
	if c == nil {
		return time.Time{}
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.catalogs[provider].LastFetched
}

// SetModels saves discovered models for a provider into memory and flushes to disk.
func (c *ModelCache) SetModels(provider string, models []ModelInfo) error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.catalogs[provider] = ProviderCatalog{
		Provider:    provider,
		Models:      models,
		LastFetched: time.Now().UTC(),
	}

	return c.saveLocked()
}

func (c *ModelCache) saveLocked() error {
	if c.filePath == "" {
		return nil
	}

	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed creating cache dir: %w", err)
	}

	data, err := json.MarshalIndent(c.catalogs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed encoding models cache: %w", err)
	}

	if err := os.WriteFile(c.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed writing models cache: %w", err)
	}

	return nil
}

// DiscoverOpenAIModels queries OpenAI GET /v1/models and filters for conversational and reasoning models.
func DiscoverOpenAIModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelInfo, error) {
	if apiKey == "" {
		return nil, errors.New("openai api key is required for model discovery")
	}
	if baseURL == "" {
		baseURL = DefaultOpenAIAPIURL
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	endpoint := strings.TrimSuffix(baseURL, "/") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai models request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading openai models response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai models error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Data []struct {
			ID      string `json:"id"`
			Created int64  `json:"created"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed decoding openai models json: %w", err)
	}

	var results []ModelInfo
	now := time.Now().UTC()

	// Filter for chat and reasoning models
	for _, item := range parsed.Data {
		id := strings.ToLower(item.ID)
		if !isRelevantOpenAIModel(id) {
			continue
		}

		displayName := formatOpenAIDisplayName(item.ID)
		results = append(results, ModelInfo{
			ID:           item.ID,
			DisplayName:  displayName,
			Provider:     ProviderOpenAI,
			DiscoveredAt: now,
		})
	}

	if len(results) == 0 {
		return DefaultModels(ProviderOpenAI), nil
	}

	// Sort prioritized models first, then alphabetical
	sort.SliceStable(results, func(i, j int) bool {
		pi := openAIModelPriority(results[i].ID)
		pj := openAIModelPriority(results[j].ID)
		if pi != pj {
			return pi < pj
		}
		return results[i].ID < results[j].ID
	})

	return results, nil
}

func isRelevantOpenAIModel(id string) bool {
	// Exclude non-text/non-chat models
	for _, prefix := range []string{"text-embedding", "whisper", "tts-", "dall-e", "babbage", "davinci", "omni-moderation", "text-moderation"} {
		if strings.HasPrefix(id, prefix) {
			return false
		}
	}
	for _, substring := range []string{"realtime", "audio", "transcription", "moderation", "search"} {
		if strings.Contains(id, substring) {
			return false
		}
	}

	// Include chat and reasoning models
	if strings.HasPrefix(id, "gpt-4") ||
		strings.HasPrefix(id, "gpt-3.5") ||
		strings.HasPrefix(id, "o1") ||
		strings.HasPrefix(id, "o3") ||
		strings.HasPrefix(id, "chatgpt-") {
		return true
	}

	return false
}

func openAIModelPriority(id string) int {
	switch {
	case id == "gpt-4o-mini":
		return 1
	case id == "gpt-4o":
		return 2
	case id == "o3-mini":
		return 3
	case id == "o1-mini":
		return 4
	case id == "o1":
		return 5
	case strings.HasPrefix(id, "gpt-4o"):
		return 6
	case strings.HasPrefix(id, "o3"):
		return 7
	case strings.HasPrefix(id, "o1"):
		return 8
	case strings.HasPrefix(id, "gpt-4"):
		return 9
	default:
		return 10
	}
}

func formatOpenAIDisplayName(id string) string {
	switch id {
	case "gpt-4o-mini":
		return "GPT-4o Mini (Fast & Affordable)"
	case "gpt-4o":
		return "GPT-4o (Flagship Multimodal)"
	case "o3-mini":
		return "o3-mini (High-Speed STEM Reasoning)"
	case "o1-mini":
		return "o1-mini (STEM Reasoning)"
	case "o1":
		return "o1 (Deep Reasoning)"
	case "gpt-4-turbo":
		return "GPT-4 Turbo"
	default:
		return id
	}
}

// DiscoverAnthropicModels queries Anthropic GET /v1/models and parses active models.
func DiscoverAnthropicModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelInfo, error) {
	if apiKey == "" {
		return nil, errors.New("anthropic api key is required for model discovery")
	}
	if baseURL == "" {
		baseURL = DefaultAnthropicURL
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	endpoint := strings.TrimSuffix(baseURL, "/") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic models request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading anthropic models response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic models error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Data []struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
			Type        string `json:"type"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed decoding anthropic models json: %w", err)
	}

	var results []ModelInfo
	now := time.Now().UTC()

	for _, item := range parsed.Data {
		displayName := item.DisplayName
		if displayName == "" {
			displayName = formatAnthropicDisplayName(item.ID)
		}
		results = append(results, ModelInfo{
			ID:           item.ID,
			DisplayName:  displayName,
			Provider:     ProviderAnthropic,
			DiscoveredAt: now,
		})
	}

	if len(results) == 0 {
		return DefaultModels(ProviderAnthropic), nil
	}

	// Sort prioritized models first
	sort.SliceStable(results, func(i, j int) bool {
		pi := anthropicModelPriority(results[i].ID)
		pj := anthropicModelPriority(results[j].ID)
		if pi != pj {
			return pi < pj
		}
		return results[i].ID < results[j].ID
	})

	return results, nil
}

func anthropicModelPriority(id string) int {
	switch {
	case strings.Contains(id, "haiku"):
		return 1
	case strings.Contains(id, "3-5-sonnet") || strings.Contains(id, "3.5-sonnet"):
		return 2
	case strings.Contains(id, "3-7-sonnet") || strings.Contains(id, "3.7-sonnet"):
		return 3
	case strings.Contains(id, "opus"):
		return 4
	default:
		return 5
	}
}

func formatAnthropicDisplayName(id string) string {
	switch {
	case strings.Contains(id, "haiku"):
		return fmt.Sprintf("Claude Haiku (%s)", id)
	case strings.Contains(id, "3-5-sonnet"):
		return fmt.Sprintf("Claude 3.5 Sonnet (%s)", id)
	case strings.Contains(id, "3-7-sonnet"):
		return fmt.Sprintf("Claude 3.7 Sonnet (%s)", id)
	case strings.Contains(id, "opus"):
		return fmt.Sprintf("Claude Opus (%s)", id)
	default:
		return id
	}
}

// DiscoverGeminiModels queries Google Gemini GET /v1beta/models and filters for generateContent support.
func DiscoverGeminiModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelInfo, error) {
	if apiKey == "" {
		return nil, errors.New("gemini api key is required for model discovery")
	}
	if baseURL == "" {
		baseURL = DefaultGeminiURL
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	endpoint := fmt.Sprintf("%s/models?key=%s", strings.TrimSuffix(baseURL, "/"), apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini models request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading gemini models response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini models error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Models []struct {
			Name                       string   `json:"name"`
			DisplayName                string   `json:"displayName"`
			Description                string   `json:"description"`
			SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed decoding gemini models json: %w", err)
	}

	var results []ModelInfo
	now := time.Now().UTC()

	for _, item := range parsed.Models {
		// Filter for generateContent support
		supportsGenerate := false
		for _, m := range item.SupportedGenerationMethods {
			if m == "generateContent" {
				supportsGenerate = true
				break
			}
		}
		if !supportsGenerate {
			continue
		}

		cleanID := strings.TrimPrefix(item.Name, "models/")
		displayName := item.DisplayName
		if displayName == "" {
			displayName = cleanID
		}

		results = append(results, ModelInfo{
			ID:           cleanID,
			DisplayName:  displayName,
			Provider:     ProviderGoogle,
			Description:  item.Description,
			DiscoveredAt: now,
		})
	}

	if len(results) == 0 {
		return DefaultModels(ProviderGoogle), nil
	}

	// Sort prioritized models first
	sort.SliceStable(results, func(i, j int) bool {
		pi := geminiModelPriority(results[i].ID)
		pj := geminiModelPriority(results[j].ID)
		if pi != pj {
			return pi < pj
		}
		return results[i].ID < results[j].ID
	})

	return results, nil
}

func geminiModelPriority(id string) int {
	switch {
	case id == "gemini-1.5-flash":
		return 1
	case id == "gemini-1.5-pro":
		return 2
	case id == "gemini-2.0-flash":
		return 3
	case id == "gemini-2.0-flash-lite":
		return 4
	case strings.HasPrefix(id, "gemini-1.5"):
		return 5
	case strings.HasPrefix(id, "gemini-2.0"):
		return 6
	default:
		return 7
	}
}

// DiscoverProviderModels fetches live models from the specified provider's API.
func DiscoverProviderModels(ctx context.Context, provider string, authStore *AuthStore, baseURL string, client *http.Client) ([]ModelInfo, error) {
	switch provider {
	case ProviderOffline:
		return DefaultModels(ProviderOffline), nil

	case ProviderAnthropic:
		key := ""
		if authStore != nil {
			key = authStore.ResolveKey(ProviderAnthropic)
		}
		if key == "" {
			return nil, errors.New("anthropic api key not configured (set ANTHROPIC_API_KEY or configure via 't')")
		}
		return DiscoverAnthropicModels(ctx, key, baseURL, client)

	case ProviderGoogle:
		key := ""
		if authStore != nil {
			key = authStore.ResolveKey(ProviderGoogle)
		}
		if key == "" {
			return nil, errors.New("gemini api key not configured (set GEMINI_API_KEY or configure via 't')")
		}
		return DiscoverGeminiModels(ctx, key, baseURL, client)

	case ProviderOpenAI:
		key := ""
		if authStore != nil {
			key = authStore.ResolveKey(ProviderOpenAI)
		}
		if key == "" {
			return nil, errors.New("openai api key not configured (set OPENAI_API_KEY or configure via 't')")
		}
		return DiscoverOpenAIModels(ctx, key, baseURL, client)

	case ProviderChatGPTPlan:
		if authStore == nil || !authStore.IsConfigured(ProviderChatGPTPlan) {
			return nil, errors.New("connect ChatGPT before discovering account models")
		}
		provider := NewOpenAIChatGPTPlanTutor(ChatGPTPlanConfig{AuthStore: authStore, HTTPClient: client})
		accessToken, err := provider.getValidAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		return DiscoverChatGPTModels(ctx, accessToken, baseURL, client)

	default:
		return nil, fmt.Errorf("unknown provider for model discovery: %s", provider)
	}
}

// RefreshProvider discovers live models from the provider API, updates the cache, and persists to disk.
func (c *ModelCache) RefreshProvider(ctx context.Context, provider string, authStore *AuthStore, baseURL string, client *http.Client) ([]ModelInfo, error) {
	models, err := DiscoverProviderModels(ctx, provider, authStore, baseURL, client)
	if err != nil {
		return nil, err
	}
	if err := c.SetModels(provider, models); err != nil {
		return models, fmt.Errorf("failed saving models to cache: %w", err)
	}
	return models, nil
}

// DiscoverChatGPTModels uses the account-specific SIWC schema, preserving server order.
func DiscoverChatGPTModels(ctx context.Context, token, baseURL string, client *http.Client) ([]ModelInfo, error) {
	if baseURL == "" {
		baseURL = DefaultOpenAIAPIURL
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(baseURL, "/")+"/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ChatGPT model discovery HTTP %d", resp.StatusCode)
	}
	var catalog struct {
		Models []struct {
			Slug        string `json:"slug"`
			DisplayName string `json:"display_name"`
			Visibility  string `json:"visibility"`
		} `json:"models"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&catalog); err != nil {
		return nil, err
	}
	var models []ModelInfo
	for _, m := range catalog.Models {
		if m.Visibility == "list" && m.Slug != "" {
			models = append(models, ModelInfo{ID: m.Slug, DisplayName: m.DisplayName, Provider: ProviderChatGPTPlan, DiscoveredAt: time.Now().UTC()})
		}
	}
	if len(models) == 0 {
		return nil, errors.New("ChatGPT returned no visible account models; the old default catalog is not an entitlement check")
	}
	return models, nil
}
