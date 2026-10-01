package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func writeTest(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}
func readTest(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func noProcesses() error { return nil }

func fixture(t *testing.T) (locations, map[string][]byte, string) {
	t.Helper()
	root := t.TempDir()
	loc := locations{filepath.Join(root, "Steam Library", "ProjectZomboid"), filepath.Join(root, "Player É", "Zomboid"), filepath.Join(root, "Steam", "userdata", "12345", "config", "localconfig.vdf")}
	writeTest(t, filepath.Join(loc.Game, "projectzomboid.jar"), []byte("test engine"))
	writeTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"), []byte(`{"mainClass":"zombie/gameStates/MainScreenState","classpath":[".","projectzomboid.jar"],"vmArgs":["-agentlib:zbNative","-Xmx3072m","-Dunknown=value"],"windows":{"10":{"vmArgs":["-XX:+UseZGC"]}},"futureField":{"keep":true}}`))
	writeTest(t, loc.SteamConfig, []byte(`"UserLocalConfigStore" { "Software" { "Valve" { "Steam" { "Apps" { "42" { "LaunchOptions" "untouched" } "108600" { "LaunchOptions" "-agentlib:zbNative=patches_jar=AftermathLHCompat4221.jar:aftermathsystems.lhcompat -- -debug -cachedir=\"`+strings.ReplaceAll(loc.Profile, `\`, `\\`)+`\"" "LastPlayed" "123" } } } } } }`))
	writeTest(t, filepath.Join(loc.Profile, "Saves", "existing.bin"), []byte("save to retain"))
	writeTest(t, filepath.Join(loc.Profile, "mods", "AnotherMod", "42", "mod.info"), []byte("id=AnotherMod\n"))
	payload := map[string][]byte{"42/mod.info": []byte("id=ZombieBuddy\nname=ZombieBuddy Community\n"), "common/.keep": {}, "libs/ZombieBuddy.jar": []byte("new framework"), "libs/ZombieBuddy.jar.zbs": []byte("test sidecar")}
	return loc, payload, digest([]byte("test engine"))
}

func TestEmbeddedRuntimeIntegrity(t *testing.T) {
	payload, err := payloadFiles(embeddedPackage)
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) != 9 || digest(payload["libs/ZombieBuddy.jar"]) != "14050862a63b55252816aba580dc4206ef540e83b494b16eae15112f62ec7445" {
		t.Fatal("wrong runtime")
	}
	modified := append([]byte(nil), embeddedPackage...)
	modified[len(modified)/2] ^= 1
	if _, err := payloadFiles(modified); err == nil {
		t.Fatal("tampered payload accepted")
	}
}

func TestMigrationFreshAndRepeatRollback(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(map[bool]string{false: "fresh", true: "original"}[existing], func(t *testing.T) {
			loc, payload, hash := fixture(t)
			if existing {
				writeTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"), []byte("original framework"))
				writeTest(t, filepath.Join(loc.Game, "zbNative.dll"), embeddedNativeLoader)
			}
			beforeJSON := readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))
			beforeSteam := readTest(t, loc.SteamConfig)
			p, err := buildPlan(loc, payload, hash)
			if err != nil {
				t.Fatal(err)
			}
			receiptPath, err := applyPlan(p, noProcesses, nil)
			if err != nil {
				t.Fatal(err)
			}
			if digest(readTest(t, filepath.Join(loc.Game, "zbNative.dll"))) != nativeLoaderHash {
				t.Fatal("native bootstrap missing or changed")
			}
			p2, err := buildPlan(loc, payload, hash)
			if err != nil {
				t.Fatal(err)
			}
			if len(p2.Changes) != 0 {
				t.Fatalf("repeat changed %d files", len(p2.Changes))
			}
			options, _, err := lookupVDF(readTest(t, loc.SteamConfig), optionsPath)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(options, "agentlib") || !strings.Contains(options, "-debug") || !strings.Contains(options, "-modfolders mods,workshop,steam") {
				t.Fatal(options)
			}
			var config map[string]any
			if err := json.Unmarshal(readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json")), &config); err != nil {
				t.Fatal(err)
			}
			if config["futureField"] == nil {
				t.Fatal("unknown JSON field lost")
			}
			if err := rollback(receiptPath, payload, noProcesses); err != nil {
				t.Fatal(err)
			}
			if err := rollback(receiptPath, payload, noProcesses); err != nil {
				t.Fatal("repeat rollback", err)
			}
			if !bytes.Equal(beforeJSON, readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))) || !bytes.Equal(beforeSteam, readTest(t, loc.SteamConfig)) {
				t.Fatal("original launch settings not restored exactly")
			}
			if existing {
				if string(readTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"))) != "original framework" {
					t.Fatal("original JAR lost")
				}
				if !bytes.Equal(readTest(t, filepath.Join(loc.Game, "zbNative.dll")), embeddedNativeLoader) {
					t.Fatal("native DLL changed")
				}
			} else if _, err := os.Stat(filepath.Join(loc.Game, "ZombieBuddy.jar")); !os.IsNotExist(err) {
				t.Fatal("fresh JAR not removed")
			}
			if !existing {
				if _, err := os.Stat(filepath.Join(loc.Game, "zbNative.dll")); !os.IsNotExist(err) {
					t.Fatal("fresh native loader not removed on rollback")
				}
			}
			if string(readTest(t, filepath.Join(loc.Profile, "Saves", "existing.bin"))) != "save to retain" {
				t.Fatal("save changed")
			}
			if _, err := os.Stat(filepath.Join(loc.Profile, "mods", "ZombieBuddy")); !os.IsNotExist(err) {
				t.Fatal("new mod directory not removed")
			}
		})
	}
}

func TestInterruptedInstallationEveryWrite(t *testing.T) {
	loc, payload, hash := fixture(t)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	for stop := 0; stop < len(p.Changes); stop++ {
		t.Run(p.Changes[stop].Area+"-"+p.Changes[stop].Name, func(t *testing.T) {
			loc, payload, hash := fixture(t)
			p, err := buildPlan(loc, payload, hash)
			if err != nil {
				t.Fatal(err)
			}
			beforeJSON := readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))
			beforeSteam := readTest(t, loc.SteamConfig)
			receiptPath, err := applyPlan(p, noProcesses, func(i int) error {
				if i == stop {
					return errors.New("simulated interruption")
				}
				return nil
			})
			if err == nil || receiptPath == "" {
				t.Fatal("interruption not recorded")
			}
			if _, err := buildPlan(loc, payload, hash); err == nil {
				t.Fatal("pending transaction ignored")
			}
			if err := rollback(receiptPath, payload, noProcesses); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(beforeJSON, readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))) || !bytes.Equal(beforeSteam, readTest(t, loc.SteamConfig)) {
				t.Fatal("interruption restoration mismatch")
			}
			if _, err := buildPlan(loc, payload, hash); err != nil {
				t.Fatal("recovery did not unblock installation", err)
			}
		})
	}
}

func TestRollbackKeepsLaterSteamActivity(t *testing.T) {
	for _, initialOptions := range []bool{true, false} {
		t.Run(map[bool]string{true: "existing-options", false: "absent-options"}[initialOptions], func(t *testing.T) {
			loc, payload, hash := fixture(t)
			if !initialOptions {
				data, err := changeVDF(readTest(t, loc.SteamConfig), optionsPath, nil)
				if err != nil {
					t.Fatal(err)
				}
				writeTest(t, loc.SteamConfig, data)
			}
			old, oldExists, _ := lookupVDF(readTest(t, loc.SteamConfig), optionsPath)
			p, err := buildPlan(loc, payload, hash)
			if err != nil {
				t.Fatal(err)
			}
			receiptPath, err := applyPlan(p, noProcesses, nil)
			if err != nil {
				t.Fatal(err)
			}
			after := strings.Replace(string(readTest(t, loc.SteamConfig)), `"LastPlayed" "123"`, `"LastPlayed" "999"`, 1)
			writeTest(t, loc.SteamConfig, []byte(after))
			if err := rollback(receiptPath, payload, noProcesses); err != nil {
				t.Fatal(err)
			}
			restored := readTest(t, loc.SteamConfig)
			now, exists, err := lookupVDF(restored, optionsPath)
			if err != nil || exists != oldExists || now != old {
				t.Fatal("old options not restored")
			}
			if !bytes.Contains(restored, []byte(`"LastPlayed" "999"`)) {
				t.Fatal("later Steam state overwritten")
			}
			other, _, _ := lookupVDF(restored, []string{"UserLocalConfigStore", "Software", "Valve", "Steam", "Apps", "42", "LaunchOptions"})
			if other != "untouched" {
				t.Fatal("other app changed")
			}
		})
	}
}

func TestConcurrentAndUnexpectedChangesArePreserved(t *testing.T) {
	t.Run("after-preview", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		p, err := buildPlan(loc, payload, hash)
		if err != nil {
			t.Fatal(err)
		}
		writeTest(t, loc.SteamConfig, append(readTest(t, loc.SteamConfig), []byte("\n// later edit")...))
		if _, err := applyPlan(p, noProcesses, nil); err == nil {
			t.Fatal("stale plan applied")
		}
		if _, err := os.Stat(filepath.Join(loc.Game, "ZombieBuddy.jar")); !os.IsNotExist(err) {
			t.Fatal("partial mutation before precheck")
		}
	})
	t.Run("after-install", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		p, err := buildPlan(loc, payload, hash)
		if err != nil {
			t.Fatal(err)
		}
		receiptPath, err := applyPlan(p, noProcesses, nil)
		if err != nil {
			t.Fatal(err)
		}
		writeTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"), []byte("user replacement"))
		if err := rollback(receiptPath, payload, noProcesses); err == nil {
			t.Fatal("user JAR overwritten")
		}
		if string(readTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"))) != "user replacement" {
			t.Fatal("lost user JAR")
		}
	})
	t.Run("changed-steam-options", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		p, err := buildPlan(loc, payload, hash)
		if err != nil {
			t.Fatal(err)
		}
		receiptPath, err := applyPlan(p, noProcesses, nil)
		if err != nil {
			t.Fatal(err)
		}
		value := "-- -debug -new-option"
		data, err := changeVDF(readTest(t, loc.SteamConfig), optionsPath, &value)
		if err != nil {
			t.Fatal(err)
		}
		writeTest(t, loc.SteamConfig, data)
		if err := rollback(receiptPath, payload, noProcesses); err == nil {
			t.Fatal("new options overwritten")
		}
		if !bytes.Equal(data, readTest(t, loc.SteamConfig)) {
			t.Fatal("Steam changed")
		}
	})
	t.Run("bad-backup", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		writeTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"), []byte("existing runtime to back up"))
		p, err := buildPlan(loc, payload, hash)
		if err != nil {
			t.Fatal(err)
		}
		receiptPath, err := applyPlan(p, noProcesses, nil)
		if err != nil {
			t.Fatal(err)
		}
		damaged := false
		for i, c := range p.Changes {
			if c.BeforeExists && c.Area != "steam" {
				writeTest(t, filepath.Join(filepath.Dir(receiptPath), "before-"+strconv.Itoa(i)), []byte("damaged backup"))
				damaged = true
				break
			}
		}
		if !damaged {
			t.Fatal("fixture did not exercise a damaged file backup")
		}
		if err := rollback(receiptPath, payload, noProcesses); err == nil {
			t.Fatal("damaged backup accepted")
		}
	})
}

func TestSafetyPreflight(t *testing.T) {
	cases := map[string]func(*testing.T, locations){
		"pending-update": func(t *testing.T, l locations) {
			writeTest(t, filepath.Join(l.Game, "ZombieBuddy.jar.new"), []byte("pending"))
		},
		"wrong-game": func(t *testing.T, l locations) {
			writeTest(t, filepath.Join(l.Game, "projectzomboid.jar"), []byte("other build"))
		},
		"unexpected-managed-file": func(t *testing.T, l locations) {
			writeTest(t, filepath.Join(l.Profile, "mods", "ZombieBuddy", "personal.txt"), []byte("keep"))
		},
		"extra-local-copy": func(t *testing.T, l locations) {
			writeTest(t, filepath.Join(l.Profile, "mods", "OldName", "42", "mod.info"), []byte("id=ZombieBuddy\n"))
		},
		"steam-wrapper": func(t *testing.T, l locations) {
			value := "custom.exe %command%"
			data, err := changeVDF(readTest(t, l.SteamConfig), optionsPath, &value)
			if err != nil {
				t.Fatal(err)
			}
			writeTest(t, l.SteamConfig, data)
		},
	}
	for name, modify := range cases {
		t.Run(name, func(t *testing.T) {
			loc, payload, hash := fixture(t)
			modify(t, loc)
			if _, err := buildPlan(loc, payload, hash); err == nil {
				t.Fatal("unsafe plan accepted")
			}
		})
	}
	t.Run("process-still-running", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		p, err := buildPlan(loc, payload, hash)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := applyPlan(p, func() error { return errors.New("Steam running") }, nil); err == nil {
			t.Fatal("guard ignored")
		}
		if _, err := os.Stat(filepath.Join(loc.Game, ".zbc-installations")); !os.IsNotExist(err) {
			t.Fatal("guard made backups/writes")
		}
	})
	t.Run("symlink", func(t *testing.T) {
		loc, payload, hash := fixture(t)
		target := filepath.Join(t.TempDir(), "mod")
		os.Mkdir(target, 0755)
		link := filepath.Join(loc.Profile, "mods", "ZombieBuddy")
		if err := os.Symlink(target, link); err != nil {
			t.Skip("OS does not allow unprivileged symlinks")
		}
		if _, err := buildPlan(loc, payload, hash); err == nil {
			t.Fatal("linked target accepted")
		}
	})
}

func TestOptionsAndJSONContracts(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "Player With Spaces")
	for _, s := range []string{"", "-debug", `-Dname="two words" -- -debug`, `-- -modfolders steam,mods,workshop -debug`} {
		result, tails, err := configureOptions(s, profile)
		if err != nil {
			t.Fatal(err)
		}
		if len(tails) != 0 || !strings.Contains(result, "-modfolders mods,") {
			t.Fatal(result)
		}
		repeated, _, err := configureOptions(result, profile)
		if err != nil || repeated != result {
			t.Fatal("not idempotent", result, repeated, err)
		}
		if strings.Contains(s, `-Dname="two words"`) && !strings.Contains(result, `-Dname="two words"`) {
			t.Fatal("quoted argument changed")
		}
	}
	for _, s := range []string{`-Xmx4g -debug`, `-- --`, `-- -modfolders`, `-- -modfolders mods,mods`, `-- -modfolders custom`, `-debug "broken`, `%command% -- -debug`, `-- -cachedir=relative`} {
		if _, _, err := configureOptions(s, profile); err == nil {
			t.Fatal("ambiguous options accepted", s)
		}
	}
	merged, err := mergeAgentOptions([]string{`policy=deny-new,config_dir=C:\Profile With Spaces`, `policy=deny-new,patches_jar=C:\Old\AftermathLHCompat4221.jar:aftermathsystems.lhcompat;C:\Keep.jar:my.patch`})
	if err != nil || strings.Contains(merged, "AftermathLHCompat") || !strings.Contains(merged, `C:\Keep.jar:my.patch`) || !strings.Contains(merged, `config_dir=C:\Profile With Spaces`) {
		t.Fatal(merged, err)
	}
	if _, err := mergeAgentOptions([]string{"policy=deny-new", "policy=allow-all"}); err == nil {
		t.Fatal("conflicting policy discarded")
	}
	data := []byte(`{"vmArgs":["-javaagent:C:\\Mods Folder\\ZombieBuddy.jar=policy=deny-new","-javaagent:Another.jar","-Xmx2g"],"windows":{"10":{"vmArgs":["-agentlib:zbNative=policy=deny-new","-XX:+UseZGC"],"unknown":"keep"}},"unknown":[1,2]}`)
	after, err := configureLauncher(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(after), "-agentlib:zbNative") != 1 || strings.Contains(string(after), "ZombieBuddy.jar") || !strings.Contains(string(after), "Another.jar") || !strings.Contains(string(after), "unknown") {
		t.Fatal(string(after))
	}
	again, err := configureLauncher(after, nil)
	if err != nil || !bytes.Equal(again, after) {
		t.Fatal("JSON update not idempotent")
	}
}

func TestVDFExactScopeEscapesAndCreation(t *testing.T) {
	input := []byte("// leading comment\r\n\"UserLocalConfigStore\" { \"Software\" { \"Valve\" { \"Steam\" { \"Apps\" {} } } } } // tail\n")
	value := `-- -cachedir="C:\Player With Spaces\Zomboid" -debug`
	after, err := changeVDF(input, optionsPath, &value)
	if err != nil {
		t.Fatal(err)
	}
	got, exists, err := lookupVDF(after, optionsPath)
	if err != nil || !exists || got != value {
		t.Fatal(got, exists, err)
	}
	if !bytes.Contains(after, []byte("// leading comment\r\n")) || !bytes.Contains(after, []byte("// tail\n")) {
		t.Fatal("comments lost")
	}
	if _, err := changeVDF([]byte(`"x" "1" "X" "2"`), []string{"x"}, &value); err == nil {
		t.Fatal("duplicate key accepted")
	}
	for _, bad := range []string{`"unterminated`, `"x" {`, `}`, `"x" }`} {
		if _, err := parseVDF(bad); err == nil {
			t.Fatal("malformed VDF accepted", bad)
		}
	}
}
