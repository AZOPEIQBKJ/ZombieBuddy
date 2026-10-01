package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type locations struct {
	Game        string `json:"game"`
	Profile     string `json:"profile"`
	SteamConfig string `json:"steamConfig"`
}
type change struct {
	Area, Name    string
	Before, After []byte
	BeforeExists  bool
}
type plan struct {
	Locations locations
	Changes   []change
	Payload   map[string][]byte
	GameSHA   string
}
type fileRecord struct {
	Area           string `json:"area"`
	Name           string `json:"name"`
	Before         string `json:"before"`
	After          string `json:"after"`
	SteamPatchHash string `json:"steamPatchHash,omitempty"`
}
type receipt struct {
	Schema      int          `json:"schema"`
	Installer   string       `json:"installerVersion"`
	Runtime     string       `json:"runtimeVersion"`
	Status      string       `json:"status"`
	Locations   locations    `json:"locations"`
	Files       []fileRecord `json:"files"`
	CreatedDirs []string     `json:"createdDirectories"`
}

func samePath(a, b string) bool { return strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) }
func within(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}
func absolute(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("empty path")
	}
	return filepath.Abs(filepath.Clean(s))
}

func safePath(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("absolute path required: %s", path)
	}
	for current := path; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil {
			if info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
				return fmt.Errorf("linked/reparse path requires separate review: %s", current)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if parent := filepath.Dir(current); parent == current {
			break
		}
	}
	return nil
}

func readFile(path string) ([]byte, bool, error) {
	if err := safePath(path); err != nil {
		return nil, false, err
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !info.Mode().IsRegular() || info.Size() > 256<<20 {
		return nil, false, fmt.Errorf("unexpected file: %s", path)
	}
	data, err := os.ReadFile(path)
	return data, true, err
}
func fileDigest(path string) (string, error) {
	data, exists, err := readFile(path)
	if err != nil || !exists {
		return "", err
	}
	return digest(data), nil
}

func targetPath(loc locations, area, name string, payload map[string][]byte) (string, error) {
	switch area {
	case "game":
		if name == "ZombieBuddy.jar" || name == "ZombieBuddy.jar.zbs" || name == "ProjectZomboid64.json" {
			return filepath.Join(loc.Game, name), nil
		}
	case "mod":
		if _, ok := payload[name]; ok {
			return filepath.Join(loc.Profile, "mods", "ZombieBuddy", filepath.FromSlash(name)), nil
		}
	case "steam":
		if name == "LaunchOptions" {
			return loc.SteamConfig, nil
		}
	}
	return "", fmt.Errorf("unmanaged transaction target %s/%s", area, name)
}

func validateLocations(loc locations) error {
	for _, p := range []string{loc.Game, loc.Profile, loc.SteamConfig} {
		if err := safePath(p); err != nil {
			return err
		}
	}
	if within(loc.Game, loc.Profile) || within(loc.Profile, loc.Game) {
		return fmt.Errorf("game and user profile must be separate directories")
	}
	if info, err := os.Stat(filepath.Dir(loc.Profile)); err != nil || !info.IsDir() {
		return fmt.Errorf("select a profile under an existing parent directory")
	}
	if filepath.Base(loc.SteamConfig) != "localconfig.vdf" || filepath.Base(filepath.Dir(loc.SteamConfig)) != "config" || filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(loc.SteamConfig)))) != "userdata" {
		return fmt.Errorf("select Steam userdata/ACCOUNT/config/localconfig.vdf")
	}
	account := filepath.Base(filepath.Dir(filepath.Dir(loc.SteamConfig)))
	if account == "" || strings.Trim(account, "0123456789") != "" {
		return fmt.Errorf("invalid Steam account directory")
	}
	return nil
}

