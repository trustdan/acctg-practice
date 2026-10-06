package tui_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/trustdan/acctg-practice/internal/tui"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

func fakeLMStudioServer(t *testing.T, modelsJSON string) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(modelsJSON))
	}))
	t.Cleanup(srv.Close)
	return srv.URL + "/v1"
}

func setupLMStudioTUI(t *testing.T) *tui.Model {
	t.Helper()
	m, db := setupTestTUI(t)
	t.Cleanup(func() { db.Close() })
	m.ModelCache, _ = tutor.NewModelCache("") // Keep discovered models out of the real user cache.
	return m
}

func TestTutorConfigLMStudioConnects(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", fakeLMStudioServer(t, `{"data":[{"id":"qwen3-8b"},{"id":"text-embedding-nomic"}]}`))
	m := setupLMStudioTUI(t)

	sendKey(m, "t")
	if !strings.Contains(m.View(), "[6] LM Studio (local)") {
		t.Fatalf("expected [6] LM Studio option in tutor settings, got:\n%s", m.View())
	}
	sendKey(m, "6")

	if got := m.AuthStore.GetConfig().ActiveProvider; got != tutor.ProviderLMStudio {
		t.Fatalf("expected lmstudio active, got %s", got)
	}
	if !strings.Contains(m.Tutor.Name(), "LMStudioTutor") {
		t.Errorf("expected LM Studio tutor, got %s", m.Tutor.Name())
	}
	if !strings.Contains(m.TutorAuthNotice, "Connected to LM Studio") || !strings.Contains(m.TutorAuthNotice, "1 model(s)") || !strings.Contains(m.TutorAuthNotice, "qwen3-8b") {
		t.Errorf("unexpected notice: %q", m.TutorAuthNotice)
	}
	if models := m.ModelCache.GetModels(tutor.ProviderLMStudio); len(models) != 1 {
		t.Errorf("expected discovered model cached for [m], got %+v", models)
	}
}

func TestTutorConfigLMStudioUnreachableShowsSetupTips(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	closed := srv.URL + "/v1"
	srv.Close()
	t.Setenv("LMSTUDIO_BASE_URL", closed)
	m := setupLMStudioTUI(t)

	sendKey(m, "t")
	sendKey(m, "6")
	if !strings.Contains(m.TutorAuthNotice, "lms server start") || !strings.Contains(m.TutorAuthNotice, "Setup:") {
		t.Errorf("expected setup guidance, got %q", m.TutorAuthNotice)
	}
	if !strings.Contains(m.View(), "Developer tab") {
		t.Errorf("expected inline guidance in tutor settings view")
	}
}

func TestTutorConfigLMStudioNoModelLoaded(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", fakeLMStudioServer(t, `{"data":[]}`))
	m := setupLMStudioTUI(t)
	sendKey(m, "t")
	sendKey(m, "6")
	if !strings.Contains(m.TutorAuthNotice, "no model is loaded") {
		t.Errorf("expected no-model guidance, got %q", m.TutorAuthNotice)
	}
}

func TestTutorConfigLMStudioURLEditing(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", "")
	t.Setenv("LMSTUDIO_ALLOW_REMOTE", "")
	m := setupLMStudioTUI(t)
	sendKey(m, "t")

	sendKey(m, "l")
	if !m.TutorInputActive || m.TutorInputBuffer != tutor.DefaultLMStudioURL {
		t.Fatalf("expected URL input prefilled with default, got active=%v buffer=%q", m.TutorInputActive, m.TutorInputBuffer)
	}

	m.TutorInputBuffer = "http://192.168.1.50:1234/v1"
	sendSpecialKey(m, tea.KeyEnter)
	if !m.TutorInputActive || !strings.Contains(m.TutorAuthNotice, "not on this machine") {
		t.Fatalf("expected remote URL rejected, notice=%q", m.TutorAuthNotice)
	}
	if got := m.AuthStore.ResolveLMStudioURL(); got != tutor.DefaultLMStudioURL {
		t.Fatalf("rejected URL was saved: %s", got)
	}

	local := fakeLMStudioServer(t, `{"data":[{"id":"llama-3.2-3b"}]}`)
	m.TutorInputBuffer = local
	sendSpecialKey(m, tea.KeyEnter)
	if m.TutorInputActive || m.AuthStore.ResolveLMStudioURL() != local {
		t.Fatalf("expected local URL saved, got %s", m.AuthStore.ResolveLMStudioURL())
	}
	if m.AuthStore.GetConfig().ActiveProvider != tutor.ProviderLMStudio || !strings.Contains(m.TutorAuthNotice, "llama-3.2-3b") {
		t.Errorf("expected LM Studio activated and checked, notice=%q", m.TutorAuthNotice)
	}
}

func TestTutorConfigLMStudioStaleCheckIgnored(t *testing.T) {
	m := setupLMStudioTUI(t)
	sendKey(m, "t")
	sendKey(m, "1")
	notice := m.TutorAuthNotice
	newM, _ := m.Update(tui.ExportedLMStudioCheckMsg("http://localhost:1234/v1", []tutor.ModelInfo{{ID: "x"}}, nil))
	*m = *newM.(*tui.Model)
	if m.TutorAuthNotice != notice {
		t.Errorf("stale LM Studio check overwrote notice: %q", m.TutorAuthNotice)
	}
}

func TestHelpScreenAdvertisesLMStudio(t *testing.T) {
	m := setupLMStudioTUI(t)
	sendKey(m, "h")
	view := m.View()
	if !strings.Contains(view, "[t] then [6]") || !strings.Contains(view, "LM Studio") {
		t.Errorf("expected help screen to advertise LM Studio option, got:\n%s", view)
	}
}
