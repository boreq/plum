package core

import (
	"testing"

	"github.com/boreq/plum/parser"
)

func TestCreateVisitHash(t *testing.T) {
	h := createVisitHash("2020-01-01", "1.2.3.4", "user agent")
	if len(h) != retainHashBytes {
		t.Fatalf("length was %d", len(h))
	}
}

func TestCreateVisitHashDependsOnPrefixAddressAndUserAgent(t *testing.T) {
	base := createVisitHash("2020-01-01", "1.2.3.4", "user agent")

	if base != createVisitHash("2020-01-01", "1.2.3.4", "user agent") {
		t.Fatal("the visit hash must be stable")
	}

	if base == createVisitHash("2020-01-02", "1.2.3.4", "user agent") {
		t.Fatal("the visit prefix must affect the visit hash")
	}

	if base == createVisitHash("2020-01-01", "5.6.7.8", "user agent") {
		t.Fatal("the remote address must affect the visit hash")
	}

	if base == createVisitHash("2020-01-01", "1.2.3.4", "other agent") {
		t.Fatal("the user agent must affect the visit hash")
	}
}

func TestInsertStoresTheBrowserOfUnclassifiedTraffic(t *testing.T) {
	testCases := []struct {
		Name     string
		Category Category
		Browser  *Browser
	}{
		{
			Name:     "unclassified",
			Category: CategoryUnclassified,
			Browser:  BrowserFirefox,
		},
		{
			Name:     "automated",
			Category: CategoryAutomated,
			Browser:  nil,
		},
		{
			Name:     "malicious",
			Category: CategoryMalicious,
			Browser:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"

			data := NewData()
			if err := data.Insert(&parser.Entry{
				RemoteAddress:  "1.2.3.4",
				UserAgent:      userAgent,
				HttpRequestURI: "/",
				Status:         "200",
			}, testCase.Category); err != nil {
				t.Fatal(err)
			}

			userAgentData := data.
				Categories[testCase.Category].
				RemoteAddresses["1.2.3.4"].
				Uris["/"].
				Statuses["200"].
				Referers[""].
				UserAgents[userAgent]

			if userAgentData.Browser != testCase.Browser {
				t.Errorf("got %v, want %v", userAgentData.Browser, testCase.Browser)
			}
		})
	}
}
