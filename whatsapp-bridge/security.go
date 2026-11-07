package main

// security.go — the safety layer added in the HK edition.
//
// Summary (details in SECURITY.md):
//   - the REST API listens on 127.0.0.1 only and needs a secret token (store/api-token)
//   - browser requests (Origin / Sec-Fetch-Site headers, foreign Host) are refused
//   - sending is OFF unless WHATSAPP_ALLOW_SEND=1; optional recipient allowlist
//   - files can only be sent from one folder (the outbox)
//   - downloaded media gets a cleaned filename and stays inside store/media
//   - store/ is private to the current user (0700 / 0600)
//   - message text is not printed to the terminal unless WHATSAPP_LOG_MESSAGES=1

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"go.mau.fi/whatsmeow/types"
)

const (
	maxRequestBody = 1 << 20   // 1 MB of JSON per API call
	maxMediaBytes  = 100 << 20 // 100 MB per file sent or downloaded
	tokenFile      = "store/api-token"
	mediaRoot      = "store/media"
)

var (
	apiToken      string
	allowSend     bool
	sendAllowlist map[string]bool // normalised JIDs; empty = anyone
	sendDir       string          // absolute, symlinks resolved
	logMessages   bool
)

func envTrue(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

// moveToBridgeDir makes the bridge work no matter which folder it is started
// from: store/ always sits next to the binary. `go run` builds into a temp
// folder, so in that case the current folder is kept.
func moveToBridgeDir() {
	if dir := os.Getenv("WHATSAPP_BRIDGE_DIR"); dir != "" {
		_ = os.Chdir(dir)
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}
	dir := filepath.Dir(exe)
	tmp, _ := filepath.EvalSymlinks(os.TempDir())
	if tmp != "" && strings.HasPrefix(dir, tmp) {
		return
	}
	if strings.Contains(dir, string(filepath.Separator)+"go-build") {
		return
	}
	_ = os.Chdir(dir)
}

// secureStore creates store/ as a private folder and tightens anything an
// older install left readable by other users of the computer.
func secureStore() error {
	if err := os.MkdirAll("store", 0700); err != nil {
		return fmt.Errorf("failed to create store directory: %v", err)
	}
	return filepath.WalkDir("store", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			_ = os.Chmod(p, 0700)
		} else {
			_ = os.Chmod(p, 0600)
		}
		return nil
	})
}

// newToken writes a fresh random store/api-token on every start, so a key
// seen by anything else while the bridge was off is useless afterwards.
func newToken() (string, error) {
	_ = os.Remove(tokenFile)
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	t := hex.EncodeToString(buf)
	f, err := os.OpenFile(tokenFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("could not write %s: %v", tokenFile, err)
	}
	defer f.Close()
	if _, err := f.WriteString(t + "\n"); err != nil {
		return "", err
	}
	return t, nil
}

// normaliseRecipient turns "+852 9123-4567" or a JID into a JID string.
func normaliseRecipient(r string) (types.JID, error) {
	r = strings.TrimSpace(r)
	if strings.Contains(r, "@") {
		jid, err := types.ParseJID(r)
		if err != nil {
			return types.JID{}, fmt.Errorf("invalid JID %q: %v", r, err)
		}
		return jid.ToNonAD(), nil
	}
	digits := strings.Map(func(c rune) rune {
		switch c {
		case '+', ' ', '-', '(', ')':
			return -1
		}
		return c
	}, r)
	if len(digits) < 6 || len(digits) > 15 || strings.Trim(digits, "0123456789") != "" {
		return types.JID{}, fmt.Errorf("invalid phone number %q: use the country code and digits only, e.g. 85291234567", r)
	}
	return types.JID{User: digits, Server: types.DefaultUserServer}, nil
}

func initSecurityConfig() error {
	allowSend = envTrue("WHATSAPP_ALLOW_SEND")
	logMessages = envTrue("WHATSAPP_LOG_MESSAGES")

	sendAllowlist = map[string]bool{}
	for _, item := range strings.Split(os.Getenv("WHATSAPP_SEND_ALLOWLIST"), ",") {
		if strings.TrimSpace(item) == "" {
			continue
		}
		jid, err := normaliseRecipient(item)
		if err != nil {
			return fmt.Errorf("WHATSAPP_SEND_ALLOWLIST: %v", err)
		}
		sendAllowlist[jid.String()] = true
	}

	dir := os.Getenv("WHATSAPP_SEND_DIR")
	if dir == "" {
		dir = filepath.Join("..", "outbox")
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(abs, 0700); err != nil {
		return fmt.Errorf("could not create send folder %s: %v", abs, err)
	}
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		abs = real
	}
	sendDir = abs

	tok, err := newToken()
	if err != nil {
		return err
	}
	apiToken = tok
	return nil
}

