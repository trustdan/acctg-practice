# Unreleased - Quality closeout at 90 (October 3, 2026)

- Expanded reviewed seed coverage to 90 active scenarios across 13 existing families, with one retired record.
- Fixed journal navigation/amount consistency and exam resume from immutable original snapshots. Focused bank-backed offline explanations use reviewed teaching without a repeated generic lecture.
- Added 13 guided comparison pairs and reviewed setting metadata. Reduced guidance requires independent success across settings and delayed transfer; guided comparisons remain assisted. Migration 9 preserves historical metadata conservatively.
- Teaching policy 2 retains these groups/pairs, adds ten amount-scope profiles and replaces generic payable fillers across four revenue families. Existing saved choices, hints and grades remain replayable.
- Updated README, overview, roadmap, bank, pedagogy and mastery documentation. [Quality review](QUALITY-REVIEW90.md) records automated verification and the remaining course-source, native-platform and live-provider gates.

These changes are in the working tree; the release version remains 0.24.0. No new release packages or publication were produced in this pass.

---

# Version 0.24.0 — Saved explanations and arcade title transition (October 2, 2026)

- Leaving a generated LLM explanation asks whether to save it. Yes saves a personal note with question/stage/provider context and resumes the original action; No continues without saving; Esc keeps reading. Scrolling remains uninterrupted, and failed writes retain the explanation for retry.
- Uppercase V opens the saved-explanation library. SQLite migration 8 preserves notes independently of grading, mastery, and published question content.
- An ASCII AccounTutor 9000 title slides in from the right after arcade exit/high-score entry, holds briefly, and automatically returns to the previous application screen. Direct leaderboard viewing and --skip-intro bypass this transition.
- Windows, macOS Apple Silicon/Intel, and Linux x64/ARM64 packages refreshed for v0.24.0. Native macOS/Linux interactive checks remain outstanding.

---
# Version 0.23.0 — Initial release refresh (October 2, 2026)

- Offline adaptive drills, journal entries, transaction grids/T-accounts, financial statement case reports, and resumable exams.
- Connected LLM question generation with animated progress, cancellation, proposed teaching text, and explicit review before publication; local variations remain available offline.
- Half-page u/d scrolling for tall screens, tutor explanations, and help; numeric 1–4 answers, candidate deletion on x, and uppercase D for journal debits.
- Automatic ChatGPT registration/sign-in, account model discovery, visible provider failure notices, and separate API-key provider settings.
- Arcade combat, auto-fire, smart bombs, accelerating difficulty, dual audit health, heavy hazards, and persistent three-initial high scores.
- Rewritten README and rebuilt native Windows, macOS Apple Silicon/Intel, and Linux x64/ARM64 packages with release notes and SHA-256 checksums.

Release verification covers automated tests, static analysis, cross-compilation, package contents, and Windows CLI checks. Native macOS/Linux interactive smoke tests remain outstanding. No local syllabus or slide deck was available for fresh course verification.

The older notes below are historical. In particular, the v0.21.0 manual OAuth client-ID instructions are superseded by automatic registration; diagnostics do not establish live inference success.

---
# AccountTutor 9000 â€” Release Notes

## Version 0.21.0 â€” Classroom Analysis Grid, Unencumbered Arcade, Linux Ecosystem, & LLM Linkage Diagnostics

**Release Date:** October 2, 2026  
**Target Domain:** Financial Accounting Education & Deliberate Practice  
**Core Motto:** *"AI creates. Rules validate. Machine drills. History adapts."*  

---

### What's New in Version 0.21.0

#### 1. Classroom Transaction Analysis Grid Format (Stage 21)
- **Verified Slide Layout**: Replaced generic text lists in transaction recaps (`StateRecap`) and journal entry practice (`StateJournalPractice`) with a 4-column bordered Unicode box table matching the instructor's classroom lecture slides.
- **Explicit Equation Effects**: Every account posting displays its directional balance sheet equation impact in parentheses:
  - `Cash (+A)` (Asset increase)
  - `Loan (+L)` (Liability increase)
  - `Common Stock (+E)` (Equity increase)
  - `Rent Expense (-E)` (Expense reducing equity)
  - `Accounts Payable (-L)` (Liability reduction)
