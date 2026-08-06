package core

import (
	"strings"
	"unicode"
)

// browserMarkers maps a token found in a user agent to a browser. Every browser
// claims to be Mozilla and most claim to be Safari or Chrome as well so the
// order matters, the most specific token wins.
var browserMarkers = []struct {
	Marker  string
	Browser *Browser
}{
	{Marker: "edg/", Browser: BrowserEdge},
	{Marker: "edga/", Browser: BrowserEdge},
	{Marker: "edgios/", Browser: BrowserEdge},
	{Marker: "opr/", Browser: BrowserOpera},
	{Marker: "opera/", Browser: BrowserOpera},
	{Marker: "vivaldi/", Browser: BrowserVivaldi},
	{Marker: "brave/", Browser: BrowserBrave},
	{Marker: "duckduckgo/", Browser: BrowserDuckDuckGo},
	{Marker: "samsungbrowser/", Browser: BrowserSamsungInternet},
	{Marker: "yabrowser/", Browser: BrowserYandex},
	{Marker: "chromium/", Browser: BrowserChromium},
	{Marker: "crios/", Browser: BrowserChrome},
	{Marker: "chrome/", Browser: BrowserChrome},
	{Marker: "fxios/", Browser: BrowserFirefox},
	{Marker: "firefox/", Browser: BrowserFirefox},
	{Marker: "safari/", Browser: BrowserSafari},
	{Marker: "msie ", Browser: BrowserInternetExplorer},
	{Marker: "trident/", Browser: BrowserInternetExplorer},
}

// UserAgent has either a name or a browser set, never both.
type UserAgent struct {
	name    string
	browser *Browser
}

func NewUserAgent(raw string) UserAgent {
	if name, ok := compatibleUserAgentName(raw); ok {
		if browser := recognizeBrowser(name); browser != nil {
			return UserAgent{browser: browser}
		}
		return UserAgent{name: name}
	}

	if browser := recognizeBrowser(raw); browser != nil {
		return UserAgent{browser: browser}
	}

	if name, ok := productUserAgentName(raw); ok {
		return UserAgent{name: name}
	}

	return UserAgent{name: raw}
}

func (u UserAgent) Browser() *Browser {
	return u.browser
}

func (u UserAgent) Name() string {
	if u.browser != nil {
		return u.browser.Name()
	}

	return u.name
}

func (u UserAgent) String() string {
	return u.Name()
}

// compatibleUserAgentName extracts the name of a bot which hides behind a
// browser user agent, for example "Mozilla/5.0 (compatible; Googlebot/2.1)".
func compatibleUserAgentName(raw string) (string, bool) {
	_, rest, found := cutFold(raw, "(compatible; ")
	if !found {
		return "", false
	}

	name, _, _ := strings.Cut(rest, "/")
	name, _, _ = strings.Cut(name, ";")
	name, _, _ = strings.Cut(name, ")")
	name = strings.TrimSpace(name)

	if !isUsableUserAgentName(name) {
		return "", false
	}

	return name, true
}

func recognizeBrowser(raw string) *Browser {
	lowered := strings.ToLower(raw)

	for _, marker := range browserMarkers {
		if strings.Contains(lowered, marker.Marker) {
			return marker.Browser
		}
	}

	return nil
}

func productUserAgentName(raw string) (string, bool) {
	name, _, _ := strings.Cut(raw, "/")
	name, _, _ = strings.Cut(name, " (")
	name = strings.TrimSpace(name)

	if !isUsableUserAgentName(name) {
		return "", false
	}

	return name, true
}

func isUsableUserAgentName(name string) bool {
	if name == "" || name == "-" {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '-' && r != '_' && r != '.' {
			return false
		}
	}

	return true
}

func cutFold(s, sep string) (before, after string, found bool) {
	index := strings.Index(strings.ToLower(s), sep)
	if index < 0 {
		return s, "", false
	}
	return s[:index], s[index+len(sep):], true
}
