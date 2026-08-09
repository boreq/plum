package http

import (
	"time"

	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/domain"
)

type RangeResult struct {
	Summary Data          `json:"summary"`
	Series  []SeriesPoint `json:"series"`
}

type SeriesPoint struct {
	Time time.Time `json:"time"`
	Metrics
	Statuses map[string]int `json:"statuses"`
}

type Data struct {
	Metrics
	Categories map[string]Metrics          `json:"categories"`
	Uris       map[string]Metrics          `json:"uris"`
	Statuses   map[string]Metrics          `json:"statuses"`
	Referers   map[string]Metrics          `json:"referers"`
	UserAgents map[string]UserAgentMetrics `json:"userAgents"`
}

type UserAgentMetrics struct {
	Metrics
	Browser string `json:"browser"`
}

type Metrics struct {
	Visits        int `json:"visits"`
	Hits          int `json:"hits"`
	BodyBytesSent int `json:"bytes"`
}

func newWebsites(websites []domain.WebsiteName) []string {
	rv := make([]string, 0, len(websites))
	for _, website := range websites {
		rv = append(rv, website.String())
	}
	return rv
}

func NewRangeResult(rangeResult app.RangeResult) RangeResult {
	series := make([]SeriesPoint, 0, len(rangeResult.Series))
	for _, v := range rangeResult.Series {
		series = append(series, newSeriesPoint(v))
	}

	return RangeResult{
		Summary: newData(rangeResult.Summary),
		Series:  series,
	}
}

func newSeriesPoint(seriesPoint app.SeriesPoint) SeriesPoint {
	statuses := make(map[string]int)
	for status, hits := range seriesPoint.Statuses {
		if status == "" {
			continue
		}
		statuses[status[:1]+"xx"] += hits
	}

	return SeriesPoint{
		Time: seriesPoint.Time,
		Metrics: Metrics{
			Visits:        seriesPoint.Visits,
			Hits:          seriesPoint.Hits,
			BodyBytesSent: seriesPoint.Bytes,
		},
		Statuses: statuses,
	}
}

func newData(summary *domain.Summary) Data {
	return Data{
		Metrics:    newMetrics(&summary.Metrics),
		Categories: newCategoryMetricsMap(summary.Categories),
		Uris:       newMetricsMap(summary.Uris),
		Statuses:   newMetricsMap(summary.Statuses),
		Referers:   newMetricsMap(summary.Referers),
		UserAgents: newUserAgentMetricsMap(summary.UserAgents),
	}
}

func newCategoryMetricsMap(source map[domain.Category]*domain.Metrics) map[string]Metrics {
	rv := make(map[string]Metrics, len(source))
	for category, metrics := range source {
		rv[category.String()] = newMetrics(metrics)
	}
	return rv
}

func newUserAgentMetricsMap(source map[string]*domain.UserAgentMetrics) map[string]UserAgentMetrics {
	rv := make(map[string]UserAgentMetrics, len(source))
	for userAgent, userAgentMetrics := range source {
		browser := ""
		if userAgentMetrics.Browser != nil {
			browser = userAgentMetrics.Browser.String()
		}

		rv[userAgent] = UserAgentMetrics{
			Metrics: newMetrics(&userAgentMetrics.Metrics),
			Browser: browser,
		}
	}
	return rv
}

func newMetricsMap(source map[string]*domain.Metrics) map[string]Metrics {
	rv := make(map[string]Metrics, len(source))
	for key, metrics := range source {
		rv[key] = newMetrics(metrics)
	}
	return rv
}

func newMetrics(metrics *domain.Metrics) Metrics {
	return Metrics{
		Visits:        metrics.Visits.Size(),
		Hits:          metrics.Hits,
		BodyBytesSent: metrics.BodyBytesSent,
	}
}
