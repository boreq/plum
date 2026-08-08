package core

import (
	"testing"
	"time"

	"github.com/boreq/plum/config"
	"github.com/boreq/plum/parser"
)

func TestInsertIgnoresOldEntries(t *testing.T) {
	r := NewRepository(config.Website{})

	old := time.Now().UTC().Add(-RetentionPeriod).Add(-time.Hour)
	if err := r.Insert(newTestEntry(old)); err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(r.data) != 0 {
		t.Fatalf("error: %v", r.data)
	}
}

func TestRemoveOldData(t *testing.T) {
	r := NewRepository(config.Website{})

	now := time.Now().UTC()
	recent := now.Add(-time.Hour)
	if err := r.Insert(newTestEntry(recent)); err != nil {
		t.Fatalf("error: %v", err)
	}

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
	r := NewRepository(config.Website{})

	now := time.Now().UTC()
	day := now.Add(-24 * time.Hour)

	entry := newTestEntry(day)
	entry.UserAgent = classifierBrowserUserAgent
	if err := r.Insert(entry); err != nil {
		t.Fatalf("error: %v", err)
	}

	summary, _ := r.RetrieveDay(day.Year(), day.Month(), day.Day(), Filter{})
	if hits := summary.Categories[CategoryUnclassified].Hits; hits != 1 {
		t.Fatalf("got %d unclassified hits", hits)
	}

	for i := 0; i <= MaliciousRequestThreshold; i++ {
		scan := newTestEntry(now)
		scan.RemoteAddress = entry.RemoteAddress
		scan.UserAgent = classifierBrowserUserAgent
		scan.HttpRequestURI = "/.env"
		if err := r.Insert(scan); err != nil {
			t.Fatalf("error: %v", err)
		}
	}

	summary, _ = r.RetrieveDay(day.Year(), day.Month(), day.Day(), Filter{})
	if hits := summary.Categories[CategoryUnclassified].Hits; hits != 0 {
		t.Fatalf("got %d unclassified hits", hits)
	}
	if hits := summary.Categories[CategoryMalicious].Hits; hits != 1 {
		t.Fatalf("got %d malicious hits", hits)
	}

	summary, _ = r.RetrieveDay(day.Year(), day.Month(), day.Day(), Filter{Category: CategoryMalicious})
	if hits := summary.Metrics.Hits; hits != 1 {
		t.Fatalf("got %d hits for the malicious filter", hits)
	}
}

func newTestEntry(t time.Time) *parser.Entry {
	return &parser.Entry{
		Time:           t,
		RemoteAddress:  "127.0.0.1",
		HttpRequestURI: "/",
		Status:         "200",
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
