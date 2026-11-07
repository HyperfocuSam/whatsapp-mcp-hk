#!/usr/bin/env bash
# WhatsApp MCP (HK edition) — installer for macOS and Linux.
# Double-click "Install.command" (Mac) or run: bash scripts/install.sh
#
# What it does, all inside this folder unless it says otherwise:
#   1. gets Go and uv (downloads official copies into .tools/ if missing; checksums verified)
#   2. builds the WhatsApp bridge
#   3. prepares the Python MCP server (its own Python, in .tools/)
#   4. asks before connecting it to the AI apps it finds (settings backed up first)
#   5. offers to start WhatsApp so you can scan the QR code
# It needs no admin password and installs nothing system-wide.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TOOLS="$ROOT/.tools"
BRIDGE="$ROOT/whatsapp-bridge"
SERVER="$ROOT/whatsapp-mcp-server"

GO_VERSION="1.27.3"
UV_VERSION="0.13.7"
GO_MIN="1.26"   # whatsapp-bridge/go.mod

say()  { printf '\n\033[1m%s\033[0m\n' "$*"; }
ok()   { printf '  ✓ %s\n' "$*"; }
info() { printf '    %s\n' "$*"; }
fail() { printf '\n  ✗ %s\n\n' "$*"; exit 1; }

case "$(uname -s)-$(uname -m)" in
  Darwin-arm64)  GO_PKG=darwin-arm64; GO_SUM=dd3908cf20959d798c1e06ebe5b72213f0505ce01373b40ce3110975a6cd99c2
                 UV_PKG=aarch64-apple-darwin; UV_SUM=81b811b57e4624df7ed3aac148e0359100d9962f3d9dce00aec533dc2e086731 ;;
  Darwin-x86_64) GO_PKG=darwin-amd64; GO_SUM=b6f8e931e4ab497fcde6da8ae4499a58f9d258e55c8a80824a516eae5950a750
                 UV_PKG=x86_64-apple-darwin; UV_SUM=20c47e4c1f630b8abaff8c291fa9b9aa300fef92c4fbff95adf198b55287cda7 ;;
  Linux-x86_64)  GO_PKG=linux-amd64; GO_SUM=65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
                 UV_PKG=x86_64-unknown-linux-gnu; UV_SUM=c5221ff41b7892bfc223543883600125f6a47fe845093e0f5bdff4f359196d83 ;;
  Linux-aarch64|Linux-arm64)
                 GO_PKG=linux-arm64; GO_SUM=a06e99309dbcbdca1a8ccd2282c3ee146e55930886e26b4f1193269587165f49
                 UV_PKG=aarch64-unknown-linux-gnu; UV_SUM=71279f8d3c67241e807b9d8db872d1bf541579fb9f08b7d1eb1682196d1e57dc ;;
  *) fail "This computer type ($(uname -s) $(uname -m)) is not supported by the installer. See README → Windows / manual install." ;;
esac

sha256() { if command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | awk '{print $1}'; else sha256sum "$1" | awk '{print $1}'; fi; }

# download URL EXPECTED_SHA256 OUTPUT
download() {
  curl -fL --retry 3 --proto '=https' --tlsv1.2 -s -o "$3" "$1" || fail "Download failed: $1. Check your internet connection and try again."
  local got; got="$(sha256 "$3")"
  if [[ "$got" != "$2" ]]; then rm -f "$3"; fail "Downloaded file did not match its official checksum. Nothing was installed. Try again later."; fi
}

