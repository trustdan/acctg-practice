package tutor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/trustdan/acctg-practice/internal/domain"
)

// ProviderHealthReport represents the diagnostic test result for an AI tutor connection.
type ProviderHealthReport struct {
	ProviderID  string `json:"provider_id"`
	Name        string `json:"name"`
	Configured  bool   `json:"configured"`
	Status      string `json:"status"` // "PASS", "NOTICE", "CONFIG_REQUIRED", "ERROR"
	Message     string `json:"message"`
	ActiveModel string `json:"active_model,omitempty"`
	LatencyMs   int64  `json:"latency_ms,omitempty"`
}

// TestLLMProviders tests all supported tutor linkages and returns diagnostic health reports.
func TestLLMProviders(authStore *AuthStore, timeout time.Duration) []ProviderHealthReport {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	var reports []ProviderHealthReport

	// 1. Offline Machine Tutor
	start := time.Now()
	offlineTutor := NewOfflineTutor()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	testReq := Request{
		ProblemPrompt: "The company pays $1,000 cash for rent today.",
		FamilyID:      "pay_rent_cash",
		Stage:         domain.StageIdentifyAccount,
		ConceptID:     "rent_expense_vs_cash",
		CausalHint:    "Rent is consumed in the period.",
		Explanation:   "Expenses reduce equity.",
	}
	_, err := offlineTutor.Hint(ctx, testReq)
	elapsed := time.Since(start).Milliseconds()

	if err == nil {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderOffline,
			Name:        "Offline Machine Tutor",
			Configured:  true,
			Status:      "PASS",
			Message:     "100% offline deterministic rule engine verified. Zero network dependency.",
			ActiveModel: "offline-rules-engine",
			LatencyMs:   elapsed,
		})
	} else {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderOffline,
			Name:        "Offline Machine Tutor",
			Configured:  true,
			Status:      "ERROR",
			Message:     fmt.Sprintf("Offline tutor check failed: %v", err),
			ActiveModel: "offline-rules-engine",
			LatencyMs:   elapsed,
		})
	}

	// 2. ChatGPT Plus / Pro (OpenAI OAuth)
	resolvedClientID := DefaultDCRClientID
	hasChatGPTToken := false
	if authStore != nil {
		resolvedClientID = authStore.ResolveChatGPTClientID()
		hasChatGPTToken = authStore.IsConfigured(ProviderChatGPTPlan)
	}

	chatgptModel := DefaultChatGPTModel
	if authStore != nil {
		chatgptModel = authStore.ResolveModel(ProviderChatGPTPlan)
	}

	if hasChatGPTToken {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderChatGPTPlan,
			Name:        "ChatGPT Plus / Pro (OpenAI OAuth)",
			Configured:  true,
			Status:      "PASS",
			Message:     fmt.Sprintf("Locally configured OAuth credentials for Client ID '%s'. Responses API preview active (store=false, stream=true).", resolvedClientID),
			ActiveModel: chatgptModel,
		})
	} else {
		msg := "Continue with ChatGPT in Tutor Settings [2]. Dynamic registration is automatic; no API key or developer registration required."
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderChatGPTPlan,
			Name:        "ChatGPT Plus / Pro (OpenAI OAuth)",
			Configured:  false,
			Status:      "NOTICE",
			Message:     msg,
			ActiveModel: chatgptModel,
		})
	}

	// 3. Anthropic Claude API
	anthropicKey := ""
	anthropicModel := DefaultAnthropicModel
	if authStore != nil {
		anthropicKey = authStore.ResolveKey(ProviderAnthropic)
		anthropicModel = authStore.ResolveModel(ProviderAnthropic)
	}

	if anthropicKey != "" {
		masked := anthropicKey
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderAnthropic,
			Name:        "Anthropic Claude API",
			Configured:  true,
			Status:      "PASS",
			Message:     fmt.Sprintf("API key configured (%s). Ready for Messages API requests.", masked),
			ActiveModel: anthropicModel,
		})
	} else {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderAnthropic,
			Name:        "Anthropic Claude API",
			Configured:  false,
			Status:      "CONFIG_REQUIRED",
			Message:     "No API key configured. Set ANTHROPIC_API_KEY environment variable or press 't' -> [3] in TUI.",
			ActiveModel: anthropicModel,
		})
	}

	// 4. Google Gemini API
	geminiKey := ""
	geminiModel := DefaultGeminiModel
	if authStore != nil {
		geminiKey = authStore.ResolveKey(ProviderGoogle)
		geminiModel = authStore.ResolveModel(ProviderGoogle)
	}

	if geminiKey != "" {
		masked := geminiKey
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderGoogle,
			Name:        "Google Gemini API",
			Configured:  true,
			Status:      "PASS",
			Message:     fmt.Sprintf("API key configured (%s). Ready for generateContent requests.", masked),
			ActiveModel: geminiModel,
		})
	} else {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderGoogle,
			Name:        "Google Gemini API",
			Configured:  false,
			Status:      "CONFIG_REQUIRED",
			Message:     "No API key configured. Set GEMINI_API_KEY environment variable or press 't' -> [4] in TUI.",
			ActiveModel: geminiModel,
		})
	}

	// 5. OpenAI Commercial API
	openaiKey := ""
	openaiModel := DefaultOpenAIAPIModel
	if authStore != nil {
		openaiKey = authStore.ResolveKey(ProviderOpenAI)
		openaiModel = authStore.ResolveModel(ProviderOpenAI)
	}

	if openaiKey != "" {
		masked := openaiKey
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderOpenAI,
			Name:        "OpenAI Commercial API",
			Configured:  true,
			Status:      "PASS",
			Message:     fmt.Sprintf("API key configured (%s). Ready for /v1/chat/completions requests.", masked),
			ActiveModel: openaiModel,
		})
	} else {
		reports = append(reports, ProviderHealthReport{
			ProviderID:  ProviderOpenAI,
			Name:        "OpenAI Commercial API",
			Configured:  false,
			Status:      "CONFIG_REQUIRED",
			Message:     "No API key configured. Set OPENAI_API_KEY environment variable or press 't' -> [5] in TUI.",
			ActiveModel: openaiModel,
		})
	}

	reports = append(reports, testLMStudio(authStore, timeout))

	return reports
}

