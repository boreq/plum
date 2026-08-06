package core

import "testing"

func TestUserAgentName(t *testing.T) {
	testCases := []struct {
		UserAgent string
		Name      string
	}{
		{
			UserAgent: "curl/7.64.0",
			Name:      "curl",
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Name:      "Mozilla",
		},
		{
			UserAgent: "moooodotfarm (+https://moooo.farm)",
			Name:      "moooodotfarm",
		},
		{
			UserAgent: "Turnitin (https://turnitin.com/robot/crawlerinfo.html)",
			Name:      "Turnitin",
		},
		{
			UserAgent: "OnlyRead",
			Name:      "OnlyRead",
		},
		{
			UserAgent: "python-requests/2.32.3",
			Name:      "python-requests",
		},
		{
			UserAgent: "Feedbin feed-id:1602368 - 4 subscribers",
			Name:      "Feedbin feed-id:1602368 - 4 subscribers",
		},
		{
			UserAgent: "http.rb/5.1.0",
			Name:      "http.rb/5.1.0",
		},
		{
			UserAgent: "com.apple.webkit.networking",
			Name:      "com.apple.webkit.networking",
		},
		{
			UserAgent: "vuln_scanner",
			Name:      "vuln_scanner",
		},
		{
			UserAgent: "",
			Name:      "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			if name := UserAgentName(testCase.UserAgent); name != testCase.Name {
				t.Errorf("got %q, want %q", name, testCase.Name)
			}
		})
	}
}

func TestClassifyUserAgent(t *testing.T) {
	testCases := []struct {
		UserAgent string
		Category  Category
	}{
		{
			UserAgent: "curl/7.64.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "curl",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Prometheus/2.48.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "-",
			Category:  CategoryUnclassified,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			if category := ClassifyUserAgent(testCase.UserAgent); category != testCase.Category {
				t.Errorf("got %q, want %q", category, testCase.Category)
			}
		})
	}
}
