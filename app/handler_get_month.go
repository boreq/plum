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

func (h *GetMonthHandler) Execute(query GetMonth) (RangeData, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return RangeData{}, ErrWebsiteNotFound
	}

	summary, ok := repository.RetrieveMonth(query.Month.Year(), query.Month.Month(), query.Filter)
	if !ok {
		return RangeData{}, errors.Wrap(ErrDataNotFound, "could not retrieve the month")
	}

	return RangeData{Time: query.Month.StartingPoint(), Data: summary}, nil
}
