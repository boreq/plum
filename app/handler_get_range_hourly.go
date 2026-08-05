package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetRangeHourly struct {
	Website string
	From    core.Hour
	To      core.Hour
	Filter  core.Filter
}

type GetRangeHourlyHandler struct {
	repositories *core.Repositories
}

func NewGetRangeHourlyHandler(repositories *core.Repositories) *GetRangeHourlyHandler {
	return &GetRangeHourlyHandler{
		repositories: repositories,
	}
}

func (h *GetRangeHourlyHandler) Execute(query GetRangeHourly) ([]RangeData, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	var rv []RangeData
	for hour := query.From; !hour.After(query.To); hour = hour.Next() {
		summary, ok := repository.RetrieveHour(hour.Year(), hour.Month(), hour.Day(), hour.Hour(), query.Filter)
		if !ok {
			return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the hour")
		}
		rv = append(rv, RangeData{Time: hour.StartingPoint(), Data: summary})
	}

	return rv, nil
}
