package adapters

import (
	"testing"
	"time"

	"github.com/boreq/plum/plum-backend/config"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
	"github.com/boreq/plum/plum-backend/domain/request"
)

const classifierBrowserUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"

// Entries which are older than the retention period are not inserted at all so
// this has to be a recent point in time.
var entryTime = time.Now().UTC().Add(-time.Hour)

var testClassifier = domain.NewTrafficClassifier()

func newRepository(t *testing.T) *Repository {
	t.Helper()

	repositories := NewRepositories()

	name, err := domain.NewWebsiteName("website")
	if err != nil {
		t.Fatalf("website name: %v", err)
	}

	repository := NewRepository(config.Website{}, NewMaliciousAddresses())
	if err := repositories.Add(name, repository); err != nil {
		t.Fatalf("add: %v", err)
	}

	return repository
}

func newTestEntry(t time.Time) *parser.Entry {
	return &parser.Entry{
		Time:           t,
		RemoteAddress:  "127.0.0.1",
		HttpRequestURI: "/",
		Status:         "200",
	}
}

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

func testEntryWithUserAgent(userAgent, remoteAddress, uri, status, referer string, bytes int) *parser.Entry {
	entry := testEntry(remoteAddress, uri, status, referer, bytes)
	entry.UserAgent = userAgent
	return entry
}

func insert(t *testing.T, r *Repository, entry *parser.Entry) {
	t.Helper()

	category := testClassifier.Classify(request.NewUri(entry.HttpRequestURI), request.NewUserAgent(entry.UserAgent), entry.Time)

	r.maliciousAddresses.Insert(entry, category)

	if err := r.Insert(entry, category); err != nil {
		t.Fatalf("insert: %v", err)
	}
}

func testRepository(t *testing.T, entries ...*parser.Entry) *Repository {
	t.Helper()

	r := newRepository(t)
	for _, entry := range entries {
		insert(t, r, entry)
	}
	return r
}

func retrieve(t *testing.T, r *Repository, filter domain.Filter) *domain.Summary {
	t.Helper()

	data, ok := r.RetrieveHour(entryTime.Year(), entryTime.Month(), entryTime.Day(), entryTime.Hour(), filter)
	if !ok {
		t.Fatal("retrieve failed")
	}
	return data
}

func TestInsertIgnoresOldEntries(t *testing.T) {
	r := newRepository(t)

	old := time.Now().UTC().Add(-RetentionPeriod).Add(-time.Hour)
	insert(t, r, newTestEntry(old))

	if len(r.data) != 0 {
		t.Fatalf("error: %v", r.data)
	}
}

func TestRemoveOldData(t *testing.T) {
	r := newRepository(t)

	now := time.Now().UTC()
	insert(t, r, newTestEntry(now.Add(-time.Hour)))

	if len(r.data) != 1 {
		t.Fatalf("error: %v", r.data)
	}

	r.RemoveOldData(now)
	if len(r.data) != 1 {
		t.Fatalf("error: %v", r.data)
	}

	r.RemoveOldData(now.Add(RetentionPeriod))
	if len(r.data) != 0 {
		t.Fatalf("error: %v", r.data)
	}
}

func TestRetrieveMarksEarlierTrafficOfMaliciousAddresses(t *testing.T) {
	r := newRepository(t)

	now := time.Now().UTC()
	day := now.Add(-24 * time.Hour)

	entry := newTestEntry(day)
	entry.UserAgent = classifierBrowserUserAgent
	insert(t, r, entry)

	summary, _ := r.RetrieveDay(day.Year(), day.Month(), day.Day(), domain.Filter{})
	if hits := summary.Categories[domain.CategoryUnclassified].Hits; hits != 1 {
		t.Fatalf("got %d unclassified hits", hits)
	}

	for i := 0; i <= domain.MaliciousRequestThreshold; i++ {
		scan := newTestEntry(now)
		scan.RemoteAddress = entry.RemoteAddress
		scan.UserAgent = classifierBrowserUserAgent
		scan.HttpRequestURI = "/.env"
		insert(t, r, scan)
	}

	summary, _ = r.RetrieveDay(day.Year(), day.Month(), day.Day(), domain.Filter{})
	if hits := summary.Categories[domain.CategoryUnclassified].Hits; hits != 0 {
		t.Fatalf("got %d unclassified hits", hits)
	}
	if hits := summary.Categories[domain.CategoryMalicious].Hits; hits != 1 {
		t.Fatalf("got %d malicious hits", hits)
	}

	summary, _ = r.RetrieveDay(day.Year(), day.Month(), day.Day(), domain.Filter{Category: domain.CategoryMalicious})
	if hits := summary.Hits; hits != 1 {
		t.Fatalf("got %d hits for the malicious filter", hits)
	}
}

func TestRetrieveIgnoresMaliciousTrafficOutsideOfTheWindow(t *testing.T) {
	r := newRepository(t)

	now := time.Now().UTC()
	day := now.Add(-domain.TrafficWindow).Add(-48 * time.Hour)

	entry := newTestEntry(day)
	entry.UserAgent = classifierBrowserUserAgent
	insert(t, r, entry)

	for i := 0; i <= domain.MaliciousRequestThreshold; i++ {
		scan := newTestEntry(now)
		scan.RemoteAddress = entry.RemoteAddress
		scan.UserAgent = classifierBrowserUserAgent
		scan.HttpRequestURI = "/.env"
		insert(t, r, scan)
	}

	summary, _ := r.RetrieveDay(day.Year(), day.Month(), day.Day(), domain.Filter{})
	if hits := summary.Categories[domain.CategoryUnclassified].Hits; hits != 1 {
		t.Fatalf("got %d unclassified hits", hits)
	}
	if hits := summary.Categories[domain.CategoryMalicious].Hits; hits != 0 {
		t.Fatalf("got %d malicious hits", hits)
	}
}

