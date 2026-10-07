package tutor

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// readSSE calls onData with each "data:" payload of a server-sent event
// stream until onData reports stop (a clean finish) or returns an error. When
// the body ends first, the stream counts as finished only if *finished is
// true; otherwise it was interrupted. At most maxChatResponseBytes are read.
func readSSE(body io.Reader, label string, finished *bool, onData func(data string) (stop bool, err error)) error {
	reader := bufio.NewReader(io.LimitReader(body, maxChatResponseBytes+1))
	read := 0
	for {
		line, readErr := reader.ReadString('\n')
		read += len(line)
		if read > maxChatResponseBytes {
			return ErrResponseTooLarge
		}
		if data, ok := strings.CutPrefix(strings.TrimSpace(line), "data:"); ok {
			stop, err := onData(strings.TrimSpace(data))
			if err != nil || stop {
				return err
			}
		}
		if readErr != nil {
			if readErr == io.EOF && *finished {
				return nil
			}
			if readErr == io.EOF {
				readErr = io.ErrUnexpectedEOF
			}
			return fmt.Errorf("%s stream interrupted: %w", label, readErr)
		}
	}
}

// streamReply accumulates one streamed reply, reports its visible text, and
// builds the final, partial or empty Response.
type streamReply struct {
	label    string
	provider string
	budget   *Budget
	onUpdate func(StreamUpdate)
	raw      strings.Builder
	usage    int // Output tokens reported by the provider, if any.
}

// add appends visible-channel text (which may still contain <think> blocks)
// and reports the update. thinking marks reasoning-only activity. Every call
// is reported, even when the visible text is unchanged, so a long reasoning
// phase still counts as activity for the idle timeout.
func (s *streamReply) add(text string, thinking bool) {
	s.raw.WriteString(text)
	if s.onUpdate == nil {
		return
	}
	s.onUpdate(StreamUpdate{
		Text:     streamVisibleText(s.raw.String()),
		Thinking: (thinking && text == "") || thinkBlockOpen(s.raw.String()),
		Provider: s.provider,
	})
}

// recordBudget records usage once; defer it so aborted streams count too.
func (s *streamReply) recordBudget() {
	if s.budget == nil {
		return
	}
	tokens := s.usage
	if tokens <= 0 {
		tokens = (s.raw.Len() + 3) / 4
	}
	_ = s.budget.RecordUsage(tokens)
}

// partial returns the visible text so far, marked Incomplete, with err.
func (s *streamReply) partial(err error) (Response, error) {
	return Response{Text: streamVisibleText(s.raw.String()), Provider: s.provider, Incomplete: true}, err
}

// result returns the completed reply, or an error when no visible text arrived.
func (s *streamReply) result() (Response, error) {
	text := streamVisibleText(s.raw.String())
	if text == "" {
		return Response{}, fmt.Errorf("empty text response from %s (a reasoning model may have spent max_tokens thinking)", s.label)
	}
	tokens := s.usage
	if tokens <= 0 {
		tokens = (len(text) + 3) / 4
	}
	return Response{
		Text:        text,
		Provider:    s.provider,
		TokensUsed:  tokens,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// streamingClient returns c without its total timeout: a streamed reply may
// legitimately run longer, and FallbackTutor bounds it with first-chunk and
// idle deadlines through the request context instead.
func streamingClient(c *http.Client) *http.Client {
	if c.Timeout == 0 {
		return c
	}
	copied := *c
	copied.Timeout = 0
	return &copied
}

// httpErrorBody reads a short error body from a failed streamed request.
func httpErrorBody(resp *http.Response) string {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return string(body)
}
