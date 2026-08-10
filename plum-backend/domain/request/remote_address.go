package request

type RemoteAddress struct {
	s string
}

func NewRemoteAddress(s string) RemoteAddress {
	return RemoteAddress{s: s}
}

func (r RemoteAddress) String() string {
	return r.s
}
