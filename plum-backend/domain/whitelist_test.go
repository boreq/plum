package domain

import (
	"testing"

	"github.com/boreq/plum/plum-backend/domain/request"
)

func TestWhitelistContains(t *testing.T) {
	testCases := []struct {
		Name          string
		Entries       []string
		RemoteAddress string
		Contains      bool
	}{
		{
			Name:          "empty whitelist",
			Entries:       nil,
			RemoteAddress: "1.2.3.4",
			Contains:      false,
		},
		{
			Name:          "exact address",
			Entries:       []string{"1.2.3.4"},
			RemoteAddress: "1.2.3.4",
			Contains:      true,
		},
		{
			Name:          "other address",
			Entries:       []string{"1.2.3.4"},
			RemoteAddress: "1.2.3.5",
			Contains:      false,
		},
		{
			Name:          "second entry matches",
			Entries:       []string{"1.2.3.4", "5.6.7.8"},
			RemoteAddress: "5.6.7.8",
			Contains:      true,
		},
		{
			Name:          "ipv6 address",
			Entries:       []string{"2001:db8::1"},
			RemoteAddress: "2001:db8::1",
			Contains:      true,
		},
		{
			Name:          "ipv6 address written differently",
			Entries:       []string{"2001:0db8:0000:0000:0000:0000:0000:0001"},
			RemoteAddress: "2001:db8::1",
			Contains:      true,
		},
		{
			Name:          "ipv4 mapped in ipv6 form",
			Entries:       []string{"1.2.3.4"},
			RemoteAddress: "::ffff:1.2.3.4",
			Contains:      true,
		},
		{
			Name:          "loopback",
			Entries:       []string{"127.0.0.1"},
			RemoteAddress: "127.0.0.1",
			Contains:      true,
		},
		{
			Name:          "garbage address",
			Entries:       []string{"1.2.3.4"},
			RemoteAddress: "not an address",
			Contains:      false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			whitelist, err := NewWhitelist(testCase.Entries)
			if err != nil {
				t.Fatalf("could not create the whitelist: %s", err)
			}

			contains := whitelist.Contains(request.NewRemoteAddress(testCase.RemoteAddress))
			if contains != testCase.Contains {
				t.Fatalf("expected %t, got %t", testCase.Contains, contains)
			}
		})
	}
}

func TestNewWhitelistRejectsInvalidEntries(t *testing.T) {
	testCases := []struct {
		Name  string
		Entry string
	}{
		{
			Name:  "blank",
			Entry: "",
		},
		{
			Name:  "not an address",
			Entry: "example.com",
		},
		{
			Name:  "cidr range",
			Entry: "10.0.0.0/8",
		},
		{
			Name:  "address with a port",
			Entry: "1.2.3.4:80",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if _, err := NewWhitelist([]string{testCase.Entry}); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
