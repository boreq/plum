package core

import (
	"crypto"
	_ "crypto/sha512"

	"github.com/boreq/plum/parser"
)

func NewData() *Data {
	return &Data{
		Uris: make(map[string]*UriData),
	}
}

type Data struct {
	Uris map[string]*UriData
}

func (d *Data) Insert(entry *parser.Entry) error {
	visit := createVisitHash(entry)

	uriData := d.getOrCreateUriData(entry.HttpRequestURI)
	statusData := uriData.getOrCreateStatusData(entry.Status)
	refererData := statusData.getOrCreateRefererData(entry.Referer)
	refererData.Insert(visit, entry.BodyBytesSent)

	return nil
}

func (d *Data) getOrCreateUriData(uri string) *UriData {
	uriData, ok := d.Uris[uri]
	if !ok {
		uriData = &UriData{
			Statuses: make(map[string]*StatusData),
		}
		d.Uris[uri] = uriData
	}
	return uriData
}

type UriData struct {
	Statuses map[string]*StatusData
}

func (b *UriData) getOrCreateStatusData(status string) *StatusData {
	statusData, ok := b.Statuses[status]
	if !ok {
		statusData = &StatusData{
			Referers: make(map[string]*RefererData),
		}
		b.Statuses[status] = statusData
	}
	return statusData
}

type StatusData struct {
	Referers map[string]*RefererData
}

func (b *StatusData) getOrCreateRefererData(referer string) *RefererData {
	refererData, ok := b.Referers[referer]
	if !ok {
		refererData = &RefererData{
			Metrics: NewMetrics(),
		}
		b.Referers[referer] = refererData
	}
	return refererData
}

type RefererData struct {
	Metrics
}

var visitHash = crypto.SHA512_256

const retainHashBytes = 8

func createVisitHash(entry *parser.Entry) string {
	h := visitHash.New()
	h.Write([]byte(entry.RemoteAddress))
	h.Write([]byte(entry.UserAgent))
	sum := h.Sum(nil)
	return string(sum)[:retainHashBytes]
}
