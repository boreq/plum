package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetHour struct {
	Website string
	Hour    core.Hour
	Filter  core.Filter
}

type GetHourHandler struct {
	repositories *core.Repositories
}

func NewGetHourHandler(repositories *core.Repositories) *GetHourHandler {
	return &GetHourHandler{
		repositories: repositories,
	}
}

func (h *GetHourHandler) Execute(query GetHour) (PointResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return PointResult{}, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveHour(query.Hour.Year(), query.Hour.Month(), query.Hour.Day(), query.Hour.Hour(), query.Filter)
	if !ok {
		return PointResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the hour")
	}

	return PointResult{Time: query.Hour.StartingPoint(), Data: summary}, nil
}
