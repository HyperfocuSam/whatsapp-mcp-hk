[繁體中文](README.md) · [简体中文](README.zh-CN.md) · **English**

# WhatsApp MCP (HK edition)

Let your AI app (Claude, Cursor, Codex…) read and search your own WhatsApp chats, for example *"Summarise what my team said in the project group this week."*
If you turn it on, it can also send messages for you, but only after it shows you the message and you say yes.
Everything runs on your own computer. This is a safer, easier version of [lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp).

**You need:** a Mac · WhatsApp on your phone · one AI app (Claude Desktop is the easiest) · about 500 MB of free space · 10 minutes.
**You don't need:** any coding, an admin password, or other software.

## Before you start: read this

> ⚠️ **Please read before installing.**
>
> - **It is unofficial.** It connects the way WhatsApp Web does, using the open-source whatsmeow library. Meta did not make or approve it. There is a small risk that WhatsApp restricts your account. Use it at your own risk, and never use it for bulk or automated messaging.
> - **Your AI company sees what the AI reads.** Your messages are stored only on your computer. But when the AI reads a chat, that text is sent to your AI company (for example Anthropic or OpenAI). Don't use it for chats you would not share with them.
> - **Anyone who can use your computer account can read your messages.** Keep your computer locked with a password.
> - **Messages can try to trick the AI.** Anyone can send you a message with hidden instructions for the AI ("prompt injection"). That is why sending is OFF unless you turn it on. Always check before you say yes to a send.

## Install (Mac)

**1. Open Terminal.** Press `Cmd` + `Space`, type `Terminal`, press `Enter`.

**2. Copy this line, paste it into Terminal, press `Enter`:**

```bash
curl -fsSL https://raw.githubusercontent.com/HyperfocuSam/whatsapp-mcp-hk/main/get.sh | bash
```

It downloads everything into a folder called **WhatsApp-MCP** in your home folder and sets it up. The first time takes about 3–5 minutes.

**3. Answer two questions.** Press `Enter` for yes each time.
- *"Connect WhatsApp to: Claude Desktop…?"* It adds WhatsApp to the AI apps it finds on your Mac. It keeps a backup of their settings first.
- *"Start WhatsApp now and show the QR code?"*

**4. Scan the QR code.** On your phone: WhatsApp → **Settings** → **Linked devices** → **Link a device**. Point the phone at the QR code in Terminal.
When it says **Connected to WhatsApp**, you're linked. **Leave this window open.** You can minimise it.

**5. Restart your AI app.** Fully quit it (`Cmd` + `Q`), then open it again.

**6. Try it.** Ask your AI app:

```text
List my 5 most recent WhatsApp chats.
```

If you see your chats, it works. Older chats can take a few minutes to appear after linking.

## Every day

1. Open the **WhatsApp-MCP** folder (Finder → **Go** → **Home** → **WhatsApp-MCP**). Tip: drag it to the Finder sidebar so it's always there.
2. Double-click **Start WhatsApp**. A Terminal window opens. Leave it open (you can minimise it).
3. Use your AI app as usual.
4. To stop, close that window.

New messages are only saved while WhatsApp is started. Your phone does not need to stay next to the computer.

## Let the AI send messages (optional)

Sending is **off** unless you start it this way:

- Double-click **Start WhatsApp - can send** (instead of **Start WhatsApp**).
- When the AI wants to send something, **the WhatsApp window beeps and shows who it goes to and the exact message**. Switch to that window. Type `y` and press `Enter` to send. Just press `Enter` to cancel. If you don't answer within 2 minutes, nothing is sent.
- Only type `y` for messages you asked for. A message from someone else can never press this key for you.

**Limit who it can message (recommended).** Open **allowed-recipients.txt** (it opens in TextEdit). Write one phone number per line, with the country code and no spaces, for example `85291234567`. Save. From then on, the AI can only message those people. Status updates and channels are never allowed unless listed. For a group, ask your AI: *"What is the ID of my WhatsApp group 'Family'?"* and add that ID (it ends in `@g.us`).

