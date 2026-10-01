package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeLaunchOptionsPreservesUserArguments(t *testing.T) {
	for _, item := range []struct{ input, want string }{
		{"", ZB_LAUNCH_OPTIONS},
		{"-debug", ZB_LAUNCH_OPTIONS + " -debug"},
		{`-cachedir="D:\My Saves" -debug`, ZB_LAUNCH_OPTIONS + ` -cachedir="D:\My Saves" -debug`},
		{`-cachedir="D:\Folder -- Name"`, ZB_LAUNCH_OPTIONS + ` -cachedir="D:\Folder -- Name"`},
		{`-javaagent:Other.jar -- -cachedir=D:\Saves`, ZB_LAUNCH_ARG + ` -javaagent:Other.jar -- -cachedir=D:\Saves`},
		{`-agentlib:zbNative=policy=deny-new -- -debug`, `-agentlib:zbNative=policy=deny-new -- -debug`},
	} {
		got, err := mergeSteamLaunchOptions(item.input)
		if err != nil || got != item.want {
			t.Fatalf("input %q: got %q, %v; want %q", item.input, got, err, item.want)
		}
		again, err := mergeSteamLaunchOptions(got)
		if err != nil || again != got {
			t.Fatalf("not idempotent: %q, %v", again, err)
		}
	}
	for _, input := range []string{`wrapper %command%`, `-javaagent:Other.jar`, `-cachedir="unfinished`, "-debug\n-x"} {
		if _, err := mergeSteamLaunchOptions(input); err == nil {
			t.Fatalf("expected refusal for %q", input)
		}
	}
}

func TestVDFRoundTripPreservesQuotedPaths(t *testing.T) {
	path := filepath.Join(t.TempDir(), "localconfig.vdf")
	data := `"UserLocalConfigStore" { "Software" { "Valve" { "Steam" { "apps" {
"108600"
{
    "LaunchOptions" "-debug"
}
} } } } }`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	want := ZB_LAUNCH_OPTIONS + ` -cachedir="D:\My Saves" -debug`
	if err := manualPatchVDF(path, "-debug", want); err != nil {
		t.Fatal(err)
	}
	got, err := readPZLaunchOptions(path)
	if err != nil || got != want {
		t.Fatalf("got %q, %v; want %q", got, err, want)
	}
}
