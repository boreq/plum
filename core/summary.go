package core

func NewSummary() *Summary {
	return &Summary{
		Metrics:  NewMetrics(),
		Uris:     make(map[string]*Metrics),
		Statuses: make(map[string]*Metrics),
		Referers: make(map[string]*Metrics),
	}
}

type Summary struct {
	Metrics
	Uris     map[string]*Metrics
	Statuses map[string]*Metrics
	Referers map[string]*Metrics
}

func (s *Summary) InsertLeaf(uri, status, referer string, metrics Metrics, visitPrefix string) {
	s.Metrics.Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Uris, uri).Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Statuses, status).Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Referers, referer).Add(metrics, visitPrefix)
}

func getOrCreateMetrics(target map[string]*Metrics, key string) *Metrics {
	metrics, ok := target[key]
	if !ok {
		created := NewMetrics()
		metrics = &created
		target[key] = metrics
	}
	return metrics
}
