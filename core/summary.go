package core

func NewSummary() *Summary {
	summary := &Summary{
		Metrics:    NewMetrics(),
		Categories: make(map[Category]*Metrics),
		Uris:       make(map[string]*Metrics),
		Statuses:   make(map[string]*Metrics),
		Referers:   make(map[string]*Metrics),
		UserAgents: make(map[string]*UserAgentMetrics),
	}
	for _, category := range Categories {
		getOrCreateMetrics(summary.Categories, category)
	}
	return summary
}

type Summary struct {
	Metrics
	Categories map[Category]*Metrics
	Uris       map[string]*Metrics
	Statuses   map[string]*Metrics
	Referers   map[string]*Metrics
	UserAgents map[string]*UserAgentMetrics
}

type UserAgentMetrics struct {
	Metrics
	Browser *Browser
}

func (s *Summary) InsertLeaf(uri, status, referer, userAgent string, browser *Browser, metrics Metrics, visitPrefix string) {
	s.Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Uris, uri).Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Statuses, status).Add(metrics, visitPrefix)
	getOrCreateMetrics(s.Referers, referer).Add(metrics, visitPrefix)
	getOrCreateUserAgentMetrics(s.UserAgents, userAgent, browser).Add(metrics, visitPrefix)
}

func getOrCreateUserAgentMetrics(target map[string]*UserAgentMetrics, userAgent string, browser *Browser) *Metrics {
	userAgentMetrics, ok := target[userAgent]
	if !ok {
		userAgentMetrics = &UserAgentMetrics{Metrics: NewMetrics()}
		target[userAgent] = userAgentMetrics
	}

	if userAgentMetrics.Browser == nil {
		userAgentMetrics.Browser = browser
	}

	return &userAgentMetrics.Metrics
}

func (s *Summary) Merge(other *Summary) {
	s.Add(other.Metrics, "")

	for category, metrics := range other.Categories {
		getOrCreateMetrics(s.Categories, category).Add(*metrics, "")
	}

	for uri, metrics := range other.Uris {
		getOrCreateMetrics(s.Uris, uri).Add(*metrics, "")
	}

	for status, metrics := range other.Statuses {
		getOrCreateMetrics(s.Statuses, status).Add(*metrics, "")
	}

	for referer, metrics := range other.Referers {
		getOrCreateMetrics(s.Referers, referer).Add(*metrics, "")
	}

	for userAgent, userAgentMetrics := range other.UserAgents {
		getOrCreateUserAgentMetrics(s.UserAgents, userAgent, userAgentMetrics.Browser).Add(userAgentMetrics.Metrics, "")
	}
}

func (s *Summary) InsertCategoryLeaf(category Category, metrics Metrics, visitPrefix string) {
	getOrCreateMetrics(s.Categories, category).Add(metrics, visitPrefix)
}

func getOrCreateMetrics[K comparable](target map[K]*Metrics, key K) *Metrics {
	metrics, ok := target[key]
	if !ok {
		created := NewMetrics()
		metrics = &created
		target[key] = metrics
	}
	return metrics
}
