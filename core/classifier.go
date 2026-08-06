package core

import (
	"strings"
	"unicode"
)

var automatedUserAgentNames = map[string]struct{}{
	"curl":       {},
	"prometheus": {},
}

func UserAgentName(userAgent string) string {
	name, _, _ := strings.Cut(userAgent, "/")
	name, _, _ = strings.Cut(name, " (")
	name = strings.TrimSpace(name)

	if !isSimpleUserAgentName(name) {
		return userAgent
	}

	return name
}

func isSimpleUserAgentName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' ' && r != '-' {
			return false
		}
	}

	return true
}

func ClassifyUserAgent(userAgent string) Category {
	name := strings.ToLower(UserAgentName(userAgent))

	if _, ok := automatedUserAgentNames[name]; ok {
		return CategoryAutomated
	}

	return CategoryUnclassified
}
