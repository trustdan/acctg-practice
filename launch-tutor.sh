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

if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
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

# Fallback: check any generic Linux binary
if [ -z "$BIN" ]; then
    for candidate in "./acctg" "./dist/acctg" "./acctg-linux-amd64" "./acctg-linux-arm64" "./dist/acctg-linux-amd64" "./dist/acctg-linux-arm64"; do
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

# Fallback: build from source if Go is installed
if [ -z "$BIN" ] && command -v go >/dev/null 2>&1; then
    echo "=================================================================="
    echo "Pre-compiled binary not found. Compiling natively with Go..."
    echo "=================================================================="
    go build -o acctg ./cmd/acctg
    BIN="./acctg"
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
    echo "2. Build natively from source (requires Go 1.24+):"
    echo "   go build -o acctg ./cmd/acctg"
    echo ""
    echo "3. Run the Windows binary using Wine (backup):"
    echo "   sudo apt install wine  # (Ubuntu/Debian)"
    echo "   wine ./acctg.exe"
    echo "=================================================================="
    read -p "Press [Enter] to exit..." dummy
    exit 1
fi

chmod +x "$BIN" 2>/dev/null || true

# If not running in a terminal (e.g. clicked from a GUI file manager), spawn a terminal
if [ ! -t 0 ] || [ ! -t 1 ]; then
    for term in x-terminal-emulator gnome-terminal konsole xfce4-terminal alacritty kitty foot terminator tilix lxterminal urxvt rxvt xterm; do
        if command -v "$term" >/dev/null 2>&1; then
            case "$term" in
                gnome-terminal|xfce4-terminal|tilix|lxterminal)
                    exec "$term" -- "$DIR/$0" "$@"
                    ;;
                konsole)
                    exec "$term" -e "$DIR/$0" "$@"
                    ;;
                alacritty|kitty|foot|terminator)
                    exec "$term" -e "$DIR/$0" "$@"
                    ;;
                x-terminal-emulator|xterm|urxvt|rxvt)
                    exec "$term" -e "$DIR/$0" "$@"
                    ;;
            esac
        fi
    done
fi

# Run the binary
exec "$BIN" "$@"
