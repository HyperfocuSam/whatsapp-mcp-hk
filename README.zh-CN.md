[繁體中文](README.md) · **简体中文** · [English](README.en.md)

# WhatsApp MCP（香港版）

让 Claude、Cursor、Codex 等 AI 应用读取并搜索你自己的 WhatsApp 聊天记录，例如 *"这周项目群里我的团队都说了什么？"*

如果开启此功能，它也可以代表你发送消息，但前提是它会先向你展示该消息，你同意后才会发出。

所有操作都在你自己的电脑上运行。这是一个更安全、更简单的 [lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp) 版本。

**你需要：**
- 一台 Mac
- 手机上的 WhatsApp
- 一个 AI 应用（Claude Desktop 最简单）
- 约 500 MB 的可用空间
- 10 分钟

**你不需要：**
- 任何编程知识
- 管理员密码
- 其他软件

## 开始前必读

> ⚠️ **安装前请阅读。**
>
> - **这是非官方版本。** 它像 WhatsApp Web 一样工作，使用了开源的 whatsmeow 库。Meta 没有制作或批准这个版本。WhatsApp 有很小的可能限制你的账号。使用此工具的风险由你自行承担。不要用它进行批量或自动化消息发送。
> - **你的 AI 公司会看到 AI 读到的内容。** 你的消息只存储在你的电脑上。但当 AI 阅读聊天记录时，那段文本会被发送给你的 AI 公司（例如 Anthropic 或 OpenAI）。不愿与他们分享的聊天，就不要让 AI 读取。
> - **任何能访问你电脑账号的人都能读取你的消息。** 请给你的电脑设置密码并锁定它。
> - **消息可能会试图欺骗 AI。** 任何人都可以给你发送带有隐藏指令的消息，让 AI 执行特定操作（这叫“提示注入”）。这就是为什么默认情况下**发送功能处于关闭状态**。在你同意发送之前，请务必仔细检查。

## 安装（Mac）

**1. 打开终端。** 按 `Cmd` + `Space`，输入 `Terminal`，然后按 `Enter`。

**2. 复制下面这行代码，粘贴到终端中，再按 `Enter`：**

```bash
curl -fsSL https://raw.githubusercontent.com/HyperfocuSam/whatsapp-mcp-hk/main/get.sh | bash
```

它会将所有文件下载到你主文件夹中的 **WhatsApp-MCP** 文件夹并自动设置好。第一次运行大约需要 3–5 分钟。

**3. 回答两个问题。** 每次看到提示时都按 `Enter` 确认“是”。
- *"将 WhatsApp 连接到：Claude Desktop…？"* 这会将 WhatsApp 添加到你的 Mac 上找到的 AI 应用中，并先备份它们的设置。
- *"现在启动 WhatsApp 并显示二维码吗？"*

**4. 扫描二维码。** 在你的手机上操作：WhatsApp → **Settings（设置）** → **Linked devices（关联设备）** → **Link a device（链接设备）**。将手机摄像头对准终端窗口中的二维码。
当提示 **Connected to WhatsApp** 时，你就已成功连接。请保持此窗口打开，你可以将其最小化。

**5. 重启你的 AI 应用。** 完全退出它（按 `Cmd` + `Q`），然后重新打开。

**6. 试试看。** 向你的 AI 应用提问：

```text
列出我最近的 5 个 WhatsApp 聊天。
```

如果你能看到聊天记录，说明成功了。旧聊天内容可能需要几分钟才会出现。

## 每天的操作

1. 打开 **WhatsApp-MCP** 文件夹（Finder → **Go** → **Home** → **WhatsApp-MCP**）。提示：把它拖到 Finder 侧边栏，这样它就一直显示在那里。
2. 双击 **Start WhatsApp**。会弹出一个终端窗口。保持它开着就行（你可以把它最小化）。
3. 像平常一样使用你的 AI 应用。
4. 要停止时，关闭那个窗口即可。

只有在 WhatsApp 运行时，新消息才会被保存。你的手机不需要一直放在电脑旁边。

## 让 AI 发送消息（可选）

默认情况下，AI **不会**自动发送消息。要开启此功能，请按以下步骤操作：

- 双击 **Start WhatsApp - can send**（而不是 **Start WhatsApp**）。
- 当 AI 想要发送消息时，**WhatsApp 窗口会发出提示音，并显示接收对象和完整内容**。请切换到该窗口，输入 `y` 并按下 `Enter` 来发送。如果只需取消操作，直接按下 `Enter` 即可。若你在 2 分钟内未回应，消息将不会被发送。
- 仅对你明确要求的内容输入 `y`。别人发来的消息永远无法替你按这个键。

