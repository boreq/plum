package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/domain"
)

type GetDay struct {
	Website domain.WebsiteName
	Day     domain.Day
	Filter  domain.Filter
}

type GetDayHandler struct {
	repositories Repositories
}

func NewGetDayHandler(repositories Repositories) *GetDayHandler {
	return &GetDayHandler{
		repositories: repositories,
	}
}

func (h *GetDayHandler) Execute(query GetDay) (*domain.Summary, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveDay(query.Day.Year(), query.Day.Month(), query.Day.Day(), query.Filter)
	if !ok {
		return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the day")
	}

	return summary, nil
}
