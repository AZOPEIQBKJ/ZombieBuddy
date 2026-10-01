//go:build !windows

package main

import (
	"fmt"
	"os"
)

func isReparse(info os.FileInfo) bool { return false }
func stoppedGuard() error             { return fmt.Errorf("this installer supports Windows only") }
func registrySteam() string           { return "" }

func acquireInstallationLock(game string) (func(), error) {
	return nil, fmt.Errorf("this installer supports Windows only")
}
