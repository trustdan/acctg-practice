#!/usr/bin/env bash
# AccountTutor 9000 — Linux Launcher & Terminal Wrapper
# Detects CPU architecture, locates native binary (or Wine fallback), and
# launches AccountTutor 9000 in an interactive terminal.

set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

# 1. Detect architecture
ARCH="$(uname -m)"
BIN=""

# A git pull updates source, not ignored executables. Always build a source
# checkout first so an older release binary cannot shadow the current code.
if [ -f "./go.mod" ] && [ -f "./cmd/acctg/main.go" ]; then
    if ! command -v go >/dev/null 2>&1; then
        echo "This is a source checkout; Go is required to build the current code." >&2
        echo "Install the Go version declared in go.mod, then run this launcher again." >&2
        echo "Alternatively, extract the latest Linux release ZIP into a separate folder." >&2
        exit 1
    fi
    echo "Building AccountTutor from the current source checkout..." >&2
    HOST_ARCH="$(go env GOHOSTARCH)"
    GOOS=linux GOARCH="$HOST_ARCH" CGO_ENABLED=0 go build -o acctg ./cmd/acctg
    BIN="./acctg"
elif [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
    # 64-bit x86 Linux
    for candidate in "./acctg-linux-amd64" "./dist/acctg-linux-amd64" "./bin/acctg-linux-amd64" "./acctg" "./dist/acctg"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    # 64-bit ARM Linux (Raspberry Pi 4/5, AWS Graviton, Asahi Linux, Chromebooks)
    for candidate in "./acctg-linux-arm64" "./dist/acctg-linux-arm64" "./bin/acctg-linux-arm64" "./acctg" "./dist/acctg"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
fi

# Fallback: a generic native binary, never a different architecture's release.
if [ -z "$BIN" ]; then
    for candidate in "./acctg" "./dist/acctg"; do
        if [ -f "$candidate" ]; then
            BIN="$candidate"
            break
        fi
    done
fi

# Fallback: check if Wine is available to run Windows binary
if [ -z "$BIN" ]; then
    for exe in "./acctg.exe" "./dist/acctg.exe" "./bin/acctg.exe"; do
        if [ -f "$exe" ] && command -v wine >/dev/null 2>&1; then
            echo "=================================================================="
            echo "Native Linux binary not found, but Wine and acctg.exe detected."
            echo "Launching via Wine backup..."
            echo "=================================================================="
            exec wine "$exe" "$@"
        fi
    done
fi

# If still not found, provide actionable guidance
if [ -z "$BIN" ]; then
    echo "=================================================================="
    echo "AccountTutor 9000 — Binary Not Found"
    echo "=================================================================="
    echo "Could not find a native Linux binary in: $DIR"
    echo "Expected one of: acctg-linux-amd64 (x86_64), acctg-linux-arm64 (ARM64), or acctg"
    echo ""
    echo "Options to get started:"
    echo "1. Download the pre-compiled binary for Linux from releases:"
    echo "   - x86_64: acctg-linux-amd64"
    echo "   - ARM64:  acctg-linux-arm64"
    echo ""
    echo "2. Build natively from a full source checkout (Go version from go.mod):"
    echo "   go build -o acctg ./cmd/acctg"
    echo ""
    echo "3. Run the Windows binary using Wine (backup):"
    echo "   sudo apt install wine  # (Ubuntu/Debian)"
    echo "   wine ./acctg.exe"
    echo "=================================================================="
    if [ -t 0 ]; then read -r -p "Press [Enter] to exit..." dummy; fi
    exit 1
fi

chmod +x "$BIN" 2>/dev/null || true
LAUNCHER="$DIR/$(basename "${BASH_SOURCE[0]}")"

# If not running in a terminal (e.g. clicked from a GUI file manager), spawn a terminal
if [ ! -t 0 ] || [ ! -t 1 ]; then
    for term in x-terminal-emulator gnome-terminal konsole xfce4-terminal alacritty kitty foot terminator tilix lxterminal urxvt rxvt xterm; do
        if command -v "$term" >/dev/null 2>&1; then
            case "$term" in
                gnome-terminal|xfce4-terminal|tilix|lxterminal)
                    exec "$term" -- bash "$LAUNCHER" "$@"
                    ;;
                konsole)
                    exec "$term" -e bash "$LAUNCHER" "$@"
                    ;;
                alacritty|kitty|foot|terminator)
                    exec "$term" -e bash "$LAUNCHER" "$@"
                    ;;
                x-terminal-emulator|xterm|urxvt|rxvt)
                    exec "$term" -e bash "$LAUNCHER" "$@"
                    ;;
            esac
        fi
    done
fi

# Run the binary
exec "$BIN" "$@"
