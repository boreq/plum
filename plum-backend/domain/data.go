package domain

import (
	"crypto"
	_ "crypto/sha512"

	"github.com/boreq/plum/plum-backend/domain/parser"
)

func NewData() *Data {
	return &Data{
		Categories: make(map[Category]*CategoryData),
	}
}

type Data struct {
	Categories map[Category]*CategoryData
}

func (d *Data) Insert(entry *parser.Entry, category Category) error {
	categoryData := d.getOrCreateCategoryData(category)
	remoteAddressData := categoryData.getOrCreateRemoteAddressData(entry.RemoteAddress)
	remoteAddressData.Insert(entry.BodyBytesSent)
	uriData := remoteAddressData.getOrCreateUriData(entry.HttpRequestURI)
	statusData := uriData.getOrCreateStatusData(entry.Status)
	refererData := statusData.getOrCreateRefererData(entry.Referer)
	userAgentData := refererData.getOrCreateUserAgentData(entry.UserAgent, category)
	userAgentData.Insert(entry.BodyBytesSent)

	return nil
}

func (d *Data) getOrCreateCategoryData(category Category) *CategoryData {
	categoryData, ok := d.Categories[category]
	if !ok {
		categoryData = &CategoryData{
			RemoteAddresses: make(map[string]*RemoteAddressData),
		}
		d.Categories[category] = categoryData
	}
	return categoryData
}

type CategoryData struct {
	RemoteAddresses map[string]*RemoteAddressData
}

func (b *CategoryData) getOrCreateRemoteAddressData(remoteAddress string) *RemoteAddressData {
	remoteAddressData, ok := b.RemoteAddresses[remoteAddress]
	if !ok {
		remoteAddressData = &RemoteAddressData{
			Uris: make(map[string]*UriData),
		}
		b.RemoteAddresses[remoteAddress] = remoteAddressData
	}
	return remoteAddressData
}

// RemoteAddressData aggregates the counters of the entire subtree so that the
// malicious requests made by an address can be counted without descending into
// it.
type RemoteAddressData struct {
	Counters
	Uris map[string]*UriData
}

func (b *RemoteAddressData) getOrCreateUriData(uri string) *UriData {
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

func (b *RefererData) getOrCreateUserAgentData(userAgent string, category Category) *UserAgentData {
	userAgentData, ok := b.UserAgents[userAgent]
	if !ok {
		var browser *Browser
		if category == CategoryUnclassified {
			browser = RecognizeBrowser(userAgent)
		}

		userAgentData = &UserAgentData{
			Browser: browser,
		}
		b.UserAgents[userAgent] = userAgentData
	}
	return userAgentData
}

type UserAgentData struct {
	Counters
	Browser *Browser
}

var visitHash = crypto.SHA512_256

const retainHashBytes = 8

func CreateVisitHash(visitPrefix, remoteAddress, userAgent string) string {
	h := visitHash.New()
	h.Write([]byte(visitPrefix))
	h.Write([]byte(remoteAddress))
	h.Write([]byte(userAgent))
	sum := h.Sum(nil)
	return string(sum)[:retainHashBytes]
}