version_ge() { [[ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -1)" == "$2" ]]; }

mkdir -p "$TOOLS"
echo
echo "  WhatsApp MCP — installer"
echo "  Folder: $ROOT"

# ---------- 1. Go ----------
say "1/5  Go (builds the WhatsApp bridge)"
GO=""
if command -v go >/dev/null 2>&1; then
  v="$(go env GOVERSION 2>/dev/null | sed 's/^go//')"
  if version_ge "$v" "$GO_MIN"; then GO="$(command -v go)"; ok "Using the Go already on this computer ($v)"; fi
fi
if [[ -z "$GO" ]]; then
  if [[ ! -x "$TOOLS/go/bin/go" || "$("$TOOLS/go/bin/go" env GOVERSION)" != "go$GO_VERSION" ]]; then
    info "Downloading Go $GO_VERSION from go.dev (about 70 MB)…"
    tmp="$TOOLS/go.tar.gz"
    download "https://go.dev/dl/go$GO_VERSION.$GO_PKG.tar.gz" "$GO_SUM" "$tmp"
    rm -rf "$TOOLS/go" && tar -xzf "$tmp" -C "$TOOLS" && rm -f "$tmp"
  fi
  GO="$TOOLS/go/bin/go"
  ok "Go $GO_VERSION (private copy in .tools/)"
fi
export GOPATH="$TOOLS/gopath" GOMODCACHE="$TOOLS/gopath/pkg/mod" GOCACHE="$TOOLS/gocache"
export GOFLAGS="-modcacherw" GOTOOLCHAIN="local" CGO_ENABLED=0

# ---------- 2. uv ----------
say "2/5  uv (runs the Python part)"
if command -v uv >/dev/null 2>&1; then
  UV="$(command -v uv)"; ok "Using the uv already on this computer"
else
  if [[ ! -x "$TOOLS/uv/uv" ]]; then
    info "Downloading uv $UV_VERSION from GitHub (about 20 MB)…"
    tmp="$TOOLS/uv.tar.gz"
    download "https://github.com/astral-sh/uv/releases/download/$UV_VERSION/uv-$UV_PKG.tar.gz" "$UV_SUM" "$tmp"
    rm -rf "$TOOLS/uv" && mkdir -p "$TOOLS/uv" && tar -xzf "$tmp" -C "$TOOLS/uv" --strip-components=1 && rm -f "$tmp"
  fi
  UV="$TOOLS/uv/uv"; ok "uv $UV_VERSION (private copy in .tools/)"
fi
export UV_PYTHON_INSTALL_DIR="$TOOLS/python" UV_CACHE_DIR="$TOOLS/uv-cache"

# ---------- 3. build ----------
say "3/5  Building (first time takes 1–3 minutes)"
LOG="$TOOLS/install.log"
: > "$LOG"
if ! ( cd "$BRIDGE" && "$GO" build -trimpath -o whatsapp-bridge . ) >>"$LOG" 2>&1; then
  tail -20 "$LOG"; fail "The bridge did not build. Send the file .tools/install.log to whoever shared this with you."
fi
# Build leftovers are large (~800 MB) and not needed to run: remove them.
chmod -R u+w "$TOOLS/gopath" "$TOOLS/gocache" 2>/dev/null || true
rm -rf "$TOOLS/gopath" "$TOOLS/gocache"
ok "WhatsApp bridge built"
if ! ( cd "$SERVER" && "$UV" sync --frozen --quiet ) >>"$LOG" 2>&1; then
  tail -20 "$LOG"; fail "Python setup failed. Send the file .tools/install.log to whoever shared this with you."
fi
rm -rf "$TOOLS/uv-cache"
ok "MCP server ready"
mkdir -p "$ROOT/outbox" && chmod 700 "$ROOT/outbox"
if [[ ! -f "$ROOT/allowed-recipients.txt" ]]; then
  cat > "$ROOT/allowed-recipients.txt" <<'EOF'
# Who the AI may message when you start "Start WhatsApp - can send".
# One phone number (with country code, e.g. 85291234567) or group ID per line.
# Lines starting with # are ignored. If no numbers are listed, anyone can be messaged
# (Status / broadcast lists / channels only if listed here).
# Every send still needs your "y" in the WhatsApp window.
EOF
fi
ok "Send folder: outbox/   ·   Allowed people: allowed-recipients.txt"

# ---------- 4. connect ----------
say "4/5  Connect to your AI apps"
PY="$SERVER/.venv/bin/python"
"$PY" "$SERVER/connect_agents.py" add || info "Could not connect automatically. See README → 'Connect by hand'."

# ---------- 5. start ----------
say "5/5  Link your phone"
if [[ -f "$BRIDGE/store/whatsapp.db" ]]; then
  info "A phone is already linked (or was linked before)."
fi
printf '  Start WhatsApp now and show the QR code? [Y/n] '
read -r answer || answer=""
if [[ -z "$answer" || "$answer" =~ ^[Yy] ]]; then
  echo
  echo "  On your phone: WhatsApp → Settings → Linked devices → Link a device, then scan."
  echo "  Keep this window open while you use your AI app. Close it to stop."
  exec bash "$ROOT/scripts/start.sh"
fi
echo
echo "  Done. Next time, double-click \"Start WhatsApp.command\"."