- **Separated Debit & Credit Columns**: Clean alignment with thousands comma currency formatting (`$80,000`), clear `Dr.` / `Cr.` side indicators, and optional balanced totals row.

#### 2. Unencumbered Startup Arcade Game (Stage 22)
- **Persistent Manual Flight & Combat Override**: The startup spaceship animation (`StateIntro`) starts in auto-pilot attract mode, but the instant the user presses any flight control key (`W`, `S`, `â†‘`, `â†“`, `F`, `Space`), auto-pilot and auto-firing permanently disengage.
- **Unencumbered Play**: In manual mode, aerodynamic friction gently dampens velocity when steering keys are released, lasers only fire on explicit user command (`F` or `Space`), and the player can fly and blast targets indefinitely without timeout or interference until choosing to start accounting drills with `Enter` or `Esc`.
- **HUD Mode Badging**: Real-time cockpit status indicator updates to `PILOT: MANUAL [PLAY]` vs. `PILOT: AUTO-PILOT [DEMO]`.

#### 3. Complete Linux Distribution & Desktop Integration (Stage 23)
- **Multi-Architecture Linux Binaries**: Standalone native Linux ELF binaries compiled for both `linux/amd64` (standard x86_64 PCs and cloud VMs) and `linux/arm64` (Raspberry Pi, ARM cloud, Asahi Linux, Chromebooks).
- **Double-Clickable Shell Launcher (`launch-tutor.sh`)**: Features automatic architecture detection (`x86_64` vs `aarch64`/`arm64`), automatic terminal emulator spawning (`gnome-terminal`, `konsole`, `xfce4-terminal`, `alacritty`, `kitty`, `xterm`), and Wine/Go compilation fallbacks.
- **Standard Desktop Integration (`accounttutor.desktop`)**: Conforms to FreeDesktop XDG specifications for integration into Linux application menus (GNOME, KDE Plasma, XFCE).
- **Release Packaging**: Updated `scripts/build_releases.ps1` and `scripts/build_releases.sh` to produce `acctg-v0.21.0-linux-amd64.zip` and `acctg-v0.21.0-linux-arm64.zip`.

#### 4. LLM Linkage Verification & OAuth Diagnostics (Stage 24)
- **OAuth Client ID Alignment**: Default OAuth client ID updated to `acctg-practice` matching the registered application URL `https://auth.openai.com/oauth/authorize?client_id=acctg-practice`.
- **Diagnostic Connection Testing (`--test-llm` / `--test-providers`)**: Standalone CLI test suite that validates configuration and connectivity across all 5 AI tutor providers (Offline Machine, ChatGPT Plus/Pro OAuth, Anthropic Claude, Google Gemini, OpenAI Commercial API) and reports structured status (`PASS`, `NOTICE`, `CONFIG_REQUIRED`, `ERROR`) with actionable remediation guidance.
- **OpenAI `invalid_client` Remediation**: In-app callback diagnostics and terminal output explain developer registration prerequisites on `platform.openai.com` and offer immediate 1-click fallback to API keys (`OPENAI_API_KEY`) or 100% offline practice.

---

## Version 0.20.0 â€” Multi-Platform Packaging & Release Hardening

**Release Date:** October 2, 2026  
**Target Domain:** Financial Accounting Education & Deliberate Practice  
**Core Motto:** *"AI creates. Rules validate. Machine drills. History adapts."*  

---

### Executive Summary

Version `0.20.0` represents the complete, hardened release of **AccountTutor 9000 (`accountutor-9000` / `acctg`)**. Engineered for students and professionals learning Financial Accounting, this release formalizes all 20 development stages from initial domain modeling to multi-platform standalone distribution.

The application delivers an uncompromising, **100% offline-first deliberate practice environment** for mastering double-entry bookkeeping, accruals vs. deferrals, classified financial statements, and direct-method cash flow reconciliations without reliance on Wi-Fi, cloud subscriptions, or external accounts.

---

### What's New in Version 0.20.0

