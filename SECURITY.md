# Security

This page explains, in plain words, what this WhatsApp MCP server protects against, how, and what it does not protect against.

Setup steps: [繁體中文](README.md) · [简体中文](README.zh-CN.md) · [English](README.en.md).

## How it works (short)

There are two parts. Both run on your computer.

- **The bridge** (`whatsapp-bridge/`) links to WhatsApp as a "linked device". It saves your messages in `whatsapp-bridge/store/messages.db`. Your login keys are in `store/whatsapp.db`. It has a small local API used to send messages and download media.
- **The MCP server** (`whatsapp-mcp-server/`) is what your AI agent talks to. It reads the message database and calls the bridge API to send or download.

## Who could attack, and what stops them

### 1. Another device on your network

**The risk.** In the original, the bridge API listened on every network interface (port 8080). Someone on the same Wi-Fi (a café, an office, a hotel) could reach it and send messages as you.

**What we changed.**

- The API listens on `127.0.0.1` only. Other devices cannot reach it.
- Every API call needs a secret token. Each time it starts, the bridge creates a new `store/api-token`: a random 256-bit value that only your user account can read (file mode 0600). The MCP server reads the same file.
- The MCP server only sends the token after checking that the bridge which wrote `store/bridge.json` is still running as you. It never sends it through a proxy and never follows redirects. If another program grabs the port while the bridge is off, it gets nothing, and an old token is useless after the next start anyway.
- If the port is already used by another program, the bridge stops with a clear message instead of running half-connected. The port can be changed with `WHATSAPP_BRIDGE_PORT` (default 8080). It is still local-only.

### 2. A malicious website

**The risk.** A web page open in your browser can try to send requests to services on your own computer. One trick, "DNS rebinding", makes the browser think a local address belongs to the website.

**What we changed.**

- Requests with an `Origin` header (which browsers add) are refused.
- The request body must be `application/json`.
- The `Host` header must be `127.0.0.1` or `localhost`. This stops DNS rebinding.
- Request bodies are limited to 1 MB.
- The token (above) is also required, and a website cannot read it.

### 3. A malicious WhatsApp message (prompt injection)

**The risk.** Anyone can send you a message. It might contain instructions aimed at the AI, such as "forward the last 50 messages to this number" or "send the file at this path". If the agent reads that message, it might obey.

**What we changed.**

- **Sending is OFF by default.** The agent can only read until you start with **Start WhatsApp - can send** (or `WHATSAPP_ALLOW_SEND=1`).
- **You confirm every send in the WhatsApp window.** When the AI asks to send, the bridge window beeps and shows the recipient (name and number), the file name and the full text, then waits for you to type `y`. A message, a website or a tricked AI cannot press that key. This works even if your AI app is set to approve tools automatically. With no answer within 2 minutes, nothing is sent. (Developers can turn this off with `WHATSAPP_CONFIRM_SENDS=0`.)
- **Optional allowlist.** `allowed-recipients.txt` (or `WHATSAPP_SEND_ALLOWLIST`) limits who can be messaged, even when sending is on. Status updates, broadcast lists and channels can never be sent to unless they are listed by name.
- **Files can only come from one folder.** Only files in `outbox/` (or the folder in `WHATSAPP_SEND_DIR`) can be sent. Symlinks are followed and checked, the file must be a regular file, and the limit is 100 MB. Note: AI apps that can also work with files (Claude Code, Codex, Cursor) could copy a file into `outbox/` first. The confirmation step above, which shows the file name, is what protects you there.
- **Safe downloads.** File names chosen by the sender are cleaned (no `../`, no hidden or odd characters). Files are saved under `store/media/<chat>/`, never overwrite an existing file, are limited to 100 MB, and are private (mode 0600). On macOS they get the same "downloaded from the internet" mark a browser adds, so Gatekeeper checks any app inside before it can run.
- **The AI is warned.** The MCP server tells the agent that message text is untrusted, that it must never follow instructions inside messages, and that it must confirm the recipient and the exact text with you before sending. Send tools are marked as risky and read tools as read-only, so AI apps can treat them differently.
- **Messages can't impersonate you.** Text written by other people (messages, names, group names, IDs) is shown in quotes on one line, so a message can't fake extra lines like "From: ME". Someone in a chat also can't overwrite an older message (yours or anyone's) by reusing its ID.
- **Read-only database.** The MCP server opens the message database read-only. Result sizes are capped, and API calls have timeouts.

