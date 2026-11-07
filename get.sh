#!/usr/bin/env bash
# One-line install for WhatsApp MCP (HK edition), macOS / Linux:
#   curl -fsSL https://raw.githubusercontent.com/HyperfocuSam/whatsapp-mcp-hk/main/get.sh | bash
# Downloads a fixed, released version into ~/WhatsApp-MCP and runs the installer.
# If it is already there, it runs the updater instead (your data is kept).
#
# Everything is inside main(), so a half-downloaded copy of this script does nothing.
main() {
  set -euo pipefail
  local REPO="HyperfocuSam/whatsapp-mcp-hk"
  local VERSION="v1.0.0"            # the released version this installer gets
  local DEST="$HOME/WhatsApp-MCP"   # not Desktop/Documents: those may sync to iCloud

  if [[ -d "$DEST" ]]; then
    echo "WhatsApp MCP is already in $DEST — updating it (your data is kept)."
    bash "$DEST/scripts/update.sh" < /dev/tty
    return
  fi

  local tmp
  tmp="$(mktemp -d)"
  echo "Downloading WhatsApp MCP $VERSION…"
  if [[ -n "${WA_MCP_ZIP:-}" ]]; then
    cp "$WA_MCP_ZIP" "$tmp/src.zip"           # local testing only
  else
    curl -fL --retry 3 --proto '=https' --tlsv1.2 -s -o "$tmp/src.zip" \
      "https://github.com/$REPO/archive/refs/tags/$VERSION.zip"
  fi
  unzip -q "$tmp/src.zip" -d "$tmp/x"
  local src
  src="$(find "$tmp/x" -mindepth 1 -maxdepth 1 -type d | head -1)"
  mv "$src" "$DEST"
  rm -rf "$tmp"
  chmod +x "$DEST"/*.command "$DEST"/scripts/*.sh
  echo "Saved in $DEST"
  if command -v open >/dev/null 2>&1; then open "$DEST" || true; fi   # show the folder in Finder
  bash "$DEST/scripts/install.sh" < /dev/tty
}
main "$@"
