package tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// maxChatResponseBytes bounds how much of a chat completion body is read.
const maxChatResponseBytes = 1 << 20

// chatCompletionsCall describes one OpenAI-format POST /chat/completions request.
// It is shared by the OpenAI API and LM Studio adapters.
type chatCompletionsCall struct {
	BaseURL    string
	Token      string // Sent as a bearer token only when non-empty.
	Model      string
	MaxTokens  int
	HTTPClient *http.Client
	Budget     *Budget
	Label      string // Error prefix, e.g. "openai api" or "lm studio".
	Provider   string // Response.Provider value.
}

func callChatCompletions(ctx context.Context, c chatCompletionsCall, req Request, isHint bool) (Response, error) {
	endpoint := strings.TrimSuffix(c.BaseURL, "/") + "/chat/completions"

	payload := map[string]interface{}{
		"model":      c.Model,
		"max_tokens": c.MaxTokens,
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
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("%s request failed: %w", c.Label, err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, maxChatResponseBytes+1))
	if err != nil {
		return Response{}, fmt.Errorf("failed reading response: %w", err)
	}
	if len(bodyBytes) > maxChatResponseBytes {
		return Response{}, ErrResponseTooLarge
	}

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("%s error (%d): %s", c.Label, httpResp.StatusCode, string(bodyBytes))
	}

	// reasoning_content (sent by some local reasoning models) is deliberately not decoded.
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
		return Response{}, fmt.Errorf("failed decoding %s response: %w", c.Label, err)
	}

	if len(chatResp.Choices) == 0 {
		return Response{}, fmt.Errorf("empty choices received from %s", c.Label)
	}

	text := strings.TrimSpace(CleanLaTeXMath(StripThinkBlocks(chatResp.Choices[0].Message.Content)))
	if text == "" {
		return Response{}, fmt.Errorf("empty text response from %s (a reasoning model may have spent max_tokens thinking)", c.Label)
	}
	tokens := chatResp.Usage.CompletionTokens
	if tokens <= 0 {
		tokens = (len(text) + 3) / 4
	}
	if c.Budget != nil {
		_ = c.Budget.RecordUsage(tokens)
	}

	return Response{
		Text:        text,
		Provider:    c.Provider,
		TokensUsed:  tokens,
		Fallback:    false,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

var (
	thinkBlockRE = regexp.MustCompile(`(?is)<think>.*?</think>`)
	thinkOpenRE  = regexp.MustCompile(`(?is)<think>.*$`)
)

// StripThinkBlocks removes <think>…</think> reasoning blocks emitted by local
// reasoning models. An unterminated block (output truncated mid-thought) is
// dropped through the end of the text.
func StripThinkBlocks(s string) string {
	s = thinkBlockRE.ReplaceAllString(s, "")
	return thinkOpenRE.ReplaceAllString(s, "")
}
