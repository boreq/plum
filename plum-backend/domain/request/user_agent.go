package request

type UserAgent struct {
	s string
}

func NewUserAgent(s string) UserAgent {
	return UserAgent{s: s}
}

func (u UserAgent) String() string {
	return u.s
}
