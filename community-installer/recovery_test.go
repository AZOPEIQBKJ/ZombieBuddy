package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestSteamBackupContainsOnlyPZEdit(t *testing.T) {
	loc, payload, hash := fixture(t)
	secretMarker := []byte("UNRELATED_ACCOUNT_DATA_MUST_STAY_IN_STEAM")
	steam := append(readTest(t, loc.SteamConfig), append([]byte("\n// "), secretMarker...)...)
	writeTest(t, loc.SteamConfig, steam)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Dir(receiptPath))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		data := readTest(t, filepath.Join(filepath.Dir(receiptPath), entry.Name()))
		if bytes.Contains(data, secretMarker) {
			t.Fatal("unrelated Steam data copied into backup")
		}
	}
	for i, c := range p.Changes {
		if c.Area == "steam" {
			path := filepath.Join(filepath.Dir(receiptPath), "steam-"+strconv.Itoa(i)+".json")
			var delta steamPatch
			if err := json.Unmarshal(readTest(t, path), &delta); err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(delta.Removed, secretMarker) || bytes.Contains(delta.Inserted, secretMarker) {
				t.Fatal("unrelated data inside edit")
			}
			writeTest(t, path, []byte("damaged edit"))
		}
	}
	if err := rollback(receiptPath, payload, noProcesses); err == nil {
		t.Fatal("corrupt Steam edit accepted")
	}
}

func TestMissingSteamAppRestoresExactFile(t *testing.T) {
	loc, payload, hash := fixture(t)
	original := []byte(`"UserLocalConfigStore" { "Software" { "Valve" { "Steam" { "Apps" { "42" { "LaunchOptions" "keep" } } } } } }`)
	writeTest(t, loc.SteamConfig, original)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := rollback(receiptPath, payload, noProcesses); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, readTest(t, loc.SteamConfig)) {
		t.Fatal("new app block not removed exactly")
	}
}

func TestInterruptedRollbackCanResume(t *testing.T) {
	loc, payload, hash := fixture(t)
	originalJSON := readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))
	originalSteam := readTest(t, loc.SteamConfig)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	err = rollback(receiptPath, payload, func() error {
		calls++
		if calls == 4 {
			return errors.New("interrupted rollback")
		}
		return nil
	})
	if err == nil {
		t.Fatal("interruption ignored")
	}
	if _, err := buildPlan(loc, payload, hash); err == nil {
		t.Fatal("incomplete rollback ignored")
	}
	if err := rollback(receiptPath, payload, noProcesses); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(originalJSON, readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))) || !bytes.Equal(originalSteam, readTest(t, loc.SteamConfig)) {
		t.Fatal("resumed rollback lost original settings")
	}
}

func TestRollbackRejectsUnmanagedReceipt(t *testing.T) {
	loc, payload, hash := fixture(t)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	var r receipt
	json.Unmarshal(readTest(t, receiptPath), &r)
	r.Files[0].Name = "../../Saves/existing.bin"
	data, _ := json.Marshal(r)
	writeTest(t, receiptPath, data)
	if err := rollback(receiptPath, payload, noProcesses); err == nil {
		t.Fatal("unmanaged path accepted")
	}
	if string(readTest(t, filepath.Join(loc.Profile, "Saves", "existing.bin"))) != "save to retain" {
		t.Fatal("save changed")
	}
}

func TestEngineAndPendingUpdateRecheckedAfterPreview(t *testing.T) {
	for _, filename := range []string{"projectzomboid.jar", "ZombieBuddy.jar.new"} {
		t.Run(filename, func(t *testing.T) {
			loc, payload, hash := fixture(t)
			p, err := buildPlan(loc, payload, hash)
			if err != nil {
				t.Fatal(err)
			}
			writeTest(t, filepath.Join(loc.Game, filename), []byte("changed after preview"))
			if _, err := applyPlan(p, noProcesses, nil); err == nil {
				t.Fatal("changed engine/pending update accepted")
			}
			if _, err := os.Stat(filepath.Join(loc.Game, "ZombieBuddy.jar")); !os.IsNotExist(err) {
				t.Fatal("partially installed despite failed recheck")
			}
		})
	}
}

func TestFrozenRuntimePayloadThroughMigration(t *testing.T) {
	loc, _, hash := fixture(t)
	payload, err := payloadFiles(embeddedPackage)
	if err != nil {
		t.Fatal(err)
	}
	writeTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"), []byte("original runtime backup"))
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range payload {
		if !bytes.Equal(data, readTest(t, filepath.Join(loc.Profile, "mods", "ZombieBuddy", filepath.FromSlash(name)))) {
			t.Fatal("deployed payload mismatch", name)
		}
	}
	if !bytes.Equal(payload["libs/ZombieBuddy.jar"], readTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"))) {
		t.Fatal("root agent mismatch")
	}
	if err := rollback(receiptPath, payload, noProcesses); err != nil {
		t.Fatal(err)
	}
	if string(readTest(t, filepath.Join(loc.Game, "ZombieBuddy.jar"))) != "original runtime backup" {
		t.Fatal("original agent not restored")
	}
}
