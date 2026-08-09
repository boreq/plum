package domain

import (
	"testing"

	"github.com/boreq/plum/plum-backend/domain/parser"
)

func testEntryWithUserAgent(userAgent, remoteAddress, uri, status, referer string, bytes int) *parser.Entry {
	entry := testEntry(remoteAddress, uri, status, referer, bytes)
	entry.UserAgent = userAgent
	return entry
}

func TestFilterCategory(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/c", "200", "example.com", 40),
	)

	testCases := []struct {
		Name          string
		Filter        Filter
		Hits          int
		Visits        int
		BodyBytesSent int
		Uris          []string
	}{
		{
			Name:          "no filter",
			Filter:        Filter{},
			Hits:          3,
			Visits:        3,
			BodyBytesSent: 70,
			Uris:          []string{"/a", "/b", "/c"},
		},
		{
			Name:          "automated",
			Filter:        Filter{Category: CategoryAutomated},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 40,
			Uris:          []string{"/c"},
		},
		{
			Name:          "unclassified",
			Filter:        Filter{Category: CategoryUnclassified},
			Hits:          2,
			Visits:        2,
			BodyBytesSent: 30,
			Uris:          []string{"/a", "/b"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			data := retrieve(t, r, testCase.Filter)

			if data.Hits != testCase.Hits {
				t.Errorf("hits: got %d, want %d", data.Hits, testCase.Hits)
			}

			if data.Visits.Size() != testCase.Visits {
				t.Errorf("visits: got %d, want %d", data.Visits.Size(), testCase.Visits)
			}

			if data.BodyBytesSent != testCase.BodyBytesSent {
				t.Errorf("bytes: got %d, want %d", data.BodyBytesSent, testCase.BodyBytesSent)
			}

			if len(data.Uris) != len(testCase.Uris) {
				t.Fatalf("uris: got %v, want %v", data.Uris, testCase.Uris)
			}

			for _, uri := range testCase.Uris {
				if _, ok := data.Uris[uri]; !ok {
					t.Errorf("uris: missing %q", uri)
				}
			}
		})
	}
}

func TestCategoryTotalsIgnoreTheCategoryFilter(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/c", "200", "example.com", 40),
	)

	for _, filter := range []Filter{
		{},
		{Category: CategoryAutomated},
		{Category: CategoryUnclassified},
	} {
		data := retrieve(t, r, filter)

		automated, ok := data.Categories[CategoryAutomated]
		if !ok {
			t.Fatalf("filter %v: missing automated category", filter)
		}

		if automated.Hits != 1 || automated.Visits.Size() != 1 || automated.BodyBytesSent != 40 {
			t.Errorf("filter %v: automated: got %d hits, %d visits, %d bytes", filter, automated.Hits, automated.Visits.Size(), automated.BodyBytesSent)
		}

		unclassified, ok := data.Categories[CategoryUnclassified]
		if !ok {
			t.Fatalf("filter %v: missing unclassified category", filter)
		}

		if unclassified.Hits != 2 || unclassified.Visits.Size() != 2 || unclassified.BodyBytesSent != 30 {
			t.Errorf("filter %v: unclassified: got %d hits, %d visits, %d bytes", filter, unclassified.Hits, unclassified.Visits.Size(), unclassified.BodyBytesSent)
		}
	}
}

func TestCategoryTotalsRespectOtherFilters(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/a", "200", "example.com", 40),
	)

	data := retrieve(t, r, Filter{Uri: "/b"})

	if automated := data.Categories[CategoryAutomated]; automated.Hits != 0 {
		t.Errorf("automated hits: got %d, want 0", automated.Hits)
	}

	if unclassified := data.Categories[CategoryUnclassified]; unclassified.Hits != 1 {
		t.Errorf("unclassified hits: got %d, want 1", unclassified.Hits)
	}
}

func TestCategoriesArePresentWithoutData(t *testing.T) {
	summary := NewSummary()

	for _, category := range Categories {
		if _, ok := summary.Categories[category]; !ok {
			t.Errorf("missing category %q", category)
		}
	}
}
