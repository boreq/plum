package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

type GetRangeMonthly struct {
	Website string
	From    core.Month
	To      core.Month
	Filter  core.Filter
}

type GetRangeMonthlyHandler struct {
	repositories *core.Repositories
}

func NewGetRangeMonthlyHandler(repositories *core.Repositories) *GetRangeMonthlyHandler {
	return &GetRangeMonthlyHandler{
		repositories: repositories,
	}
}

func (h *GetRangeMonthlyHandler) Execute(query GetRangeMonthly) (RangeResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return RangeResult{}, ErrWebsiteNotFound
	}

	summary := core.NewSummary()

	var series []SeriesPoint
	for month := query.From; !month.After(query.To); month = month.Next() {
		monthSummary, ok := repository.RetrieveMonth(month.Year(), month.Month(), query.Filter)
		if !ok {
			return RangeResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the month")
		}

		series = append(series, NewSeriesPoint(month.StartingPoint(), monthSummary))

		summary.Merge(monthSummary)
	}

	return RangeResult{Summary: summary, Series: series}, nil
}
