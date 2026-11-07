#!/usr/bin/env bash
# Start the WhatsApp bridge. Keep this window open while you use your AI app.
#   scripts/start.sh                 read-only (default, safest)
#   scripts/start.sh --allow-send    the AI may also send (it must ask you first)
#   --only 85291234567,...           limit sending to these numbers / groups
#                                    (default: the list in allowed-recipients.txt)
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/whatsapp-bridge/whatsapp-bridge"
[[ -x "$BIN" ]] || { echo "Not installed yet. Double-click Install.command first."; exit 1; }

export WHATSAPP_ALLOW_SEND="${WHATSAPP_ALLOW_SEND:-0}"
ONLY=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --allow-send) export WHATSAPP_ALLOW_SEND=1 ;;
    --only) shift; ONLY="${1:-}" ;;
    --log-messages) export WHATSAPP_LOG_MESSAGES=1 ;;
    -h|--help) sed -n '2,7p' "$0"; exit 0 ;;
    *) echo "Unknown option: $1 (try --help)"; exit 1 ;;
  esac
  shift
done

if [[ -z "$ONLY" && -f "$ROOT/allowed-recipients.txt" ]]; then
  ONLY="$(grep -v '^[[:space:]]*#' "$ROOT/allowed-recipients.txt" | tr -d ' \t\r+' | grep -v '^$' | paste -sd, - || true)"
fi
[[ -n "$ONLY" ]] && export WHATSAPP_SEND_ALLOWLIST="$ONLY"

printf '\033]0;WhatsApp bridge — keep this window open\007'
cd "$ROOT/whatsapp-bridge"
exec "$BIN"   # full path, so uninstall/update can stop exactly this bridge