// writeBridgeInfo tells the MCP server which port and send folder this bridge
// uses, so the agent config never needs extra settings.
func writeBridgeInfo(port int) error {
	b, _ := json.Marshal(map[string]interface{}{"port": port, "send_dir": sendDir, "pid": os.Getpid()})
	_ = os.WriteFile("store/bridge.pid", []byte(fmt.Sprintf("%d\n", os.Getpid())), 0600)
	tmp := "store/bridge.json.tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, "store/bridge.json")
}

// clearBridgeInfo runs on exit, so the MCP server knows the bridge is gone.
func clearBridgeInfo() {
	_ = os.Remove("store/bridge.json")
	_ = os.Remove("store/bridge.pid")
}

func recipientAllowed(jid types.JID) bool {
	key := jid.ToNonAD().String()
	if jid.Server == types.BroadcastServer || jid.Server == types.NewsletterServer {
		return sendAllowlist[key] // Status / broadcast / channels: only if listed by name
	}
	return len(sendAllowlist) == 0 || sendAllowlist[key]
}

// within reports whether path is inside root (both absolute, cleaned).
func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || filepath.IsAbs(rel) {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// resolveSendPath only accepts regular files inside the send folder.
// A bare file name ("photo.jpg") is looked up in the send folder.
func resolveSendPath(p string) (string, error) {
	if strings.TrimSpace(p) == "" {
		return "", errors.New("no file given")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(sendDir, p)
	}
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		return "", fmt.Errorf("file not found: %s", p)
	}
	real, _ = filepath.Abs(real)
	if !within(sendDir, real) {
		return "", fmt.Errorf("files can only be sent from the send folder %s — copy the file there first", sendDir)
	}
	fi, err := os.Stat(real)
	if err != nil || !fi.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file: %s", p)
	}
	if fi.Size() > maxMediaBytes {
		return "", fmt.Errorf("file is larger than %d MB", maxMediaBytes>>20)
	}
	return real, nil
}

var unsafeNameRunes = `<>:"/\|?*`

func isUnsafeRune(r rune) bool {
	switch {
	case r < 0x20 || r == 0x7f:
		return true
	case r == 0x061c, r >= 0x200b && r <= 0x200f, r >= 0x202a && r <= 0x202e, r == 0x2060, r >= 0x2066 && r <= 0x2069, r == 0xfeff:
		return true // direction overrides that hide the real extension
	case r == utf8.RuneError:
		return true
	}
	return strings.ContainsRune(unsafeNameRunes, r)
}

// safeFilename cleans a name chosen by the sender of a WhatsApp message:
// no folders, no "..", no hidden files, no control characters, bounded length.
// Chinese and other non-Latin names are kept.
func safeFilename(name, fallback string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.Map(func(r rune) rune {
		if isUnsafeRune(r) {
			return '_'
		}
		return r
	}, name)
	name = strings.TrimLeft(strings.TrimSpace(name), ". ")
	name = strings.TrimRight(name, ". ")
	if name == "" {
		name = fallback
	}
	const maxLen = 150
	if len(name) > maxLen {
		ext := filepath.Ext(name)
		if len(ext) > 16 {
			ext = ""
		}
		base := name[:maxLen-len(ext)]
		for !utf8.ValidString(base) {
			base = base[:len(base)-1]
		}
		name = base + ext
	}
	return name
}

// safeSegment keeps a JID or message ID usable as a single folder/file part.
func safeSegment(s string) string {
	out := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '@', r == '-', r == '_', r == '.':
			return r
		}
		return '_'
	}, s)
	out = strings.TrimLeft(out, ".")
	if out == "" {
		out = "unknown"
	}
	if len(out) > 100 {
		out = out[:100]
	}
	return out
}

// guard wraps every API handler: POST only, localhost only, no browsers,
// JSON only, valid token, bounded body.
func guard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Origin") != "" {
			http.Error(w, "Browser requests are not allowed", http.StatusForbidden)
			return
		}
		if s := r.Header.Get("Sec-Fetch-Site"); s != "" && s != "none" {
			http.Error(w, "Browser requests are not allowed", http.StatusForbidden)
			return
		}
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		if host != "127.0.0.1" && host != "localhost" {
			http.Error(w, "Forbidden host", http.StatusForbidden)
			return
		}
		if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		auth := r.Header.Get("Authorization")
		tok := strings.TrimPrefix(auth, "Bearer ")
		if !strings.HasPrefix(auth, "Bearer ") || subtle.ConstantTimeCompare([]byte(tok), []byte(apiToken)) != 1 {
			http.Error(w, "Unauthorized: missing or wrong API token (store/api-token)", http.StatusUnauthorized)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// logMsg prints personal content (message text, names) only when the user
// asked for it with WHATSAPP_LOG_MESSAGES=1.
func logMsg(format string, args ...interface{}) {
	if logMessages {
		fmt.Printf(format+"\n", args...)
	}
}
