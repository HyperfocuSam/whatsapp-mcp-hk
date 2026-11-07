package main

// confirm.go — every send must be approved by the person at the computer, in
// the bridge's own window. A WhatsApp message (or a tricked AI) can't press a
// key in that window, so this holds even if the AI app auto-approves tools.
// Developers can turn it off with WHATSAPP_CONFIRM_SENDS=0.

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const confirmTimeout = 110 * time.Second

var (
	confirmSends = true
	stdinLines   chan string // nil when there is no keyboard (not started from a window)
	confirmMu    sync.Mutex
)

func startConfirmReader() {
	confirmSends = strings.TrimSpace(os.Getenv("WHATSAPP_CONFIRM_SENDS")) != "0"
	if !confirmSends {
		return
	}
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return
	}
	stdinLines = make(chan string)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			stdinLines <- sc.Text()
		}
	}()
}

// terminalSafe shows text exactly but harmlessly: control and direction
// characters can't move the cursor or fake the prompt.
func terminalSafe(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' {
			return r
		}
		if r == '\t' {
			return ' '
		}
		if isUnsafeRune(r) || (r >= 0x80 && r <= 0x9f) {
			return '?'
		}
		return r
	}, s)
	return strings.ReplaceAll(s, "\n", "\n          ")
}

// confirmSend asks in the bridge window. Returns ok and, if not ok, a reason for the AI.
func confirmSend(to, text, file string) (bool, string) {
	if !confirmSends {
		return true, ""
	}
	if stdinLines == nil {
		return false, "Sending needs a confirmation in the WhatsApp window, but the bridge was not started from a window. Ask the user to start it with 'Start WhatsApp - can send'."
	}
	confirmMu.Lock()
	defer confirmMu.Unlock()
	for drained := false; !drained; { // ignore anything typed before this question
		select {
		case <-stdinLines:
		default:
			drained = true
		}
	}
	fmt.Print("\a\n──────────────────────────────────────────────\n")
	fmt.Println("  The AI wants to SEND this on WhatsApp:")
	fmt.Printf("  To:     %s\n", terminalSafe(to))
	if file != "" {
		fmt.Printf("  File:   %s\n", terminalSafe(filepath.Base(file)))
	}
	if text != "" {
		fmt.Printf("  Text:   %s\n", terminalSafe(text))
	}
	fmt.Print("  Type y then Enter to send. Just Enter cancels. > ")
	select {
	case line := <-stdinLines:
		if a := strings.ToLower(strings.TrimSpace(line)); a == "y" || a == "yes" {
			fmt.Println("  → Sending.")
			return true, ""
		}
		fmt.Println("  → Cancelled.")
		return false, "The user said no in the WhatsApp window. Nothing was sent."
	case <-time.After(confirmTimeout):
		fmt.Println("\n  → No answer in time. Cancelled.")
		return false, "No confirmation in the WhatsApp window within 2 minutes. Nothing was sent."
	}
}
