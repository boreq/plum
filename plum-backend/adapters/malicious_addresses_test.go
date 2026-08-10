package adapters

import (
	"testing"
	"time"

	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/request"
)

func scan(t time.Time, remoteAddress string) request.Request {
	return newTestRequest(remoteAddress, classifierBrowserUserAgent, t, "/.env", "404", "", 0)
}

func insertScans(m *MaliciousAddresses, t time.Time, remoteAddress string, requests int) {
	for i := 0; i < requests; i++ {
		m.Insert(scan(t, remoteAddress), domain.CategoryMalicious)
	}
}

func TestIsIpMaliciousMarksAddresses(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold+1)

	if !m.IsIpMalicious(now, request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should be malicious")
	}

	if m.IsIpMalicious(now, request.NewRemoteAddress("2.2.2.2")) {
		t.Error("the address should not be malicious")
	}
}

func TestIsIpMaliciousRequiresTheThreshold(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold)

	if m.IsIpMalicious(now, request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should not be malicious")
	}
}

func TestIsIpMaliciousIgnoresOtherCategories(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	for i := 0; i <= domain.MaliciousRequestThreshold; i++ {
		m.Insert(scan(now, "1.1.1.1"), domain.CategoryUnclassified)
	}

	if m.IsIpMalicious(now, request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should not be malicious")
	}
}

func TestIsIpMaliciousSumsDaysWithinTheWindow(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	for i := 0; i <= domain.MaliciousRequestThreshold; i++ {
		insertScans(m, now.AddDate(0, 0, -i), "1.1.1.1", 1)
	}

	if !m.IsIpMalicious(now, request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should be malicious")
	}
}

func TestIsIpMaliciousIgnoresTrafficOutsideOfTheWindow(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold+1)

	if m.IsIpMalicious(now.Add(-2*domain.TrafficWindow), request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should not be malicious")
	}

	if !m.IsIpMalicious(now.Add(-domain.TrafficWindow).AddDate(0, 0, 1), request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the address should be malicious")
	}
}

func TestInsertIgnoresOldScans(t *testing.T) {
	now := time.Now().UTC()
	old := now.Add(-RetentionPeriod).Add(-time.Hour)

	m := NewMaliciousAddresses(domain.Whitelist{})
	insertScans(m, old, "1.1.1.1", domain.MaliciousRequestThreshold+1)

	if len(m.hits) != 0 {
		t.Fatalf("error: %v", m.hits)
	}
}

func TestMaliciousAddressesRemoveOldData(t *testing.T) {
	now := time.Now().UTC()

	m := NewMaliciousAddresses(domain.Whitelist{})
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold+1)

	m.RemoveOldData(now)
	if len(m.hits) != 1 {
		t.Fatalf("error: %v", m.hits)
	}

	m.RemoveOldData(now.Add(RetentionPeriod).AddDate(0, 0, 1))
	if len(m.hits) != 0 {
		t.Fatalf("error: %v", m.hits)
	}
}

func TestIsIpMaliciousSkipsWhitelistedAddresses(t *testing.T) {
	now := time.Now().UTC()

	whitelist, err := domain.NewWhitelist([]string{"1.1.1.1"})
	if err != nil {
		t.Fatalf("could not create the whitelist: %s", err)
	}

	m := NewMaliciousAddresses(whitelist)
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold+1)
	insertScans(m, now, "2.2.2.2", domain.MaliciousRequestThreshold+1)

	if m.IsIpMalicious(now, request.NewRemoteAddress("1.1.1.1")) {
		t.Error("the whitelisted address should not be malicious")
	}

	if !m.IsIpMalicious(now, request.NewRemoteAddress("2.2.2.2")) {
		t.Error("the other address should be malicious")
	}
}

func TestWhitelistedAddressesAreStillCounted(t *testing.T) {
	now := time.Now().UTC()

	whitelist, err := domain.NewWhitelist([]string{"1.1.1.1"})
	if err != nil {
		t.Fatalf("could not create the whitelist: %s", err)
	}

	m := NewMaliciousAddresses(whitelist)
	insertScans(m, now, "1.1.1.1", domain.MaliciousRequestThreshold+1)

	hits := m.hits[m.createKey(now)]["1.1.1.1"]
	if hits != domain.MaliciousRequestThreshold+1 {
		t.Errorf("expected %d hits, got %d", domain.MaliciousRequestThreshold+1, hits)
	}
}
