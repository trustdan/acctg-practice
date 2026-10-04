# AccountTutor 9000

```text
   +--------------------------------------------------------+
   |             A C C O U N T T U T O R   9 0 0 0          |
   |                                                        |
   |       DEBIT             |             CREDIT           |
   |       ----------        |        ----------            |
   |           $             |             $                |
   |                                                        |
   |         [ DR ] ======== BALANCED ======== [ CR ]        |
   |                                                        |
   |          LEVEL UP YOUR LEDGER. KEEP IT BALANCED.        |
   +--------------------------------------------------------+
```

Financial accounting practice in your terminal: work through a business event, build a balanced entry, and see how it changes the accounting equation. AccountTutor runs offline, saves your progress locally, and offers optional AI help when you connect a provider.

![AccountTutor 9000 animated terminal demo](https://github.com/trustdan/acctg-practice/blob/main/acctg-practice.gif?raw=true)

**Current version: 0.24.0.** Native releases support Windows x64, macOS Apple Silicon and Intel, and Linux x64 and ARM64. No Go installation, external SQLite library, or separate curriculum download is needed to run a release binary.

> AI creates. Rules validate. Machine drills. History adapts.

## Contents

- [Install and run](#install-and-run)
- [Practice and assessment](#practice-and-assessment)
- [AI tutor](#ai-tutor)
- [Generate and review new questions](#generate-and-review-new-questions)
- [Arcade flight](#arcade-flight)
- [Keyboard reference](#keyboard-reference)
- [Command line](#command-line)
- [Local data and privacy](#local-data-and-privacy)
- [Accounting scope](#accounting-scope)
- [Build and release](#build-and-release)

## Install and run

Download the matching ZIP from [Releases](https://github.com/trustdan/acctg-practice/releases), extract all its files into one folder, and use a terminal window. Start with about 80 columns by 24 rows; long pages support scrolling.

| Computer | Release ZIP | Executable | Optional launcher |
|---|---|---|---|
| Windows x64 | `acctg-v0.24.0-windows-amd64.zip` | `acctg.exe` | Run the executable in Windows Terminal |
| Mac with Apple Silicon | `acctg-v0.24.0-macos-arm64.zip` | `acctg-mac-arm64` | `Launch-Tutor.command` |
| Intel Mac | `acctg-v0.24.0-macos-amd64.zip` | `acctg-mac-amd64` | `Launch-Tutor.command` |
| Linux x64 | `acctg-v0.24.0-linux-amd64.zip` | `acctg-linux-amd64` | `launch-tutor.sh` |
| Linux ARM64 | `acctg-v0.24.0-linux-arm64.zip` | `acctg-linux-arm64` | `launch-tutor.sh` |

The macOS classmate bundle contains both Mac executables and the same launcher; it is a ZIP with two architectures, not a universal executable.

### Windows

Open PowerShell in the extracted folder:

```powershell
.\acctg.exe
.\acctg.exe --skip-intro
```

### macOS

For Apple Silicon, open Terminal in the extracted folder:

```sh
chmod +x acctg-mac-arm64 Launch-Tutor.command
./acctg-mac-arm64
```

On Intel, substitute `acctg-mac-amd64`. After making the launcher executable, you can double-click `Launch-Tutor.command` in Finder. Keep it beside the matching binary. The current launcher attempts to remove the binary's quarantine attribute; these builds are unsigned and have not been notarized or smoke-tested on macOS. Use your normal macOS security approval process for a download you trust.

### Linux

For x64:

```sh
chmod +x acctg-linux-amd64 launch-tutor.sh
./launch-tutor.sh
# Or run directly:
./acctg-linux-amd64 --skip-intro
```

On ARM64, substitute `acctg-linux-arm64`. The shell launcher looks beside itself and in `dist/` and `bin/`, and attempts to open a terminal when started from a desktop file manager. ZIP extraction on some systems loses executable permissions, so run `chmod` first.

**Running a Git checkout:** compiled binaries and `dist/` are ignored by Git, so `git pull` updates only the source. In a source checkout, `launch-tutor.sh` builds the current code into `./acctg` on each launch before running it, using Go's build cache. This requires the Go version declared in `go.mod`. A failed build stops the launcher instead of running an old binary. Extracted release packages run their included binary without requiring Go.

To refresh a checkout using an older launcher immediately:

```sh
CGO_ENABLED=0 go build -o acctg ./cmd/acctg
./acctg --version
./acctg --skip-intro
```

After pulling the launcher fix, run `chmod +x launch-tutor.sh` if needed, then `./launch-tutor.sh --version` to build and verify the current version.

`accounttutor.desktop` is optional. Before installing it in `~/.local/share/applications/`, set `Exec` to the absolute path to your launcher and `Path` to its containing directory. Copying the unmodified desktop file into that menu directory will not locate a launcher stored elsewhere.

### Windows binary through Wine

The launchers include a Wine fallback if `acctg.exe` is present and Wine is already installed. Native binaries are available for every platform above. Wine compatibility has not been verified in this release.

## Practice and assessment

- **Progressive drills:** identify accounts, classify them, choose increase/decrease and debit/credit, then reconcile the entry with the accounting equation. Errors receive a causal hint, one retry, and an explanation. Reviewed pairs can then provide a guided comparison within the session.
- **Adaptive scaffolding:** practice moves from seven steps to four and then two as independent evidence supports it. Mistakes restore guidance. Reduced guidance requires independent success in two reviewed reasoning settings and a delayed success that changes settings. Guided comparisons do not count as independent evidence.
- **Journal entry practice:** enter multiple debit/credit lines, including split postings. Grading compares account totals independently of line order. Recaps include the transaction analysis grid, T-accounts, and equation effects derived from the same postings.
- **Financial statements:** inspect the Pioneer Consulting accounting-cycle case, trial balances, income statement, retained earnings, balance sheet, cash flows, and reconciliation proofs.
- **Exams:** take timed or untimed assessments with assistance and feedback withheld until completion. Saved exams support resumption, history, and reports with error patterns and question reviews.
- **Mastery:** see independent and assisted evidence, delayed retrieval, scaffolding, and time decay. Displayed retention is a scheduling estimate, not a validated probability or medical measurement.

## AI tutor

Offline help works without an account. Press `t` to open Tutor Settings and choose offline, ChatGPT subscription sign-in, Anthropic, Google Gemini, or OpenAI API.

For ChatGPT, choose `2` to open the browser sign-in flow. Registration is automatic; no manual OAuth client ID is needed. After sign-in, choose an available model. In settings, `m` opens model selection and `r` requests a fresh provider catalog; cached models support later selection. Use `a` in settings to disconnect the saved ChatGPT connection and sign in with another account.

For API providers, enter a key in settings or supply `ANTHROPIC_API_KEY`, `GEMINI_API_KEY`, or `OPENAI_API_KEY`. API billing is separate from subscription access. The application does not silently switch between these routes.

During practice, `?` requests a hint and `e` requests an explanation. Requests run asynchronously with a rotating loading indicator; `Esc` cancels. Provider failures display a notice and fall back to offline help. Tutor prose cannot change answer keys, grades, or mastery. The default request timeout is 60 seconds and the default session budget is 20 requests.

When you leave an LLM explanation (Esc, another hotkey, a mouse click, or quit), the app asks: **Would you like to save this explanation in the database?** Press `y` to save and continue the original action, `n` to continue without saving, or Esc to keep reading. Scrolling does not trigger the prompt. Save failures keep the text available for retry. Press uppercase `V` to browse saved explanations and their original question/stage/provider; use `n`/`p` to browse and `u`/`d` to scroll. These are personal advisory notes, not approved question content or grading evidence. Offline hints and fallback explanations do not trigger this prompt.

The available subscription permissions and models depend on the connected account. `--test-llm` reports local configuration diagnostics; it does not prove that live inference will succeed.

## Generate and review new questions

Press `n` from a drill, feedback screen, or question review to request a new question from your connected LLM. The request includes the current transaction family/concept and up to two same-family examples, without learner history. The loading indicator stays visible while it runs; `Esc` cancels. Offline fallback prose cannot become a question.

Press `p` to open **Review New Questions**, a holding area for proposed questions saved locally before they join your practice bank. You can also press `g` inside review to create a local variation without a provider.

Review the scenario, parameters, locally derived journal entry, proposed hints/explanations, validation results, and provenance. Check that the wording actually describes the displayed accounting event: balanced entries and automated validation alone do not establish accounting correctness.

| Review key | Action |
|---|---|
| `j` / `k`, arrows | Browse proposals |
| `n` | Request a new question from the selected LLM |
| `g` | Generate a local variation |
| `a` | Explicitly approve and publish the selected proposal |
| `r` | Reject the selected proposal |
| `x` / Delete | Delete the selected proposal |
| `u` / `d` | Scroll up/down half a page |
| `p` / Esc | Return to practice |

Only explicit approval adds a validated proposal to the active practice bank. Reviewed teaching text is preserved with published content and question snapshots. Failed proposals stay outside the active bank. Reviewing content does not award grades or mastery, and successful learner answers never approve a proposal automatically. From practice, press `h` or F1 for the in-app review walkthrough; use `u`/`d` to read the full Help page.

For a first try: connect a tutor with `t`, return to practice, press `n`, and read the proposed scenario, answer, and every hint/explanation. Press `a` only when you approve the content. Each request creates one proposal; repeat `n` for more variations.

## Arcade flight

Startup includes an accounting-themed spaceship arcade with an autopilot demo and manual flight. Steer with arrows or `w`/`s`, fire with `f`/Space, toggle auto-fire with `g`, and use `b` for a smart bomb. Manual controls take over from autopilot. Speed and difficulty rise during play. Press `p` to pause and `L` to inspect the leaderboard. Survive four quarters to close the fiscal year, earn a bonus for remaining audit health, and choose `y` to continue into a harder year or `n` to retire with a win. After game over, `r` restarts the flight.

The HUD tracks Internal Audit shields and External Audit integrity. Collisions damage your shields; escaped hazards damage global integrity. Heavy hazards such as fraud and insider trading require multiple hits and can fragment into smaller targets. New all-time records prompt for three initials and persist in SQLite; arcade scores do not affect accounting mastery.

Press Enter or Esc to move into practice. From practice, `A` opens the arcade and `L` opens high scores. When you leave the game, after any required high-score entry, a cyan ASCII **AccounTutor 9000** title slides in from the right for one second and holds for two seconds before automatically returning to the application. Narrow terminals use a compact title. Closing a leaderboard opened directly from practice returns immediately. Use `--skip-intro` to start directly in drills.

## Keyboard reference

Keys depend on the current screen; text fields accept ordinary typing. Uppercase shortcuts matter.

| Key | Practice action |
|---|---|
| `j` / `k`, arrows | Move selection |
| `1`-`4` | Select answer; `a`/`b`/`c` remain aliases for the first three |
| Enter | Submit selection or continue |
| Space | Continue after feedback |
| `u` / `d`, Page Up / Page Down | Scroll tall pages, including tutor responses and help |
| `?` / `e` | Hint / explanation |
| `h` / F1 | Help and accounting reference |
| `s` | Mastery dashboard |
| `V` | Read saved explanations |
| `t` | Tutor settings |
| `p` / `n` | Review questions / request a new LLM question |
| `J` | Journal entry practice |
| `F` | Financial statements |
| `E` | Exam mode |
| `[` / `]` | Adjust session size |
| `i` | Cycle practice intensity |
| `A` / `L` | Arcade / high scores |
| Esc | Cancel request or dismiss current screen |
| `q` / Ctrl+C | Quit |

Recaps and statement views also support `j`/`k`, arrows, and `g`/`G` scrolling. In journal practice, use uppercase `D` for debit so lowercase `d` remains available for scrolling. Opening a reference or requesting help is tracked as assistance; exams suppress these aids.

## Command line

Examples use the Windows filename. On macOS or Linux, substitute your executable.

```powershell
.\acctg.exe --help
.\acctg.exe --version
.\acctg.exe --skip-intro --size=10 --intensity=spaced
.\acctg.exe --practice-journal
.\acctg.exe --statements
.\acctg.exe --reconcile-all
.\acctg.exe --exam --exam-time=15m
.\acctg.exe --resume-exam
.\acctg.exe --exam-history
.\acctg.exe --exam-report=SESSION_ID
.\acctg.exe --mastery
.\acctg.exe --high-scores
.\acctg.exe --list-models
.\acctg.exe --fetch-models
.\acctg.exe --test-llm
.\acctg.exe --tutor=openai --tutor-model=MODEL_ID
.\acctg.exe --generate-candidate --tutor=chatgpt-plan --candidate-family=customer_advance
.\acctg.exe --generate-candidate --tutor=openai --candidate-family=cash_service
.\acctg.exe --generate-candidate --weakest-concept
.\acctg.exe --candidates
.\acctg.exe --preview-candidate=CANDIDATE_ID
.\acctg.exe --approve-candidate=CANDIDATE_ID --reviewer="Your name" --notes="Reviewed wording and entry"
.\acctg.exe --active-bank
.\acctg.exe --export-bank
.\acctg.exe --export-json
```

The ChatGPT generation command uses your saved connection; API providers use their configured keys. Omitting `--tutor` generates offline. Generated content is saved for review; approve it separately with `--approve-candidate`, your reviewer name, and review notes.

Use `--questions` or `--size` for session length; intensities are `standard`, `spaced`, `intensive`, and `transfer`. `--seed` makes randomized selection reproducible. `--data-dir` overrides application storage; `--db` overrides only the database path. `--help` lists candidate generation, repair, rejection, retirement, exam abandonment, provider budgets, and aliases.

## Local data and privacy

| Platform | Default application directory |
|---|---|
| Windows | `%APPDATA%\acctg-practice` |
| macOS | `~/Library/Application Support/acctg-practice` |
| Linux | `$XDG_CONFIG_HOME/acctg-practice`, or `~/.config/acctg-practice` |

`acctg_practice.db` stores attempts, question snapshots, local published content, exams, arcade scores, and learner-saved explanations. `tutor_auth.json` stores credentials and provider settings; `models_cache.json` stores discovered models. Credentials use local JSON file storage with restrictive permission requests, not an encrypted OS keychain. Protect that directory and exclude it from releases and commits.

Schema upgrades create database backups. For a manual file backup, close the application first and copy the database along with any remaining SQLite sidecar files. JSON export provides practice history, not a complete replacement for a database backup. Optional provider requests send the context needed for the selected action; offline drills require no network.

## Accounting scope

The embedded seed bank contains **90 active reviewed scenarios across 13 families**, plus one retired record. It covers cash services, credit services, customer advances and their fulfillment, receivable collection, cash expenses, prepaid purchases/consumption, equipment purchases, borrowing, principal repayment, stock issuance, and dividends. Reviewed local questions extend these supported families. Guided contrasts follow eligible mistakes; amount choices distinguish today's transaction from original totals, prior entries and per-recipient amounts.

The current 90-scenario quality pass covers journal and exam replay, focused offline help, delayed transfer evidence and stronger distractors. See [the quality review](docs/QUALITY-REVIEW90.md) for coverage and remaining validation gates. Expansion beyond 90 has not started.

Money uses integer minor units. Debit means left and credit means right; increase/decrease follows account metadata. Revenue records earning, and cash receipt alone does not establish revenue. Journal entries, T-accounts, and statement effects share deterministic rules.

See [the course evidence map](https://github.com/trustdan/acctg-practice/blob/main/docs/COURSE-MAP.md) for source status and limitations. No local syllabus or slide deck is currently available; course-specific coverage and exam rules should be checked against your actual course materials.

## Build and release

Source builds require **Go 1.27.1 or newer**, as declared in [go.mod](https://github.com/trustdan/acctg-practice/blob/main/go.mod). SQLite uses a pure-Go driver; cross-compilation uses `CGO_ENABLED=0`.

```sh
go build -o acctg ./cmd/acctg
gofmt -s -d .
go test ./...
go vet ./...
```

From PowerShell, run `./scripts/build_releases.ps1`. From a Bash environment with Go and ZIP tools, run `./scripts/build_releases.sh`. These build all five targets, package launchers and documentation, and generate `dist/checksums.sha256`.

### Initial release attachments

Upload these files from `dist/` to the release:

- `acctg-v0.24.0-windows-amd64.zip`
- `acctg-v0.24.0-macos-arm64.zip`
- `acctg-v0.24.0-macos-amd64.zip`
- `acctg-v0.24.0-linux-amd64.zip`
- `acctg-v0.24.0-linux-arm64.zip`
- `checksums.sha256`

Optionally add `acctg-v0.24.0-macos-classmate-bundle.zip` for classmates who want one download for either Mac architecture. Bash builds additionally produce Linux `.tar.gz` archives with executable permissions.

The ZIPs already contain the corresponding executable, launcher where applicable, README, and release notes. Linux ZIPs also include the optional desktop template. If distributing loose files instead, add `acctg-mac-arm64`, `acctg-mac-amd64`, `acctg-linux-amd64`, `acctg-linux-arm64`, `Launch-Tutor.command`, `README.md`, and `RELEASE-NOTES.md` alongside the Windows `.exe` and Linux `.sh`; the shell script cannot run without a matching binary. The duplicate `acctg-windows-amd64.exe` is an alternative filename for `acctg.exe` and need not be uploaded separately.

The curriculum and account catalog are embedded. Do not ship learner databases, authentication files, model caches, API keys, or your local data directory. No standalone DLLs, SQLite installer, Go runtime, or source tree are required. Cross-compilation alone does not verify interactive behavior on macOS/Linux; test extracted packages on those platforms before claiming native runtime verification.

See [CONTRIBUTING.md](https://github.com/trustdan/acctg-practice/blob/main/CONTRIBUTING.md) for contribution guidance and [release notes](https://github.com/trustdan/acctg-practice/blob/main/docs/RELEASE-NOTES.md) for changes.
