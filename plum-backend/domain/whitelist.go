package domain

import (
	"net/netip"

	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/domain/request"
)

type Whitelist struct {
	addresses map[netip.Addr]struct{}
}

func NewWhitelist(entries []string) (Whitelist, error) {
	addresses := make(map[netip.Addr]struct{})

	for _, entry := range entries {
		if entry == "" {
			return Whitelist{}, errors.New("blank whitelist entry")
		}

		address, err := netip.ParseAddr(entry)
		if err != nil {
			return Whitelist{}, errors.Wrapf(err, "invalid whitelist entry '%s'", entry)
		}

		addresses[address.Unmap()] = struct{}{}
	}

	return Whitelist{addresses: addresses}, nil
}

func (w Whitelist) Contains(remoteAddress request.RemoteAddress) bool {
	if len(w.addresses) == 0 {
		return false
	}

	address, err := netip.ParseAddr(remoteAddress.String())
	if err != nil {
		return false
	}

	_, ok := w.addresses[address.Unmap()]
	return ok
}
