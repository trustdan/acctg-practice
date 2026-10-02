#!/usr/bin/env bash
# AccountTutor 9000 Ã¢â‚¬â€ Multi-Platform Release Packaging (Bash)
# Builds standalone native binaries for Windows, macOS (Apple Silicon & Intel), and Linux.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${ROOT_DIR}"

echo "=================================================================="
echo " AccountTutor 9000 Ã¢â‚¬â€ Multi-Platform Build & Packaging"
echo "=================================================================="

# 1. Detect Version
VERSION=$(grep -E 'const AppVersion = "[^"]+"' cmd/acctg/main.go | head -n 1 | cut -d '"' -f 2)
if [ -z "${VERSION}" ]; then
    echo "ERROR: Could not determine release version" >&2
    exit 1
fi
echo "Target Version: v${VERSION}"

# 2. Automated Quality Gates
echo ""
echo "[1/5] Running automated quality gates..."

echo "  -> gofmt formatting check..."
DIFF=$(gofmt -s -d .)
if [ -n "${DIFF}" ]; then
    echo "ERROR: gofmt check failed! Please run 'gofmt -s -w .' before releasing."
    exit 1
fi

echo "  -> go vet static analysis..."
go vet ./...

echo "  -> go test full test suite..."
go test ./...
echo "  All quality gates PASSED."

# 3. Clean and prepare dist/
echo ""
echo "[2/5] Preparing output directory (dist/)..."
DIST_DIR="${ROOT_DIR}/dist"
# Preserve previous release artifacts; overwrite only generated files.
mkdir -p "${DIST_DIR}"

export CGO_ENABLED=0
LDFLAGS="-s -w"

# 4. Cross-compiling release binaries
echo ""
echo "[3/5] Cross-compiling standalone release binaries..."

echo "  -> Compiling Windows AMD64 (dist/acctg-windows-amd64.exe)..."
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "${DIST_DIR}/acctg-windows-amd64.exe" ./cmd/acctg
cp "${DIST_DIR}/acctg-windows-amd64.exe" "${DIST_DIR}/acctg.exe"

echo "  -> Compiling macOS Apple Silicon ARM64 (dist/acctg-mac-arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o "${DIST_DIR}/acctg-mac-arm64" ./cmd/acctg

echo "  -> Compiling macOS Intel AMD64 (dist/acctg-mac-amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "${DIST_DIR}/acctg-mac-amd64" ./cmd/acctg

echo "  -> Compiling Linux AMD64 (dist/acctg-linux-amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "${DIST_DIR}/acctg-linux-amd64" ./cmd/acctg

echo "  -> Compiling Linux ARM64 (dist/acctg-linux-arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o "${DIST_DIR}/acctg-linux-arm64" ./cmd/acctg

# Copy macOS & Linux Launchers and desktop integration
cp "${ROOT_DIR}/Launch-Tutor.command" "${DIST_DIR}/Launch-Tutor.command"
chmod +x "${DIST_DIR}/Launch-Tutor.command" 2>/dev/null || true
cp "${ROOT_DIR}/launch-tutor.sh" "${DIST_DIR}/launch-tutor.sh"
chmod +x "${DIST_DIR}/launch-tutor.sh" 2>/dev/null || true
cp "${ROOT_DIR}/accounttutor.desktop" "${DIST_DIR}/accounttutor.desktop"

cp "${ROOT_DIR}/docs/RELEASE-NOTES.md" "${DIST_DIR}/RELEASE-NOTES.md"
# 5. Package release bundles
echo ""
echo "[4/5] Creating distribution archive packages..."

# Windows zip
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-windows-amd64.zip" acctg.exe RELEASE-NOTES.md -j "${ROOT_DIR}/README.md" "${ROOT_DIR}/OVERVIEW.md")

# macOS Apple Silicon zip
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-macos-arm64.zip" acctg-mac-arm64 Launch-Tutor.command RELEASE-NOTES.md -j "${ROOT_DIR}/README.md")

# macOS Intel zip
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-macos-amd64.zip" acctg-mac-amd64 Launch-Tutor.command RELEASE-NOTES.md -j "${ROOT_DIR}/README.md")

# macOS Universal Classmate bundle
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-macos-classmate-bundle.zip" acctg-mac-arm64 acctg-mac-amd64 Launch-Tutor.command RELEASE-NOTES.md -j "${ROOT_DIR}/README.md")

# Linux AMD64 tar.gz & zip
(cd "${DIST_DIR}" && tar -czf "acctg-v${VERSION}-linux-amd64.tar.gz" acctg-linux-amd64 launch-tutor.sh accounttutor.desktop RELEASE-NOTES.md -C "${ROOT_DIR}" README.md OVERVIEW.md)
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-linux-amd64.zip" acctg-linux-amd64 launch-tutor.sh accounttutor.desktop RELEASE-NOTES.md -j "${ROOT_DIR}/README.md" "${ROOT_DIR}/OVERVIEW.md")

# Linux ARM64 tar.gz & zip
(cd "${DIST_DIR}" && tar -czf "acctg-v${VERSION}-linux-arm64.tar.gz" acctg-linux-arm64 launch-tutor.sh accounttutor.desktop RELEASE-NOTES.md -C "${ROOT_DIR}" README.md OVERVIEW.md)
(cd "${DIST_DIR}" && zip -q "acctg-v${VERSION}-linux-arm64.zip" acctg-linux-arm64 launch-tutor.sh accounttutor.desktop RELEASE-NOTES.md -j "${ROOT_DIR}/README.md" "${ROOT_DIR}/OVERVIEW.md")

# 6. Checksums
echo ""
echo "[5/5] Generating SHA256 checksums..."
cd "${DIST_DIR}"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum acctg.exe acctg-windows-amd64.exe acctg-mac-* acctg-linux-* acctg-v"${VERSION}"-* > checksums.sha256
elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 acctg.exe acctg-windows-amd64.exe acctg-mac-* acctg-linux-* acctg-v"${VERSION}"-* > checksums.sha256
fi

echo "=================================================================="
echo " Build & Packaging Complete! Release Artifacts in dist/:"
echo "=================================================================="
ls -lh "${DIST_DIR}"
