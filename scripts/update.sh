#!/usr/bin/env bash
# Update to the newest released version. Keeps your linked phone, saved
# messages, outbox/ and allowed-recipients.txt. Then rebuilds.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO="HyperfocuSam/whatsapp-mcp-hk"
# shellcheck source=lib.sh
source "$ROOT/scripts/lib.sh"

echo
echo "  WhatsApp MCP — update"
if [[ -d "$ROOT/.git" ]] && command -v git >/dev/null 2>&1; then
  git -C "$ROOT" pull --ff-only
else
  tag="$(curl -fsSL --proto '=https' --tlsv1.2 "https://api.github.com/repos/$REPO/releases/latest" \
        | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)"
  [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "  Could not find the newest version. Try again later."; exit 1; }
  echo "  Newest version: $tag"
  tmp="$(mktemp -d)"
  curl -fL --retry 3 --proto '=https' --tlsv1.2 -s -o "$tmp/src.zip" \
    "https://github.com/$REPO/archive/refs/tags/$tag.zip"
  unzip -q "$tmp/src.zip" -d "$tmp/x"
  src="$(find "$tmp/x" -mindepth 1 -maxdepth 1 -type d | head -1)"
  if ! stop_bridge "$ROOT"; then echo "  Could not stop the running WhatsApp window. Close it and try again."; exit 1; fi
  rsync -a --exclude '.tools/' --exclude 'whatsapp-bridge/store/' --exclude 'outbox/' \
        --exclude 'allowed-recipients.txt' --exclude 'whatsapp-mcp-server/.venv/' "$src/" "$ROOT/"
  rm -rf "$tmp"
  chmod +x "$ROOT"/*.command "$ROOT"/scripts/*.sh
fi
stop_bridge "$ROOT" || true
bash "$ROOT/scripts/install.sh"
