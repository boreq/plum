package core

type Filter struct {
	Uri     string
	Status  string
	Referer string
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
