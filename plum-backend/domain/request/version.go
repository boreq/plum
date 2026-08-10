package request

import "github.com/boreq/errors"

type Version struct {
	s string
}

func NewVersion(s string) (Version, error) {
	if s == "" {
		return Version{}, errors.New("version is empty")
	}

	return Version{s: s}, nil
}

func (v Version) String() string {
	return v.s
}