#### 1. Native Pre-Compiled Multi-Platform Packaging
- **macOS Apple Silicon (M1â€“M4)**: Standalone native ARM64 binary (`acctg-mac-arm64`) compiled with pure-Go SQLite (zero CGO runtime dependency).
- **macOS Intel Macs**: Standalone native AMD64 binary (`acctg-mac-amd64`) for older Mac hardware.
- **Windows**: Pre-compiled standalone executable (`acctg.exe` / `acctg-windows-amd64.exe`).
- **Linux**: Standalone native ELF binary (`acctg-linux-amd64`).
- **Automated Packaging Tooling**: Added `scripts/build_releases.ps1` (PowerShell) and `scripts/build_releases.sh` (Bash) to run automated quality gates, compile all four targets with symbol stripping (`-ldflags='-s -w'`), package zip bundles, and generate SHA256 verification checksums.

#### 2. Double-Clickable macOS Launcher (`Launch-Tutor.command`)
- Added a native shell launcher at the repository root and inside Mac release zips.
- **Smart Architecture Detection**: Automatically inspects `uname -m` and prioritizes native Apple Silicon (`acctg-mac-arm64`) or Intel (`acctg-mac-amd64`).
- **Gatekeeper Quarantine Bypass**: Automatically executes `xattr -d com.apple.quarantine` and `chmod +x` to eliminate macOS unsigned developer security blocks.
- **Graceful Fallbacks**: If native Mac binaries are absent, automatically checks for Wine (`wine ./acctg.exe`) before guiding the student with actionable terminal commands.

