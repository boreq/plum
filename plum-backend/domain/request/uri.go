package request

type Uri struct {
	s string
}

func NewUri(s string) Uri {
	return Uri{s: s}
}

func (u Uri) String() string {
	return u.s
}
