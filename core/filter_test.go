package core

import (
	"testing"
	"time"

	"github.com/boreq/plum/config"
	"github.com/boreq/plum/parser"
)

var entryTime = time.Date(2019, time.February, 28, 13, 30, 0, 0, time.UTC)

func testEntry(remoteAddress, uri, status, referer string, bytes int) *parser.Entry {
	return &parser.Entry{
		RemoteAddress:  remoteAddress,
		UserAgent:      "user agent",
		Time:           entryTime,
		HttpRequestURI: uri,
		Status:         status,
		Referer:        referer,
		BodyBytesSent:  bytes,
	}
}

func testRepository(t *testing.T, entries ...*parser.Entry) *Repository {
	t.Helper()

	r := NewRepository(config.Website{})
	for _, entry := range entries {
		if err := r.Insert(entry); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}
	return r
}

func retrieve(t *testing.T, r *Repository, filter Filter) *Summary {
	t.Helper()

	data, ok := r.RetrieveHour(entryTime.Year(), entryTime.Month(), entryTime.Day(), entryTime.Hour(), filter)
	if !ok {
		t.Fatal("retrieve failed")
	}
	return data
}

func TestFilter(t *testing.T) {
	r := testRepository(t,
		testEntry("1.1.1.1", "/a", "200", "example.com", 10),
		testEntry("2.2.2.2", "/a", "404", "example.com", 20),
		testEntry("3.3.3.3", "/b", "200", "other.com", 40),
		testEntry("1.1.1.1", "/b", "200", "example.com", 80),
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
			Hits:          4,
			Visits:        3,
			BodyBytesSent: 150,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "uri",
			Filter:        Filter{Uri: "/a"},
			Hits:          2,
			Visits:        2,
			BodyBytesSent: 30,
			Uris:          []string{"/a"},
		},
		{
			Name:          "status",
			Filter:        Filter{Status: "200"},
			Hits:          3,
			Visits:        2,
			BodyBytesSent: 130,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "referer",
			Filter:        Filter{Referer: "example.com"},
			Hits:          3,
			Visits:        2,
			BodyBytesSent: 110,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "uri and status",
			Filter:        Filter{Uri: "/a", Status: "404"},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 20,
			Uris:          []string{"/a"},
		},
		{
			Name:          "all dimensions",
			Filter:        Filter{Uri: "/b", Status: "200", Referer: "other.com"},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 40,
			Uris:          []string{"/b"},
		},
		{
			Name:          "no match",
			Filter:        Filter{Uri: "/a", Referer: "other.com"},
			Hits:          0,
			Visits:        0,
			BodyBytesSent: 0,
			Uris:          nil,
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

func TestFilterAggregatesFromLeaves(t *testing.T) {
	r := testRepository(t,
		testEntry("1.1.1.1", "/a", "200", "example.com", 10),
		testEntry("2.2.2.2", "/a", "200", "other.com", 20),
	)

	data := retrieve(t, r, Filter{Referer: "example.com"})

	uriMetrics, ok := data.Uris["/a"]
	if !ok {
		t.Fatal("missing uri")
	}

	if uriMetrics.Hits != 1 {
		t.Errorf("uri hits: got %d, want 1", uriMetrics.Hits)
	}

	if uriMetrics.BodyBytesSent != 10 {
		t.Errorf("uri bytes: got %d, want 10", uriMetrics.BodyBytesSent)
	}

	statusMetrics, ok := data.Statuses["200"]
	if !ok {
		t.Fatal("missing status")
	}

	if statusMetrics.Hits != 1 {
		t.Errorf("status hits: got %d, want 1", statusMetrics.Hits)
	}

	if len(data.Referers) != 1 {
		t.Errorf("referers: got %v, want only example.com", data.Referers)
	}
}
