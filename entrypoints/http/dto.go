package http

import (
	"time"

	"github.com/boreq/plum/app"
	"github.com/boreq/plum/core"
)

type RangeData struct {
	Time time.Time `json:"time"`
	Data Data      `json:"data"`
}

type Data struct {
	Metrics
	Uris     map[string]Metrics `json:"uris"`
	Statuses map[string]Metrics `json:"statuses"`
	Referers map[string]Metrics `json:"referers"`
}

type Metrics struct {
	Visits        int `json:"visits"`
	Hits          int `json:"hits"`
	BodyBytesSent int `json:"bytes"`
}

func NewRangeData(rangeData app.RangeData) RangeData {
	return RangeData{
		Time: rangeData.Time,
		Data: newData(rangeData.Data),
	}
}

func NewRangeDataSlice(rangeData []app.RangeData) []RangeData {
	rv := make([]RangeData, 0, len(rangeData))
	for _, v := range rangeData {
		rv = append(rv, NewRangeData(v))
	}
	return rv
}

func newData(summary *core.Summary) Data {
	return Data{
		Metrics:  newMetrics(&summary.Metrics),
		Uris:     newMetricsMap(summary.Uris),
		Statuses: newMetricsMap(summary.Statuses),
		Referers: newMetricsMap(summary.Referers),
	}
}

func newMetricsMap(source map[string]*core.Metrics) map[string]Metrics {
	rv := make(map[string]Metrics, len(source))
	for key, metrics := range source {
		rv[key] = newMetrics(metrics)
	}
	return rv
}

func newMetrics(metrics *core.Metrics) Metrics {
	return Metrics{
		Visits:        metrics.Visits.Size(),
		Hits:          metrics.Hits,
		BodyBytesSent: metrics.BodyBytesSent,
	}
}
