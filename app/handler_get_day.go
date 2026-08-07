package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetDay struct {
	Website string
	Day     core.Day
	Filter  core.Filter
}

type GetDayHandler struct {
	repositories *core.Repositories
}

func NewGetDayHandler(repositories *core.Repositories) *GetDayHandler {
	return &GetDayHandler{
		repositories: repositories,
	}
}

func (h *GetDayHandler) Execute(query GetDay) (PointResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return PointResult{}, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveDay(query.Day.Year(), query.Day.Month(), query.Day.Day(), query.Filter)
	if !ok {
		return PointResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the day")
	}

	return PointResult{Time: query.Day.StartingPoint(), Data: summary}, nil
}
