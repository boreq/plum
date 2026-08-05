package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetRangeDaily struct {
	Website string
	From    core.Day
	To      core.Day
	Filter  core.Filter
}

type GetRangeDailyHandler struct {
	repositories *core.Repositories
}

func NewGetRangeDailyHandler(repositories *core.Repositories) *GetRangeDailyHandler {
	return &GetRangeDailyHandler{
		repositories: repositories,
	}
}

func (h *GetRangeDailyHandler) Execute(query GetRangeDaily) ([]RangeData, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	var rv []RangeData
	for day := query.From; !day.After(query.To); day = day.Next() {
		summary, ok := repository.RetrieveDay(day.Year(), day.Month(), day.Day(), query.Filter)
		if !ok {
			return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the day")
		}
		rv = append(rv, RangeData{Time: day.StartingPoint(), Data: summary})
	}

	return rv, nil
}
