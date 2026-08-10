package domain

import (
	"crypto"
	_ "crypto/sha512"

	"github.com/boreq/plum/plum-backend/domain/request"
)

func NewData() *Data {
	return &Data{
		Categories: make(map[Category]*CategoryData),
	}
}

type Data struct {
	Categories map[Category]*CategoryData
}

func (d *Data) Insert(r request.Request, category Category) error {
	bodyBytesSent := r.BodyBytesSent().Int()

	categoryData := d.getOrCreateCategoryData(category)
	remoteAddressData := categoryData.getOrCreateRemoteAddressData(r.RemoteAddress())
	remoteAddressData.Insert(bodyBytesSent)
	uriData := remoteAddressData.getOrCreateUriData(r.Uri())
	statusData := uriData.getOrCreateStatusData(r.Status())
	refererData := statusData.getOrCreateRefererData(r.Referer())
	userAgentData := refererData.getOrCreateUserAgentData(r.UserAgent(), category)
	userAgentData.Insert(bodyBytesSent)

	return nil
}

func (d *Data) getOrCreateCategoryData(category Category) *CategoryData {
	categoryData, ok := d.Categories[category]
	if !ok {
		categoryData = &CategoryData{
			RemoteAddresses: make(map[request.RemoteAddress]*RemoteAddressData),
		}
		d.Categories[category] = categoryData
	}
	return categoryData
}

type CategoryData struct {
	RemoteAddresses map[request.RemoteAddress]*RemoteAddressData
}

func (b *CategoryData) getOrCreateRemoteAddressData(remoteAddress request.RemoteAddress) *RemoteAddressData {
	remoteAddressData, ok := b.RemoteAddresses[remoteAddress]
	if !ok {
		remoteAddressData = &RemoteAddressData{
			Uris: make(map[request.Uri]*UriData),
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
	Uris map[request.Uri]*UriData
}

func (b *RemoteAddressData) getOrCreateUriData(uri request.Uri) *UriData {
	uriData, ok := b.Uris[uri]
	if !ok {
		uriData = &UriData{
			Statuses: make(map[request.Status]*StatusData),
		}
		b.Uris[uri] = uriData
	}
	return uriData
}

type UriData struct {
	Statuses map[request.Status]*StatusData
}

func (b *UriData) getOrCreateStatusData(status request.Status) *StatusData {
	statusData, ok := b.Statuses[status]
	if !ok {
		statusData = &StatusData{
			Referers: make(map[request.Referer]*RefererData),
		}
		b.Statuses[status] = statusData
	}
	return statusData
}

type StatusData struct {
	Referers map[request.Referer]*RefererData
}

func (b *StatusData) getOrCreateRefererData(referer request.Referer) *RefererData {
	refererData, ok := b.Referers[referer]
	if !ok {
		refererData = &RefererData{
			UserAgents: make(map[request.UserAgent]*UserAgentData),
		}
		b.Referers[referer] = refererData
	}
	return refererData
}

type RefererData struct {
	UserAgents map[request.UserAgent]*UserAgentData
}

func (b *RefererData) getOrCreateUserAgentData(userAgent request.UserAgent, category Category) *UserAgentData {
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

func CreateVisitHash(visitPrefix string, remoteAddress request.RemoteAddress, userAgent request.UserAgent) string {
	h := visitHash.New()
	h.Write([]byte(visitPrefix))
	h.Write([]byte(remoteAddress.String()))
	h.Write([]byte(userAgent.String()))
	sum := h.Sum(nil)
	return string(sum)[:retainHashBytes]
}
