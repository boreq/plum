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

type PointResult struct {
	Time time.Time
	Data *core.Summary
}

type RangeResult struct {
	Summary *core.Summary
	Series  []SeriesPoint
}

type SeriesPoint struct {
	Time     time.Time
	Visits   int
	Hits     int
	Bytes    int
	Statuses map[string]int
}

func NewSeriesPoint(t time.Time, summary *core.Summary) SeriesPoint {
	statuses := make(map[string]int, len(summary.Statuses))
	for status, metrics := range summary.Statuses {
		statuses[status] = metrics.Hits
	}

	return SeriesPoint{
		Time:     t,
		Visits:   summary.Visits.Size(),
		Hits:     summary.Hits,
		Bytes:    summary.BodyBytesSent,
		Statuses: statuses,
	}
}