func buildPlan(loc locations, payload map[string][]byte, expectedGame string) (plan, error) {
	p := plan{Locations: loc, Payload: payload, GameSHA: expectedGame}
	if err := validateLocations(loc); err != nil {
		return p, err
	}
	if hash, err := fileDigest(filepath.Join(loc.Game, "projectzomboid.jar")); err != nil || hash != expectedGame {
		return p, fmt.Errorf("game JAR is not the audited PZ 42.21 build")
	}
	if _, exists, err := readFile(filepath.Join(loc.Game, "ZombieBuddy.jar.new")); err != nil || exists {
		return p, fmt.Errorf("pending or unreadable ZombieBuddy.jar.new; resolve it before installation")
	}
	if pending, err := pendingReceipts(loc.Game); err != nil {
		return p, err
	} else if len(pending) > 0 {
		return p, fmt.Errorf("unfinished installation: roll back %s before retrying", pending[0])
	}
	modDir := filepath.Join(loc.Profile, "mods", "ZombieBuddy")
	if err := inspectLocalMods(loc.Profile, payload); err != nil {
		return p, err
	}
	if err := filepath.WalkDir(modDir, func(path string, d os.DirEntry, walkErr error) error {
		if errors.Is(walkErr, os.ErrNotExist) && samePath(path, modDir) {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}
		if err := safePath(path); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(modDir, path)
		if _, ok := payload[filepath.ToSlash(rel)]; !ok {
			return fmt.Errorf("unexpected file in managed mod directory: %s", path)
		}
		return nil
	}); err != nil {
		return p, err
	}
	add := func(area, name string, after []byte) error {
		path, err := targetPath(loc, area, name, payload)
		if err != nil {
			return err
		}
		before, exists, err := readFile(path)
		if err != nil {
			return err
		}
		if exists && bytes.Equal(before, after) {
			return nil
		}
		p.Changes = append(p.Changes, change{area, name, before, after, exists})
		return nil
	}
	for _, name := range sortedKeys(payload) {
		if err := add("mod", name, payload[name]); err != nil {
			return p, err
		}
	}
	for _, name := range []string{"ZombieBuddy.jar", "ZombieBuddy.jar.zbs"} {
		if err := add("game", name, payload["libs/"+name]); err != nil {
			return p, err
		}
	}
	steam, exists, err := readFile(loc.SteamConfig)
	if err != nil || !exists {
		return p, fmt.Errorf("cannot read selected Steam account configuration")
	}
	// Require the selected account's Steam namespace, rather than accepting an arbitrary VDF.
	nodes, err := parseVDF(string(steam))
	if err != nil {
		return p, err
	}
	for _, key := range optionsPath[:4] {
		n, e := vdfChild(nodes, key)
		if e != nil || n == nil || n.scalar {
			return p, fmt.Errorf("invalid Steam account namespace")
		}
		nodes = n.children
	}
	current, _, err := lookupVDF(steam, optionsPath)
	if err != nil {
		return p, err
	}
	updated, tails, err := configureOptions(current, loc.Profile)
	if err != nil {
		return p, err
	}
	steamAfter, err := changeVDF(steam, optionsPath, &updated)
	if err != nil {
		return p, err
	}
	if err := add("steam", "LaunchOptions", steamAfter); err != nil {
		return p, err
	}
	launcher, exists, err := readFile(filepath.Join(loc.Game, "ProjectZomboid64.json"))
	if err != nil || !exists {
		return p, fmt.Errorf("normal Windows launcher JSON is missing")
	}
	launcherAfter, err := configureLauncher(launcher, tails)
	if err != nil {
		return p, err
	}
	if err := add("game", "ProjectZomboid64.json", launcherAfter); err != nil {
		return p, err
	}
	return p, nil
}

func inspectLocalMods(profile string, payload map[string][]byte) error {
	root := filepath.Join(profile, "mods")
	if err := safePath(root); err != nil {
		return err
	}
	dirs, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range dirs {
		dir := filepath.Join(root, entry.Name())
		if err := safePath(dir); err != nil {
			return err
		}
		if !entry.IsDir() || strings.EqualFold(entry.Name(), "ZombieBuddy") {
			continue
		}
		versions, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		candidates := []string{filepath.Join(dir, "mod.info")}
		for _, v := range versions {
			if v.IsDir() {
				candidates = append(candidates, filepath.Join(dir, v.Name(), "mod.info"))
			}
		}
		for _, path := range candidates {
			data, exists, err := readFile(path)
			if err != nil {
				return err
			}
			if !exists {
				continue
			}
			for _, line := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(strings.TrimSuffix(line, "\r")) == "id=ZombieBuddy" {
					return fmt.Errorf("another local ZombieBuddy copy exists: %s; preserve and review it first", dir)
				}
			}
		}
	}
	return nil
}

func atomicWrite(path string, data []byte) error {
	if err := safePath(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".zbc-write-*")
	if err != nil {
		return err
	}
	temp := f.Name()
	defer os.Remove(temp)
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(temp, path)
}

func saveReceipt(dir string, r receipt) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(dir, "receipt.json"), append(data, '\n'))
}

func pendingReceipts(game string) ([]string, error) {
	root := filepath.Join(game, ".zbc-installations")
	if err := safePath(root); err != nil {
		return nil, err
	}
	dirs, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range dirs {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(root, entry.Name(), "receipt.json")
		data, exists, err := readFile(path)
		if err != nil {
			return nil, err
		}
		if !exists {
			continue
		}
		var r receipt
		if err = json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("unreadable installation receipt: %s", path)
		}
		if r.Status != "INSTALLED" && r.Status != "ROLLED_BACK" {
			result = append(result, path)
		}
	}
	return result, nil
}

