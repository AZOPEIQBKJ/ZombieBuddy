package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type launchToken struct{ raw, value string }

func launchTokens(s string) ([]launchToken, error) {
	if strings.ContainsAny(s, "\r\n") || strings.Contains(strings.ToLower(s), "%command%") {
		return nil, fmt.Errorf("custom Steam command wrapper requires separate review")
	}
	var out []launchToken
	start, slashes := -1, 0
	quoted := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' && slashes%2 == 0 {
			quoted = !quoted
		}
		if (c == ' ' || c == '\t') && !quoted {
			if start >= 0 {
				raw := s[start:i]
				out = append(out, launchToken{raw, strings.ReplaceAll(raw, `"`, "")})
				start = -1
			}
		} else if start < 0 {
			start = i
		}
		if c == '\\' {
			slashes++
		} else {
			slashes = 0
		}
	}
	if quoted {
		return nil, fmt.Errorf("unbalanced quotes in Steam launch options")
	}
	if start >= 0 {
		raw := s[start:]
		out = append(out, launchToken{raw, strings.ReplaceAll(raw, `"`, "")})
	}
	return out, nil
}

func zbArgument(arg string) (string, bool) {
	head, tail, _ := strings.Cut(arg, "=")
	lower := strings.ToLower(head)
	if lower == "-agentlib:zbnative" {
		return tail, true
	}
	if strings.HasPrefix(lower, "-javaagent:") && strings.EqualFold(filepath.Base(strings.ReplaceAll(head[len("-javaagent:"):], "\\", "/")), "ZombieBuddy.jar") {
		return tail, true
	}
	if strings.HasPrefix(lower, "-agentpath:") {
		name := filepath.Base(strings.ReplaceAll(head[len("-agentpath:"):], "\\", "/"))
		if strings.EqualFold(name, "zbNative.dll") {
			return tail, true
		}
	}
	return "", false
}

func mergeAgentOptions(tails []string) (string, error) {
	values := map[string]string{}
	var order []string
	for _, tail := range tails {
		for _, part := range strings.Split(tail, ",") {
			if part == "" {
				continue
			}
			key, value, ok := strings.Cut(part, "=")
			if !ok || key == "" {
				return "", fmt.Errorf("unrecognized ZombieBuddy option: %s", key)
			}
			if key == "patches_jar" {
				var kept []string
				for _, entry := range strings.Split(value, ";") {
					colon := strings.LastIndex(entry, ":")
					if colon > 0 && entry[colon+1:] == "aftermathsystems.lhcompat" && strings.EqualFold(filepath.Base(strings.ReplaceAll(entry[:colon], "\\", "/")), "AftermathLHCompat4221.jar") {
						continue
					}
					kept = append(kept, entry)
				}
				value = strings.Join(kept, ";")
				if value == "" {
					continue
				}
			}
			if previous, ok := values[key]; ok {
				if previous != value {
					return "", fmt.Errorf("conflicting ZombieBuddy options for %s; installation unchanged", key)
				}
			} else {
				order = append(order, key)
			}
			values[key] = value
		}
	}
	var result []string
	for _, key := range order {
		result = append(result, key+"="+values[key])
	}
	return strings.Join(result, ","), nil
}

func vmArguments(raw json.RawMessage) ([]string, error) {
	var args []string
	err := json.Unmarshal(raw, &args)
	return args, err
}

func cleanVMArgs(args []string) ([]string, []string, error) {
	kept := []string{}
	var tails []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-Duser.home=") {
			return nil, nil, fmt.Errorf("custom JVM user.home requires separate profile review")
		}
		if tail, ok := zbArgument(arg); ok {
			tails = append(tails, tail)
		} else {
			kept = append(kept, arg)
		}
	}
	return kept, tails, nil
}