**建议限制可发消息的对象**。打开 **allowed-recipients.txt**（它会在 TextEdit 中打开）。每行写一个电话号码，包含国家代码且无空格，例如 `85291234567`。保存后，AI 只能向这些人发送消息。状态更新和频道除非明确列出，否则不可用。对于群组，请询问你的 AI：“我的 WhatsApp 群组‘Family’的 ID 是什么？”并添加该 ID（它通常以 `@g.us` 结尾）。

**发送文件**。AI 只能从 **outbox** 文件夹中发送文件。先将文件放入该文件夹，然后让 AI 发送（最大限制为 100 MB）。

## WhatsApp-MCP 文件夹里有什么

| 项目 | 说明 |
|---|---|
| **Start WhatsApp** | 每天启动。AI 只能读取消息。 |
| **Start WhatsApp - can send** | 允许发送消息（AI 必须先问你）。 |
| **allowed-recipients.txt** | AI 可以发消息的对象。留空 = 任何人都可以。 |
| **outbox** | AI 能从中发送文件的唯一文件夹。 |
| **Update** | 获取最新版本，保留你的聊天记录和登录状态。 |
| **Uninstall** | 从你的 AI 应用中移除它，并可删除你保存的聊天记录。 |
| **Install** | 重新运行设置程序（例如移动文件夹后）。 |
| whatsapp-bridge/store | 你的聊天记录和 WhatsApp 登录凭证。私密信息，不要分享。 |

把这个文件夹放在你的主目录里。别把它放到桌面、文档或 iCloud Drive，因为这些位置可能会同步到 iCloud。如果你移动了它，请双击 **Install** 重新运行设置。

## 更新

双击 **Update**。你的关联手机、已保存的消息、发件箱和允许列表都会保留。然后重新启动 WhatsApp，并重启你的 AI 应用。

## 卸载

1. 双击 **Uninstall**。这会从你的 AI 应用中移除 WhatsApp，删除已保存的消息前会先询问你（输入 `DELETE` 确认）。
2. 在手机上：WhatsApp → **设置（Settings）** → **关联设备（Linked devices）** → 点击这台电脑 → **登出（Log out）**。
3. 将 **WhatsApp-MCP** 文件夹移到废纸篓。

## 如果出问题了

