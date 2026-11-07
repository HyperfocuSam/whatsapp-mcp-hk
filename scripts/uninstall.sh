#!/usr/bin/env bash
# Disconnect WhatsApp from your AI apps and (if you agree) delete the saved messages.
set -uo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PY="$ROOT/whatsapp-mcp-server/.venv/bin/python"
# shellcheck source=lib.sh
source "$ROOT/scripts/lib.sh"

echo
echo "  WhatsApp MCP — uninstall"
echo
echo "  1) Removing it from your AI apps"
if [[ -x "$PY" ]]; then "$PY" "$ROOT/whatsapp-mcp-server/connect_agents.py" remove --yes
else echo "    (Python part not installed — nothing to remove)"; fi

echo
echo "  2) Saved messages and WhatsApp login on this computer"
if [[ -d "$ROOT/whatsapp-bridge/store" ]]; then
  printf '  Delete them now? This cannot be undone. Type DELETE to confirm: '
  read -r answer || answer=""
  if [[ "$answer" == "DELETE" ]]; then
    if stop_bridge "$ROOT"; then
      rm -rf "$ROOT/whatsapp-bridge/store" && echo "    Deleted."
    else
      echo "    WhatsApp is still running. Close its window, then run Uninstall again. Nothing deleted."
    fi
  else
    echo "    Kept."
  fi
else
  echo "    Nothing saved."
fi

echo
echo "  3) On your phone: WhatsApp → Settings → Linked devices → tap this computer → Log out."
echo "  4) Then you can move this whole folder to the Trash:"
echo "     $ROOT"
echo
