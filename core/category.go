package core

type Category struct {
	s string
}

var (
	CategoryMalicious    = Category{"malicious"}
	CategoryAutomated    = Category{"automated"}
	CategoryUnclassified = Category{"unclassified"}
)

var Categories = []Category{
	CategoryMalicious,
	CategoryAutomated,
	CategoryUnclassified,
}

func (c Category) IsZero() bool {
	return c == Category{}
}

func (c Category) String() string {
	return c.s
}
