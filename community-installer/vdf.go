package main

import (
	"fmt"
	"strings"
)

// Byte spans let us edit one Steam value without reserializing the account file.
type vdfNode struct {
	key, value                              string
	scalar                                  bool
	children                                []*vdfNode
	start, valueStart, valueEnd, end, close int
}
type vdfToken struct {
	text       string
	start, end int
	structural bool
}

func vdfTokens(s string) ([]vdfToken, error) {
	var out []vdfToken
	for i := 0; i < len(s); {
		if i == 0 && strings.HasPrefix(s, "\xef\xbb\xbf") {
			i = 3
			continue
		}
		if strings.ContainsRune(" \t\r\n", rune(s[i])) {
			i++
			continue
		}
		if strings.HasPrefix(s[i:], "//") {
			for i < len(s) && s[i] != '\n' {
				i++
			}
			continue
		}
		start := i
		if s[i] == '{' || s[i] == '}' {
			out = append(out, vdfToken{string(s[i]), i, i + 1, true})
			i++
			continue
		}
		var b strings.Builder
		if s[i] == '"' {
			i++
			closed := false
			for i < len(s) {
				c := s[i]
				i++
				if c == '"' {
					closed = true
					break
				}
				if c == '\\' && i < len(s) && (s[i] == '\\' || s[i] == '"') {
					c = s[i]
					i++
				}
				b.WriteByte(c)
			}
			if !closed {
				return nil, fmt.Errorf("unterminated Steam configuration string")
			}
		} else {
			for i < len(s) && !strings.ContainsRune(" \t\r\n{}", rune(s[i])) {
				b.WriteByte(s[i])
				i++
			}
		}
		out = append(out, vdfToken{b.String(), start, i, false})
		if len(out) > 500000 {
			return nil, fmt.Errorf("Steam configuration too large")
		}
	}
	return out, nil
}

func parseVDF(s string) ([]*vdfNode, error) {
	tokens, err := vdfTokens(s)
	if err != nil {
		return nil, err
	}
	i := 0
	var block func(int) ([]*vdfNode, error)
	block = func(depth int) ([]*vdfNode, error) {
		if depth > 32 {
			return nil, fmt.Errorf("Steam configuration nesting too deep")
		}
		var nodes []*vdfNode
		for i < len(tokens) && !(tokens[i].structural && tokens[i].text == "}") {
			key := tokens[i]
			i++
			if key.structural || i >= len(tokens) {
				return nil, fmt.Errorf("invalid Steam key/value")
			}
			value := tokens[i]
			i++
			n := &vdfNode{key: key.text, start: key.start, valueStart: value.start, valueEnd: value.end, end: value.end}
			if value.structural {
				if value.text != "{" {
					return nil, fmt.Errorf("invalid Steam block")
				}
				n.children, err = block(depth + 1)
				if err != nil {
					return nil, err
				}
				if i >= len(tokens) || tokens[i].text != "}" {
					return nil, fmt.Errorf("unclosed Steam block")
				}
				n.close = tokens[i].start
				n.end = tokens[i].end
				i++
			} else {
				n.scalar = true
				n.value = value.text
			}
			nodes = append(nodes, n)
		}
		return nodes, nil
	}
	nodes, err := block(0)
	if err == nil && i != len(tokens) {
		err = fmt.Errorf("unexpected Steam closing brace")
	}
	return nodes, err
}

func vdfChild(nodes []*vdfNode, key string) (*vdfNode, error) {
	var found *vdfNode
	for _, n := range nodes {
		if strings.EqualFold(n.key, key) {
			if found != nil {
				return nil, fmt.Errorf("ambiguous duplicate Steam key: %s", key)
			}
			found = n
		}
	}
	return found, nil
}

var optionsPath = []string{"UserLocalConfigStore", "Software", "Valve", "Steam", "Apps", "108600", "LaunchOptions"}

func lookupVDF(data []byte, path []string) (string, bool, error) {
	nodes, err := parseVDF(string(data))
	if err != nil {
		return "", false, err
	}
	for i, key := range path {
		n, err := vdfChild(nodes, key)
		if err != nil {
			return "", false, err
		}
		if n == nil {
			return "", false, nil
		}
		if i == len(path)-1 {
			if !n.scalar {
				return "", false, fmt.Errorf("expected Steam string")
			}
			return n.value, true, nil
		}
		if n.scalar {
			return "", false, fmt.Errorf("expected Steam container")
		}
		nodes = n.children
	}
	return "", false, fmt.Errorf("empty Steam path")
}

func quoteVDF(s string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s) + `"`
}

func changeVDF(data []byte, path []string, value *string) ([]byte, error) {
	s := string(data)
	nodes, err := parseVDF(s)
	if err != nil {
		return nil, err
	}
	closeAt := len(s)
	for i, key := range path {
		n, err := vdfChild(nodes, key)
		if err != nil {
			return nil, err
		}
		if n == nil {
			if value == nil {
				return data, nil
			}
			tail := quoteVDF(path[len(path)-1]) + "\t" + quoteVDF(*value) + "\n"
			for j := len(path) - 2; j >= i; j-- {
				tail = quoteVDF(path[j]) + "\n{\n" + tail + "}\n"
			}
			return []byte(s[:closeAt] + "\n" + tail + s[closeAt:]), nil
		}
		if i == len(path)-1 {
			if !n.scalar {
				return nil, fmt.Errorf("expected Steam string")
			}
			if value == nil {
				return []byte(s[:n.start] + s[n.end:]), nil
			}
			return []byte(s[:n.valueStart] + quoteVDF(*value) + s[n.valueEnd:]), nil
		}
		if n.scalar {
			return nil, fmt.Errorf("expected Steam container")
		}
		nodes = n.children
		closeAt = n.close
	}
	return nil, fmt.Errorf("empty Steam path")
}
