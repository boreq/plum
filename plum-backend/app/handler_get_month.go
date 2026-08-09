package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/domain"
)

type GetMonth struct {
	Website domain.WebsiteName
	Month   domain.Month
	Filter  domain.Filter
}

type GetMonthHandler struct {
	repositories Repositories
}

func NewGetMonthHandler(repositories Repositories) *GetMonthHandler {
	return &GetMonthHandler{
		repositories: repositories,
	}
}

func (h *GetMonthHandler) Execute(query GetMonth) (*domain.Summary, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveMonth(query.Month.Year(), query.Month.Month(), query.Filter)
	if !ok {
		return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the month")
	}

	return summary, nil
}
