package core

import "testing"

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
