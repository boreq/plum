package request

import "github.com/boreq/errors"

type BodyBytesSent struct {
	n int
}

func NewBodyBytesSent(n int) (BodyBytesSent, error) {
	if n < 0 {
		return BodyBytesSent{}, errors.New("body bytes sent is negative")
	}

	return BodyBytesSent{n: n}, nil
}

func (b BodyBytesSent) Int() int {
	return b.n
}
