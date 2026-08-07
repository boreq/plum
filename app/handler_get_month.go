package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetMonth struct {
	Website string
	Month   core.Month
	Filter  core.Filter
}

type GetMonthHandler struct {
	repositories *core.Repositories
}

func NewGetMonthHandler(repositories *core.Repositories) *GetMonthHandler {
	return &GetMonthHandler{
		repositories: repositories,
	}
}

func (h *GetMonthHandler) Execute(query GetMonth) (PointResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return PointResult{}, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveMonth(query.Month.Year(), query.Month.Month(), query.Filter)
	if !ok {
		return PointResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the month")
	}

	return PointResult{Time: query.Month.StartingPoint(), Data: summary}, nil
}
