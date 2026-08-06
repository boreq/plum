package http

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/boreq/plum/app"
	"github.com/boreq/plum/config"
	"github.com/boreq/plum/core"
	"github.com/boreq/plum/parser"
)

// Entries which are older than the retention period are not inserted at all so
// this has to be a recent point in time.
var entryTime = time.Now().UTC().Add(-time.Hour)

func TestNewRangeDataJSON(t *testing.T) {
	repository := core.NewRepository(config.Website{})

	entries := []*parser.Entry{
		{
			RemoteAddress:  "1.2.3.4",
			UserAgent:      "user agent",
			Time:           entryTime,
			HttpRequestURI: "/index.html",
			Status:         "200",
			BodyBytesSent:  100,
			Referer:        "example.com",
		},
		{
			RemoteAddress:  "1.2.3.4",
			UserAgent:      "user agent",
			Time:           entryTime,
			HttpRequestURI: "/index.html",
			Status:         "404",
			BodyBytesSent:  50,
			Referer:        "example.com",
		},
		{
			RemoteAddress:  "5.6.7.8",
			UserAgent:      "other agent",
			Time:           entryTime,
			HttpRequestURI: "/index.html",
			Status:         "200",
			BodyBytesSent:  100,
			Referer:        "other.example.com",
		},
	}

	for _, entry := range entries {
		if err := repository.Insert(entry); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	data, ok := repository.RetrieveHour(entryTime.Year(), entryTime.Month(), entryTime.Day(), entryTime.Hour(), core.Filter{})
	if !ok {
		t.Fatal("retrieve failed")
	}

	j, err := json.Marshal(NewRangeData(app.RangeData{
		Time: entryTime.Truncate(time.Hour),
		Data: data,
	}))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	expected := fmt.Sprintf(`{"time":"%s","data":{"visits":2,"hits":3,"bytes":250,"categories":{"automated":{"visits":0,"hits":0,"bytes":0},"malicious":{"visits":0,"hits":0,"bytes":0},"unclassified":{"visits":2,"hits":3,"bytes":250}},"uris":{"/index.html":{"visits":2,"hits":3,"bytes":250}},"statuses":{"200":{"visits":2,"hits":2,"bytes":200},"404":{"visits":1,"hits":1,"bytes":50}},"referers":{"example.com":{"visits":1,"hits":2,"bytes":150},"other.example.com":{"visits":1,"hits":1,"bytes":100}},"userAgents":{"other agent":{"visits":1,"hits":1,"bytes":100},"user agent":{"visits":1,"hits":2,"bytes":150}}}}`,
		entryTime.Truncate(time.Hour).Format(time.RFC3339))

	if string(j) != expected {
		t.Fatalf("\n got: %s\nwant: %s", string(j), expected)
	}
}
