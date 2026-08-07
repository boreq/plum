package core

import (
	"testing"

	"github.com/boreq/plum/parser"
)

func TestCreateVisitHash(t *testing.T) {
	entry := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "user agent",
	}
	h := createVisitHash(entry)
	if len(h) != retainHashBytes {
		t.Fatalf("length was %d", len(h))
	}
}

func TestCreateVisitHashDependsOnAddressAndUserAgent(t *testing.T) {
	base := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "user agent",
	}

	sameVisitor := &parser.Entry{
		RemoteAddress:  "1.2.3.4",
		UserAgent:      "user agent",
		HttpRequestURI: "/other",
	}

	otherAddress := &parser.Entry{
		RemoteAddress: "5.6.7.8",
		UserAgent:     "user agent",
	}

	otherUserAgent := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "other agent",
	}

	if createVisitHash(base) != createVisitHash(sameVisitor) {
		t.Fatal("the uri must not affect the visit hash")
	}

	if createVisitHash(base) == createVisitHash(otherAddress) {
		t.Fatal("the remote address must affect the visit hash")
	}

	if createVisitHash(base) == createVisitHash(otherUserAgent) {
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
