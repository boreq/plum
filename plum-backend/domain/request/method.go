package request

import "github.com/boreq/errors"

type Method struct {
	s string
}

func NewMethod(s string) (Method, error) {
	if s == "" {
		return Method{}, errors.New("method is empty")
	}

	return Method{s: s}, nil
}

func (m Method) String() string {
	return m.s
}
