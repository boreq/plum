package domain

import "github.com/boreq/plum/plum-backend/domain/request"

type Filter struct {
	Category  Category
	Uri       *request.Uri
	Status    *request.Status
	Referer   *request.Referer
	UserAgent *request.UserAgent
}

func (f Filter) MatchesUserAgent(userAgent request.UserAgent) bool {
	return f.UserAgent == nil || *f.UserAgent == userAgent
}

func (f Filter) MatchesCategory(category Category) bool {
	return f.Category.IsZero() || f.Category == category
}

func (f Filter) MatchesUri(uri request.Uri) bool {
	return f.Uri == nil || *f.Uri == uri
}

func (f Filter) MatchesStatus(status request.Status) bool {
	return f.Status == nil || *f.Status == status
}

func (f Filter) MatchesReferer(referer request.Referer) bool {
	return f.Referer == nil || *f.Referer == referer
}
