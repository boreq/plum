package request

type Status struct {
	s string
}

func NewStatus(s string) Status {
	return Status{s: s}
}

func (s Status) String() string {
	return s.s
}
