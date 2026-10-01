package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNativeLoaderIntegrityAndUnknownExistingRefusal(t *testing.T) {
	if err := verifyNativeLoader(embeddedNativeLoader); err != nil {
		t.Fatal(err)
	}
	modified := append([]byte(nil), embeddedNativeLoader...)
	modified[len(modified)/2] ^= 1
	if verifyNativeLoader(modified) == nil {
		t.Fatal("modified embedded native bootstrap accepted")
	}
	loc, payload, hash := fixture(t)
	writeTest(t, filepath.Join(loc.Game, "zbNative.dll"), modified)
	if _, err := buildPlan(loc, payload, hash); err == nil || !strings.Contains(err.Error(), "existing zbNative.dll differs") {
		t.Fatal("unknown existing bootstrap accepted", err)
	}
	if _, err := os.Stat(filepath.Join(loc.Game, "ZombieBuddy.jar")); !os.IsNotExist(err) {
		t.Fatal("refusal changed installed framework")
	}
}

func TestNativeLoaderChangedAfterPreviewRefusesWrites(t *testing.T) {
	loc, payload, hash := fixture(t)
	writeTest(t, filepath.Join(loc.Game, "zbNative.dll"), embeddedNativeLoader)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	writeTest(t, filepath.Join(loc.Game, "zbNative.dll"), []byte("unexpected later change"))
	if _, err := applyPlan(p, noProcesses, nil); err == nil {
		t.Fatal("changed native loader accepted")
	}
	if _, err := os.Stat(filepath.Join(loc.Game, "ZombieBuddy.jar")); !os.IsNotExist(err) {
		t.Fatal("partially installed after native loader changed")
	}
}

func TestPreview2RepairAndExactRollback(t *testing.T) {
	loc, payload, hash := fixture(t)
	before := []byte(`{"vmArgs":["-javaagent:ZombieBuddy.jar=config_dir=C:/QA,policy=deny-new","-javaagent:Other.jar","-Xmx3g"],"keep":true}`)
	writeTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"), before)
	writeTest(t, filepath.Join(loc.Game, "zbNative.dll"), embeddedNativeLoader)
	p, err := buildPlan(loc, payload, hash)
	if err != nil {
		t.Fatal(err)
	}
	receiptPath, err := applyPlan(p, noProcesses, nil)
	if err != nil {
		t.Fatal(err)
	}
	after := string(readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json")))
	if !strings.Contains(after, "-agentlib:zbNative=config_dir=C:/QA,policy=deny-new") ||
		strings.Contains(after, "-javaagent:ZombieBuddy.jar") || !strings.Contains(after, "-javaagent:Other.jar") {
		t.Fatal("repair lost configuration or kept failed route", after)
	}
	if err := rollback(receiptPath, payload, noProcesses); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, readTest(t, filepath.Join(loc.Game, "ProjectZomboid64.json"))) {
		t.Fatal("preview.2 baseline not restored exactly")
	}
}
