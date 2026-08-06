package core

import "testing"

func TestNewUserAgent(t *testing.T) {
	testCases := []struct {
		UserAgent string
		Name      string
		Browser   *Browser
	}{
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
			Browser:   BrowserChrome,
			Name:      "Chrome",
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0",
			Browser:   BrowserEdge,
			Name:      "Edge",
		},
		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:153.0) Gecko/20100101 Firefox/153.0",
			Browser:   BrowserFirefox,
			Name:      "Firefox",
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			Browser:   BrowserSafari,
			Name:      "Safari",
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Name:      "Googlebot",
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.7871.186 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Name:      "Googlebot",
		},
		{
			UserAgent: "curl/7.64.0",
			Name:      "curl",
		},
		{
			UserAgent: "moooodotfarm (+https://moooo.farm)",
			Name:      "moooodotfarm",
		},
		{
			UserAgent: "http.rb/5.1.1 (Mastodon/4.2.17; +https://mastodon.social/)",
			Name:      "http.rb",
		},
		{
			UserAgent: "vuln_scanner/3.1.0 (CVE-2026-4020)",
			Name:      "vuln_scanner",
		},
		{
			UserAgent: "",
			Name:      "",
		},
		{
			UserAgent: "-",
			Name:      "-",
		},
		{
			UserAgent: "Ditto <https://ditto.pub>",
			Name:      "Ditto <https://ditto.pub>",
		},
		{
			UserAgent: "Feedbin feed-id:1602368 - 4 subscribers",
			Name:      "Feedbin feed-id:1602368 - 4 subscribers",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			userAgent := NewUserAgent(testCase.UserAgent)

			if userAgent.Name() != testCase.Name {
				t.Errorf("name: got %q, want %q", userAgent.Name(), testCase.Name)
			}

			if userAgent.Browser() != testCase.Browser {
				t.Errorf("browser: got %v, want %v", userAgent.Browser(), testCase.Browser)
			}
		})
	}
}
