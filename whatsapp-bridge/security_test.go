package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeFilename(t *testing.T) {
	cases := map[string]string{
		"../../Library/LaunchAgents/x.plist": "x.plist",
		`..\..\evil.bat`:                     "evil.bat",
		".zshrc":                             "zshrc",
		"..":                                 "fallback",
		"":                                   "fallback",
		"報價單 2026.pdf":                       "報價單 2026.pdf",
		"invoice\u202Efdp.exe":               "invoice_fdp.exe",
		"a\x00b.txt":                         "a_b.txt",
		"/etc/passwd":                        "passwd",
	}
	for in, want := range cases {
		if got := safeFilename(in, "fallback"); got != want {
			t.Errorf("safeFilename(%q) = %q, want %q", in, got, want)
		}
	}
	long := strings.Repeat("長", 200) + ".pdf"
	got := safeFilename(long, "f")
	if len(got) > 150 || !strings.HasSuffix(got, ".pdf") {
		t.Errorf("long name not bounded: %d bytes %q", len(got), got[len(got)-8:])
	}
}

func TestSafeSegment(t *testing.T) {
	for _, in := range []string{"../../x", "..", "a/b", `a\b`} {
		got := safeSegment(in)
		if strings.ContainsAny(got, `/\`) || strings.HasPrefix(got, ".") {
			t.Errorf("safeSegment(%q) = %q", in, got)
		}
	}
}

func TestResolveSendPath(t *testing.T) {
	root := t.TempDir()
	outbox := filepath.Join(root, "outbox")
	os.MkdirAll(outbox, 0700)
	sendDir, _ = filepath.EvalSymlinks(outbox)
	secret := filepath.Join(root, "secret.txt")
	os.WriteFile(secret, []byte("x"), 0600)
	ok := filepath.Join(outbox, "photo.jpg")
	os.WriteFile(ok, []byte("x"), 0600)
	os.Symlink(secret, filepath.Join(outbox, "link.txt"))
	os.MkdirAll(filepath.Join(outbox, "dir"), 0700)

	if _, err := resolveSendPath(ok); err != nil {
		t.Errorf("file in outbox refused: %v", err)
	}
	if _, err := resolveSendPath("photo.jpg"); err != nil {
		t.Errorf("bare name refused: %v", err)
	}
	for _, bad := range []string{secret, filepath.Join(outbox, "..", "secret.txt"),
		filepath.Join(outbox, "link.txt"), filepath.Join(outbox, "dir"), outbox, "/etc/hosts", ""} {
		if _, err := resolveSendPath(bad); err == nil {
			t.Errorf("resolveSendPath(%q) allowed, want refused", bad)
		}
	}
}

func TestNormaliseRecipient(t *testing.T) {
	j, err := normaliseRecipient("+852 9123-4567")
	if err != nil || j.String() != "85291234567@s.whatsapp.net" {
		t.Errorf("got %v %v", j, err)
	}
	for _, bad := range []string{"abc", "123", "85291234567; rm -rf"} {
		if _, err := normaliseRecipient(bad); err == nil {
			t.Errorf("normaliseRecipient(%q) accepted", bad)
		}
	}
}

func TestGuard(t *testing.T) {
	apiToken = strings.Repeat("a", 64)
	h := guard(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	mk := func(mod func(r *http.Request)) int {
		r := httptest.NewRequest("POST", "http://127.0.0.1:8080/api/send", strings.NewReader("{}"))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+apiToken)
		if mod != nil {
			mod(r)
		}
		w := httptest.NewRecorder()
		h(w, r)
		return w.Code
	}
	if c := mk(nil); c != 200 {
		t.Fatalf("valid request got %d", c)
	}
	checks := map[string]func(r *http.Request){
		"no token":       func(r *http.Request) { r.Header.Del("Authorization") },
		"wrong token":    func(r *http.Request) { r.Header.Set("Authorization", "Bearer nope") },
		"browser origin": func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") },
		"cross-site":     func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "cross-site") },
		"text/plain":     func(r *http.Request) { r.Header.Set("Content-Type", "text/plain") },
		"rebinding host": func(r *http.Request) { r.Host = "evil.example:8080" },
		"GET":            func(r *http.Request) { r.Method = "GET" },
	}
	for name, mod := range checks {
		if c := mk(mod); c == 200 {
			t.Errorf("%s: allowed, want refused", name)
		}
	}
}
