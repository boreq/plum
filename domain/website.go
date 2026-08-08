package domain

import "github.com/boreq/errors"

type WebsiteName struct {
	name string
}

func NewWebsiteName(name string) (WebsiteName, error) {
	if name == "" {
		return WebsiteName{}, errors.New("website name is empty")
	}

	return WebsiteName{name: name}, nil
}

func (w WebsiteName) String() string {
	return w.name
}
