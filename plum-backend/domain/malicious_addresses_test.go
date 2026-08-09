package domain

import (
	"testing"

	"github.com/boreq/plum/plum-backend/domain/parser"
)

func scannedData(t *testing.T, remoteAddress string, requests int) *Data {
	data := NewData()
	for i := 0; i < requests; i++ {
		if err := data.Insert(&parser.Entry{
			RemoteAddress:  remoteAddress,
			UserAgent:      classifierBrowserUserAgent,
			HttpRequestURI: "/.env",
			Status:         "404",
		}, CategoryMalicious); err != nil {
			t.Fatal(err)
		}
	}
	return data
}

func TestIsMaliciousMarksAddresses(t *testing.T) {
	m := NewMaliciousAddresses()
	m.Insert(scannedData(t, "1.1.1.1", MaliciousRequestThreshold+1))

	if !m.IsMalicious("1.1.1.1") {
		t.Error("the address should be malicious")
	}

	if m.IsMalicious("2.2.2.2") {
		t.Error("the address should not be malicious")
	}
}

func TestIsMaliciousIgnoresOtherCategories(t *testing.T) {
	data := NewData()
	for i := 0; i <= MaliciousRequestThreshold; i++ {
		if err := data.Insert(&parser.Entry{
			RemoteAddress:  "1.1.1.1",
			UserAgent:      classifierBrowserUserAgent,
			HttpRequestURI: "/",
			Status:         "200",
		}, CategoryUnclassified); err != nil {
			t.Fatal(err)
		}
	}

	m := NewMaliciousAddresses()
	m.Insert(data)

	if m.IsMalicious("1.1.1.1") {
		t.Error("the address should not be malicious")
	}
}

func TestIsMaliciousRequiresTheThreshold(t *testing.T) {
	m := NewMaliciousAddresses()
	m.Insert(scannedData(t, "1.1.1.1", MaliciousRequestThreshold))

	if m.IsMalicious("1.1.1.1") {
		t.Error("the address should not be malicious")
	}
}

func TestIsMaliciousSumsInsertedData(t *testing.T) {
	m := NewMaliciousAddresses()
	for i := 0; i <= MaliciousRequestThreshold; i++ {
		m.Insert(scannedData(t, "1.1.1.1", 1))
	}

	if !m.IsMalicious("1.1.1.1") {
		t.Error("the address should be malicious")
	}
}