**Sending files.** The AI can only send files from the **outbox** folder. Put the file there first, then ask the AI to send it (100 MB maximum).

## What's in the WhatsApp-MCP folder

| Item | What it is |
|---|---|
| **Start WhatsApp** | Daily start. The AI can read only. |
| **Start WhatsApp - can send** | Start with sending allowed (the AI must ask you first). |
| **allowed-recipients.txt** | Who the AI may message. Empty = anyone. |
| **outbox** | The only folder the AI can send files from. |
| **Update** | Get the newest version. Keeps your messages and link. |
| **Uninstall** | Removes it from your AI apps and can delete your saved messages. |
| **Install** | Runs the setup again (for example after moving the folder). |
| whatsapp-bridge/store | Your saved messages and WhatsApp login. Private. Never share it. |

Keep the folder in your home folder. Don't put it in Desktop, Documents or iCloud Drive, because those can sync to iCloud. If you move it, double-click **Install** again.

## Update

Double-click **Update**. Your linked phone, saved messages, outbox and allowed list stay. Then start WhatsApp again and restart your AI app.

## Uninstall

1. Double-click **Uninstall**. It removes WhatsApp from your AI apps. It asks before deleting your saved messages (type `DELETE` to confirm).
2. On your phone: WhatsApp → **Settings** → **Linked devices** → tap this computer → **Log out**.
3. Move the **WhatsApp-MCP** folder to the Trash.

## If something goes wrong

