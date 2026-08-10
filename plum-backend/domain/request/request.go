package request

import "time"

type Request struct {
	remoteAddress RemoteAddress
	timestamp     time.Time
	method        Method
	uri           Uri
	version       Version
	status        Status
	bodyBytesSent BodyBytesSent
	referer       Referer
	userAgent     UserAgent
}

func NewRequest(
	remoteAddress RemoteAddress,
	timestamp time.Time,
	method Method,
	uri Uri,
	version Version,
	status Status,
	bodyBytesSent BodyBytesSent,
	referer Referer,
	userAgent UserAgent,
) Request {
	return Request{
		remoteAddress: remoteAddress,
		timestamp:     timestamp,
		method:        method,
		uri:           uri,
		version:       version,
		status:        status,
		bodyBytesSent: bodyBytesSent,
		referer:       referer,
		userAgent:     userAgent,
	}
}

func (r Request) RemoteAddress() RemoteAddress {
	return r.remoteAddress
}

func (r Request) Timestamp() time.Time {
	return r.timestamp
}

func (r Request) Method() Method {
	return r.method
}

func (r Request) Uri() Uri {
	return r.uri
}

func (r Request) Version() Version {
	return r.version
}

func (r Request) Status() Status {
	return r.status
}

func (r Request) BodyBytesSent() BodyBytesSent {
	return r.bodyBytesSent
}

func (r Request) Referer() Referer {
	return r.referer
}

func (r Request) UserAgent() UserAgent {
	return r.userAgent
}
