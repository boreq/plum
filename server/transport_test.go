package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/boreq/plum/config"
	"github.com/boreq/plum/core"
	"github.com/boreq/plum/parser"
)

var entryTime = time.Date(2019, time.February, 28, 13, 30, 0, 0, time.UTC)

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

	j, err := json.Marshal(NewRangeData(entryTime.Truncate(time.Hour), data))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	expected := `{"time":"2019-02-28T13:00:00Z","data":{"visits":2,"hits":3,"bytes":250,"uris":{"/index.html":{"visits":2,"hits":3,"bytes":250}},"statuses":{"200":{"visits":2,"hits":2,"bytes":200},"404":{"visits":1,"hits":1,"bytes":50}},"referers":{"example.com":{"visits":1,"hits":2,"bytes":150},"other.example.com":{"visits":1,"hits":1,"bytes":100}}}}`

	if string(j) != expected {
		t.Fatalf("\n got: %s\nwant: %s", string(j), expected)
	}
}
