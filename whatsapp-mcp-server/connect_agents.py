"""Connect (or disconnect) this WhatsApp MCP server to the AI apps on this computer.

    python connect_agents.py add      # asks first, then adds to every app it finds
    python connect_agents.py remove   # removes it again
    python connect_agents.py status   # shows what is connected
    add --yes                         # no questions

Every settings file is backed up (name.bak-YYYYMMDD-HHMMSS) before it is changed.
A settings file that can't be read safely is left alone.
"""
import json
import os
import platform
import shutil
import subprocess
import sys
import time
from pathlib import Path

NAME = "whatsapp"
HERE = Path(__file__).resolve().parent
PYTHON = str(Path(sys.executable))
MAIN = str(HERE / "main.py")
ENTRY = {"command": PYTHON, "args": [MAIN]}
HOME = Path.home()
MARK_START = "# >>> whatsapp-mcp-hk (added by connect_agents.py)"
MARK_END = "# <<< whatsapp-mcp-hk"


def claude_desktop_path():
    system = platform.system()
    if system == "Darwin":
        return HOME / "Library/Application Support/Claude/claude_desktop_config.json"
    if system == "Windows":
        return Path(os.environ.get("APPDATA", HOME)) / "Claude/claude_desktop_config.json"
    return HOME / ".config/Claude/claude_desktop_config.json"


def claude_code_cli():
    found = shutil.which("claude")
    if found:
        return found
    for p in (HOME / ".claude/local/claude", HOME / ".local/bin/claude"):
        if p.exists():
            return str(p)
    return None


ALT_NAME = "whatsapp-hk"   # used when the user already has a different "whatsapp" server


def is_ours(entry) -> bool:
    return isinstance(entry, dict) and MAIN in [str(a) for a in entry.get("args", [])]


def real_path(path: Path) -> Path:
    """Edit the file a symlink points to (dotfile setups), not the link itself."""
    return path.resolve() if path.is_symlink() else path


def backup(path: Path):
    if path.exists():
        dst = path.with_name(path.name + time.strftime(".bak-%Y%m%d-%H%M%S"))
        shutil.copy2(path, dst)
        return dst
    return None


def write_atomic(path: Path, text: str):
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_name(path.name + ".tmp")
    fd = os.open(tmp, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0o600)
    with os.fdopen(fd, "w", encoding="utf-8") as f:
        f.write(text)
    if path.exists():
        shutil.copymode(path, tmp)   # keep the user's file permissions
    os.replace(tmp, path)


# ---------- JSON apps (Claude Desktop, Cursor) ----------

def json_load(path: Path):
    if not path.exists() or not path.read_text(encoding="utf-8").strip():
        return {}
    data = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise ValueError("not a JSON object")
    return data


def json_status(path: Path):
    try:
        return any(is_ours(e) for e in json_load(path).get("mcpServers", {}).values())
    except (ValueError, OSError, AttributeError):
        return False


def json_change(path: Path, add: bool):
    path = real_path(path)
    try:
        data = json_load(path)
        servers = data.setdefault("mcpServers", {})
        if not isinstance(servers, dict):
            raise ValueError("mcpServers is not an object")
    except (ValueError, OSError) as e:
        return f"skipped: could not read {path} safely ({e}). Add it by hand (see README)."
    mine = [k for k, v in servers.items() if is_ours(v)]
    if add:
        if mine and all(servers[k] == ENTRY for k in mine):
            return "already connected"
        for k in mine:
            del servers[k]
        name = NAME if NAME not in servers else ALT_NAME
        servers[name] = ENTRY
    else:
        if not mine:
            return "was not connected"
        for k in mine:
            del servers[k]
    b = backup(path)
    write_atomic(path, json.dumps(data, indent=2, ensure_ascii=False) + "\n")
    return ("connected" if add else "removed") + (f" (backup: {b.name})" if b else "")


# ---------- Codex (TOML) ----------

def toml_str(s: str) -> str:
    return json.dumps(s, ensure_ascii=False)   # valid TOML basic string, keeps emoji/Chinese as-is


def codex_block(name: str):
    return "\n".join([MARK_START, f"[mcp_servers.{json.dumps(name)}]",
                      f"command = {toml_str(PYTHON)}", f"args = [{toml_str(MAIN)}]", MARK_END, ""])


def codex_strip(text: str):
    """Remove our marked block. Returns None if a start marker has no end marker."""
    out, skip = [], False
    for line in text.splitlines(keepends=True):
        if line.strip() == MARK_START:
            skip = True
            continue
        if skip and line.strip() == MARK_END:
            skip = False
            continue
        if not skip:
            out.append(line)
    return None if skip else "".join(out)


