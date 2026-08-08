package core

func NewMetrics() Metrics {
	return Metrics{
		Visits: NewSet(),
	}
}

type Metrics struct {
	Visits        Set
	Hits          int
	BodyBytesSent int
}

type Counters struct {
	Hits          int
	BodyBytesSent int
}

func (c *Counters) Insert(bodyBytesSent int) {
	c.Hits++
	c.BodyBytesSent += bodyBytesSent
}

func (m *Metrics) Insert(visit string, counters Counters) {
	m.Visits.Add(visit)
	m.Hits += counters.Hits
	m.BodyBytesSent += counters.BodyBytesSent
}

func (m *Metrics) Add(other Metrics, visitPrefix string) {
	for visit := range other.Visits {
		m.Visits.Add(visitPrefix + visit)
	}
	m.Hits += other.Hits
	m.BodyBytesSent += other.BodyBytesSent
}