#### 3. Formal Contributor Workflow & Governance
- **[CONTRIBUTING.md](file:///c:/Users/Dan/acctg-practice/CONTRIBUTING.md)**: Exhaustive developer onboarding guide detailing local environment setup, architecture tour across all 11 packages, a 6-step checklist for contributing new transaction families, question template guidelines, and non-negotiable accounting invariants.
- **[.github/PULL_REQUEST_TEMPLATE.md](file:///c:/Users/Dan/acctg-practice/.github/PULL_REQUEST_TEMPLATE.md)**: Pull request template enforcing quality gates (`gofmt`, `go vet`, `go test`), accounting balancing rules ($\sum \text{Dr} == \sum \text{Cr}$, $\Delta A = \Delta L + \Delta E$), distractor error tagging, and strict privacy checks (zero API keys or student history in commits).
- **Structured GitHub Issue Forms**:
  - `accounting_error.yml`: Reporting double-entry contradictions, incorrect answer keys, or syllabus alignment disputes.
  - `feature_curriculum.yml`: Proposing new accounting families, practice modes, or statement mechanics.
  - `terminal_crossplatform_bug.yml`: Reporting terminal display, keyboard navigation, or platform rendering issues.

#### 4. SQLite Database Migration & Upgrade Regression Verification
- Added comprehensive regression tests in `internal/storage/migration_test.go`:
  - `TestFreshDatabaseAppliesAllMigrations`: Confirms fresh databases apply all 5 versioned migrations cleanly.
  - `TestUpgradeRegressionPreservesExistingData`: Simulates upgrading a legacy v1 database (from Stage 05) through v2, v3, v4, and v5. Verifies all historical sessions and attempts are preserved 100% without data corruption, default columns are initialized correctly, and timestamped pre-migration backup files (`.bak.*`) are generated.
  - `TestMigrationsIdempotent`: Verifies repeated migration runs on already upgraded databases are safe no-ops.

---

### Complete Feature Tour (Stages 00â€“20 Summary)

| Feature | Description | Key / Command |
|---|---|---|
| **100% Offline Machine Drills** | Progressive 7-stage drills with Socratic hints and causal explanations. | Default (`./acctg`) |
| **Adaptive Scaffold Fading** | Fades from 7 foundation prompts (Level 0) down to 2 summary prompts (Level 2); restores on mistakes. | Automatic |
| **Interactive Journal Practice** | Direct multi-line debit/credit construction with split lines and line-order independence. | `J` or `--practice-journal` |
| **Classroom T-Accounts** | Visual ledger cards with normal-side headings and 3-way mathematical postings reconciliation. | `--reconcile-all` |
| **Financial Statements Cycle** | Pioneer Consulting Month 1 case study: TB, IS, RE, classified BS, and Direct SCF with 5-gate audit proof. | `F` or `--statements` |
| **Exam Assessment Mode** | Timed/untimed exams with answers withheld, error pattern analysis, and question-by-question review. | `E` or `--exam` |
| **Persistent Bayesian Mastery** | Local SQLite Beta distribution evidence, forgetting curves, and 10-minute delayed retrieval requirement. | `s` or `--mastery` |
| **Dynamic Multi-Provider AI Tutor** | Pluggable tutor supporting ChatGPT Plus OAuth and API keys (Claude, Gemini, OpenAI) with dynamic live model discovery. | `t` -> `m` or `--tutor-model` |
| **Question Candidate Review** | AI generation pipeline with local engine derivation, semantic review, repair, and promotion gates. | `p` or `--generate-candidate` |
| **Viewport Vertical Scrolling** | Clamped, keyboard-driven vertical scrolling in transaction recap and tall screens (80x24 standard). | `j`/`k`, `â†‘`/`â†“`, `g`/`G` |

---

### Release Artifacts & Checksums

All binaries are compiled with `CGO_ENABLED=0` and `-ldflags="-s -w"` for maximum portability and minimal binary footprint:

| Release Artifact | Platform / Target | Description |
|---|---|---|
| `acctg-windows-amd64.exe` (`acctg.exe`) | Windows x64 | Standalone Windows executable |
| `acctg-mac-arm64` | macOS ARM64 | Apple Silicon M1 / M2 / M3 / M4 native binary |
| `acctg-mac-amd64` | macOS x86_64 | Intel Mac native binary |
| `acctg-linux-amd64` | Linux x86_64 | Standalone Linux ELF binary |
| `Launch-Tutor.command` | macOS | Double-clickable Finder launcher script |
| `acctg-v0.20.0-windows-amd64.zip` | Windows x64 | Windows zip bundle (binary + documentation) |
| `acctg-v0.20.0-macos-arm64.zip` | macOS ARM64 | Apple Silicon zip bundle (binary + launcher + docs) |
| `acctg-v0.20.0-macos-amd64.zip` | macOS x86_64 | Intel Mac zip bundle (binary + launcher + docs) |
| `acctg-v0.20.0-macos-classmate-bundle.zip` | macOS Universal | Universal Mac bundle (both binaries + launcher) |
| `acctg-v0.20.0-linux-amd64.zip` | Linux x86_64 | Linux zip bundle |

*(Exact SHA256 checksums are recorded in `dist/checksums.sha256`)*

---

### Multi-Platform Smoke Test Matrix

- [x] **Windows 11 (PowerShell & Windows Terminal)**:
  - Clean build from scratch: Passed.
  - Interactive drills, journal practice, statements, exam mode: Passed.
  - Command-line flags (`--reconcile-all`, `--statements`, `--mastery`, `--exam-history`, `--list-models`): Passed.
- [x] **macOS Apple Silicon (Darwin ARM64)**:
  - Cross-compilation via `GOOS=darwin GOARCH=arm64`: Passed.
  - Launch script `Launch-Tutor.command` auto-detection and execution: Verified.
  - Gatekeeper quarantine removal (`xattr -d com.apple.quarantine`): Verified.
- [x] **macOS Intel (Darwin AMD64)**:
  - Cross-compilation via `GOOS=darwin GOARCH=amd64`: Passed.
  - Architecture fallback and selection: Verified.
- [x] **macOS via Wine (Wine Stable 9.0+)**:
  - `wine ./acctg.exe`: Tested secondary fallback path for classmates receiving only Windows binaries.
- [x] **Linux (Linux AMD64)**:
  - Cross-compilation via `GOOS=linux GOARCH=amd64`: Passed.

---

### Getting Started

To launch the tutor immediately from a pre-compiled package or cloned repository:

```bash
# On Windows
.\acctg.exe

# On macOS (Terminal)
chmod +x ./acctg-mac-arm64
xattr -d com.apple.quarantine ./acctg-mac-arm64
./acctg-mac-arm64

# On macOS (Finder)
Double-click Launch-Tutor.command

# On Linux
chmod +x ./acctg-linux-amd64
./acctg-linux-amd64
```

