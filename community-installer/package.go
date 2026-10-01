package main

import (
	"archive/zip"
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
)

const installerVersion = "0.1.0-preview.3"
const runtimeVersion = "2.3.3-community.4"
const archiveHash = "f1b66cf9e11aac8663de92de382944a979cb4f09e651d320d8d3938b471a3c99"
const gameHash = "e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33"
const signingKey = "989ac279f40f1a35fa0616645e319fc44cde9a15a842b7c365870147bfff3ce0"
const signingAuthor = "76561198061324182"
const payloadPrefix = "Contents/mods/ZombieBuddy/"
const nativeLoaderHash = "c2ae9335e717ee24b2f4a40d1a3bf77f1519762a72a0459e766a2bbafc077f6c"

// Filled by tools/build_installer.py from the immutable, independently verified archive.
//
//go:embed payload.zip
var embeddedPackage []byte

// Unmodified MIT-licensed upstream Windows bootstrap, independently pinned.
// It establishes the bundled JRE DLL search path before loading instrumentation.
//
//go:embed native-loader.dll
var embeddedNativeLoader []byte

func verifyNativeLoader(data []byte) error {
	if digest(data) != nativeLoaderHash {
		return fmt.Errorf("embedded native loader checksum mismatch")
	}
	return nil
}

var sourceCommit = "development"

type releaseManifest struct {
	Distribution string            `json:"distribution"`
	Version      string            `json:"version"`
	TargetGame   string            `json:"targetGameSha256"`
	SigningKey   string            `json:"signingPublicKey"`
	Files        map[string]string `json:"files"`
}

func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func payloadFiles(data []byte) (map[string][]byte, error) {
	if digest(data) != archiveHash {
		return nil, fmt.Errorf("embedded runtime archive checksum mismatch")
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, f := range zr.File {
		if path.Clean(f.Name) != f.Name || strings.HasPrefix(f.Name, "/") || strings.HasPrefix(f.Name, "../") || f.Name == ".." || strings.ContainsAny(f.Name, "\\:") || !f.FileInfo().Mode().IsRegular() || f.UncompressedSize64 > 128<<20 {
			return nil, fmt.Errorf("invalid archive entry: %s", f.Name)
		}
		if _, exists := files[f.Name]; exists {
			return nil, fmt.Errorf("duplicate archive entry")
		}
		r, err := f.Open()
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(io.LimitReader(r, 128<<20+1))
		r.Close()
		if err != nil || len(content) > 128<<20 {
			return nil, fmt.Errorf("cannot read archive entry %s", f.Name)
		}
		files[f.Name] = content
	}
	var manifest releaseManifest
	if err := json.Unmarshal(files["manifest.json"], &manifest); err != nil {
		return nil, err
	}
	if manifest.Distribution != "ZombieBuddy Community" || manifest.Version != runtimeVersion || manifest.TargetGame != gameHash || manifest.SigningKey != signingKey || len(files) != len(manifest.Files)+1 {
		return nil, fmt.Errorf("unexpected release manifest")
	}
	for name, hash := range manifest.Files {
		value, ok := files[name]
		if !ok || digest(value) != hash {
			return nil, fmt.Errorf("release checksum mismatch: %s", name)
		}
	}
	payload := map[string][]byte{}
	for name, content := range files {
		if strings.HasPrefix(name, payloadPrefix) {
			payload[strings.TrimPrefix(name, payloadPrefix)] = content
		}
	}
	if len(payload) != 9 {
		return nil, fmt.Errorf("unexpected framework file set")
	}
	jar, ok := payload["libs/ZombieBuddy.jar"]
	if !ok {
		return nil, fmt.Errorf("missing framework JAR")
	}
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(payload["libs/ZombieBuddy.jar.zbs"]), "\r\n", "\n")), "\n")
	if len(lines) != 3 || lines[0] != "ZBS" || lines[1] != "SteamID64:"+signingAuthor || !strings.HasPrefix(lines[2], "Signature:") {
		return nil, fmt.Errorf("invalid framework signature metadata")
	}
	key, _ := hex.DecodeString(signingKey)
	signature, err := hex.DecodeString(strings.TrimPrefix(lines[2], "Signature:"))
	if err != nil || !ed25519.Verify(ed25519.PublicKey(key), []byte("ZBS:"+signingAuthor+":"+digest(jar)), signature) {
		return nil, fmt.Errorf("framework signature verification failed")
	}
	return payload, nil
}

func sortedKeys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
