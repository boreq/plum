package request

import "github.com/boreq/errors"

type Status struct {
	s string
}

func NewStatus(s string) (Status, error) {
	if s == "" {
		return Status{}, errors.New("status is empty")
	}

	return Status{s: s}, nil
}

func (s Status) String() string {
	return s.s
}