func configureLauncher(data []byte, steamTails []string) ([]byte, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(bytes.TrimPrefix(data, []byte{239, 187, 191}), &root); err != nil {
		return nil, err
	}
	args, err := vmArguments(root["vmArgs"])
	if err != nil {
		return nil, err
	}
	kept, tails, err := cleanVMArgs(args)
	if err != nil {
		return nil, err
	}
	if raw, ok := root["windows"]; ok {
		var windows map[string]map[string]json.RawMessage
		if err := json.Unmarshal(raw, &windows); err != nil {
			return nil, err
		}
		for _, key := range sortedKeys(windows) {
			block := windows[key]
			if value, ok := block["vmArgs"]; ok {
				args, err := vmArguments(value)
				if err != nil {
					return nil, err
				}
				cleaned, extra, err := cleanVMArgs(args)
				if err != nil {
					return nil, err
				}
				tails = append(tails, extra...)
				block["vmArgs"], _ = json.Marshal(cleaned)
			}
		}
		root["windows"], _ = json.Marshal(windows)
	}
	tail, err := mergeAgentOptions(append(tails, steamTails...))
	if err != nil {
		return nil, err
	}
	agent := "-javaagent:ZombieBuddy.jar"
	if tail != "" {
		agent += "=" + tail
	}
	root["vmArgs"], _ = json.Marshal(append([]string{agent}, kept...))
	output, err := json.MarshalIndent(root, "", "\t")
	if err != nil {
		return nil, err
	}
	// Avoid needless rewrites when the parsed configuration already has this meaning.
	var a, b any
	json.Unmarshal(bytes.TrimPrefix(data, []byte{239, 187, 191}), &a)
	json.Unmarshal(output, &b)
	aa, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	if bytes.Equal(aa, bb) {
		return data, nil
	}
	return append(output, '\n'), nil
}

func configureOptions(current string, profile string) (string, []string, error) {
	tokens, err := launchTokens(current)
	if err != nil {
		return "", nil, err
	}
	separator := -1
	for i, t := range tokens {
		if t.raw == "--" {
			if separator >= 0 {
				return "", nil, fmt.Errorf("multiple launch separators")
			}
			separator = i
		}
	}
	var before, after, tails []string
	order := "workshop,steam,mods"
	seenOrder := false
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t.raw == "--" {
			continue
		}
		if tail, ok := zbArgument(t.value); ok {
			tails = append(tails, tail)
			continue
		}
		if t.value == "-modfolders" {
			if seenOrder || i+1 >= len(tokens) || separator >= 0 && i < separator {
				return "", nil, fmt.Errorf("ambiguous mod folder arguments")
			}
			seenOrder = true
			i++
			order = tokens[i].value
			continue
		}
		if strings.HasPrefix(t.value, "-cachedir=") {
			candidate := strings.TrimPrefix(t.value, "-cachedir=")
			if !filepath.IsAbs(candidate) || !samePath(candidate, profile) {
				return "", nil, fmt.Errorf("Steam -cachedir differs from selected profile; select that exact profile")
			}
		}
		if separator >= 0 && i < separator {
			before = append(before, t.raw)
		} else {
			if separator < 0 && (strings.HasPrefix(t.value, "-D") || strings.HasPrefix(t.value, "-X") || strings.HasPrefix(t.value, "-agent") || strings.HasPrefix(t.value, "-javaagent")) {
				return "", nil, fmt.Errorf("JVM options without -- require separate review")
			}
			after = append(after, t.raw)
		}
	}
	reordered := []string{"mods"}
	seen := map[string]bool{}
	for _, part := range strings.Split(order, ",") {
		if seen[part] || (part != "mods" && part != "workshop" && part != "steam") {
			return "", nil, fmt.Errorf("unknown or repeated mod folder source")
		}
		seen[part] = true
		if part != "mods" {
			reordered = append(reordered, part)
		}
	}
	after = append(after, "-modfolders", strings.Join(reordered, ","))
	result := append(before, "--")
	result = append(result, after...)
	return strings.Join(result, " "), tails, nil
}
