package request

type Version struct {
	s string
}

func NewVersion(s string) Version {
	return Version{s: s}
}

func (v Version) String() string {
	return v.s
}