| 你看到的现象 | 该怎么做 |
|---|---|
| AI 提示 WhatsApp 未启动 | 双击 **Start WhatsApp**，并保持窗口打开。 |
| 二维码已过期 | 关闭窗口，再次双击 **Start WhatsApp**。 |
| AI 里根本没有 WhatsApp | 完全退出 AI 应用（`Cmd` + `Q`），然后重新打开它。如果还是找不到：再双击 **Install**。 |
| Mac 提示“无法验证”你双击的文件 | **设置（System Settings）** → **隐私与安全性（Privacy & Security）** → 向下滚动 → **继续运行（Open Anyway）**。只有用浏览器下载 ZIP 文件时才会出现。 |
| "Sending is turned off"（发送已关闭） | 关闭窗口，改用 **Start WhatsApp - can send**。 |
| AI 提示发送被取消或未确认 | 切换到 WhatsApp 窗口，当它询问时输入 `y` + `Enter`。 |
| "port 8080 is already used by another program"（8080 端口已被占用） | 有其他应用占用了该端口。重启 Mac，或联系分享给你的人。 |
| "not in WHATSAPP_SEND_ALLOWLIST"（不在允许列表内） | 那个人不在 **allowed-recipients.txt**里。添加他/她，或者清空列表。 |
| 文件无法发送 | 先把它放到 **outbox**文件夹中。最大限制为 100 MB。 |
| "The bridge refused the access key"（访问密钥被拒） | 关闭 WhatsApp 窗口，重新启动它，然后重启 AI 应用。 |
| 旧消息不见了 | 只有当 WhatsApp 运行时才会保存消息。部分历史记录会在你首次链接时到达。 |
| 长时间中断后停止工作 | 如果手机约 20 天没有上线，关联设备会自动登出。启动 WhatsApp 并扫描新的二维码。 |
| 其他情况 | 将文件 `.tools/install.log`（位于该文件夹内）发送给分享给你的人，或 [提交 issue](https://github.com/HyperfocuSam/whatsapp-mcp-hk/issues)。切勿发送你的 `store`文件夹。

## 常见问题

**分享者能看到我的消息吗？**  
不会。没有任何内容会发送给维护者，也没有追踪机制。你的消息只会像往常一样发送到 WhatsApp，以及在 AI 阅读它们时传送到你的 AI 公司。

**安装程序会下载什么？**  
它会从 go.dev 下载 Go，并从 GitHub 下载 uv，每个文件都会通过官方校验和进行验证。此外还会从官方的 Go 和 Python 包服务器下载代码库。所有内容都存放在 WhatsApp-MCP 文件夹内。不会在全系统范围内安装任何东西，也绝不会要求你输入密码。

**支持两个 WhatsApp 号码吗？**  
可以。只需在另一个文件夹中创建一份副本即可。详见下方的“开发者”部分。

**Windows 能用吗？**  
目前 Windows 还没有便捷的安装包。开发者可以通过手动设置来完成（见下文）。

## 面向开发者

<details>
<summary>手动安装、其他应用、设置及工作原理</summary>

### 手动安装（macOS / Linux）

```bash
git clone https://github.com/HyperfocuSam/whatsapp-mcp-hk.git ~/WhatsApp-MCP
cd ~/WhatsApp-MCP && bash scripts/install.sh
```

`scripts/install.sh` 会使用你自己的 `go`（版本需 ≥ 1.26）和 `uv`（如果有的话），否则使用 `.tools/` 中的私有副本。桥接程序完全用 Go 编写（`CGO_ENABLED=0`、modernc SQLite），因此不需要 C 编译器。Linux 的操作方式也相同。

### 手动连接 MCP 客户端

MCP 服务器通过标准输入输出（stdio）运行。命令：`<folder>/whatsapp-mcp-server/.venv/bin/python`，参数：`<folder>/whatsapp-mcp-server/main.py`。

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

- Claude Desktop: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Cursor: `~/.cursor/mcp.json`
- Claude Code: `claude mcp add --scope user whatsapp -- <python> <main.py>`
- Codex (`~/.codex/config.toml`)：

```toml
[mcp_servers.whatsapp]
command = "/Users/you/WhatsApp-MCP/whatsapp-mcp-server/.venv/bin/python"
args = ["/Users/you/WhatsApp-MCP/whatsapp-mcp-server/main.py"]
```

`whatsapp-mcp-server/connect_agents.py add|remove|status` 会自动完成此操作，并支持备份。

### 设置（桥接的环境变量）

| 变量 | 含义 |
|---|---|
| `WHATSAPP_ALLOW_SEND=1` | 允许发送消息（默认关闭）。`scripts/start.sh --allow-send` 会开启它。 |
| `WHATSAPP_SEND_ALLOWLIST` | 逗号分隔的电话号码 / JID。`start.sh` 会从 `allowed-recipients.txt` 填充，或从 `--only a,b` 获取。 |
| `WHATSAPP_SEND_DIR` | 发送文件夹（默认 `outbox/`）。MCP 服务器会从 `store/bridge.json` 读取它。 |
| `WHATSAPP_BRIDGE_PORT` | 本地 API 端口（默认 8080）。用于运行第二个副本并配合另一个号码。 |
| `WHATSAPP_LOG_MESSAGES=1` | 将消息文本打印到终端（默认关闭）。 |
| `WHATSAPP_CONFIRM_SENDS=0` | 在桥接窗口每次发送前不询问 `y`（默认：会询问）。 |

### 它是如何工作的

- `whatsapp-bridge/`（Go，[whatsmeow](https://github.com/tulir/whatsmeow)）是一个 WhatsApp 关联设备。它将消息存储在 `store/messages.db` 中，并提供一个本地 API（127.0.0.1，令牌在 `store/api-token`）用于发送和下载媒体。
- `whatsapp-mcp-server/`（Python，MCP）以只读方式读取 `messages.db`，并调用桥接 API。
- 安全设计和威胁模型：[SECURITY.md](SECURITY.md)。
- 测试：`cd whatsapp-bridge && go test ./...` 和 `cd whatsapp-mcp-server && uv run python -m unittest test_security`。

### Windows（手动，未测试）

安装 Go 和 uv。然后：

```powershell
cd whatsapp-bridge; $env:CGO_ENABLED="0"; go build -o whatsapp-bridge.exe .
cd ..\whatsapp-mcp-server; uv sync
cd ..\whatsapp-bridge; .\whatsapp-bridge.exe
```

发送消息前：`$env:WHATSAPP_ALLOW_SEND="1"`。用 `whatsapp-mcp-server\.venv\Scripts\python.exe` 和 `main.py` 连接你的客户端。

### 工具

`search_contacts`、`list_chats`、`get_chat`、`list_messages`、`get_message_context`、`get_direct_chat_by_contact`、`get_contact_chats`、`get_last_interaction`、`download_media`，以及发送工具 `send_message`、`send_file` 和 `send_audio_message`（仅在开启发送功能时使用；语音消息需要 ffmpeg）。

</details>

## 致谢

- 原始项目：[lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp)，作者 Luke Harries（采用 MIT 许可证）。
- WhatsApp 库：[whatsmeow](https://github.com/tulir/whatsmeow)，作者 Tulir Asokan。
- 本版本由 Sam Wong ([HyperfocuSam](https://github.com/HyperfocuSam)) 维护。

## 许可证

MIT。详见 [LICENSE](LICENSE)。
