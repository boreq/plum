package core

import (
	"testing"

	"github.com/boreq/plum/parser"
)

func TestCreateVisitHash(t *testing.T) {
	entry := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "user agent",
	}
	h := createVisitHash(entry)
	if len(h) != retainHashBytes {
		t.Fatalf("length was %d", len(h))
	}
}

func TestCreateVisitHashDependsOnAddressAndUserAgent(t *testing.T) {
	base := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "user agent",
	}

	sameVisitor := &parser.Entry{
		RemoteAddress:  "1.2.3.4",
		UserAgent:      "user agent",
		HttpRequestURI: "/other",
	}

	otherAddress := &parser.Entry{
		RemoteAddress: "5.6.7.8",
		UserAgent:     "user agent",
	}

	otherUserAgent := &parser.Entry{
		RemoteAddress: "1.2.3.4",
		UserAgent:     "other agent",
	}

	if createVisitHash(base) != createVisitHash(sameVisitor) {
		t.Fatal("the uri must not affect the visit hash")
	}

	if createVisitHash(base) == createVisitHash(otherAddress) {
		t.Fatal("the remote address must affect the visit hash")
	}

	if createVisitHash(base) == createVisitHash(otherUserAgent) {
		t.Fatal("the user agent must affect the visit hash")
	}
}
