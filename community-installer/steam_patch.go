package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Back up only the PZ edit, never the rest of a Steam account's configuration.
type steamPatch struct {
	Offset    int    `json:"offset"`
	Removed   []byte `json:"removed"`
	Inserted  []byte `json:"inserted"`
	OldValue  string `json:"oldValue"`
	OldExists bool   `json:"oldExists"`
	NewValue  string `json:"newValue"`
}

func makeSteamPatch(before, after []byte) ([]byte, error) {
	old, exists, err := lookupVDF(before, optionsPath)
	if err != nil {
		return nil, err
	}
	value, present, err := lookupVDF(after, optionsPath)
	if err != nil || !present {
		return nil, fmt.Errorf("missing new Steam options")
	}
	prefix := 0
	for prefix < len(before) && prefix < len(after) && before[prefix] == after[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(before)-prefix && suffix < len(after)-prefix && before[len(before)-1-suffix] == after[len(after)-1-suffix] {
		suffix++
	}
	return json.Marshal(steamPatch{prefix, before[prefix : len(before)-suffix], after[prefix : len(after)-suffix], old, exists, value})
}

func restoreSteamPatch(current, encoded []byte, beforeHash, afterHash string) ([]byte, error) {
	var patch steamPatch
	if err := json.Unmarshal(encoded, &patch); err != nil {
		return nil, err
	}
	if digest(current) == afterHash {
		end := patch.Offset + len(patch.Inserted)
		if patch.Offset < 0 || end < patch.Offset || end > len(current) || !bytes.Equal(current[patch.Offset:end], patch.Inserted) {
			return nil, fmt.Errorf("invalid Steam edit backup")
		}
		restored := append([]byte{}, current[:patch.Offset]...)
		restored = append(restored, patch.Removed...)
		restored = append(restored, current[end:]...)
		if digest(restored) != beforeHash {
			return nil, fmt.Errorf("Steam edit backup does not restore the recorded hash")
		}
		return restored, nil
	}
	value, exists, err := lookupVDF(current, optionsPath)
	if err != nil {
		return nil, err
	}
	if value == patch.OldValue && exists == patch.OldExists {
		return current, nil
	}
	if !exists || value != patch.NewValue {
		return nil, fmt.Errorf("Steam launch options changed since installation; preserve them and review rollback")
	}
	var old *string
	if patch.OldExists {
		old = &patch.OldValue
	}
	return changeVDF(current, optionsPath, old)
}
