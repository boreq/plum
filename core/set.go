package core

func NewSet() Set {
	return make(Set)
}

type Set map[string]struct{}

func (s *Set) Add(value string) {
	(*s)[value] = struct{}{}
}

func (s Set) Size() int {
	return len(s)
}
