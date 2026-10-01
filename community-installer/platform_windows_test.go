package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindowsJunctionRefused(t *testing.T) {
	loc, payload, hash := fixture(t)
	destination := filepath.Join(t.TempDir(), "preserve")
	if err := os.Mkdir(destination, 0755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(loc.Profile, "mods", "ZombieBuddy")
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
	command := "New-Item -ItemType Junction -Path " + quote(link) + " -Value " + quote(destination) + " -ErrorAction Stop | Out-Null"
	if output, err := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command).CombinedOutput(); err != nil {
		t.Fatalf("cannot create fixture junction: %v %s", err, output)
	}
	if _, err := buildPlan(loc, payload, hash); err == nil {
		t.Fatal("junction target accepted")
	}
	if _, err := os.Stat(destination); err != nil {
		t.Fatal("junction target lost")
	}
}

func TestConcurrentInstallerLock(t *testing.T) {
	game := t.TempDir()
	release, err := acquireInstallationLock(game)
	if err != nil {
		t.Fatal(err)
	}
	if second, err := acquireInstallationLock(game); err == nil {
		second()
		release()
		t.Fatal("second installer acquired the same directory")
	}
	release()
	after, err := acquireInstallationLock(game)
	if err != nil {
		t.Fatal("lock not released", err)
	}
	after()
}

func TestWindowsProcessEnumeration(t *testing.T) {
	names, err := runningProcessNames()
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if strings.EqualFold(name, filepath.Base(executable)) {
			return
		}
	}
	t.Fatal("native snapshot did not include this test process")
}

func TestProcessGuardNames(t *testing.T) {
	for _, name := range []string{"steam.exe", "STEAM.EXE", "java.exe", "javaw.exe", "ProjectZomboid64.exe"} {
		if checkProcessNames([]string{name}) == nil {
			t.Fatal("running application accepted", name)
		}
	}
	if checkProcessNames(nil) == nil {
		t.Fatal("empty process list accepted")
	}
	if err := checkProcessNames([]string{"explorer.exe", "ZombieBuddyCommunityInstaller.exe"}); err != nil {
		t.Fatal(err)
	}
}
