package core

import "testing"

func TestRecognizeBrowser(t *testing.T) {
	testCases := []struct {
		Name      string
		UserAgent string
		Browser   *Browser
	}{
		{
			Name:      "chrome",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
			Browser:   BrowserChrome,
		},
		{
			Name:      "chrome on ios",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/99.0.4844.47 Mobile/15E148 Safari/604.1",
			Browser:   BrowserChrome,
		},
		{
			Name:      "chromium",
			UserAgent: "Mozilla/5.0 (X11; Linux armv7l) AppleWebKit/537.36 (KHTML, like Gecko) Raspbian Chromium/72.0.3626.121 Chrome/72.0.3626.121 Safari/537.36",
			Browser:   BrowserChromium,
		},
		{
			Name:      "firefox",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
			Browser:   BrowserFirefox,
		},
		{
			Name:      "firefox on ios",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/112.2  Mobile/15E148 Safari/605.1.15",
			Browser:   BrowserFirefox,
		},
		{
			Name:      "firefox with a comment at the end",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; WOW64; rv:41.0) Gecko/20100101 Firefox/140.0.2 (x64 de)",
			Browser:   BrowserFirefox,
		},
		{
			Name:      "safari",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			Browser:   BrowserSafari,
		},
		{
			Name:      "safari without a version",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Safari",
			Browser:   BrowserSafari,
		},
		{
			Name:      "browser based on chromium",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36 Edg/151.0.0.0",
			Browser:   BrowserChrome,
		},
		{
			Name:      "browser based on firefox",
			UserAgent: "Mozilla/5.0 (Windows NT 5.1; rv:38.0) Gecko/20100101 Firefox/38.0 SeaMonkey/2.35",
			Browser:   BrowserFirefox,
		},
		{
			Name:      "crawler using a browser user agent",
			UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/136.0.0.0 Safari/537.36",
			Browser:   BrowserChrome,
		},
		{
			Name:      "crawler pretending to be chrome on android",
			UserAgent: "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.7871.186 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Browser:   BrowserChrome,
		},
		{
			Name:      "headless chrome",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/145.0.7632.6 Safari/537.36",
			Browser:   BrowserChrome,
		},
		{
			Name:      "mail client based on firefox",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Thunderbird/102.12.0",
			Browser:   nil,
		},
		{
			Name:      "sync client",
			UserAgent: "Mozilla/5.0 (Android) Nextcloud-android/3.24.2",
			Browser:   nil,
		},
		{
			Name:      "desktop sync client",
			UserAgent: "Mozilla/5.0 (Linux) mirall/3.10.0git (Nextcloud, arch-6.3.8-arch1-1 ClientArchitecture: x86_64 OsArchitecture: x86_64)",
			Browser:   nil,
		},
		{
			Name:      "command line program",
			UserAgent: "curl/8.7.1",
			Browser:   nil,
		},
		{
			Name:      "scanner",
			UserAgent: "vuln_scanner/3.1.0 (CVE-2026-4020)",
			Browser:   nil,
		},
		{
			Name:      "missing user agent",
			UserAgent: "",
			Browser:   nil,
		},
		{
			Name:      "empty user agent",
			UserAgent: "-",
			Browser:   nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if browser := RecognizeBrowser(testCase.UserAgent); browser != testCase.Browser {
				t.Errorf("got %v, want %v", browser, testCase.Browser)
			}
		})
	}
}
