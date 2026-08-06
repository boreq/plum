package core

import "strings"

var automatedUserAgentNames = map[string]struct{}{
	"curl":       {},
	"prometheus": {},
}

func ClassifyUserAgent(userAgent string) Category {
	name, _, _ := strings.Cut(userAgent, "/")
	name = strings.ToLower(strings.TrimSpace(name))

	if _, ok := automatedUserAgentNames[name]; ok {
		return CategoryAutomated
	}

	return CategoryUnclassified
}
