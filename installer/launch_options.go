package main

import (
	"fmt"
	"strings"
)

// Keep the user's argument text verbatim. Ambiguous wrappers need manual review.
func mergeSteamLaunchOptions(current string) (string, error) {
	if strings.ContainsAny(current, "\r\n") || strings.Contains(current, "%command%") {
		return "", fmt.Errorf("custom Steam launch wrapper: preserve it and configure the agent manually")
	}
	if hasZombieBuddyLaunchOptions(current) {
		return current, nil
	}
	tokens, err := launchTokens(current)
	if err != nil {
		return "", err
	}
	hasSeparator := false
	hasJVMOption := false
	for _, token := range tokens {
		if token == "--" {
			hasSeparator = true
		}
		if strings.HasPrefix(token, "-javaagent:") || strings.HasPrefix(token, "-agentlib:") ||
			strings.HasPrefix(token, "-agentpath:") || strings.HasPrefix(token, "-X") ||
			strings.HasPrefix(token, "-D") {
			hasJVMOption = true
		}
	}
	if hasSeparator {
		return ZB_LAUNCH_ARG + " " + current, nil
	}
	if hasJVMOption {
		return "", fmt.Errorf("JVM options without --: review the launch separator before adding ZombieBuddy")
	}
	if strings.TrimSpace(current) == "" {
		return ZB_LAUNCH_OPTIONS, nil
	}
	return ZB_LAUNCH_OPTIONS + " " + current, nil
}

func launchTokens(options string) ([]string, error) {
	var result []string
	start, slashes := -1, 0
	quoted := false
	for i := 0; i < len(options); i++ {
		ch := options[i]
		if ch == '"' && slashes%2 == 0 {
			quoted = !quoted
		}
		if (ch == ' ' || ch == '\t') && !quoted {
			if start >= 0 {
				result = append(result, options[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
		if ch == '\\' {
			slashes++
		} else {
			slashes = 0
		}
	}
	if quoted {
		return nil, fmt.Errorf("unbalanced quotes in Steam launch options; no changes applied")
	}
	if start >= 0 {
		result = append(result, options[start:])
	}
	return result, nil
}
