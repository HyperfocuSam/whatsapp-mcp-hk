//go:build darwin

package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

// markQuarantined gives downloaded attachments the same "downloaded from the
// internet" mark a browser would, so macOS Gatekeeper checks any app inside
// before it can run.
func markQuarantined(path string) {
	v := fmt.Sprintf("0081;%x;WhatsApp MCP;", time.Now().Unix())
	_ = unix.Setxattr(path, "com.apple.quarantine", []byte(v), 0)
}
