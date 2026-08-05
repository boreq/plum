package app

import (
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/plum/core"
)

var (
	ErrWebsiteNotFound = errors.New("website not found")
	ErrDataNotFound    = errors.New("data not found")
)

type Application struct {
	GetWebsites     *GetWebsitesHandler
	GetHour         *GetHourHandler
	GetDay          *GetDayHandler
	GetMonth        *GetMonthHandler
	GetRangeHourly  *GetRangeHourlyHandler
	GetRangeDaily   *GetRangeDailyHandler
	GetRangeMonthly *GetRangeMonthlyHandler
	RemoveOldData   *RemoveOldDataHandler
}

func New(repositories *core.Repositories) *Application {
	return &Application{
		GetWebsites:     NewGetWebsitesHandler(repositories),
		GetHour:         NewGetHourHandler(repositories),
		GetDay:          NewGetDayHandler(repositories),
		GetMonth:        NewGetMonthHandler(repositories),
		GetRangeHourly:  NewGetRangeHourlyHandler(repositories),
		GetRangeDaily:   NewGetRangeDailyHandler(repositories),
		GetRangeMonthly: NewGetRangeMonthlyHandler(repositories),
		RemoveOldData:   NewRemoveOldDataHandler(repositories),
	}
}

type RangeData struct {
	Time time.Time
	Data *core.Summary
}
