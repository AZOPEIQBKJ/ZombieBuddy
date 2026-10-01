//go:build windows

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func isReparse(info os.FileInfo) bool {
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && attr.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0
}
func hiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func stoppedGuard() error {
	output, err := hiddenCommand("tasklist.exe", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return fmt.Errorf("could not check running processes; no installation changes are allowed")
	}
	rows, err := csv.NewReader(strings.NewReader(string(output))).ReadAll()
	if err != nil || len(rows) == 0 {
		return fmt.Errorf("could not parse running processes")
	}
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		name := strings.ToLower(row[0])
		if name == "steam.exe" || name == "java.exe" || name == "javaw.exe" || strings.HasPrefix(name, "projectzomboid") {
			return fmt.Errorf("close Steam, Project Zomboid and Java processes normally before installing or rolling back (still running: %s)", row[0])
		}
	}
	for _, name := range []string{"JAVA_TOOL_OPTIONS", "_JAVA_OPTIONS", "JDK_JAVA_OPTIONS"} {
		value := strings.ToLower(os.Getenv(name))
		if strings.Contains(value, "zombiebuddy") || strings.Contains(value, "zbnative") {
			return fmt.Errorf("%s also configures ZombieBuddy; review it before installation", name)
		}
	}
	return nil
}

func registrySteam() string {
	output, err := hiddenCommand("reg.exe", "query", `HKCU\Software\Valve\Steam`, "/v", "SteamPath").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(output), "\n") {
		if _, value, ok := strings.Cut(line, "REG_SZ"); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// The handle keeps a named object alive; process exit automatically releases it.
// No owned Win32 mutex is used, so goroutine thread scheduling cannot lose ownership.
func acquireInstallationLock(game string) (func(), error) {
	name, err := syscall.UTF16PtrFromString(`Local\ZombieBuddyCommunity-` + digest([]byte(strings.ToLower(filepath.Clean(game))))[:24])
	if err != nil {
		return nil, err
	}
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("CreateMutexW")
	handle, _, callErr := proc.Call(0, 0, uintptr(unsafe.Pointer(name)))
	if handle == 0 {
		return nil, fmt.Errorf("cannot acquire installation lock: %v", callErr)
	}
	closeHandle := func() { _ = syscall.CloseHandle(syscall.Handle(handle)) }
	if callErr == syscall.ERROR_ALREADY_EXISTS {
		closeHandle()
		return nil, fmt.Errorf("another Community installer is using this game directory")
	}
	return closeHandle, nil
}
