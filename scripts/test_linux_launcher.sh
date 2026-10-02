#!/usr/bin/env bash
# Isolated launcher regressions: no real compiler, GUI, network, or learner data.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
TEST_ROOT="$(mktemp -d "$ROOT/scripts/.launcher-tests.XXXXXX")"
case "$TEST_ROOT" in "$ROOT"/scripts/.launcher-tests.*) ;; *) exit 1 ;; esac
trap 'rm -rf -- "$TEST_ROOT"' EXIT
MOCKS="$TEST_ROOT/tools"
mkdir -p "$MOCKS"
for command_name in dirname basename chmod; do
    command_path="$(command -v "$command_name")"
    printf '#!/bin/bash\nexec %q "$@"\n' "$command_path" > "$MOCKS/$command_name"
    chmod +x "$MOCKS/$command_name"
done
cat > "$MOCKS/uname" <<'MOCK'
#!/bin/bash
printf '%s\n' "${MOCK_ARCH:-x86_64}"
MOCK
cat > "$MOCKS/go" <<'MOCK'
#!/bin/bash
if [ "$1" = env ]; then printf '%s\n' "${MOCK_GO_ARCH:-amd64}"; exit; fi
if [ "${FAIL_BUILD:-0}" = 1 ]; then echo 'BUILD FAILED' >&2; exit 2; fi
if [ "$GOOS" != linux ] || [ "$GOARCH" != "${MOCK_GO_ARCH:-amd64}" ] || [ "$CGO_ENABLED" != 0 ]; then exit 3; fi
printf '#!/bin/bash\nprintf "CURRENT\\n"\nprintf "<%%s>\\n" "$@"\n' > acctg
chmod +x acctg
MOCK
chmod +x "$MOCKS/uname" "$MOCKS/go"

fixture() {
    CASE_DIR="$TEST_ROOT/$1 with spaces"
    mkdir -p "$CASE_DIR"
    cp "$ROOT/launch-tutor.sh" "$CASE_DIR/launch-tutor.sh"
}
binary() {
    printf '#!/bin/bash\nprintf "%s\\n"\nprintf "<%%s>\\n" "$@"\n' "$2" > "$CASE_DIR/$1"
    chmod +x "$CASE_DIR/$1"
}
source_tree() {
    mkdir -p "$CASE_DIR/cmd/acctg"
    touch "$CASE_DIR/go.mod" "$CASE_DIR/cmd/acctg/main.go"
}
run_launcher() { PATH="$MOCKS" /bin/bash "$CASE_DIR/launch-tutor.sh" "$@"; }

fixture checkout
source_tree
binary acctg-linux-amd64 STALE_RELEASE
binary acctg STALE_LOCAL
output="$(GOOS=windows GOARCH=arm64 run_launcher --version 'two words')"
[[ "$output" == *CURRENT* && "$output" != *STALE* && "$output" == *'<two words>'* ]]
echo 'PASS: source checkout rebuilds and wins over stale binaries; args and native target preserved'

if FAIL_BUILD=1 run_launcher > "$TEST_ROOT/failure.log" 2>&1; then exit 1; fi
if grep -q CURRENT "$TEST_ROOT/failure.log"; then exit 1; fi
echo 'PASS: build failure does not launch stale cached code'

mv "$MOCKS/go" "$MOCKS/go.disabled"
if run_launcher > "$TEST_ROOT/no-go.log" 2>&1; then exit 1; fi
grep -q 'Go is required' "$TEST_ROOT/no-go.log"
mv "$MOCKS/go.disabled" "$MOCKS/go"
echo 'PASS: source checkout without Go explains how to get current code'

fixture release
binary acctg-linux-amd64 RELEASE_X64
binary acctg-linux-arm64 RELEASE_ARM
output="$(run_launcher --version)"
[[ "$output" == *RELEASE_X64* && "$output" != *CURRENT* ]]
output="$(MOCK_ARCH=aarch64 run_launcher --version)"
[[ "$output" == *RELEASE_ARM* ]]
echo 'PASS: extracted packages select their native architecture without compiling'

fixture wrong-architecture
binary acctg-linux-arm64 WRONG_ARCH
if run_launcher > "$TEST_ROOT/wrong-arch.log" 2>&1; then exit 1; fi
echo 'PASS: incompatible architecture is never selected as fallback'

fixture desktop
binary acctg-linux-amd64 RELEASE_X64
cat > "$MOCKS/x-terminal-emulator" <<'MOCK'
#!/bin/bash
printf '<%s>\n' "$@"
MOCK
chmod +x "$MOCKS/x-terminal-emulator"
output="$(run_launcher 'two words')"
[[ "$output" == *"<$CASE_DIR/launch-tutor.sh>"* && "$output" == *'<two words>'* ]]
echo 'PASS: desktop relaunch uses the correct absolute script path and preserves spaces'
