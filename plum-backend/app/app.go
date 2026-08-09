package app

import (
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/domain"
)

var (
	ErrWebsiteNotFound = errors.New("website not found")
	ErrDataNotFound    = errors.New("data not found")
)

type Repositories interface {
	Get(name domain.WebsiteName) (*domain.Repository, bool)
	Names() []domain.WebsiteName
	RemoveOldData(now time.Time)
}

type Application struct {
	GetWebsites     *GetWebsitesHandler
	GetHour         *GetHourHandler
	GetDay          *GetDayHandler
	GetMonth        *GetMonthHandler
	GetRangeHourly  *GetRangeHourlyHandler
	GetRangeDaily   *GetRangeDailyHandler
	GetRangeMonthly *GetRangeMonthlyHandler
	RemoveOldData   *RemoveOldDataHandler
	AddRequest      *AddRequestHandler
}

func New(repositories Repositories) *Application {
	return &Application{
		GetWebsites:     NewGetWebsitesHandler(repositories),
		GetHour:         NewGetHourHandler(repositories),
		GetDay:          NewGetDayHandler(repositories),
		GetMonth:        NewGetMonthHandler(repositories),
		GetRangeHourly:  NewGetRangeHourlyHandler(repositories),
		GetRangeDaily:   NewGetRangeDailyHandler(repositories),
		GetRangeMonthly: NewGetRangeMonthlyHandler(repositories),
		RemoveOldData:   NewRemoveOldDataHandler(repositories),
		AddRequest:      NewAddRequestHandler(repositories),
	}
}

type RangeResult struct {
	Summary *domain.Summary
	Series  []SeriesPoint
}

type SeriesPoint struct {
	Time     time.Time
	Visits   int
	Hits     int
	Bytes    int
	Statuses map[string]int
}

func NewSeriesPoint(t time.Time, summary *domain.Summary) SeriesPoint {
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
