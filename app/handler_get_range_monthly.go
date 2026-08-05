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

func (h *GetRangeMonthlyHandler) Execute(query GetRangeMonthly) ([]RangeData, error) {
	repository, ok := h.repositories.Get(query.Website)
	if !ok {
		return nil, ErrWebsiteNotFound
	}

	var rv []RangeData
	for month := query.From; !month.After(query.To); month = month.Next() {
		summary, ok := repository.RetrieveMonth(month.Year(), month.Month(), query.Filter)
		if !ok {
			return nil, errors.Wrap(ErrDataNotFound, "could not retrieve the month")
		}
		rv = append(rv, RangeData{Time: month.StartingPoint(), Data: summary})
	}

	return rv, nil
}
