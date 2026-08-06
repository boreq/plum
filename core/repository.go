package core

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/boreq/plum/config"
	"github.com/boreq/plum/logging"
	"github.com/boreq/plum/parser"
)

const entryKeyFormat = "2006-01-02 15"

// visitPrefixFormat is used to generate a visit prefix which prevents
// identical visits from different days from getting merged.
const visitPrefixFormat = "2006-01-02"

// RetentionPeriod specifies how old the stored data is allowed to be. Older
// entries are not inserted and the data which becomes too old is eventually
// discarded by RemoveOldData.
const RetentionPeriod = 365 * 24 * time.Hour

type Repository struct {
	data      map[string]*Data
	dataMutex sync.Mutex
	conf      config.Website
	log       logging.Logger
}

func NewRepository(conf config.Website) *Repository {
	rv := &Repository{
		data: make(map[string]*Data),
		log:  logging.New("repository"),
		conf: conf,
	}
	return rv
}

func (r *Repository) Insert(entry *parser.Entry) error {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	if entry.Time.Before(retentionCutoff(time.Now())) {
		return nil
	}

	r.normalize(entry)

	key := r.createKey(entry.Time)
	data, ok := r.data[key]
	if !ok {
		data = NewData()
		r.data[key] = data
	}
	return data.Insert(entry)
}

func (r *Repository) RetrieveHour(year int, month time.Month, day int, hour int, filter Filter) (*Summary, bool) {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	target := NewSummary()

	t := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
	key := r.createKey(t)
	if d, ok := r.data[key]; ok {
		visitPrefix := t.Format(visitPrefixFormat)
		mergeData(target, d, visitPrefix, filter)
	}
	return target, true
}

func (r *Repository) RetrieveDay(year int, month time.Month, day int, filter Filter) (*Summary, bool) {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	target := NewSummary()

	for _, t := range iterateDay(year, month, day) {
		key := r.createKey(t)
		if d, ok := r.data[key]; ok {
			visitPrefix := t.Format(visitPrefixFormat)
			mergeData(target, d, visitPrefix, filter)
		}
	}
	return target, true
}

func (r *Repository) RetrieveMonth(year int, month time.Month, filter Filter) (*Summary, bool) {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	target := NewSummary()

	for _, t := range iterateMonth(year, month) {
		key := r.createKey(t)
		if d, ok := r.data[key]; ok {
			visitPrefix := t.Format(visitPrefixFormat)
			mergeData(target, d, visitPrefix, filter)
		}
	}
	return target, true
}

// RemoveOldData discards the data which is older than the retention period.
func (r *Repository) RemoveOldData(now time.Time) {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	cutoff := retentionCutoff(now)

	for key := range r.data {
		t, err := time.ParseInLocation(entryKeyFormat, key, time.UTC)
		if err != nil {
			r.log.Error("could not parse a key", "err", err, "key", key)
			continue
		}

		if t.Before(cutoff) {
			delete(r.data, key)
		}
	}
}

func retentionCutoff(now time.Time) time.Time {
	return now.UTC().Add(-RetentionPeriod)
}

func (r *Repository) createKey(date time.Time) string {
	return date.UTC().Format(entryKeyFormat)
}

func (r *Repository) normalize(entry *parser.Entry) {
	if r.conf.NormalizeQuery {
		u, err := url.ParseRequestURI(entry.HttpRequestURI)
		if err != nil {
			if entry.Status != "400" {
				r.log.Warn("query normalization failed", "err", err, "entry", entry)
			}
		} else {
			u.RawQuery = ""
			entry.HttpRequestURI = u.String()
		}
	}
	if r.conf.NormalizeSlash {
		if len(entry.HttpRequestURI) > 1 {
			entry.HttpRequestURI = strings.TrimRight(entry.HttpRequestURI, "/")
		}
	}
	if r.conf.StripRefererProtocol {
		entry.Referer = strings.TrimPrefix(entry.Referer, "http://")
		entry.Referer = strings.TrimPrefix(entry.Referer, "https://")
	}
}

func iterateDay(year int, month time.Month, day int) []time.Time {
	var result []time.Time
	start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)
	for t := start; t.Before(end); t = t.Add(time.Hour) {
		result = append(result, t)
	}
	return result
}

func iterateMonth(year int, month time.Month) []time.Time {
	var result []time.Time
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	for t := start; t.Before(end); t = t.Add(time.Hour) {
		result = append(result, t)
	}
	return result
}

func mergeData(target *Summary, source *Data, visitPrefix string, filter Filter) {
	for category, categoryData := range source.Categories {
		categoryMatches := filter.MatchesCategory(category)

		for uri, uriData := range categoryData.Uris {
			if !filter.MatchesUri(uri) {
				continue
			}

			for status, statusData := range uriData.Statuses {
				if !filter.MatchesStatus(status) {
					continue
				}

				for referer, refererData := range statusData.Referers {
					if !filter.MatchesReferer(referer) {
						continue
					}

					for userAgent, userAgentData := range refererData.UserAgents {
						if !filter.MatchesUserAgent(userAgent) {
							continue
						}

						target.InsertCategoryLeaf(category, userAgentData.Metrics, visitPrefix)

						if categoryMatches {
							target.InsertLeaf(uri, status, referer, userAgent, userAgentData.Metrics, visitPrefix)
						}
					}
				}
			}
		}
	}
}
