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

func (h *GetRangeHourlyHandler) Execute(query GetRangeHourly) (RangeResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return RangeResult{}, ErrWebsiteNotFound
	}

	summary := core.NewSummary()

	var series []SeriesPoint
	for hour := query.From; !hour.After(query.To); hour = hour.Next() {
		hourSummary, ok := repository.RetrieveHour(hour.Year(), hour.Month(), hour.Day(), hour.Hour(), query.Filter)
		if !ok {
			return RangeResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the hour")
		}

		series = append(series, NewSeriesPoint(hour.StartingPoint(), hourSummary))

		summary.Merge(hourSummary)
	}

	return RangeResult{Summary: summary, Series: series}, nil
}
