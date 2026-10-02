# ACCTG 502 practice tutor — planning starter

This is a planning and handoff package, not a runnable application. Copy its contents into the course repository without overwriting existing instructions or course materials. Merge root instruction files when they already exist. Use uppercase `CLAUDE.md` and `AGENTS.md` for conventional discovery; avoid case-only duplicate filenames.

Start by reading [OVERVIEW.md](OVERVIEW.md), [PLAN.md](PLAN.md), and [docs/AGENT-CONTRACT.md](docs/AGENT-CONTRACT.md). The course syllabus lives in the destination repository and has not been inspected when preparing this package. Course scope, exact six-square layout, and assessment conventions remain verification tasks.

Suggested first prompt to Antigravity:

> Read AGENTS.md and the documents it references. Perform Stage 00 and Stage 01 of PLAN.md using the syllabus and course materials in this repository. Record evidence and unresolved assumptions in docs/COURSE-MAP.md. Then implement the offline vertical slice through Stage 07, one stage at a time, with the required checks and a concise handoff after each stage. Do not build later features before the offline slice works. Preserve existing repository files and never invent course requirements.

## Contents

- OVERVIEW.md: purpose, learner needs, decisions, boundaries.
- PLAN.md: sequenced implementation stages and acceptance gates.
- AGENTS.md / CLAUDE.md / GEMINI.md: thin entries to one shared contract.
- .agent/rules/accounting-tutor.md: optional Antigravity workspace-rule entry; verify discovery in the installed version.
- docs/: architecture, pedagogy, mastery, validation, course map, handoff.
- curriculum/: starter account catalog and semantic question fixtures; provisional course scope.
- config/example.toml: intended configuration contract, not an implemented parser.

No API key, database, executable, dependency lockfile, or generated course claims are included. The seed fixtures are independently specified elementary examples, not extracted assignments. Keep actual learner history and secrets out of git. No new agent should treat this package as proof that a particular course topic has been covered.
