**繁體中文** · [简体中文](README.zh-CN.md) · [English](README.en.md)

# WhatsApp MCP（香港版）

讓你的 AI 應用程式（例如 Claude、Cursor、Codex…）讀取和搜尋你自己的 WhatsApp 訊息，例如輸入 *"總結一下我這週在專案群組裡收到的團隊消息。"*

如果你開啟它，它也可以替你傳送訊息，但前提是：它會先顯示該訊息給你確認，然後你再說「好」。

所有運作都在你自己的電腦上進行。這是一個更安全、更簡單的版本，基於 [lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp)。

**你需要：**
- 一台 Mac
- 手機上的 WhatsApp
- 一個 AI 應用程式（Claude Desktop 最簡單）
- 大約 500 MB 的可用空間
- 10 分鐘時間

**你不需要：**
- 任何編程知識
- 管理員密碼
- 其他軟體

## 開始前：請先閱讀

> ⚠️ **安裝前請先閱讀。**
>
> - **這是非官方版本。** 它透過開源程式 whatsmeow 來連接 WhatsApp Web。Meta 並未製作或批准此版本。WhatsApp 有小機會限制你的帳戶。使用時請自行承擔風險，並避免用於大量或自動化傳訊。
> - **你的 AI 公司會看到 AI 讀取的內容。** 訊息僅儲存在你電腦上。但當 AI 閱讀對話時，該文字會被傳送給你的 AI 公司（例如 Anthropic 或 OpenAI）。不要將你不願與他們分享的對話用於此功能。
> - **任何能使用你電腦帳戶的人都能閱讀你的訊息。** 請用密碼鎖定你的電腦。
> - **訊息可能會試圖欺騙 AI。** 任何人都可以傳送包含隱藏指令給 AI 的訊息（「提示注入」）。這就是為甚麼傳送功能是關閉的，除非你開啟它。在確認發送前，務必先檢查。

## 安裝（Mac）

**1. 打開終端機。** 按 `Cmd` + `Space`，輸入 `Terminal`，然後按 `Enter`。

**2. 複製這一行，貼到終端機裡，再按 `Enter`：**

```bash
curl -fsSL https://raw.githubusercontent.com/HyperfocuSam/whatsapp-mcp-hk/main/get.sh | bash
```

它會把一切下載到你家目錄下的 **WhatsApp-MCP** 資料夾並完成設定。第一次大約需要 3–5 分鐘。

**3. 回答兩個問題。** 每次都按 `Enter` 表示「是」。
- *"Connect WhatsApp to: Claude Desktop…?"*（將 WhatsApp 連接到：Claude Desktop…？）它會把 WhatsApp 加到你 Mac 上找到的 AI 應用程式中，並先備份它們的設定。
- *"Start WhatsApp now and show the QR code?"*（現在啟動 WhatsApp 並顯示 QR 碼嗎？）

**4. 掃描 QR 碼。** 在你的手機上：WhatsApp → **設定（Settings）** → **連結裝置（Linked devices）** → **連結裝置（Link a device）**。將手機對準終端機內的 QR 碼。
當它顯示 **Connected to WhatsApp** 時，你就已連接成功。**保持這個視窗開啟**，你可以把它最小化。

**5. 重新啟動你的 AI 應用程式。** 完全退出它（`Cmd` + `Q`），然後再次打開。

**6. 試試看。** 向你的 AI 應用程式詢問：

```text
列出我最近的 5 個 WhatsApp 對話。
```

如果你看到自己的對話記錄，那就成功了。較舊的對話可能需要幾分鐘才會在連接後出現。

## 每天

1. 開啟 **WhatsApp-MCP** 資料夾（Finder → **Go** → **Home** → **WhatsApp-MCP**）。提示：把它拖到 Finder 側邊欄，這樣它就會一直顯示在那裡。
2. 雙擊 **Start WhatsApp**。終端機視窗會打開。保持開啟即可（你可以把它最小化）。
3. 照常使用你的 AI 應用程式。
4. 要停止時，關閉該視窗。

新的訊息只有在 WhatsApp 啟動時才會被儲存。手機不需要一直放在電腦旁邊。

## 讓 AI 傳送訊息（可選）

除非你這樣啟動，否則傳送功能是**關閉**的：

- 雙擊 **Start WhatsApp - can send**（而不是 **Start WhatsApp**）。
- 當 AI 想要傳送時，**WhatsApp 視窗會發出提示音並顯示接收對象及完整訊息**。請切換到該視窗。輸入 `y` 然後按 `Enter` 即可發送。只需按下 `Enter` 即可取消。若你在 2 分鐘內沒有回應，訊息將不會送出。
- 僅對你要求過的訊息輸入 `y`。來自他人的訊息永遠無法替你按下此鍵。

