package request

type Method struct {
	s string
}

func NewMethod(s string) Method {
	return Method{s: s}
}

func (m Method) String() string {
	return m.s
}
