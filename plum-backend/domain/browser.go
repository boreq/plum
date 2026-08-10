package domain

import (
	"strings"

	"github.com/boreq/plum/plum-backend/domain/request"
)

type Browser struct {
	id   string
	name string
}

var (
	BrowserChrome   = &Browser{"chrome", "Chrome"}
	BrowserChromium = &Browser{"chromium", "Chromium"}
	BrowserFirefox  = &Browser{"firefox", "Firefox"}
	BrowserSafari   = &Browser{"safari", "Safari"}
)

var browserMarkers = []struct {
	Marker  string
	Browser *Browser
}{
	{Marker: "chromium/", Browser: BrowserChromium},
	{Marker: "crios/", Browser: BrowserChrome},
	{Marker: "chrome/", Browser: BrowserChrome},
	{Marker: "fxios/", Browser: BrowserFirefox},
	{Marker: "firefox/", Browser: BrowserFirefox},
	{Marker: "safari", Browser: BrowserSafari},
}

func RecognizeBrowser(userAgent request.UserAgent) *Browser {
	raw := strings.ToLower(userAgent.String())

	for _, marker := range browserMarkers {
		if strings.Contains(raw, marker.Marker) {
			return marker.Browser
		}
	}

	return nil
}

func (b *Browser) Name() string {
	return b.name
}

func (b *Browser) String() string {
	return b.id
}