func applyPlan(p plan, guard func() error, beforeWrite func(int) error) (string, error) {
	release, err := acquireInstallationLock(p.Locations.Game)
	if err != nil {
		return "", err
	}
	defer release()
	if err := guard(); err != nil {
		return "", err
	}
	if len(p.Changes) == 0 {
		return "", nil
	}
	if hash, err := fileDigest(filepath.Join(p.Locations.Game, "projectzomboid.jar")); err != nil || hash != p.GameSHA {
		return "", fmt.Errorf("game build changed since preview")
	}
	if _, exists, err := readFile(filepath.Join(p.Locations.Game, "ZombieBuddy.jar.new")); err != nil || exists {
		return "", fmt.Errorf("pending framework replacement appeared since preview")
	}
	if pending, err := pendingReceipts(p.Locations.Game); err != nil {
		return "", err
	} else if len(pending) > 0 {
		return "", fmt.Errorf("another installation requires recovery: %s", pending[0])
	}
	for _, c := range p.Changes {
		path, _ := targetPath(p.Locations, c.Area, c.Name, p.Payload)
		actual, exists, err := readFile(path)
		if err != nil || exists != c.BeforeExists || !bytes.Equal(actual, c.Before) {
			return "", fmt.Errorf("file changed since preview: %s", path)
		}
	}
	root := filepath.Join(p.Locations.Game, ".zbc-installations")
	if err := safePath(root); err != nil {
		return "", err
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	dir, err := os.MkdirTemp(root, time.Now().UTC().Format("20060102T150405Z-")+"*")
	if err != nil {
		return "", err
	}
	r := receipt{Schema: 1, Installer: installerVersion, Runtime: runtimeVersion, Status: "PREPARED", Locations: p.Locations}
	created := map[string]bool{}
	for i, c := range p.Changes {
		record := fileRecord{Area: c.Area, Name: c.Name, After: digest(c.After)}
		if c.Area == "steam" {
			patch, err := makeSteamPatch(c.Before, c.After)
			if err != nil {
				return dir, err
			}
			record.Before = digest(c.Before)
			record.SteamPatchHash = digest(patch)
			if err := atomicWrite(filepath.Join(dir, fmt.Sprintf("steam-%d.json", i)), patch); err != nil {
				return dir, err
			}
		} else {
			if c.BeforeExists {
				record.Before = digest(c.Before)
				if err := atomicWrite(filepath.Join(dir, fmt.Sprintf("before-%d", i)), c.Before); err != nil {
					return dir, err
				}
			}
			if err := atomicWrite(filepath.Join(dir, fmt.Sprintf("after-%d", i)), c.After); err != nil {
				return dir, err
			}
		}
		r.Files = append(r.Files, record)
		target, _ := targetPath(p.Locations, c.Area, c.Name, p.Payload)
		for parent := filepath.Dir(target); ; parent = filepath.Dir(parent) {
			_, err := os.Stat(parent)
			if err == nil {
				break
			}
			if !errors.Is(err, os.ErrNotExist) {
				return dir, err
			}
			created[parent] = true
			if filepath.Dir(parent) == parent {
				return dir, fmt.Errorf("no existing parent")
			}
		}
	}
	r.CreatedDirs = sortedKeys(created)
	if err := saveReceipt(dir, r); err != nil {
		return dir, err
	}
	for i, c := range p.Changes {
		if beforeWrite != nil {
			err = beforeWrite(i)
		}
		if err == nil {
			err = guard()
		}
		path, _ := targetPath(p.Locations, c.Area, c.Name, p.Payload)
		if err == nil {
			actual, exists, e := readFile(path)
			if e != nil || exists != c.BeforeExists || !bytes.Equal(actual, c.Before) {
				err = fmt.Errorf("concurrent edit: %s", path)
			}
		}
		if err == nil {
			err = atomicWrite(path, c.After)
		}
		if err != nil {
			return filepath.Join(dir, "receipt.json"), fmt.Errorf("installation incomplete: %w; use Roll back with %s", err, filepath.Join(dir, "receipt.json"))
		}
	}
	for _, c := range p.Changes {
		path, _ := targetPath(p.Locations, c.Area, c.Name, p.Payload)
		hash, e := fileDigest(path)
		if e != nil || hash != digest(c.After) {
			return filepath.Join(dir, "receipt.json"), fmt.Errorf("installed file verification failed: %s", path)
		}
	}
	r.Status = "INSTALLED"
	err = saveReceipt(dir, r)
	return filepath.Join(dir, "receipt.json"), err
}

func rollback(receiptPath string, payload map[string][]byte, guard func() error) error {
	if err := guard(); err != nil {
		return err
	}
	data, exists, err := readFile(receiptPath)
	if err != nil || !exists {
		return fmt.Errorf("cannot read receipt")
	}
	var r receipt
	if err = json.Unmarshal(data, &r); err != nil {
		return err
	}
	if r.Schema != 1 || r.Runtime != runtimeVersion || r.Installer != installerVersion {
		return fmt.Errorf("unsupported receipt")
	}
	if err = validateLocations(r.Locations); err != nil {
		return err
	}
	release, err := acquireInstallationLock(r.Locations.Game)
	if err != nil {
		return err
	}
	defer release()
	dir := filepath.Dir(receiptPath)
	if filepath.Base(receiptPath) != "receipt.json" || !samePath(filepath.Dir(dir), filepath.Join(r.Locations.Game, ".zbc-installations")) {
		return fmt.Errorf("receipt is outside its recorded installation")
	}
	if r.Status == "ROLLED_BACK" {
		return nil
	}
	if r.Status != "INSTALLED" && r.Status != "PREPARED" && r.Status != "ROLLING_BACK" {
		return fmt.Errorf("unsupported transaction status")
	}
	var changes []change
	targets := map[string]bool{}
	for i, f := range r.Files {
		path, err := targetPath(r.Locations, f.Area, f.Name, payload)
		if err != nil {
			return err
		}
		if targets[path] {
			return fmt.Errorf("duplicate receipt target")
		}
		targets[path] = true
		current, currentExists, err := readFile(path)
		if err != nil {
			return err
		}
		currentHash := ""
		if currentExists {
			currentHash = digest(current)
		}
		if f.Area == "steam" {
			patch, ok, e := readFile(filepath.Join(dir, fmt.Sprintf("steam-%d.json", i)))
			if e != nil || !ok || digest(patch) != f.SteamPatchHash {
				return fmt.Errorf("Steam patch checksum mismatch")
			}
			if currentHash == f.Before {
				continue
			}
			restore, e := restoreSteamPatch(current, patch, f.Before, f.After)
			if e != nil {
				return e
			}
			if !bytes.Equal(current, restore) {
				changes = append(changes, change{f.Area, f.Name, current, restore, true})
			}
			continue
		}
		before, beforeExists, err := readFile(filepath.Join(dir, fmt.Sprintf("before-%d", i)))
		if err != nil {
			return err
		}
		if (f.Before != "") != beforeExists || beforeExists && digest(before) != f.Before {
			return fmt.Errorf("backup checksum mismatch: %s", path)
		}
		after, afterExists, err := readFile(filepath.Join(dir, fmt.Sprintf("after-%d", i)))
		if err != nil || !afterExists || digest(after) != f.After {
			return fmt.Errorf("staged checksum mismatch: %s", path)
		}
		if currentHash == f.Before {
			continue
		}
		restore := before
		if currentHash != f.After {
			return fmt.Errorf("refusing to overwrite a changed file: %s", path)
		}
		changes = append(changes, change{f.Area, f.Name, current, restore, beforeExists})
	}
	for _, path := range r.CreatedDirs {
		if !within(filepath.Join(r.Locations.Profile, "mods"), path) && !samePath(path, r.Locations.Profile) {
			return fmt.Errorf("unexpected created directory")
		}
		if err := safePath(path); err != nil {
			return err
		}
	}
	r.Status = "ROLLING_BACK"
	if err = saveReceipt(dir, r); err != nil {
		return err
	}
	for i := len(changes) - 1; i >= 0; i-- {
		c := changes[i]
		path, _ := targetPath(r.Locations, c.Area, c.Name, payload)
		if err = guard(); err != nil {
			return err
		}
		actual, _, e := readFile(path)
		if e != nil || !bytes.Equal(actual, c.Before) {
			return fmt.Errorf("concurrent edit during rollback: %s", path)
		}
		if c.BeforeExists {
			err = atomicWrite(path, c.After)
		} else {
			err = os.Remove(path)
		}
		if err != nil {
			return err
		}
		restored, exists, verifyErr := readFile(path)
		if verifyErr != nil || exists != c.BeforeExists || !bytes.Equal(restored, c.After) {
			return fmt.Errorf("restored file verification failed: %s", path)
		}
	}
	for i := len(r.CreatedDirs) - 1; i >= 0; i-- {
		_ = os.Remove(r.CreatedDirs[i])
	} // Empty directories only; retain newly added user content.
	r.Status = "ROLLED_BACK"
	return saveReceipt(dir, r)
}
