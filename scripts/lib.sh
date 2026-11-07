# Shared helpers for the scripts in this folder. Sourced, not run.

# stop_bridge ROOT — stops THIS copy's bridge (by the PID it wrote), nothing else.
# Returns 1 if it is still running afterwards.
stop_bridge() {
  local root="$1" pidfile="$1/whatsapp-bridge/store/bridge.pid" pid cmd
  [[ -f "$pidfile" ]] || return 0
  pid="$(tr -cd '0-9' < "$pidfile")"
  [[ -n "$pid" ]] || return 0
  cmd="$(ps -p "$pid" -o comm= 2>/dev/null || true)"
  [[ "$cmd" == *whatsapp-bridge* ]] || return 0      # not ours (or already gone)
  kill "$pid" 2>/dev/null || true
  for _ in 1 2 3 4 5 6 7 8 9 10; do
    kill -0 "$pid" 2>/dev/null || return 0
    sleep 0.5
  done
  kill -9 "$pid" 2>/dev/null || true
  sleep 0.5
  ! kill -0 "$pid" 2>/dev/null
}
