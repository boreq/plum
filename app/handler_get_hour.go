package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/domain"
)

type GetHour struct {
	Website string
	Hour    domain.Hour
	Filter  domain.Filter
}

type GetHourHandler struct {
	repositories *domain.Repositories
}

func NewGetHourHandler(repositories *domain.Repositories) *GetHourHandler {
	return &GetHourHandler{
		repositories: repositories,
	}
}

func (h *GetHourHandler) Execute(query GetHour) (*domain.Summary, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveHour(query.Hour.Year(), query.Hour.Month(), query.Hour.Day(), query.Hour.Hour(), query.Filter)
	if !ok {
		return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the hour")
	}

	return summary, nil
}
