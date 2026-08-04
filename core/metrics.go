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

func (m *Metrics) Insert(visit string, bodyBytesSent int) {
	m.Visits.Add(visit)
	m.Hits++
	m.BodyBytesSent += bodyBytesSent
}

func (m *Metrics) Add(other Metrics, visitPrefix string) {
	for visit := range other.Visits {
		m.Visits.Add(visitPrefix + visit)
	}
	m.Hits += other.Hits
	m.BodyBytesSent += other.BodyBytesSent
}
