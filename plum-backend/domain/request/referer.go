package request

type Referer struct {
	s string
}

func NewReferer(s string) Referer {
	return Referer{s: s}
}

func (r Referer) String() string {
	return r.s
}
