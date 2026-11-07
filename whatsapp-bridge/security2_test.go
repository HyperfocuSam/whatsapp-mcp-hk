package main

import (
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"go.mau.fi/whatsmeow/types"
)

func TestBroadcastNeedsExplicitListing(t *testing.T) {
	sendAllowlist = map[string]bool{}
	status, _ := types.ParseJID("status@broadcast")
	channel, _ := types.ParseJID("120363000000000000@newsletter")
	person, _ := normaliseRecipient("85291234567")
	if recipientAllowed(status) || recipientAllowed(channel) {
		t.Error("Status/channel allowed with an empty allowlist")
	}
	if !recipientAllowed(person) {
		t.Error("person refused with an empty allowlist")
	}
	sendAllowlist = map[string]bool{"status@broadcast": true}
	if !recipientAllowed(status) || recipientAllowed(person) {
		t.Error("explicit allowlist not applied")
	}
	sendAllowlist = map[string]bool{}
}

func TestTerminalSafe(t *testing.T) {
	got := terminalSafe("hi\x1b[2K\rfake prompt‮\x9b")
	if strings.ContainsAny(got, "\x1b\r‮\x9b") {
		t.Errorf("control characters survived: %q", got)
	}
}

func TestConfirmWithoutKeyboardRefuses(t *testing.T) {
	confirmSends, stdinLines = true, nil
	if ok, _ := confirmSend("x", "y", ""); ok {
		t.Error("send allowed without a window to confirm in")
	}
}

func TestConfirmYesAndNo(t *testing.T) {
	confirmSends = true
	stdinLines = make(chan string, 1)
	stdinLines <- "y" // typed BEFORE the question: must be ignored
	go func() { time.Sleep(50 * time.Millisecond); stdinLines <- "n" }()
	if ok, _ := confirmSend("x", "hello", ""); ok {
		t.Error("stale 'y' or 'n' was accepted as yes")
	}
	go func() { time.Sleep(50 * time.Millisecond); stdinLines <- "y" }()
	if ok, _ := confirmSend("x", "hello", ""); !ok {
		t.Error("'y' was not accepted")
	}
	stdinLines = nil
}

func TestReusedIDCannotRewriteHistory(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)
	ms, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	defer ms.Close()
	ms.StoreChat("c@g.us", "Group", time.Now())
	ms.StoreMessage("ID1", "c@g.us", "me", "original", time.Now(), true, "", "", "", nil, nil, nil, 0)
	ms.StoreMessage("ID1", "c@g.us", "attacker", "fake", time.Now(), false, "", "", "", nil, nil, nil, 0)
	var content, sender string
	var fromMe bool
	if err := ms.db.QueryRow("SELECT content, sender, is_from_me FROM messages WHERE id='ID1'").Scan(&content, &sender, &fromMe); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if content != "original" || sender != "me" || !fromMe {
		t.Errorf("history rewritten: %q %q %v", content, sender, fromMe)
	}
	ms.StoreMessage("ID1", "c@g.us", "me", "edited", time.Now(), true, "", "", "", nil, nil, nil, 0)
	ms.db.QueryRow("SELECT content FROM messages WHERE id='ID1'").Scan(&content)
	if content != "edited" {
		t.Error("same sender could not update its own message")
	}
}

func TestTokenChangesEveryStart(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)
	os.MkdirAll("store", 0700)
	a, _ := newToken()
	b, _ := newToken()
	if a == b || len(a) < 64 {
		t.Error("token not renewed")
	}
}
