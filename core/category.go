package core

type Category struct {
	s string
}

var (
	CategoryAutomated    = Category{"automated"}
	CategoryUnclassified = Category{"unclassified"}
)

var Categories = []Category{
	CategoryAutomated,
	CategoryUnclassified,
}

func (c Category) IsZero() bool {
	return c == Category{}
}

func (c Category) String() string {
	return c.s
}
