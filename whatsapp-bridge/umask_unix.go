//go:build !windows

package main

import "syscall"

// New files (databases, media, logs) are private to the current user.
func setPrivateUmask() { syscall.Umask(0o077) }
