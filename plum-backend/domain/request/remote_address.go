package request

import "github.com/boreq/errors"

type RemoteAddress struct {
	s string
}

func NewRemoteAddress(s string) (RemoteAddress, error) {
	if s == "" {
		return RemoteAddress{}, errors.New("remote address is empty")
	}

	return RemoteAddress{s: s}, nil
}

func (r RemoteAddress) String() string {
	return r.s
}
