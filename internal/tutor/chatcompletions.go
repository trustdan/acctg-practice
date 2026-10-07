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

// newRequest builds the POST /chat/completions request, optionally streamed.
func (c chatCompletionsCall) newRequest(ctx context.Context, req Request, isHint, stream bool) (*http.Request, error) {
	payload := map[string]interface{}{
		"model":      c.Model,
		"max_tokens": c.MaxTokens,
		"messages": []map[string]string{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": FormatUserPrompt(req, isHint)},
		},
	}
	if stream {
		payload["stream"] = true
		payload["stream_options"] = map[string]bool{"include_usage": true}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed encoding request: %w", err)
	}

	endpoint := strings.TrimSuffix(c.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func callChatCompletions(ctx context.Context, c chatCompletionsCall, req Request, isHint bool) (Response, error) {
	httpReq, err := c.newRequest(ctx, req, isHint, false)
	if err != nil {
		return Response{}, err
	}

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

// streamChatCompletions sends a streamed (SSE) chat completion and reports the
// visible text after every chunk. Reasoning (<think> blocks or
// reasoning_content deltas) is never shown; it only sets StreamUpdate.Thinking.
func streamChatCompletions(ctx context.Context, c chatCompletionsCall, req Request, isHint bool, onUpdate func(StreamUpdate)) (Response, error) {
	httpReq, err := c.newRequest(ctx, req, isHint, true)
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("%s request failed: %w", c.Label, err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("%s error (%d): %s", c.Label, httpResp.StatusCode, httpErrorBody(httpResp))
	}

	reply := &streamReply{label: c.Label, provider: c.Provider, budget: c.Budget, onUpdate: onUpdate}
	defer reply.recordBudget()
	finished := false
	err = readSSE(httpResp.Body, c.Label, &finished, func(data string) (bool, error) {
		if data == "[DONE]" {
			return true, nil
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content          string `json:"content"`
					ReasoningContent string `json:"reasoning_content"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				CompletionTokens int `json:"completion_tokens"`
			} `json:"usage"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return false, fmt.Errorf("failed decoding %s stream: %w", c.Label, err)
		}
		if chunk.Error != nil {
			return false, fmt.Errorf("%s stream error: %s", c.Label, chunk.Error.Message)
		}
		if chunk.Usage != nil {
			reply.usage = chunk.Usage.CompletionTokens
		}
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			if chunk.Choices[0].FinishReason != "" {
				finished = true
			}
			if delta.Content != "" || delta.ReasoningContent != "" {
				reply.add(delta.Content, delta.ReasoningContent != "")
			}
		}
		return false, nil
	})
	if err != nil {
		return reply.partial(err)
	}
	return reply.result()
}

// streamVisibleText removes reasoning from accumulated stream text, including a
// "<think>" tag still split across chunks, and cleans LaTeX for display.
func streamVisibleText(raw string) string {
	s := StripThinkBlocks(raw)
	lower := strings.ToLower(s)
	for i := len("<think>") - 1; i > 0; i-- {
		if strings.HasSuffix(lower, "<think>"[:i]) {
			s = s[:len(s)-i]
			break
		}
	}
	return strings.TrimSpace(CleanLaTeXMath(s))
}

// thinkBlockOpen reports whether raw ends inside an unterminated <think> block.
func thinkBlockOpen(raw string) bool {
	return thinkOpenRE.MatchString(thinkBlockRE.ReplaceAllString(raw, ""))
}