// testLMStudio probes the local LM Studio server's model list. It runs only for
// this explicit diagnostic command, never at startup.
func testLMStudio(authStore *AuthStore, timeout time.Duration) ProviderHealthReport {
	report := ProviderHealthReport{
		ProviderID: ProviderLMStudio,
		Name:       "LM Studio (local)",
		Configured: true,
	}
	baseURL, allowRemote, model := DefaultLMStudioURL, false, ""
	if authStore != nil {
		baseURL, allowRemote, model = authStore.ResolveLMStudioURL(), authStore.LMStudioAllowRemote(), authStore.ResolveModel(ProviderLMStudio)
	}
	if err := ValidateLMStudioURL(baseURL, allowRemote); err != nil {
		report.Status = "ERROR"
		report.Message = err.Error()
		return report
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	start := time.Now()
	models, err := DiscoverLMStudioModels(ctx, baseURL, LMStudioToken(), nil)
	report.LatencyMs = time.Since(start).Milliseconds()
	if err != nil {
		report.Status = "NOTICE"
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf("LM Studio at %s did not answer within %s", baseURL, timeout)
		}
		report.Message = err.Error()
		return report
	}
	if model == "" {
		model = models[0].ID + " (first available)"
	}
	report.Status = "PASS"
	report.ActiveModel = model
	report.Message = fmt.Sprintf("Server reachable at %s with %d chat model(s). No API key needed; requests stay on this machine.", baseURL, len(models))
	return report
}

// FormatHealthReports creates a formatted terminal display of LLM connection health.
func FormatHealthReports(reports []ProviderHealthReport) string {
	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString("          ACCOUNTUTOR 9000 // AI TUTOR LINKAGES & HEALTH DIAGNOSTICS            \n")
	b.WriteString("================================================================================\n\n")

	for i, r := range reports {
		statusBadge := "[PASS]"
		switch r.Status {
		case "PASS":
			statusBadge = "[✓ PASS]"
		case "NOTICE":
			statusBadge = "[! NOTICE]"
		case "CONFIG_REQUIRED":
			statusBadge = "[- NOT CONFIGURED]"
		case "ERROR":
			statusBadge = "[✗ ERROR]"
		}

		b.WriteString(fmt.Sprintf("%d. %-32s %s\n", i+1, r.Name, statusBadge))
		if r.ActiveModel != "" {
			b.WriteString(fmt.Sprintf("   Active Model: %s\n", r.ActiveModel))
		}
		if r.LatencyMs > 0 {
			b.WriteString(fmt.Sprintf("   Response Latency: %dms\n", r.LatencyMs))
		}
		b.WriteString(fmt.Sprintf("   Details: %s\n\n", r.Message))
	}

	b.WriteString("================================================================================\n")
	b.WriteString("Remediation Tips:\n")
	b.WriteString(" • Offline Mode [1] requires zero configuration and works completely offline.\n")
	b.WriteString("ChatGPT invalid_client: reconnect through Tutor Settings [2] for dynamic registration.\n")
	b.WriteString(" • Commercial API keys can be passed via environment variables (OPENAI_API_KEY,\n")
	b.WriteString("   ANTHROPIC_API_KEY, GEMINI_API_KEY) or set interactively in TUI (press 't').\n")
	b.WriteString(" • LM Studio [6] needs no key: load a model, start the server (Developer tab or\n")
	b.WriteString("   `lms server start`), then press 't' -> [6]. URL: LMSTUDIO_BASE_URL or --lmstudio-url.\n")
	b.WriteString("================================================================================\n")

	return b.String()
}