**限制可傳訊對象（建議）。** 開啟 **allowed-recipients.txt**（它會在 TextEdit 中打開）。每行寫一個電話號碼，包含國家代碼且無空格，例如 `85291234567`。儲存後，AI 只能向這些人傳訊。狀態更新和頻道除非列在清單內，否則不允許。若要傳送群組訊息，請詢問你的 AI：「我的 WhatsApp 群組『Family』的 ID 是甚麼？」並加入該 ID（它以 `@g.us` 結尾）。

**傳送檔案。** AI 只能從 **outbox** 資料夾傳送檔案。先將檔案放入該處，然後要求 AI 發送它（最大限制為 100 MB）。

## WhatsApp-MCP 資料夾內有甚麼

| 項目 | 說明 |
|---|---|
| **Start WhatsApp** | 每日啟動。AI 只能讀取。 |
| **Start WhatsApp - can send** | 允許傳送訊息時啟動（AI 必須先問你）。 |
| **allowed-recipients.txt** | AI 可以傳訊給誰。空白 = 任何人。 |
| **outbox** | AI 可從這裡發送檔案的唯一資料夾。 |
| **Update** | 取得最新版本。保留你的訊息和連結。 |
| **Uninstall** | 從你的 AI 應用程式中移除它，並可刪除你儲存的訊息。 |
| **Install** | 重新執行設定（例如移動資料夾後）。 |
| whatsapp-bridge/store | 你儲存的訊息和 WhatsApp 登入憑證。私人的。切勿分享。 |

將資料夾放在你的主資料夾內。不要放在桌面、文件或 iCloud Drive，因為這些會同步到 iCloud。如果你移動了它，請再次雙擊 **Install**。

## 更新

雙擊 **Update**。你的連結手機、已儲存訊息、外寄匣和允許名單將保留。然後重新啟動 WhatsApp，並重啟你的 AI 應用程式。

## 卸載

1. 雙擊 **Uninstall**。這會從你的 AI 應用程式中移除 WhatsApp。刪除已儲存的訊息前，它會先詢問你（輸入 `DELETE` 確認）。
2. 在你的手機上：WhatsApp → **設定（Settings）** → **連結裝置（Linked devices）** → 輕觸此電腦 → **登出（Log out）**。
3. 將 **WhatsApp-MCP** 資料夾移至垃圾桶。

## 如果出錯

