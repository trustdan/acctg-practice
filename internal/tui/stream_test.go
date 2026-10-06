package tui_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/trustdan/acctg-practice/internal/tui"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

func sseContent(text string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"choices": []map[string]interface{}{{"delta": map[string]string{"content": text}}},
	})
	return "data: " + string(b) + "\n\n"
}

// streamingTutorFor serves one SSE chunk, then either finishes, drops the
// connection, or holds it open until the client goes away.
func streamingTutorFor(t *testing.T, first, ending string) tutor.Tutor {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(sseContent(first)))
		w.(http.Flusher).Flush()
		switch ending {
		case "done":
			_, _ = w.Write([]byte("data: [DONE]\n\n"))
		case "hold":
			select {
			case <-r.Context().Done():
			case <-time.After(10 * time.Second):
			}
		}
	}))
	t.Cleanup(srv.Close)
	lm := tutor.NewLMStudioTutor(tutor.LMStudioConfig{BaseURL: srv.URL + "/v1", Model: "m"})
	return tutor.NewFallbackTutor(lm, nil, 5*time.Second)
}

// startStream presses key and returns the stream runner and first update wait.
func startStream(t *testing.T, m *tui.Model, key string) (run, wait tea.Cmd) {
	t.Helper()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	batch, ok := cmd().(tea.BatchMsg)
	if !ok || len(batch) != 3 {
		t.Fatalf("expected stream runner, update wait and loading tick, got %#v", cmd())
	}
	return batch[0], batch[1]
}

func TestStreamingReplyStaysResponsiveAndStopKeepsPartial(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = streamingTutorFor(t, "Cash is an asset; paying rent", "hold")

	run, wait := startStream(t, m, "e")
	done := make(chan tea.Msg, 1)
	go func() { done <- run() }()

	m.Update(wait())
	if m.TutorStreamText != "Cash is an asset; paying rent" || !m.TutorActive {
		t.Fatalf("expected streamed text while active, got %q active=%v", m.TutorStreamText, m.TutorActive)
	}
	view := m.View()
	if !strings.Contains(view, "streaming") || !strings.Contains(view, "Esc stops") {
		t.Errorf("streaming reply not shown:\n%s", view)
	}
	if !strings.Contains(view, "[Esc] Stop reply") {
		t.Error("status bar does not advertise the stop key while streaming")
	}

	// Keys are handled immediately while the stream is still open.
	start := time.Now()
	sendKey(m, "j")
	sendKey(m, "d")
	if time.Since(start) > time.Second || !m.TutorActive {
		t.Fatalf("key input blocked or ended the stream (%s)", time.Since(start))
	}

	sendSpecialKey(m, tea.KeyEsc)
	if m.TutorActive || !m.ShowHint || m.TutorResponse == nil || !m.TutorResponse.Incomplete {
		t.Fatalf("expected stopped partial reply kept visible, got active=%v resp=%+v", m.TutorActive, m.TutorResponse)
	}
	if m.PendingExplanation != nil {
		t.Error("partial reply was offered for saving")
	}
	if !strings.Contains(m.View(), "incomplete, not saved") {
		t.Errorf("partial reply not marked incomplete:\n%s", m.View())
	}

	select {
	case msg := <-done:
		m.Update(msg) // Late result of the cancelled request is ignored.
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not stop after Esc")
	}
	if m.PendingExplanation != nil || !m.TutorResponse.Incomplete {
		t.Error("cancelled stream result replaced the stopped partial reply")
	}
	sendSpecialKey(m, tea.KeyEsc)
	if m.ShowHint || m.ExplanationSavePrompt {
		t.Error("expected Esc to dismiss the partial reply without a save prompt")
	}
}

func TestStreamingMidStreamFailureIsNotSaved(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = streamingTutorFor(t, "Revenue records earning", "drop")

	run, _ := startStream(t, m, "e")
	m.Update(run())
	if m.TutorActive || m.TutorResponse == nil || !m.TutorResponse.Incomplete {
		t.Fatalf("expected incomplete reply, got %+v", m.TutorResponse)
	}
	if m.PendingExplanation != nil {
		t.Error("incomplete reply was offered for saving")
	}
	if !strings.Contains(m.TutorError, "not saved") {
		t.Errorf("expected not-saved notice, got %q", m.TutorError)
	}
}

func TestStreamingCompletedExplanationCanBeSaved(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	m.Tutor = streamingTutorFor(t, "Rent used up this month is an expense.", "done")

	run, _ := startStream(t, m, "e")
	m.Update(run())
	if m.TutorResponse == nil || m.TutorResponse.Incomplete || m.PendingExplanation == nil {
		t.Fatalf("expected a complete, saveable explanation, got %+v", m.TutorResponse)
	}
	if m.TutorStreamText != "" {
		t.Error("stream text not cleared after completion")
	}
}

func TestHelpScreenExplainsStreamingStop(t *testing.T) {
	m, db := setupTestTUI(t)
	defer db.Close()
	sendKey(m, "h")
	view := m.View()
	for _, want := range []string{"Stop a reply while it streams", "incomplete", "thinking…"} {
		if !strings.Contains(view, want) {
			t.Errorf("help screen missing %q", want)
		}
	}
}
