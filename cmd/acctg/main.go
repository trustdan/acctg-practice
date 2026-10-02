package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"encoding/json"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/candidate"
	"github.com/trustdan/acctg-practice/internal/domain"
	"github.com/trustdan/acctg-practice/internal/drill"
	"github.com/trustdan/acctg-practice/internal/engine"
	"github.com/trustdan/acctg-practice/internal/exam"
	"github.com/trustdan/acctg-practice/internal/mastery"
	"github.com/trustdan/acctg-practice/internal/statements"
	"github.com/trustdan/acctg-practice/internal/storage"
	"github.com/trustdan/acctg-practice/internal/tui"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

const AppVersion = "0.24.0"

func main() {
	var (
		questionsFlag         = flag.Int("questions", 10, "number of questions in practice session")
		sizeFlag              = flag.Int("size", 0, "session size (5, 10, 15, 20 questions)")
		intensityFlag         = flag.String("intensity", "standard", "practice intensity: standard, spaced, intensive, transfer")
		seedFlag              = flag.Int64("seed", 0, "deterministic random seed (default: current time)")
		providerFlag          = flag.String("provider", "", "alias for --tutor")
		tutorFlag             = flag.String("tutor", "offline", "tutor provider: offline, chatgpt-plan, anthropic, google, openai, simulated, simulated-slow, simulated-error")
		tutorModelFlag        = flag.String("tutor-model", "", "select active AI tutor model (e.g. gpt-4o, claude-3-5-sonnet-20241022, gemini-1.5-flash)")
		tutorTimeoutFlag      = flag.Duration("tutor-timeout", 60*time.Second, "tutor request timeout")
		tutorBudgetFlag       = flag.Int("tutor-budget", 20, "max tutor requests per session (0 = unlimited)")
		listModelsFlag        = flag.Bool("list-models", false, "display discovered and default AI models for configured providers and exit")
		fetchModelsFlag       = flag.Bool("fetch-models", false, "discover live models from configured provider APIs and save to local cache")
		refreshModelsFlag     = flag.Bool("refresh-models", false, "alias for --fetch-models")
		oauthClientIDFlag     = flag.String("oauth-client-id", "", "set custom OpenAI OAuth client ID for ChatGPT Plus plan")
		dataDirFlag           = flag.String("data-dir", "", "override user data directory for storage")
		dbPathFlag            = flag.String("db", "", "override SQLite database file path directly")
		exportJSON            = flag.Bool("export-json", false, "export all practice sessions and attempts as JSON and exit")
		masteryFlag           = flag.Bool("mastery", false, "print current learner mastery projections and exit")
		generateCandidateFlag = flag.Bool("generate-candidate", false, "generate a creative question candidate outside active bank and exit")
		candidateFamilyFlag   = flag.String("candidate-family", "", "family ID for candidate generation")
		weakestConceptFlag    = flag.Bool("weakest-concept", false, "target learner's weakest concept for candidate generation")
		candidatesFlag        = flag.Bool("candidates", false, "list all stored candidate questions and exit")
		previewCandidateFlag  = flag.String("preview-candidate", "", "display full preview for candidate question with given ID and exit")
		candidateFailedFlag   = flag.Bool("candidate-failed", false, "filter for failed proposals when listing candidates")
		approveCandidateFlag  = flag.String("approve-candidate", "", "approve a candidate question by ID into active bank")
		rejectCandidateFlag   = flag.String("reject-candidate", "", "reject a candidate question by ID")
		repairCandidateFlag   = flag.String("repair-candidate", "", "repair a candidate question by ID")
		retireQuestionFlag    = flag.String("retire-question", "", "retire an active question by ID from practice")
		reviewerFlag          = flag.String("reviewer", "", "reviewer identity for approval/rejection/repair/retirement")
		notesFlag             = flag.String("notes", "", "notes or reason for review action")
		scenarioTemplateFlag  = flag.String("scenario-template", "", "new scenario template text for repair-candidate")
		approvalHistoryFlag   = flag.Bool("approval-history", false, "display approval events audit history and exit")
		activeBankFlag        = flag.Bool("active-bank", false, "list all active question bank items (seed + published) and exit")
		exportBankFlag        = flag.Bool("export-bank", false, "export active question bank as JSON to stdout and exit")
		practiceJournalFlag   = flag.Bool("practice-journal", false, "start interactive multi-line journal entry practice mode directly")
		reconcileAllFlag      = flag.Bool("reconcile-all", false, "verify and print postings reconciliation across all active bank questions and exit")
		statementsFlag        = flag.Bool("statements", false, "display comprehensive case financial statements, trial balances, and audit proof and exit")
		caseStudyFlag         = flag.Bool("case-study", false, "alias for --statements")
		examFlag              = flag.Bool("exam", false, "start interactive exam mode (assessment conditions: hints, references, and feedback suppressed)")
		examTimeFlag          = flag.Duration("exam-time", 0, "exam time limit (e.g. 10m, 15m, 600s; default: 0 = untimed)")
		timerFlag             = flag.Duration("timer", 0, "alias for --exam-time")
		resumeExamFlag        = flag.Bool("resume-exam", false, "resume latest interrupted exam without prompting")
		abandonExamFlag       = flag.String("abandon-exam", "", "abandon/close an interrupted exam by ID and exit")
		closeExamFlag         = flag.String("close-exam", "", "alias for --abandon-exam")
		examHistoryFlag       = flag.Bool("exam-history", false, "display all past exam sessions and scores and exit")
		examReportFlag        = flag.String("exam-report", "", "display full report for a completed or interrupted exam by ID and exit")
		skipIntroFlag         = flag.Bool("skip-intro", false, "skip startup spaceship animation and enter practice directly")
		noIntroFlag           = flag.Bool("no-intro", false, "alias for --skip-intro")
		introFlag             = flag.Bool("intro", false, "launch startup spaceship animation showcase directly")
		highScoresFlag        = flag.Bool("high-scores", false, "display top arcade flight scores from database and exit")
		arcadeScoresFlag      = flag.Bool("arcade-scores", false, "alias for --high-scores")
		testLLMFlag           = flag.Bool("test-llm", false, "test and diagnose connections across all configured AI tutor providers and exit")
		testProvidersFlag     = flag.Bool("test-providers", false, "alias for --test-llm")
		versionFlag           = flag.Bool("version", false, "display version and exit")
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("acctg version %s (offline keyboard drill with arcade combat upgrades, smart bombs, dual audit HP, and high score hall of fame)\n", AppVersion)
		os.Exit(0)
	}

	if *providerFlag != "" && *tutorFlag == "offline" {
		*tutorFlag = *providerFlag
	}

	totalQuestions := *questionsFlag
	if *sizeFlag > 0 {
		totalQuestions = *sizeFlag
	}

	// Determine user data directory
	var dataDir string
	if *dataDirFlag != "" {
		dataDir = *dataDirFlag
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to create data directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		defaultDir, err := storage.DefaultDataDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not determine default data directory: %v\n", err)
			os.Exit(1)
		}
		dataDir = defaultDir
	}

	// Determine SQLite DB path
	dbPath := *dbPathFlag
	if dbPath == "" {
		dbPath = filepath.Join(dataDir, "acctg_practice.db")
	}

	// Load tutor credentials store (0600 permissions in data directory)
	authPath := filepath.Join(dataDir, "tutor_auth.json")
	authStore, err := tutor.NewAuthStore(authPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load tutor auth store: %v\n", err)
	}

	// Load model cache (0600 permissions in data directory)
	cachePath := filepath.Join(dataDir, "models_cache.json")
	modelCache, err := tutor.NewModelCache(cachePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load model cache: %v\n", err)
	}

	// Handle --provider standalone switch
	if *providerFlag != "" && authStore != nil {
		_ = authStore.SetActiveProvider(*providerFlag)
		if *tutorModelFlag == "" && !*examFlag && !*practiceJournalFlag && !*statementsFlag && !*listModelsFlag && !*fetchModelsFlag && !*refreshModelsFlag {
			fmt.Printf("✓ Active tutor provider set to: %s\n", strings.ToUpper(*providerFlag))
			os.Exit(0)
		}
	}

	if *oauthClientIDFlag != "" {
		fmt.Fprintln(os.Stderr, "--oauth-client-id is obsolete. Continue with ChatGPT in Tutor Settings [2]; registration is automatic.")
	}

	// Handle --test-llm / --test-providers
	if *testLLMFlag || *testProvidersFlag {
		reports := tutor.TestLLMProviders(authStore, *tutorTimeoutFlag)
		fmt.Print(tutor.FormatHealthReports(reports))
		os.Exit(0)
	}

	// Handle --tutor-model
	if *tutorModelFlag != "" && authStore != nil {
		active := authStore.GetConfig().ActiveProvider
		targetProvider := active
		if *tutorFlag != "offline" {
			targetProvider = *tutorFlag
		}
		if targetProvider == tutor.ProviderOffline {
			// Auto-detect provider from model ID prefix
			mLower := strings.ToLower(*tutorModelFlag)
			switch {
			case strings.HasPrefix(mLower, "claude"):
				targetProvider = tutor.ProviderAnthropic
			case strings.HasPrefix(mLower, "gemini"):
				targetProvider = tutor.ProviderGoogle
			case strings.HasPrefix(mLower, "gpt-") || strings.HasPrefix(mLower, "o1") || strings.HasPrefix(mLower, "o3") || strings.HasPrefix(mLower, "chatgpt"):
				if authStore.IsConfigured(tutor.ProviderChatGPTPlan) {
					targetProvider = tutor.ProviderChatGPTPlan
				} else {
					targetProvider = tutor.ProviderOpenAI
				}
			}
		}
		if targetProvider != tutor.ProviderOffline {
			if err := authStore.SetSelectedModel(targetProvider, *tutorModelFlag); err != nil {
				fmt.Fprintf(os.Stderr, "error setting active model: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Active model for provider %s set to: %s\n", strings.ToUpper(targetProvider), *tutorModelFlag)
			if *tutorFlag != "offline" {
				_ = authStore.SetActiveProvider(targetProvider)
			}
		} else {
			fmt.Println("Notice: Active provider is Offline mode. Model selection will apply when an AI provider is active.")
		}
		if *tutorFlag == "offline" && !*examFlag && !*practiceJournalFlag && !*statementsFlag && !*listModelsFlag && !*fetchModelsFlag && !*refreshModelsFlag {
			os.Exit(0)
		}
	}

	// Handle --fetch-models / --refresh-models
	if *fetchModelsFlag || *refreshModelsFlag {
		fmt.Println("==================================================================================")
		fmt.Println("DISCOVERING LIVE MODELS ACROSS CONFIGURED AI PROVIDERS")
		fmt.Println("==================================================================================")
		providers := []string{tutor.ProviderAnthropic, tutor.ProviderGoogle, tutor.ProviderOpenAI, tutor.ProviderChatGPTPlan}
		discoveredCount := 0
		for _, p := range providers {
			isConfigured := authStore != nil && authStore.IsConfigured(p)
			if !isConfigured {
				fmt.Printf("[%s] Skipped (credentials not configured)\n", strings.ToUpper(p))
				continue
			}
			fmt.Printf("[%s] Fetching models from API...\n", strings.ToUpper(p))
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			models, err := modelCache.RefreshProvider(ctx, p, authStore, "", nil)
			cancel()
			if err != nil {
				fmt.Printf("  ⚠️ Discovery failed: %v (using fallback catalog)\n", err)
			} else {
				fmt.Printf("  ✓ Successfully discovered and cached %d models.\n", len(models))
				discoveredCount += len(models)
			}
		}
		fmt.Printf("\nDiscovery complete (%d models). Cache saved to: %s\n", discoveredCount, cachePath)
		fmt.Println("==================================================================================")
		if !*listModelsFlag {
			os.Exit(0)
		}
	}

	// Handle --list-models
	if *listModelsFlag {
		activeProvider := tutor.ProviderOffline
		if authStore != nil {
			activeProvider = authStore.GetConfig().ActiveProvider
		}
		fmt.Println("==================================================================================")
		fmt.Printf("TUTOR MODEL CATALOG & SELECTION STATUS (Active Provider: %s)\n", strings.ToUpper(activeProvider))
		fmt.Println("==================================================================================")

		providers := []struct {
			id    string
			label string
		}{
			{tutor.ProviderAnthropic, "Anthropic Claude (Messages API)"},
			{tutor.ProviderGoogle, "Google Gemini (Gemini API)"},
			{tutor.ProviderOpenAI, "OpenAI API (Chat Completions)"},
			{tutor.ProviderChatGPTPlan, "ChatGPT Plus Plan (OpenAI OAuth)"},
		}

		for _, p := range providers {
			isConfigured := authStore != nil && authStore.IsConfigured(p.id)
			cfgBadge := "[Not Configured]"
			if isConfigured {
				cfgBadge = "[Ready]"
			}
			if p.id == activeProvider {
				cfgBadge += " [ACTIVE PROVIDER]"
			}

			activeModel := ""
			if authStore != nil {
				activeModel = authStore.ResolveModel(p.id)
			}

			fmt.Printf("\nPROVIDER: %s %s\n", strings.ToUpper(p.id), cfgBadge)
			fmt.Printf("Description: %s\n", p.label)
			if p.id == tutor.ProviderChatGPTPlan && authStore != nil {
				fmt.Printf("OAuth Client ID: %s\n", authStore.ResolveChatGPTClientID())
			}
			lastFetched := modelCache.GetLastFetched(p.id)
			if !lastFetched.IsZero() {
				fmt.Printf("Last Fetched from API: %s\n", lastFetched.Format(time.RFC3339))
			} else {
				fmt.Println("Last Fetched from API: Never (using standard defaults)")
			}
			fmt.Println("----------------------------------------------------------------------------------")
			fmt.Printf("  %-32s %-40s %s\n", "MODEL ID", "DISPLAY NAME", "STATUS")
			fmt.Println("  --------------------------------------------------------------------------------")

			models := modelCache.GetModels(p.id)
			for _, m := range models {
				status := ""
				if m.ID == activeModel {
					status = "★ [ACTIVE MODEL]"
				}
				fmt.Printf("  %-32s %-40s %s\n", m.ID, m.DisplayName, status)
			}
		}
		fmt.Println("\nTo change models: use --tutor-model=<id> or open TUI Tutor Settings ('t') and press 'm'.")
		fmt.Println("To discover live models: use --fetch-models or press 'r' in TUI Tutor Settings.")
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Open durable SQLite storage
	db, err := storage.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open database at %s: %v\n", dbPath, err)
		os.Exit(1)
	}
	defer db.Close()

	// Handle --high-scores / --arcade-scores
	if *highScoresFlag || *arcadeScoresFlag {
		scores, err := db.ListTopArcadeHighScores(10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving arcade high scores: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("==================================================================================")
		fmt.Println("🚀 ACCOUNTUTOR 9000 // ORBITAL AUDIT DEFENSE — HALL OF FAME")
		fmt.Printf("Database: %s\n", dbPath)
		fmt.Println("==================================================================================")
		fmt.Printf("%-6s %-10s %-12s %-10s %-14s %-20s\n", "Rank", "Initials", "Score", "Blasted", "Flight Time", "Date Recorded")
		fmt.Println("----------------------------------------------------------------------------------")
		for i, s := range scores {
			fmt.Printf("#%-5d %-10s %-12d %-10d %-14s %-20s\n",
				i+1,
				fmt.Sprintf("[%s]", s.Initials),
				s.Score,
				s.BlastedCount,
				fmt.Sprintf("%ds", s.SurvivalSeconds),
				s.CreatedAt.Format("2006-01-02 15:04"),
			)
		}
		if len(scores) == 0 {
			fmt.Println("No pilot flight scores recorded yet. Launch the arcade intro to set an all-time high score!")
		}
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --export-json
	if *exportJSON {
		if err := db.ExportJSON(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error exporting JSON: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle --mastery (terminal text output)
	if *masteryFlag {
		attempts, err := db.GetAllAttempts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving attempts: %v\n", err)
			os.Exit(1)
		}
		projections := mastery.RebuildProjections(attempts)
		now := time.Now().UTC()

		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Learner Mastery Projections (as of %s)\n", now.Format(time.RFC3339))
		fmt.Printf("Database: %s\n", dbPath)
		fmt.Println("==================================================================================")
		fmt.Printf("%-32s %-12s %-10s %-8s %-10s %-12s\n", "Concept", "Attempts", "Retention", "Delayed", "Scaffold", "Status")
		fmt.Println("----------------------------------------------------------------------------------")

		for cid, stats := range projections {
			indep := fmt.Sprintf("%d / %d", stats.IndependentSuccesses, stats.IndependentAttempts)
			ret := fmt.Sprintf("%.0f%%", stats.RetentionFactor(now)*100.0)
			delayedStr := fmt.Sprintf("%d", stats.DelayedSuccesses)
			if stats.HasDelayedRetrieval {
				delayedStr += " ✓"
			}
			status := "Learning"
			if stats.EffectiveScore(now) >= 0.80 && stats.HasDelayedRetrieval {
				status = "Mastered"
			} else if stats.EffectiveScore(now) >= 0.80 && !stats.HasDelayedRetrieval {
				status = "Needs Spacing"
			} else if stats.EffectiveScore(now) < 0.50 {
				status = "Needs Review"
			}
			fmt.Printf("%-32s %-12s %-10s %-8s %-10s %-12s\n", cid, indep, ret, delayedStr, stats.ScaffoldLevel.String(), status)
		}
		if len(projections) == 0 {
			fmt.Println("No practice attempts recorded yet. Start practicing to generate mastery evidence!")
		}
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --candidates (terminal text summary of candidate store)
	if *candidatesFlag {
		filter := storage.CandidateFilter{}
		if *candidateFailedFlag {
			filter.ValidationStatus = string(candidate.ValidationFailed)
		}
		if *candidateFamilyFlag != "" {
			filter.FamilyID = *candidateFamilyFlag
		}
		cands, err := db.ListCandidates(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing candidates: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Stored Question Candidates (Outside Active Bank)\n")
		fmt.Printf("Total: %d candidates found\n", len(cands))
		fmt.Println("==================================================================================")
		fmt.Printf("%-24s %-22s %-10s %-8s %-16s\n", "Candidate ID", "Family", "Status", "Valid", "Source")
		fmt.Println("----------------------------------------------------------------------------------")
		for _, c := range cands {
			valMark := "✓"
			if c.ValidationStatus == candidate.ValidationFailed {
				valMark = "✗ (failed)"
			}
			fmt.Printf("%-24s %-22s %-10s %-8s %-16s\n", c.ID, c.FamilyID, c.Status, valMark, c.Provenance.Source)
		}
		if len(cands) == 0 {
			fmt.Println("No candidate questions currently stored. Run with --generate-candidate to create one!")
		}
		fmt.Println("==================================================================================")
		fmt.Println("Use --preview-candidate=<id> to view full wording, parameters, and derived postings.")
		os.Exit(0)
	}

	// Handle --preview-candidate
	if *previewCandidateFlag != "" {
		cand, err := db.GetCandidate(*previewCandidateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving candidate %q: %v\n", *previewCandidateFlag, err)
			os.Exit(1)
		}
		fmt.Print(cand.FormatPreview())
		os.Exit(0)
	}

	// Handle --approval-history
	if *approvalHistoryFlag {
		events, err := db.ListApprovalEvents()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving approval events: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Content Review & Promotion Audit Log\n")
		fmt.Printf("Total Events: %d recorded\n", len(events))
		fmt.Println("==================================================================================")
		fmt.Printf("%-20s %-8s %-24s %-16s %-10s %s\n", "Timestamp (UTC)", "Action", "Question/Cand ID", "Reviewer", "RuleVer", "Notes")
		fmt.Println("----------------------------------------------------------------------------------")
		for _, e := range events {
			notes := e.Notes
			if len(notes) > 30 {
				notes = notes[:27] + "..."
			}
			fmt.Printf("%-20s %-8s %-24s %-16s v%-9d %s\n",
				e.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
				e.Action,
				e.QuestionID,
				e.Reviewer,
				e.RuleVersion,
				notes,
			)
		}
		if len(events) == 0 {
			fmt.Println("No approval or review events recorded yet.")
		}
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --approve-candidate
	if *approveCandidateFlag != "" {
		if strings.TrimSpace(*reviewerFlag) == "" {
			fmt.Fprintf(os.Stderr, "error: --reviewer=<name> is required to approve a candidate question\n")
			os.Exit(1)
		}

		cand, err := db.GetCandidate(*approveCandidateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving candidate %q: %v\n", *approveCandidateFlag, err)
			os.Exit(1)
		}

		notes := *notesFlag
		if notes == "" {
			notes = "Approved via CLI review command"
		}

		pubQ, apprEvent, err := candidate.ApproveCandidate(cand, *reviewerFlag, notes, time.Now().UTC())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Approval rejected by review gate:\n%v\n", err)
			os.Exit(1)
		}

		if err := db.PublishQuestion(*pubQ, cand.ID, *apprEvent); err != nil {
			fmt.Fprintf(os.Stderr, "error publishing approved question to database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("✓ CANDIDATE APPROVED AND PROMOTED TO ACTIVE BANK\n")
		fmt.Println("==================================================================================")
		fmt.Printf("Candidate ID:      %s\n", cand.ID)
		fmt.Printf("Published ID:      %s\n", pubQ.ID)
		fmt.Printf("Family ID:         %s (Rule Version %d)\n", pubQ.FamilyID, pubQ.RuleVersion)
		fmt.Printf("Reviewer:          %s\n", *reviewerFlag)
		fmt.Printf("Approved At:       %s\n", *pubQ.Review.ApprovedAt)
		fmt.Printf("Status:            %s\n", pubQ.Status)
		fmt.Printf("Audit Event ID:    %s\n", apprEvent.ID)
		fmt.Println("----------------------------------------------------------------------------------")
		fmt.Printf("Scenario Template: %s\n", pubQ.ScenarioTemplate)
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --reject-candidate
	if *rejectCandidateFlag != "" {
		if strings.TrimSpace(*reviewerFlag) == "" {
			fmt.Fprintf(os.Stderr, "error: --reviewer=<name> is required to reject a candidate question\n")
			os.Exit(1)
		}

		cand, err := db.GetCandidate(*rejectCandidateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving candidate %q: %v\n", *rejectCandidateFlag, err)
			os.Exit(1)
		}

		reason := *notesFlag
		if reason == "" {
			reason = "Rejected during human review"
		}

		apprEvent, err := candidate.RejectCandidate(cand, *reviewerFlag, reason, time.Now().UTC())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Rejection error: %v\n", err)
			os.Exit(1)
		}

		if err := db.SaveCandidate(*cand); err != nil {
			fmt.Fprintf(os.Stderr, "error updating candidate: %v\n", err)
			os.Exit(1)
		}
		if err := db.RecordApprovalEvent(*apprEvent); err != nil {
			fmt.Fprintf(os.Stderr, "error recording rejection event: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("✗ CANDIDATE REJECTED\n")
		fmt.Println("==================================================================================")
		fmt.Printf("Candidate ID:      %s\n", cand.ID)
		fmt.Printf("Reviewer:          %s\n", *reviewerFlag)
		fmt.Printf("Reason:            %s\n", reason)
		fmt.Printf("Audit Event ID:    %s\n", apprEvent.ID)
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --repair-candidate
	if *repairCandidateFlag != "" {
		if strings.TrimSpace(*reviewerFlag) == "" {
			fmt.Fprintf(os.Stderr, "error: --reviewer=<name> is required to repair a candidate question\n")
			os.Exit(1)
		}
		if strings.TrimSpace(*scenarioTemplateFlag) == "" {
			fmt.Fprintf(os.Stderr, "error: --scenario-template=<text> is required to repair a candidate question\n")
			os.Exit(1)
		}

		catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading curriculum accounts: %v\n", err)
			os.Exit(1)
		}
		eng := engine.NewEngine(catalog)

		cand, err := db.GetCandidate(*repairCandidateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving candidate %q: %v\n", *repairCandidateFlag, err)
			os.Exit(1)
		}

		notes := *notesFlag
		if notes == "" {
			notes = "Repaired scenario template wording via CLI"
		}

		apprEvent, err := candidate.RepairCandidate(cand, *scenarioTemplateFlag, nil, eng, catalog, *reviewerFlag, notes, time.Now().UTC())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Repair validation error: %v\n", err)
			os.Exit(1)
		}

		if err := db.SaveCandidate(*cand); err != nil {
			fmt.Fprintf(os.Stderr, "error saving repaired candidate: %v\n", err)
			os.Exit(1)
		}
		if err := db.RecordApprovalEvent(*apprEvent); err != nil {
			fmt.Fprintf(os.Stderr, "error recording repair event: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("✓ CANDIDATE REPAIRED\n")
		fmt.Println("==================================================================================")
		fmt.Printf("Candidate ID:      %s\n", cand.ID)
		fmt.Printf("Reviewer:          %s\n", *reviewerFlag)
		fmt.Printf("Validation:        %s\n", cand.ValidationStatus)
		if cand.RejectionReason != "" {
			fmt.Printf("Warning:           %s\n", cand.RejectionReason)
		}
		fmt.Printf("Audit Event ID:    %s\n", apprEvent.ID)
		fmt.Printf("New Template:      %s\n", cand.ScenarioTemplate)
		fmt.Println("==================================================================================")
		if cand.ValidationStatus == candidate.ValidationValid {
			fmt.Println("Candidate passed validation! To approve, run:")
			fmt.Printf("  acctg --approve-candidate=%s --reviewer=%s\n", cand.ID, *reviewerFlag)
		}
		os.Exit(0)
	}

	// Handle --retire-question
	if *retireQuestionFlag != "" {
		if strings.TrimSpace(*reviewerFlag) == "" {
			fmt.Fprintf(os.Stderr, "error: --reviewer=<name> is required to retire a question\n")
			os.Exit(1)
		}

		reason := *notesFlag
		if reason == "" {
			reason = "Retired during curriculum review"
		}

		apprEvent, err := db.RetireQuestion(*retireQuestionFlag, *reviewerFlag, reason, time.Now().UTC())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retiring question: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("✓ QUESTION RETIRED FROM ACTIVE PRACTICE\n")
		fmt.Println("==================================================================================")
		fmt.Printf("Question ID:       %s\n", *retireQuestionFlag)
		fmt.Printf("Reviewer:          %s\n", *reviewerFlag)
		fmt.Printf("Reason:            %s\n", reason)
		fmt.Printf("Audit Event ID:    %s\n", apprEvent.ID)
		fmt.Println("Historical learner attempts remain completely intact and replayable.")
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --active-bank
	if *activeBankFlag {
		catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading accounts: %v\n", err)
			os.Exit(1)
		}
		qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading seed bank: %v\n", err)
			os.Exit(1)
		}
		activeQuestions, err := db.GetActiveBankQuestions(qBank.Questions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting active bank: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Active Question Bank (Eligible for Practice)\n")
		fmt.Printf("Total Active Questions: %d\n", len(activeQuestions))
		fmt.Println("==================================================================================")
		fmt.Printf("%-28s %-4s %-22s %-16s %s\n", "Question ID", "Ver", "Family", "Status", "Reviewer")
		fmt.Println("----------------------------------------------------------------------------------")
		for _, q := range activeQuestions {
			rev := "unreviewed"
			if q.Review.Reviewer != nil {
				rev = *q.Review.Reviewer
			}
			fmt.Printf("%-28s v%-3d %-22s %-16s %s\n", q.ID, q.Version, q.FamilyID, q.Status, rev)
		}
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --export-bank
	if *exportBankFlag {
		catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading accounts: %v\n", err)
			os.Exit(1)
		}
		qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading seed bank: %v\n", err)
			os.Exit(1)
		}
		activeQuestions, err := db.GetActiveBankQuestions(qBank.Questions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting active bank: %v\n", err)
			os.Exit(1)
		}

		bankFile := bank.QuestionBankFile{
			SchemaVersion: 1,
			Questions:     activeQuestions,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(bankFile); err != nil {
			fmt.Fprintf(os.Stderr, "error encoding bank: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle --generate-candidate
	if *generateCandidateFlag {
		catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading curriculum accounts: %v\n", err)
			os.Exit(1)
		}
		eng := engine.NewEngine(catalog)

		targetFamily := *candidateFamilyFlag
		targetConcept := ""

		if *weakestConceptFlag {
			attempts, _ := db.GetAllAttempts()
			projections := mastery.RebuildProjections(attempts)
			weakConcept, weakFam := candidate.SelectWeakestConcept(projections, time.Now().UTC())
			targetConcept = weakConcept
			if targetFamily == "" {
				targetFamily = weakFam
			}
			fmt.Printf("Targeting learner's weakest concept: %s (Family: %s)\n\n", targetConcept, targetFamily)
		}

		if targetFamily == "" {
			targetFamily = bank.FamilyCustomerAdvance
		}

		var gen candidate.CandidateGenerator
		if *tutorFlag != "offline" {
			budget := tutor.NewBudget(*tutorBudgetFlag, 0)
			var tut tutor.Tutor
			switch *tutorFlag {
			case "anthropic":
				tut = tutor.BuildTutor(tutor.FactoryOptions{Provider: tutor.ProviderAnthropic, Timeout: *tutorTimeoutFlag, Budget: budget, AuthStore: authStore})
			case "google", "gemini":
				tut = tutor.BuildTutor(tutor.FactoryOptions{Provider: tutor.ProviderGoogle, Timeout: *tutorTimeoutFlag, Budget: budget, AuthStore: authStore})
			case "openai":
				tut = tutor.BuildTutor(tutor.FactoryOptions{Provider: tutor.ProviderOpenAI, Timeout: *tutorTimeoutFlag, Budget: budget, AuthStore: authStore})
			case "chatgpt-plan", "chatgpt_plan":
				tut = tutor.BuildTutor(tutor.FactoryOptions{Provider: tutor.ProviderChatGPTPlan, Timeout: *tutorTimeoutFlag, Budget: budget, AuthStore: authStore})
			default:
				tut = tutor.NewOfflineTutor()
			}
			gen = candidate.NewProviderCandidateGenerator(tut, eng, catalog)
		} else {
			gen = candidate.NewOfflineCandidateGenerator(eng, catalog)
		}

		fmt.Printf("Generating candidate with generator: %s...\n", gen.Name())
		cand, genErr := gen.Generate(context.Background(), candidate.GenerateRequest{
			FamilyID:      targetFamily,
			TargetConcept: targetConcept,
			Seed:          *seedFlag,
		})

		if cand != nil {
			_ = db.SaveCandidate(*cand)
			fmt.Println(cand.FormatPreview())
		}

		if genErr != nil {
			fmt.Fprintf(os.Stderr, "Candidate generation error: %v\n", genErr)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Load curriculum
	catalog, _, err := bank.LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load accounts curriculum: %v\n", err)
		os.Exit(1)
	}

	qBank, err := bank.LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load seed question bank: %v\n", err)
		os.Exit(1)
	}

	activeQuestions, err := db.GetActiveBankQuestions(qBank.Questions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not merge active questions from db: %v\n", err)
		activeQuestions = qBank.Questions
	}

	// Handle --reconcile-all
	if *reconcileAllFlag {
		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Active Bank Postings Reconciliation Proof\n")
		fmt.Printf("[%s]\n", domain.GenericVisualNotice)
		fmt.Printf("Total Active Question Templates: %d\n", len(activeQuestions))
		fmt.Println("==================================================================================")

		reconciledCount := 0
		eng := engine.NewEngine(catalog)

		for i, q := range activeQuestions {
			param := map[string]int64{"amount_minor_units": 150000}
			ev := engine.TransactionEvent{
				FamilyID:   q.FamilyID,
				Parameters: param,
			}
			processed, err := eng.ProcessEvent(ev)
			if err != nil {
				fmt.Printf("[%2d] ✗ %-28s Family: %-24s Error: %v\n", i+1, q.ID, q.FamilyID, err)
				continue
			}

			recon, err := domain.ReconcileTransaction(processed.Entry, catalog)
			if err != nil || !recon.Reconciled {
				fmt.Printf("[%2d] ✗ %-28s Family: %-24s FAILED RECONCILIATION: %v\n", i+1, q.ID, q.FamilyID, err)
				continue
			}

			reconciledCount++
			fmt.Printf("[%2d] ✓ %-28s Family: %-24s Dr=%s Cr=%s ΔA=%s ΔL=%s ΔE=%s\n",
				i+1, q.ID, q.FamilyID,
				recon.TotalDr.FormatDollars(), recon.TotalCr.FormatDollars(),
				recon.DeltaAssets.FormatDollars(), recon.DeltaLiabilities.FormatDollars(), recon.DeltaEquity.FormatDollars())
		}

		fmt.Println("==================================================================================")
		fmt.Printf("Reconciliation Result: %d / %d active questions 100%% verified (Journal == T-Accounts == Equation)\n",
			reconciledCount, len(activeQuestions))
		fmt.Println("==================================================================================")
		os.Exit(0)
	}

	// Handle --statements and --case-study
	if *statementsFlag || *caseStudyFlag {
		eng := engine.NewEngine(catalog)
		c := statements.CanonicalCasePioneerConsulting()
		report, err := statements.BuildAccountingCycleReport(c, catalog, eng)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error building accounting cycle report: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(report.FormatReport())
		os.Exit(0)
	}

	// Handle --abandon-exam / --close-exam
	if *abandonExamFlag != "" || *closeExamFlag != "" {
		abandonID := *abandonExamFlag
		if abandonID == "" {
			abandonID = *closeExamFlag
		}
		if err := db.AbandonExamSession(abandonID); err != nil {
			fmt.Fprintf(os.Stderr, "error abandoning exam session: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Exam session %s abandoned.\n", abandonID)
		os.Exit(0)
	}

	// Handle --exam-history
	if *examHistoryFlag {
		sessions, err := db.ListExamSessions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing exam sessions: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("==================================================================================")
		fmt.Printf("ACCTG PRACTICE — Exam Assessment Sessions History\n")
		fmt.Printf("Total: %d exam sessions recorded\n", len(sessions))
		fmt.Println("==================================================================================")
		fmt.Printf("%-24s %-20s %-12s %-10s %-8s %-10s %s\n", "Exam Session ID", "Started (UTC)", "Duration", "Status", "Items", "Score", "Grade")
		fmt.Println("----------------------------------------------------------------------------------")
		for _, s := range sessions {
			durStr := fmt.Sprintf("%02dm %02ds", s.ElapsedSeconds/60, s.ElapsedSeconds%60)
			grade := "—"
			scoreStr := "—"
			if s.Status == exam.ExamStatusCompleted {
				scoreStr = fmt.Sprintf("%.1f%%", s.Score)
				switch {
				case s.Score >= 90.0:
					grade = "A"
				case s.Score >= 80.0:
					grade = "B"
				case s.Score >= 70.0:
					grade = "C"
				case s.Score >= 60.0:
					grade = "D"
				default:
					grade = "F"
				}
			}
			itemsStr := fmt.Sprintf("%d/%d", s.CorrectCount, s.TotalAttempts)
			fmt.Printf("%-24s %-20s %-12s %-10s %-8s %-10s %s\n",
				s.ID,
				s.StartedAt.UTC().Format("2006-01-02 15:04:05"),
				durStr,
				s.Status,
				itemsStr,
				scoreStr,
				grade,
			)
		}
		if len(sessions) == 0 {
			fmt.Println("No exam sessions recorded yet. Run with --exam to take an assessment!")
		}
		fmt.Println("==================================================================================")
		fmt.Println("Use --exam-report=<id> to view full skill breakdowns, error patterns, and reviews.")
		os.Exit(0)
	}

	// Handle --exam-report
	if *examReportFlag != "" {
		sess, err := db.GetExamSession(*examReportFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving exam session %q: %v\n", *examReportFlag, err)
			os.Exit(1)
		}
		atts, err := db.GetExamAttempts(*examReportFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving exam attempts for %q: %v\n", *examReportFlag, err)
			os.Exit(1)
		}

		eng := engine.NewEngine(catalog)
		gen := drill.NewGenerator(catalog, eng)

		runner, err := exam.ResumeExamRunner(*sess, atts, activeQuestions, catalog, eng, gen)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reconstructing exam report: %v\n", err)
			os.Exit(1)
		}
		report := runner.GenerateReport()
		fmt.Print(report.FormatText())
		os.Exit(0)
	}

	// Seed selection
	seed := *seedFlag
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	// Tutor provider selection
	var tutorProvider tutor.Tutor
	budget := tutor.NewBudget(*tutorBudgetFlag, 0)
	offline := tutor.NewOfflineTutor()

	switch *tutorFlag {
	case "simulated":
		sim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
			Delay:  300 * time.Millisecond,
			Budget: budget,
		})
		tutorProvider = tutor.NewFallbackTutor(sim, offline, *tutorTimeoutFlag)
	case "simulated-slow":
		sim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
			Delay:  5 * time.Second,
			Budget: budget,
		})
		tutorProvider = tutor.NewFallbackTutor(sim, offline, *tutorTimeoutFlag)
	case "simulated-error":
		sim := tutor.NewSimulatedProvider(tutor.SimulatedConfig{
			FailWith: errors.New("simulated upstream tutor provider failure"),
			Budget:   budget,
		})
		tutorProvider = tutor.NewFallbackTutor(sim, offline, *tutorTimeoutFlag)
	case "chatgpt-plan", "chatgpt_plan":
		tutorProvider = tutor.BuildTutor(tutor.FactoryOptions{
			Provider:  tutor.ProviderChatGPTPlan,
			Model:     *tutorModelFlag,
			Timeout:   *tutorTimeoutFlag,
			Budget:    budget,
			AuthStore: authStore,
		})
	case "anthropic":
		tutorProvider = tutor.BuildTutor(tutor.FactoryOptions{
			Provider:  tutor.ProviderAnthropic,
			Model:     *tutorModelFlag,
			Timeout:   *tutorTimeoutFlag,
			Budget:    budget,
			AuthStore: authStore,
		})
	case "google", "gemini":
		tutorProvider = tutor.BuildTutor(tutor.FactoryOptions{
			Provider:  tutor.ProviderGoogle,
			Model:     *tutorModelFlag,
			Timeout:   *tutorTimeoutFlag,
			Budget:    budget,
			AuthStore: authStore,
		})
	case "openai":
		tutorProvider = tutor.BuildTutor(tutor.FactoryOptions{
			Provider:  tutor.ProviderOpenAI,
			Model:     *tutorModelFlag,
			Timeout:   *tutorTimeoutFlag,
			Budget:    budget,
			AuthStore: authStore,
		})
	default: // "offline"
		tutorProvider = tutor.BuildTutor(tutor.FactoryOptions{
			Provider:  tutor.ProviderOffline,
			Timeout:   *tutorTimeoutFlag,
			Budget:    budget,
			AuthStore: authStore,
		})
	}

	// Initialize TUI Model
	examTime := *examTimeFlag
	if examTime == 0 && *timerFlag > 0 {
		examTime = *timerFlag
	}
	isExam := *examFlag || *resumeExamFlag || examTime > 0
	showIntro := !*skipIntroFlag && !*noIntroFlag && !isExam && !*practiceJournalFlag
	if *introFlag {
		showIntro = true
	}

	cfg := tui.Config{
		DB:             db,
		Catalog:        catalog,
		Questions:      activeQuestions,
		TotalQuestions: totalQuestions,
		Intensity:      mastery.SessionIntensity(*intensityFlag),
		Seed:           seed,
		Clock:          mastery.RealClock{},
		Tutor:          tutorProvider,
		TutorTimeout:   *tutorTimeoutFlag,
		AuthStore:      authStore,
		ModelCache:     modelCache,
		TutorModel:     *tutorModelFlag,
		OAuthClientID:  *oauthClientIDFlag,
		ExamMode:       isExam,
		ExamTimeLimit:  examTime,
		ResumeExam:     *resumeExamFlag,
		Intro:          showIntro,
	}

	model, err := tui.NewModel(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to initialize drill model: %v\n", err)
		os.Exit(1)
	}

	if *practiceJournalFlag {
		model.PreviousState = tui.StateDrill
		model.State = tui.StateJournalPractice
	}

	// Launch Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running application: %v\n", err)
		os.Exit(1)
	}
}