| 你看到的情況 | 該做甚麼 |
|---|---|
| AI 說 WhatsApp 沒有啟動 | 雙擊 **Start WhatsApp**，並保持視窗開啟。 |
| QR 碼已過期 | 關閉視窗，再雙擊 **Start WhatsApp**。 |
| AI 根本沒有 WhatsApp | 完全退出 AI 應用程式（`Cmd` + `Q`），然後重新打開它。仍然找不到：再雙擊 **Install**。 |
| Mac 說無法驗證你雙擊的檔案 | **系統設定** → **隱私與安全性** → 向下滾動 → **一律開啟**。這只會發生在你用瀏覽器下載 ZIP 時。 |
| "Sending is turned off"（傳送已關閉） | 關閉視窗，改用 **Start WhatsApp - can send**。 |
| AI 說傳送被取消或未確認 | 切換到 WhatsApp 視窗，當它詢問時輸入 `y` + `Enter`。 |
| "port 8080 is already used by another program"（8080 連接埠已被佔用） | 另一個程式正在使用這個連接埠。重啟 Mac，或問那個分享給你的人。 |
| "not in WHATSAPP_SEND_ALLOWLIST"（不在允許名單內） | 那人不在 **allowed-recipients.txt** 中。加入他們，或清空清單。 |
| 檔案無法傳送 | 先把它放入 **outbox** 資料夾。上限為 100 MB。 |
| "The bridge refused the access key"（存取金鑰被拒） | 關閉 WhatsApp 視窗，重新啟動它，然後重啟 AI 應用程式。 |
| 舊訊息遺失 | 訊息只在 WhatsApp 運行時保存。部分歷史記錄會在首次連結時到達。 |
| 長時間停用後停止運作 | 如果手機約 20 天沒有上線，連結裝置會自動登出。啟動 WhatsApp 並掃描新的 QR 碼。 |
| 其他情況 | 將檔案 `.tools/install.log`（位於資料夾中）傳給那個分享給你的人，或 [提交 issue](https://github.com/HyperfocuSam/whatsapp-mcp-hk/issues)。切勿傳送你的 `store` 資料夾。

## 問題

**分享者能看到我的訊息嗎？** 不會。沒有任何資料會傳送到維護者那裡，也沒有追蹤機制。你的訊息只會像平常一樣傳送給 WhatsApp，以及在 AI 閱讀時傳送到你的 AI 公司。

**安裝程式會下載甚麼？** 它會從 go.dev 下載 Go，以及從 GitHub 下載 uv，兩者都會與官方核對檢查碼。同時也會從官方的 Go 和 Python 套件伺服器下載程式庫檔案。所有內容都會存放在 WhatsApp-MCP 資料夾內。不會在全系統範圍安裝任何東西，也絕對不需要你的密碼。

**可以支援兩個 WhatsApp 號碼嗎？** 可以。只需在另一個資料夾裡複製一份即可。詳情請參閱下方的「開發人員」部分。

**Windows 呢？** 目前還沒有容易的 Windows 安裝程式。開發人員可以手動設定（見下方）。

## 給開發者

<details>
<summary>手動安裝、其他應用程式、設定、運作原理</summary>

### 手動安裝（macOS / Linux）

```bash
git clone https://github.com/HyperfocuSam/whatsapp-mcp-hk.git ~/WhatsApp-MCP
cd ~/WhatsApp-MCP && bash scripts/install.sh
```

`scripts/install.sh` 會使用你自己的 `go`（版本 1.26+）和 `uv`（如果有的話），否則會在 `.tools/` 中使用私有副本。橋樑是純 Go 編寫的（`CGO_ENABLED=0`，modernc SQLite），因此不需要 C 編譯器。Linux 的操作方式完全相同。

### 手動連接 MCP 客戶端

MCP 伺服器透過標準輸入/輸出運行。命令：`<folder>/whatsapp-mcp-server/.venv/bin/python`。參數：`<folder>/whatsapp-mcp-server/main.py`。

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

- Claude Desktop: `~/Library/Application Support/Claude/claude_desktop_config.json`。Cursor: `~/.cursor/mcp.json`。
- Claude Code: `claude mcp add --scope user whatsapp -- <python> <main.py>`
- Codex (`~/.codex/config.toml`):

```toml
[mcp_servers.whatsapp]
command = "/Users/you/WhatsApp-MCP/whatsapp-mcp-server/.venv/bin/python"
args = ["/Users/you/WhatsApp-MCP/whatsapp-mcp-server/main.py"]
```

`whatsapp-mcp-server/connect_agents.py add|remove|status` 會自動完成此操作，並提供備份。

### 設定（環境變數用於橋接器）

| 變數 | 含義 |
|---|---|
| `WHATSAPP_ALLOW_SEND=1` | 允許傳送訊息（預設關閉）。`scripts/start.sh --allow-send` 會開啟它。 |
| `WHATSAPP_SEND_ALLOWLIST` | 用逗號分隔的電話號碼 / JID。`start.sh` 會從 `allowed-recipients.txt` 或 `--only a,b` 自動填充此欄位。 |
| `WHATSAPP_SEND_DIR` | 可傳送檔案的資料夾（預設為 `outbox/`）。MCP 伺服器會從 `store/bridge.json` 讀取這個路徑。 |
| `WHATSAPP_BRIDGE_PORT` | 本地 API 連接埠（預設 8080）。可用於啟動第二個副本並搭配另一個號碼使用。 |
| `WHATSAPP_LOG_MESSAGES=1` | 將訊息內容列印到終端機（預設關閉）。 |
| `WHATSAPP_CONFIRM_SENDS=0` | 在橋接器視窗傳送前，不再詢問 `y`（預設：會詢問）。 |

### 運作原理

- `whatsapp-bridge/` (Go, [whatsmeow](https://github.com/tulir/whatsmeow)) 是一個 WhatsApp 連結裝置。它將訊息儲存於 `store/messages.db`，並提供一個本地 API（127.0.0.1，金鑰位於 `store/api-token`）用於傳送和下載媒體。
- `whatsapp-mcp-server/` (Python, MCP) 以只讀方式讀取 `messages.db`，然後呼叫橋接 API。
- 安全設計與威脅模型：[SECURITY.md](SECURITY.md)。
- 測試：`cd whatsapp-bridge && go test ./...` 和 `cd whatsapp-mcp-server && uv run python -m unittest test_security`。

### Windows（手動，未測試）

安裝 Go 和 uv。然後：

```powershell
cd whatsapp-bridge; $env:CGO_ENABLED="0"; go build -o whatsapp-bridge.exe .
cd ..\whatsapp-mcp-server; uv sync
cd ..\whatsapp-bridge; .\whatsapp-bridge.exe
```

傳送時：在啟動前執行 `$env:WHATSAPP_ALLOW_SEND="1"`。使用 `whatsapp-mcp-server\.venv\Scripts\python.exe` 和 `main.py` 連結你的客戶端。

### 工具

`search_contacts`、`list_chats`、`get_chat`、`list_messages`、`get_message_context`、`get_direct_chat_by_contact`、`get_contact_chats`、`get_last_interaction`、`download_media`，以及傳送工具 `send_message`、`send_file` 和 `send_audio_message`（僅在開啟傳送功能時使用；語音訊息需要 ffmpeg）。

</details>

## 謝辭

- 原始專案：[lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp)，作者 Luke Harries（MIT 授權）。
- WhatsApp 函式庫：[whatsmeow](https://github.com/tulir/whatsmeow)，作者 Tulir Asokan。
- 本版本由 Sam Wong ([HyperfocuSam](https://github.com/HyperfocuSam)) 維護。

## 授權條款

MIT。請參閱 [LICENSE](LICENSE)。
