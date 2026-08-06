package core

type Filter struct {
	Category  Category
	Uri       string
	Status    string
	Referer   string
	UserAgent string
}

func (f Filter) MatchesUserAgent(userAgent UserAgent) bool {
	return f.UserAgent == "" || f.UserAgent == userAgent.Name()
}

func (f Filter) MatchesCategory(category Category) bool {
	return f.Category.IsZero() || f.Category == category
}

func (f Filter) MatchesUri(uri string) bool {
	return f.Uri == "" || f.Uri == uri
}

func (f Filter) MatchesStatus(status string) bool {
	return f.Status == "" || f.Status == status
}

func (f Filter) MatchesReferer(referer string) bool {
	return f.Referer == "" || f.Referer == referer
}
