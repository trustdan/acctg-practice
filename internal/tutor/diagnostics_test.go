package tutor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/trustdan/acctg-practice/internal/tutor"
)

func TestLLMProvidersDiagnosticsOfflineOnly(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", closedServerURL())
	tempDir, err := os.MkdirTemp("", "diag_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")
	store, err := tutor.NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed creating auth store: %v", err)
	}

	reports := tutor.TestLLMProviders(store, 2*time.Second)
	if len(reports) != 6 {
		t.Fatalf("expected 6 provider reports, got %d", len(reports))
	}
	if lm := reports[5]; lm.ProviderID != tutor.ProviderLMStudio || lm.Status != "NOTICE" || !strings.Contains(lm.Message, "lms server start") {
		t.Errorf("expected LM Studio server-not-running notice, got %+v", lm)
	}

	// 1. Offline should PASS
	if reports[0].ProviderID != tutor.ProviderOffline || reports[0].Status != "PASS" {
		t.Errorf("expected offline tutor to pass, got %+v", reports[0])
	}

	// 2. ChatGPT plan should have client ID 'acctg-practice' and NOTICE
	if reports[1].ProviderID != tutor.ProviderChatGPTPlan {
		t.Errorf("expected chatgpt plan, got %+v", reports[1])
	}
	if !strings.Contains(reports[1].Message, "Dynamic registration") {
		t.Errorf("expected client ID acctg-practice in message, got: %s", reports[1].Message)
	}

	// 3. FormatHealthReports output
	formatted := tutor.FormatHealthReports(reports)
	if !strings.Contains(formatted, "ACCOUNTUTOR 9000 // AI TUTOR LINKAGES") {
		t.Errorf("expected header banner in formatted report, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "Offline Machine Tutor") {
		t.Errorf("expected Offline Machine Tutor in formatted report, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "invalid_client") {
		t.Errorf("expected invalid_client remediation guidance, got:\n%s", formatted)
	}
}

func TestLLMProvidersDiagnosticsWithConfiguredKey(t *testing.T) {
	t.Setenv("LMSTUDIO_BASE_URL", closedServerURL())
	tempDir, err := os.MkdirTemp("", "diag_key_test_*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	authPath := filepath.Join(tempDir, "tutor_auth.json")
	store, err := tutor.NewAuthStore(authPath)
	if err != nil {
		t.Fatalf("failed creating auth store: %v", err)
	}

	_ = store.SetAPIKey(tutor.ProviderAnthropic, "sk-ant-api03-testkey12345678")
	_ = store.SetChatGPTClientID("acctg-practice")

	reports := tutor.TestLLMProviders(store, 2*time.Second)

	// Anthropic should now be PASS
	var anthropicReport *tutor.ProviderHealthReport
	for i := range reports {
		if reports[i].ProviderID == tutor.ProviderAnthropic {
			anthropicReport = &reports[i]
			break
		}
	}

	if anthropicReport == nil {
		t.Fatalf("anthropic report not found")
	}
	if anthropicReport.Status != "PASS" {
		t.Errorf("expected Anthropic to PASS with configured key, got %s (%s)", anthropicReport.Status, anthropicReport.Message)
	}
	if !strings.Contains(anthropicReport.Message, "sk-a...5678") {
		t.Errorf("expected masked key in report, got: %s", anthropicReport.Message)
	}
}

func TestLLMProvidersDiagnosticsLMStudio(t *testing.T) {
	store, _ := tutor.NewAuthStore("")

	t.Setenv("LMSTUDIO_BASE_URL", startFake(t, &fakeLMStudio{models: []string{"qwen3-8b"}}))
	lm := tutor.TestLLMProviders(store, 2*time.Second)[5]
	if lm.Status != "PASS" || !strings.Contains(lm.ActiveModel, "qwen3-8b") || !strings.Contains(lm.Message, "No API key needed") {
		t.Errorf("expected LM Studio PASS, got %+v", lm)
	}

	t.Setenv("LMSTUDIO_BASE_URL", startFake(t, &fakeLMStudio{}))
	lm = tutor.TestLLMProviders(store, 2*time.Second)[5]
	if lm.Status != "NOTICE" || !strings.Contains(lm.Message, "no model is loaded") {
		t.Errorf("expected no-model notice, got %+v", lm)
	}

	t.Setenv("LMSTUDIO_BASE_URL", "http://10.1.2.3:1234/v1")
	lm = tutor.TestLLMProviders(store, 2*time.Second)[5]
	if lm.Status != "ERROR" || !strings.Contains(lm.Message, "not on this machine") {
		t.Errorf("expected remote-host rejection, got %+v", lm)
	}
}