func TestRetrieveFilter(t *testing.T) {
	r := testRepository(t,
		testEntry("1.1.1.1", "/a", "200", "example.com", 10),
		testEntry("2.2.2.2", "/a", "404", "example.com", 20),
		testEntry("3.3.3.3", "/b", "200", "other.com", 40),
		testEntry("1.1.1.1", "/b", "200", "example.com", 80),
	)

	testCases := []struct {
		Name          string
		Filter        domain.Filter
		Hits          int
		Visits        int
		BodyBytesSent int
		Uris          []string
	}{
		{
			Name:          "no filter",
			Filter:        domain.Filter{},
			Hits:          4,
			Visits:        3,
			BodyBytesSent: 150,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "uri",
			Filter:        domain.Filter{Uri: "/a"},
			Hits:          2,
			Visits:        2,
			BodyBytesSent: 30,
			Uris:          []string{"/a"},
		},
		{
			Name:          "status",
			Filter:        domain.Filter{Status: "200"},
			Hits:          3,
			Visits:        2,
			BodyBytesSent: 130,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "referer",
			Filter:        domain.Filter{Referer: "example.com"},
			Hits:          3,
			Visits:        2,
			BodyBytesSent: 110,
			Uris:          []string{"/a", "/b"},
		},
		{
			Name:          "uri and status",
			Filter:        domain.Filter{Uri: "/a", Status: "404"},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 20,
			Uris:          []string{"/a"},
		},
		{
			Name:          "all dimensions",
			Filter:        domain.Filter{Uri: "/b", Status: "200", Referer: "other.com"},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 40,
			Uris:          []string{"/b"},
		},
		{
			Name:          "no match",
			Filter:        domain.Filter{Uri: "/a", Referer: "other.com"},
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

func TestRetrieveAggregatesFromLeaves(t *testing.T) {
	r := testRepository(t,
		testEntry("1.1.1.1", "/a", "200", "example.com", 10),
		testEntry("2.2.2.2", "/a", "200", "other.com", 20),
	)

	data := retrieve(t, r, domain.Filter{Referer: "example.com"})

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

func TestRetrieveCategoryFilter(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent(classifierBrowserUserAgent, "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent(classifierBrowserUserAgent, "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/c", "200", "example.com", 40),
	)

	testCases := []struct {
		Name          string
		Filter        domain.Filter
		Hits          int
		Visits        int
		BodyBytesSent int
		Uris          []string
	}{
		{
			Name:          "no filter",
			Filter:        domain.Filter{},
			Hits:          3,
			Visits:        3,
			BodyBytesSent: 70,
			Uris:          []string{"/a", "/b", "/c"},
		},
		{
			Name:          "automated",
			Filter:        domain.Filter{Category: domain.CategoryAutomated},
			Hits:          1,
			Visits:        1,
			BodyBytesSent: 40,
			Uris:          []string{"/c"},
		},
		{
			Name:          "unclassified",
			Filter:        domain.Filter{Category: domain.CategoryUnclassified},
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

func TestRetrieveCategoryTotalsIgnoreTheCategoryFilter(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent(classifierBrowserUserAgent, "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent(classifierBrowserUserAgent, "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/c", "200", "example.com", 40),
	)

	for _, filter := range []domain.Filter{
		{},
		{Category: domain.CategoryAutomated},
		{Category: domain.CategoryUnclassified},
	} {
		data := retrieve(t, r, filter)

		automated, ok := data.Categories[domain.CategoryAutomated]
		if !ok {
			t.Fatalf("filter %v: missing automated category", filter)
		}

		if automated.Hits != 1 || automated.Visits.Size() != 1 || automated.BodyBytesSent != 40 {
			t.Errorf("filter %v: automated: got %d hits, %d visits, %d bytes", filter, automated.Hits, automated.Visits.Size(), automated.BodyBytesSent)
		}

		unclassified, ok := data.Categories[domain.CategoryUnclassified]
		if !ok {
			t.Fatalf("filter %v: missing unclassified category", filter)
		}

		if unclassified.Hits != 2 || unclassified.Visits.Size() != 2 || unclassified.BodyBytesSent != 30 {
			t.Errorf("filter %v: unclassified: got %d hits, %d visits, %d bytes", filter, unclassified.Hits, unclassified.Visits.Size(), unclassified.BodyBytesSent)
		}
	}
}

func TestRetrieveCategoryTotalsRespectOtherFilters(t *testing.T) {
	r := testRepository(t,
		testEntryWithUserAgent(classifierBrowserUserAgent, "1.1.1.1", "/a", "200", "example.com", 10),
		testEntryWithUserAgent(classifierBrowserUserAgent, "2.2.2.2", "/b", "200", "example.com", 20),
		testEntryWithUserAgent("curl/7.64.0", "3.3.3.3", "/a", "200", "example.com", 40),
	)

	data := retrieve(t, r, domain.Filter{Uri: "/b"})

	if automated := data.Categories[domain.CategoryAutomated]; automated.Hits != 0 {
		t.Errorf("automated hits: got %d, want 0", automated.Hits)
	}

	if unclassified := data.Categories[domain.CategoryUnclassified]; unclassified.Hits != 1 {
		t.Errorf("unclassified hits: got %d, want 1", unclassified.Hits)
	}
}

func TestIterateDay(t *testing.T) {
	r := iterateDay(2019, time.February, 28)
	if len(r) != 24 {
		t.Fatalf("error: %v", r)
	}
	for i := 0; i < 24; i++ {
		if !r[i].Equal(time.Date(2019, time.February, 28, i, 0, 0, 0, time.UTC)) {
			t.Fatalf("error [%d]: %v", i, r[i])
		}
	}
}
