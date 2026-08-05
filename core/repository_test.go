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
