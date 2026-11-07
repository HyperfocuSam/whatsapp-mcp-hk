//go:build windows

package main

// Windows has no umask; files inherit the user profile's permissions.
func setPrivateUmask() {}