def codex_status(path: Path):
    return path.exists() and MARK_START in path.read_text(encoding="utf-8")


def codex_change(path: Path, add: bool):
    import tomllib
    path = real_path(path)
    text = path.read_text(encoding="utf-8") if path.exists() else ""
    stripped = codex_strip(text)
    if stripped is None:
        return f"skipped: the WhatsApp section in {path} looks edited by hand. Fix it by hand."
    try:
        servers = tomllib.loads(stripped).get("mcp_servers", {})
    except tomllib.TOMLDecodeError as e:
        return f"skipped: {path} is not valid TOML ({e}). Nothing changed."
    if not add and stripped == text:
        return "was not connected"
    new = stripped
    if add:
        name = NAME if NAME not in servers else ALT_NAME
        if name in servers:
            return f"skipped: {path} already has servers named {NAME} and {ALT_NAME}. Add it by hand."
        if new and not new.endswith("\n"):
            new += "\n"
        new += ("\n" if new else "") + codex_block(name)
    try:
        tomllib.loads(new)
    except tomllib.TOMLDecodeError as e:
        return f"skipped: the change would break {path} ({e}). Nothing changed."
    b = backup(path)
    write_atomic(path, new)
    return ("connected" if add else "removed") + (f" (backup: {b.name})" if b else "")


# ---------- Claude Code (CLI) ----------

def cc_run(cli, *args):
    return subprocess.run([cli, "mcp", *args], capture_output=True, text=True, timeout=60)


def cc_ours(cli, name) -> bool:
    try:
        r = cc_run(cli, "get", name)
        return r.returncode == 0 and MAIN in r.stdout
    except Exception:
        return False


def cc_exists(cli, name) -> bool:
    try:
        return cc_run(cli, "get", name).returncode == 0
    except Exception:
        return False


def cc_status(cli):
    return cc_ours(cli, NAME) or cc_ours(cli, ALT_NAME)


def cc_change(cli, add: bool):
    try:
        removed = False
        for name in (NAME, ALT_NAME):
            if cc_ours(cli, name):
                removed = cc_run(cli, "remove", "--scope", "user", name).returncode == 0 or removed
        if not add:
            return "removed" if removed else "was not connected"
        name = NAME if not cc_exists(cli, NAME) else ALT_NAME
        if cc_exists(cli, name):
            return f"skipped: Claude Code already has servers named {NAME} and {ALT_NAME}."
        r = cc_run(cli, "add", "--scope", "user", name, "--", PYTHON, MAIN)
        return f"connected as '{name}'" if r.returncode == 0 else f"failed: {(r.stderr or r.stdout).strip()[:200]}"
    except Exception as e:
        return f"failed: {e}"


# ---------- main ----------

def find_apps():
    apps = []
    cd = claude_desktop_path()
    if cd.parent.exists():
        apps.append(("Claude Desktop", lambda a: json_change(cd, a), lambda: json_status(cd)))
    cur = HOME / ".cursor/mcp.json"
    if cur.parent.exists():
        apps.append(("Cursor", lambda a: json_change(cur, a), lambda: json_status(cur)))
    cx = HOME / ".codex/config.toml"
    if cx.parent.exists():
        apps.append(("Codex", lambda a: codex_change(cx, a), lambda: codex_status(cx)))
    cli = claude_code_cli()
    if cli:
        apps.append(("Claude Code", lambda a: cc_change(cli, a), lambda: cc_status(cli)))
    return apps


def main():
    action = sys.argv[1] if len(sys.argv) > 1 else "status"
    yes = "--yes" in sys.argv
    apps = find_apps()
    if not apps:
        print("  No supported AI app found (Claude Desktop, Claude Code, Cursor, Codex).")
        print("  Install one, then run this again. Other apps: see README, 'Other MCP apps'.")
        return
    if action == "status":
        for name, _, status in apps:
            print(f"  {name}: {'connected' if status() else 'not connected'}")
        return
    add = action == "add"
    names = ", ".join(n for n, _, _ in apps)
    if not yes:
        verb = "Connect WhatsApp to" if add else "Remove WhatsApp from"
        try:
            answer = input(f"  {verb}: {names}? [Y/n] ").strip().lower()
        except EOFError:
            answer = ""
        if answer not in ("", "y", "yes"):
            print("  Nothing changed.")
            return
    for name, change, _ in apps:
        print(f"  {name}: {change(add)}")
    if add:
        print("\n  Fully quit and reopen those apps (Mac: Cmd+Q) so they load WhatsApp.")


if __name__ == "__main__":
    main()
