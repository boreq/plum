package http

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/boreq/plum/plum-backend/adapters"
	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/config"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/request"
)

var rangeEntryTime = time.Now().UTC().Truncate(24 * time.Hour).Add(-12 * time.Hour)

func TestNewRangeResultJSON(t *testing.T) {
	repositories := adapters.NewRepositories()

	websiteName, err := domain.NewWebsiteName("website")
	if err != nil {
		t.Fatalf("website name: %v", err)
	}

	repository := adapters.NewRepository(config.Website{}, adapters.NewMaliciousAddresses())
	if err := repositories.Add(websiteName, repository); err != nil {
		t.Fatalf("add: %v", err)
	}

	entries := []request.Request{
		newTestRequest(rangeEntryTime, "/index.html", "200", 100),
		newTestRequest(rangeEntryTime.Add(time.Hour), "/other.html", "404", 50),
	}

	classifier := domain.NewTrafficClassifier()

	for _, entry := range entries {
		if err := repository.Insert(entry, classifier.Classify(entry.Uri(), entry.UserAgent(), entry.Timestamp())); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	var series []app.SeriesPoint
	summary := domain.NewSummary()

	for _, t2 := range []time.Time{rangeEntryTime, rangeEntryTime.Add(time.Hour)} {
		data, ok := repository.RetrieveHour(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), domain.Filter{})
		if !ok {
			t.Fatal("retrieve failed")
		}
		series = append(series, app.NewSeriesPoint(t2.Truncate(time.Hour), data))
		summary.Merge(data)
	}

	j, err := json.Marshal(NewRangeResult(app.RangeResult{Summary: summary, Series: series}))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	expected := fmt.Sprintf(`{"summary":{"categories":{"automated":{"visits":0,"hits":0,"bytes":0},"malicious":{"visits":0,"hits":0,"bytes":0},"possibly-automated":{"visits":0,"hits":0,"bytes":0},"unclassified":{"visits":1,"hits":2,"bytes":150}},"uris":{"/index.html":{"visits":1,"hits":1,"bytes":100},"/other.html":{"visits":1,"hits":1,"bytes":50}},"statuses":{"200":{"visits":1,"hits":1,"bytes":100},"404":{"visits":1,"hits":1,"bytes":50}},"referers":{"example.com":{"visits":1,"hits":2,"bytes":150}},"userAgents":{"User Agent":{"visits":1,"hits":2,"bytes":150,"browser":""}}},"series":[{"time":"%s","visits":1,"hits":1,"bytes":100,"statuses":{"2xx":1}},{"time":"%s","visits":1,"hits":1,"bytes":50,"statuses":{"4xx":1}}]}`,
		rangeEntryTime.Truncate(time.Hour).Format(time.RFC3339),
		rangeEntryTime.Add(time.Hour).Truncate(time.Hour).Format(time.RFC3339))

	if string(j) != expected {
		t.Fatalf("\n got: %s\nwant: %s", string(j), expected)
	}
}

func newTestRequest(t time.Time, uri, status string, bytes int) request.Request {
	bodyBytesSent, err := request.NewBodyBytesSent(bytes)
	if err != nil {
		panic(err)
	}

	return request.NewRequest(
		request.NewRemoteAddress("1.2.3.4"),
		t,
		request.NewMethod("GET"),
		request.NewUri(uri),
		request.NewVersion("HTTP/1.1"),
		request.NewStatus(status),
		bodyBytesSent,
		request.NewReferer("example.com"),
		request.NewUserAgent("User Agent"),
	)
}
