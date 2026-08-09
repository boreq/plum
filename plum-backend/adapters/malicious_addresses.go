package adapters

import (
	"sync"
	"time"

	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
	"github.com/boreq/plum/plum-backend/logging"
)

const maliciousAddressesKeyFormat = "2006-01-02"

type MaliciousAddresses struct {
	hits      map[string]map[string]int
	hitsMutex sync.Mutex
	log       logging.Logger
}

func NewMaliciousAddresses() *MaliciousAddresses {
	return &MaliciousAddresses{
		hits: make(map[string]map[string]int),
		log:  logging.New("malicious_addresses"),
	}
}

func (m *MaliciousAddresses) Insert(entry *parser.Entry, category domain.Category) {
	if category != domain.CategoryMalicious {
		return
	}

	if entry.Time.Before(retentionCutoff(time.Now())) {
		return
	}

	m.hitsMutex.Lock()
	defer m.hitsMutex.Unlock()

	key := m.createKey(entry.Time)

	dayHits, ok := m.hits[key]
	if !ok {
		dayHits = make(map[string]int)
		m.hits[key] = dayHits
	}

	dayHits[entry.RemoteAddress]++
}

func (m *MaliciousAddresses) IsIpMalicious(t time.Time, remoteAddress string) bool {
	m.hitsMutex.Lock()
	defer m.hitsMutex.Unlock()

	to := t.Add(domain.TrafficWindow)

	var hits int

	for day := startOfDay(t.Add(-domain.TrafficWindow)); !day.After(to); day = day.AddDate(0, 0, 1) {
		hits += m.hits[m.createKey(day)][remoteAddress]

		if hits > domain.MaliciousRequestThreshold {
			return true
		}
	}

	return false
}

func (m *MaliciousAddresses) RemoveOldData(now time.Time) {
	m.hitsMutex.Lock()
	defer m.hitsMutex.Unlock()

	cutoff := retentionCutoff(now)

	for key := range m.hits {
		t, err := time.ParseInLocation(maliciousAddressesKeyFormat, key, time.UTC)
		if err != nil {
			m.log.Error("could not parse a key", "err", err, "key", key)
			continue
		}

		if t.AddDate(0, 0, 1).Before(cutoff) {
			delete(m.hits, key)
		}
	}
}

func (m *MaliciousAddresses) createKey(date time.Time) string {
	return date.UTC().Format(maliciousAddressesKeyFormat)
}

func startOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
