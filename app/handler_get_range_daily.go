package app

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/domain"
)

type GetRangeDaily struct {
	Website domain.WebsiteName
	From    domain.Day
	To      domain.Day
	Filter  domain.Filter
}

type GetRangeDailyHandler struct {
	repositories Repositories
}

func NewGetRangeDailyHandler(repositories Repositories) *GetRangeDailyHandler {
	return &GetRangeDailyHandler{
		repositories: repositories,
	}
}

func (h *GetRangeDailyHandler) Execute(query GetRangeDaily) (RangeResult, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return RangeResult{}, ErrWebsiteNotFound
	}

	summary := domain.NewSummary()

	var series []SeriesPoint
	for day := query.From; !day.After(query.To); day = day.Next() {
		daySummary, ok := repository.RetrieveDay(day.Year(), day.Month(), day.Day(), query.Filter)
		if !ok {
			return RangeResult{}, errors.Wrap(ErrDataNotFound, "could not retrieve the day")
		}

		series = append(series, NewSeriesPoint(day.StartingPoint(), daySummary))

		summary.Merge(daySummary)
	}

	return RangeResult{Summary: summary, Series: series}, nil
}