| What you see | What to do |
|---|---|
| The AI says WhatsApp is not started | Double-click **Start WhatsApp** and keep its window open. |
| The QR code expired | Close the window, double-click **Start WhatsApp** again. |
| The AI doesn't have WhatsApp at all | Fully quit the AI app (`Cmd` + `Q`) and reopen it. Still missing: double-click **Install** again. |
| Mac says it "cannot verify" a file you double-clicked | System Settings → **Privacy & Security** → scroll down → **Open Anyway**. This only happens if you downloaded the ZIP with a browser. |
| "Sending is turned off" | Close the window and use **Start WhatsApp - can send**. |
| The AI says the send was cancelled or not confirmed | Switch to the WhatsApp window and type `y` + `Enter` when it asks. |
| "port 8080 is already used by another program" | Another app is using that number. Restart the Mac, or ask whoever shared this with you. |
| "not in WHATSAPP_SEND_ALLOWLIST" | That person isn't in **allowed-recipients.txt**. Add them, or empty the list. |
| A file can't be sent | Put it in the **outbox** folder first. 100 MB maximum. |
| "The bridge refused the access key" | Close the WhatsApp window, start it again, then restart the AI app. |
| Old messages are missing | Messages are only saved while WhatsApp is started. Some history arrives when you first link. |
| It stopped working after a long break | Linked devices log out after about 20 days without your phone online. Start WhatsApp and scan a new QR code. |
| Something else | Send the file `.tools/install.log` (in the folder) to whoever shared this with you, or [open an issue](https://github.com/HyperfocuSam/whatsapp-mcp-hk/issues). Never send your `store` folder. |

## Questions

**Can the person who shared this see my messages?** No. Nothing is sent to the maintainer. There is no tracking. Your messages go only to WhatsApp (as usual) and to your AI company when the AI reads them.

**What does the installer download?** Go (from go.dev) and uv (from GitHub), each checked against its official checksum. It also downloads code libraries from the official Go and Python package servers. Everything goes inside the WhatsApp-MCP folder. Nothing is installed system-wide, and it never asks for your password.

**Does it work with two WhatsApp numbers?** Yes, with a second copy in another folder. See "For developers" below.

**Windows?** There is no easy installer for Windows yet. A developer can set it up by hand (see below).

## For developers

<details>
<summary>Manual install, other apps, settings, how it works</summary>

### Manual install (macOS / Linux)

```bash
git clone https://github.com/HyperfocuSam/whatsapp-mcp-hk.git ~/WhatsApp-MCP
cd ~/WhatsApp-MCP && bash scripts/install.sh
```

`scripts/install.sh` uses your own `go` (1.26+) and `uv` if present, otherwise private copies in `.tools/`. The bridge is pure Go (`CGO_ENABLED=0`, modernc SQLite), so no C compiler is needed. Linux works the same way.

### Connect an MCP client by hand

The MCP server runs over stdio. Command: `<folder>/whatsapp-mcp-server/.venv/bin/python`. Argument: `<folder>/whatsapp-mcp-server/main.py`.

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "/Users/you/WhatsApp-MCP/whatsapp-mcp-server/.venv/bin/python",
      "args": ["/Users/you/WhatsApp-MCP/whatsapp-mcp-server/main.py"]
    }
  }
}
```

- Claude Desktop: `~/Library/Application Support/Claude/claude_desktop_config.json`. Cursor: `~/.cursor/mcp.json`.
- Claude Code: `claude mcp add --scope user whatsapp -- <python> <main.py>`
- Codex (`~/.codex/config.toml`):

```toml
[mcp_servers.whatsapp]
command = "/Users/you/WhatsApp-MCP/whatsapp-mcp-server/.venv/bin/python"
args = ["/Users/you/WhatsApp-MCP/whatsapp-mcp-server/main.py"]
```

`whatsapp-mcp-server/connect_agents.py add|remove|status` does this automatically, with backups.

### Settings (environment variables for the bridge)

| Variable | Meaning |
|---|---|
| `WHATSAPP_ALLOW_SEND=1` | Allow sending (default off). `scripts/start.sh --allow-send` sets it. |
| `WHATSAPP_SEND_ALLOWLIST` | Comma-separated numbers / JIDs. `start.sh` fills it from `allowed-recipients.txt`, or from `--only a,b`. |
| `WHATSAPP_SEND_DIR` | Send folder (default `outbox/`). The MCP server reads it from `store/bridge.json`. |
| `WHATSAPP_BRIDGE_PORT` | Local API port (default 8080). Use it for a second copy with another number. |
| `WHATSAPP_LOG_MESSAGES=1` | Print message text to the terminal (default off). |
| `WHATSAPP_CONFIRM_SENDS=0` | Don't ask for `y` in the bridge window before each send (default: ask). |

### How it works

- `whatsapp-bridge/` (Go, [whatsmeow](https://github.com/tulir/whatsmeow)) is a WhatsApp linked device. It stores messages in `store/messages.db` and serves a local API (127.0.0.1, token in `store/api-token`) for sending and media downloads.
- `whatsapp-mcp-server/` (Python, MCP) reads `messages.db` read-only and calls the bridge API.
- Security design and threat model: [SECURITY.md](SECURITY.md).
- Tests: `cd whatsapp-bridge && go test ./...` and `cd whatsapp-mcp-server && uv run python -m unittest test_security`.

### Windows (manual, untested)

Install Go and uv. Then:

```powershell
cd whatsapp-bridge; $env:CGO_ENABLED="0"; go build -o whatsapp-bridge.exe .
cd ..\whatsapp-mcp-server; uv sync
cd ..\whatsapp-bridge; .\whatsapp-bridge.exe
```

For sending: `$env:WHATSAPP_ALLOW_SEND="1"` before starting. Connect your client with `whatsapp-mcp-server\.venv\Scripts\python.exe` and `main.py`.

### Tools

`search_contacts`, `list_chats`, `get_chat`, `list_messages`, `get_message_context`, `get_direct_chat_by_contact`, `get_contact_chats`, `get_last_interaction`, `download_media`, and the send tools `send_message`, `send_file` and `send_audio_message` (only when sending is on; voice notes need ffmpeg).

</details>

## Credits

- Original project: [lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp) by Luke Harries (MIT License).
- WhatsApp library: [whatsmeow](https://github.com/tulir/whatsmeow) by Tulir Asokan.
- This edition is maintained by Sam Wong ([HyperfocuSam](https://github.com/HyperfocuSam)).

## License

MIT. See [LICENSE](LICENSE).
