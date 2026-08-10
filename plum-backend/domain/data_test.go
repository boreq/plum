package domain

import (
	"testing"
	"time"

	"github.com/boreq/plum/plum-backend/domain/request"
)

func TestCreateVisitHash(t *testing.T) {
	h := CreateVisitHash("2020-01-01", request.NewRemoteAddress("1.2.3.4"), request.NewUserAgent("user agent"))
	if len(h) != retainHashBytes {
		t.Fatalf("length was %d", len(h))
	}
}

func TestCreateVisitHashDependsOnPrefixAddressAndUserAgent(t *testing.T) {
	base := CreateVisitHash("2020-01-01", request.NewRemoteAddress("1.2.3.4"), request.NewUserAgent("user agent"))

	if base != CreateVisitHash("2020-01-01", request.NewRemoteAddress("1.2.3.4"), request.NewUserAgent("user agent")) {
		t.Fatal("the visit hash must be stable")
	}

	if base == CreateVisitHash("2020-01-02", request.NewRemoteAddress("1.2.3.4"), request.NewUserAgent("user agent")) {
		t.Fatal("the visit prefix must affect the visit hash")
	}

	if base == CreateVisitHash("2020-01-01", request.NewRemoteAddress("5.6.7.8"), request.NewUserAgent("user agent")) {
		t.Fatal("the remote address must affect the visit hash")
	}

	if base == CreateVisitHash("2020-01-01", request.NewRemoteAddress("1.2.3.4"), request.NewUserAgent("other agent")) {
		t.Fatal("the user agent must affect the visit hash")
	}
}

func TestInsertStoresTheBrowserOfUnclassifiedTraffic(t *testing.T) {
	testCases := []struct {
		Name     string
		Category Category
		Browser  *Browser
	}{
		{
			Name:     "unclassified",
			Category: CategoryUnclassified,
			Browser:  BrowserFirefox,
		},
		{
			Name:     "automated",
			Category: CategoryAutomated,
			Browser:  nil,
		},
		{
			Name:     "malicious",
			Category: CategoryMalicious,
			Browser:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"

			bodyBytesSent, err := request.NewBodyBytesSent(0)
			if err != nil {
				t.Fatal(err)
			}

			data := NewData()
			if err := data.Insert(request.NewRequest(
				request.NewRemoteAddress("1.2.3.4"),
				time.Time{},
				request.NewMethod("GET"),
				request.NewUri("/"),
				request.NewVersion("HTTP/1.1"),
				request.NewStatus("200"),
				bodyBytesSent,
				request.NewReferer(""),
				request.NewUserAgent(userAgent),
			), testCase.Category); err != nil {
				t.Fatal(err)
			}

			userAgentData := data.
				Categories[testCase.Category].
				RemoteAddresses[request.NewRemoteAddress("1.2.3.4")].
				Uris[request.NewUri("/")].
				Statuses[request.NewStatus("200")].
				Referers[request.NewReferer("")].
				UserAgents[request.NewUserAgent(userAgent)]

			if userAgentData.Browser != testCase.Browser {
				t.Errorf("got %v, want %v", userAgentData.Browser, testCase.Browser)
			}
		})
	}
}