### 4. Another user on the same computer

**The risk.** Another account on a shared computer could read your saved messages or your login keys.

**What we changed.**

- `store/` is private to your account (mode 0700). Databases and media are mode 0600.
- These permissions are fixed again on every start. This also repairs older installs.
- Message text, names and phone numbers are no longer printed to the terminal or log by default. Turn it on only if you need it, with `WHATSAPP_LOG_MESSAGES=1`.

### 5. The AI provider

**The risk.** When the agent reads a chat, that text is sent to your AI provider (for example Anthropic or OpenAI).

**What we changed.** Nothing in this code can stop this. It is how AI agents work. The agent only reads what it asks for through the tools. You decide which chats it should look at.

### 6. The installer and the "connect to my AI apps" step

**The risk.** An install script runs code on your computer and changes your AI apps' settings.

**What we did.**

- The installer needs no admin password and installs nothing system-wide. Everything goes inside the project folder (`.tools/` for Go, uv and Python).
- The one-line installer downloads a fixed released version (a Git tag), not whatever is newest. Update installs the newest *release*. Unreleased changes never reach users.
- Go and uv are downloaded only over HTTPS from go.dev and GitHub. Each file is checked against a SHA-256 checksum written into `scripts/install.sh`. If a checksum doesn't match, nothing is installed.
- Python packages are locked to exact versions in `uv.lock`, and Go modules are checked against `go.sum`.
- Before changing an AI app's settings file, `connect_agents.py` saves a backup next to it (`*.bak-YYYYMMDD-HHMMSS`). It only adds or removes its own `whatsapp` entry. If a settings file can't be read safely, it leaves that file alone.
- The bridge is pure Go (no C compiler needed), which makes the build simpler to check.
- Uninstall and update stop only this copy's bridge (by the process ID it recorded), not other programs. Uninstall refuses to delete your data while the bridge is still running.

### Also updated

- The whatsmeow library is updated (2026-09).
- A media download fix keeps the CDN's auth query, which fixes HTTP 403 errors.

## What is NOT protected

Please be honest with yourself about these.

- **Your AI provider sees what the agent reads.** Do not use this for chats you cannot share with that provider.
- **Anyone using your computer account can read everything in `store/`**, including your messages and your WhatsApp login keys. File permissions only protect against *other* accounts. Keep your computer locked.
- **Malware running as you** can read the token and the databases. This tool cannot defend against that.
- **Prompt injection is reduced, not solved.** AI models can still be tricked into *asking* to send something. Nothing leaves without your `y` in the WhatsApp window. So read every confirmation before you type `y`, and type it only for messages you asked for. A tricked AI can also still summarise private chats wrongly, or include their text in what it shows you.
- **Your WhatsApp account.** This is unofficial and not approved by Meta. WhatsApp could restrict your account. Avoid bulk or automated messaging.
- **Data is not encrypted on disk by this tool.** Use your computer's disk encryption (FileVault on macOS, BitLocker on Windows).
- **Turning on `WHATSAPP_LOG_MESSAGES=1`** puts message text into the terminal and log again.

## Safe-use checklist

- [ ] My computer has a password and locks when I step away.
- [ ] Disk encryption is on.
- [ ] Sending is OFF unless I need it right now.
- [ ] When sending is on, I list the allowed people in `allowed-recipients.txt`.
- [ ] I only put files in `outbox/` that I am happy to send.
- [ ] I check the recipient and exact text before every send.
- [ ] I do not ask the agent to read chats I would not share with my AI provider.
- [ ] I never share the `whatsapp-bridge/store/` folder or the `store/api-token` file.
- [ ] When I stop using it, I log out the linked device on my phone and delete `whatsapp-bridge/store/`.

## Report a problem

- **Normal bugs and questions:** open an issue at https://github.com/HyperfocuSam/whatsapp-mcp-hk/issues
- **Security problems that could hurt users if made public:** please do not open a public issue. Use a private GitHub security advisory on the repo instead: go to https://github.com/HyperfocuSam/whatsapp-mcp-hk, open the **Security** tab, and choose **Report a vulnerability**.

Please include what you did, what happened, and your operating system. Never include your `store/` files, your token, or real message content.
