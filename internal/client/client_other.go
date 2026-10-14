//go:build !windows

package client

import "fmt"

// Run reports that this OS build supports only the server mode.
func Run() error { return fmt.Errorf("desktop client requires Windows") }
