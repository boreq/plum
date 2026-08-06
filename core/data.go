package core

import (
	"crypto"
	_ "crypto/sha512"

	"github.com/boreq/plum/parser"
)

func NewData() *Data {
	return &Data{
		Categories: make(map[Category]*CategoryData),
	}
}

type Data struct {
	Categories map[Category]*CategoryData
}

func (d *Data) Insert(entry *parser.Entry) error {
	visit := createVisitHash(entry)
	category := ClassifyUserAgent(entry.UserAgent)

	categoryData := d.getOrCreateCategoryData(category)
	uriData := categoryData.getOrCreateUriData(entry.HttpRequestURI)
	statusData := uriData.getOrCreateStatusData(entry.Status)
	refererData := statusData.getOrCreateRefererData(entry.Referer)
	userAgentData := refererData.getOrCreateUserAgentData(UserAgentName(entry.UserAgent))
	userAgentData.Insert(visit, entry.BodyBytesSent)

	return nil
}

func (d *Data) getOrCreateCategoryData(category Category) *CategoryData {
	categoryData, ok := d.Categories[category]
	if !ok {
		categoryData = &CategoryData{
			Uris: make(map[string]*UriData),
		}
		d.Categories[category] = categoryData
	}
	return categoryData
}

type CategoryData struct {
	Uris map[string]*UriData
}

func (b *CategoryData) getOrCreateUriData(uri string) *UriData {
	uriData, ok := b.Uris[uri]
	if !ok {
		uriData = &UriData{
			Statuses: make(map[string]*StatusData),
		}
		b.Uris[uri] = uriData
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
			UserAgents: make(map[string]*UserAgentData),
		}
		b.Referers[referer] = refererData
	}
	return refererData
}

type RefererData struct {
	UserAgents map[string]*UserAgentData
}

func (b *RefererData) getOrCreateUserAgentData(userAgent string) *UserAgentData {
	userAgentData, ok := b.UserAgents[userAgent]
	if !ok {
		userAgentData = &UserAgentData{
			Metrics: NewMetrics(),
		}
		b.UserAgents[userAgent] = userAgentData
	}
	return userAgentData
}

type UserAgentData struct {
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
