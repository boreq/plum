package app

import (
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
)

var (
	ErrWebsiteNotFound = errors.New("website not found")
	ErrDataNotFound    = errors.New("data not found")
)

type Repository interface {
	Insert(entry *parser.Entry, category domain.Category) error
	RetrieveHour(year int, month time.Month, day int, hour int, filter domain.Filter) (*domain.Summary, bool)
	RetrieveDay(year int, month time.Month, day int, filter domain.Filter) (*domain.Summary, bool)
	RetrieveMonth(year int, month time.Month, filter domain.Filter) (*domain.Summary, bool)
}

type Repositories interface {
	Get(name domain.WebsiteName) (Repository, bool)
	Names() []domain.WebsiteName
	RemoveOldData(now time.Time)
}

type MaliciousAddresses interface {
	Insert(entry *parser.Entry, category domain.Category)
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

func New(repositories Repositories, maliciousAddresses MaliciousAddresses) *Application {
	return &Application{
		GetWebsites:     NewGetWebsitesHandler(repositories),
		GetHour:         NewGetHourHandler(repositories),
		GetDay:          NewGetDayHandler(repositories),
		GetMonth:        NewGetMonthHandler(repositories),
		GetRangeHourly:  NewGetRangeHourlyHandler(repositories),
		GetRangeDaily:   NewGetRangeDailyHandler(repositories),
		GetRangeMonthly: NewGetRangeMonthlyHandler(repositories),
		RemoveOldData:   NewRemoveOldDataHandler(repositories, maliciousAddresses),
		AddRequest:      NewAddRequestHandler(repositories, maliciousAddresses, domain.NewTrafficClassifier()),
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
