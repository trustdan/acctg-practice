#!/bin/sh
# AccountTutor 9000 — macOS Double-Click Launcher
# This script launches the terminal accounting drill inside Terminal.app on macOS.

DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

# Detect Mac architecture (arm64 for Apple Silicon M1-M4, x86_64 for Intel)
ARCH="$(uname -m)"
BIN=""

if [ "$ARCH" = "arm64" ]; then
    # Apple Silicon priority
    for candidate in "./acctg-mac-arm64" "./dist/acctg-mac-arm64" "./bin/acctg-mac-arm64" "./acctg-mac" "./dist/acctg-mac" "./acctg-mac-universal" "./acctg" "./dist/acctg"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
else
    # Intel Mac priority
    for candidate in "./acctg-mac-amd64" "./dist/acctg-mac-amd64" "./bin/acctg-mac-amd64" "./acctg-mac" "./dist/acctg-mac" "./acctg-mac-universal" "./acctg" "./dist/acctg"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
fi

# Fallback: check any native binary if preferred arch binary wasn't found
if [ -z "$BIN" ]; then
    for candidate in "./acctg-mac-arm64" "./acctg-mac-amd64" "./acctg-mac" "./acctg-mac-universal" "./acctg" "./dist/acctg-mac-arm64" "./dist/acctg-mac-amd64" "./dist/acctg-mac"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
fi

# Secondary Fallback: Wine execution if Windows binary exists and Wine is installed
if [ -z "$BIN" ]; then
    for exe in "./acctg.exe" "./dist/acctg.exe" "./bin/acctg.exe"; do
        if [ -f "$exe" ] && command -v wine >/dev/null 2>&1; then
            echo "=================================================================="
            echo "Native Mac binary not found, but Wine and acctg.exe detected."
            echo "Launching via Wine backup..."
            echo "=================================================================="
            exec wine "$exe" "$@"
        fi
    done
fi

# If still not found, display actionable setup guidance
if [ -z "$BIN" ]; then
    echo "=================================================================="
    echo "AccountTutor 9000 — Binary Not Found"
    echo "=================================================================="
    echo "Could not find a native macOS binary in: $DIR"
    echo "Expected one of: acctg-mac-arm64 (Apple Silicon), acctg-mac-amd64 (Intel), or acctg"
    echo ""
    echo "Options to get started:"
    echo "1. Download the pre-compiled binary for your Mac:"
    echo "   - Apple Silicon (M1/M2/M3/M4): acctg-mac-arm64"
    echo "   - Intel Mac:                   acctg-mac-amd64"
    echo ""
    echo "2. Build natively from source (requires Go):"
    echo "   go build -o acctg ./cmd/acctg"
    echo ""
    echo "3. Run the Windows binary using Wine (backup):"
    echo "   brew install --cask wine-stable"
    echo "   wine ./acctg.exe"
    echo "=================================================================="
    read -p "Press [Enter] to exit..." dummy
    exit 1
fi

# Ensure executable permissions and remove macOS Gatekeeper quarantine attribute
chmod +x "$BIN" 2>/dev/null
xattr -d com.apple.quarantine "$BIN" 2>/dev/null

# Launch the interactive terminal UI
exec "$BIN" "$@"
