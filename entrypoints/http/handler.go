package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/plum/app"
	"github.com/boreq/plum/core"
	"github.com/boreq/plum/entrypoints/http/api"
	_ "github.com/boreq/plum/statik"
	"github.com/julienschmidt/httprouter"
	"github.com/rakyll/statik/fs"
)

type Handler struct {
	app    *app.Application
	router *httprouter.Router
}

func NewHandler(app *app.Application) (*Handler, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, errors.Wrap(err, "could not create the statik filesystem")
	}

	h := &Handler{
		app:    app,
		router: httprouter.New(),
	}

	// List websites
	h.router.GET("/api/websites", api.Wrap(h.websites))

	// Discrete endpoints
	h.router.GET("/api/websites/:name/hour/:year/:month/:day/:hour", api.Wrap(h.hour))
	h.router.GET("/api/websites/:name/day/:year/:month/:day", api.Wrap(h.day))
	h.router.GET("/api/websites/:name/month/:year/:month", api.Wrap(h.month))

	// Range endpoints
	h.router.GET("/api/websites/:name/range/hourly/:yearFrom/:monthFrom/:dayFrom/:hourFrom/:yearTo/:monthTo/:dayTo/:hourTo", api.Wrap(h.rangeHourly))
	h.router.GET("/api/websites/:name/range/daily/:yearFrom/:monthFrom/:dayFrom/:yearTo/:monthTo/:dayTo", api.Wrap(h.rangeDaily))
	h.router.GET("/api/websites/:name/range/monthly/:yearFrom/:monthFrom/:yearTo/:monthTo", api.Wrap(h.rangeMonthly))

	// Frontend
	h.router.NotFound = http.FileServer(statikFS)

	return h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) websites(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	websites, err := h.app.GetWebsites.Execute(app.GetWebsites{})
	if err != nil {
		return nil, mapError(err)
	}
	return websites, nil
}

func (h *Handler) hour(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	hour, err := getHour(ps, "year", "month", "day", "hour")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetHour.Execute(app.GetHour{
		Website: ps.ByName("name"),
		Hour:    hour,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeData(rangeData), nil
}

func (h *Handler) day(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	day, err := getDay(ps, "year", "month", "day")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetDay.Execute(app.GetDay{
		Website: ps.ByName("name"),
		Day:     day,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeData(rangeData), nil
}

func (h *Handler) month(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	month, err := getMonth(ps, "year", "month")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetMonth.Execute(app.GetMonth{
		Website: ps.ByName("name"),
		Month:   month,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeData(rangeData), nil
}

func (h *Handler) rangeHourly(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	from, err := getHour(ps, "yearFrom", "monthFrom", "dayFrom", "hourFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	to, err := getHour(ps, "yearTo", "monthTo", "dayTo", "hourTo")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetRangeHourly.Execute(app.GetRangeHourly{
		Website: ps.ByName("name"),
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeDataSlice(rangeData), nil
}

func (h *Handler) rangeDaily(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	from, err := getDay(ps, "yearFrom", "monthFrom", "dayFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	to, err := getDay(ps, "yearTo", "monthTo", "dayTo")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetRangeDaily.Execute(app.GetRangeDaily{
		Website: ps.ByName("name"),
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeDataSlice(rangeData), nil
}

func (h *Handler) rangeMonthly(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	from, err := getMonth(ps, "yearFrom", "monthFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	to, err := getMonth(ps, "yearTo", "monthTo")
	if err != nil {
		return nil, api.BadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return nil, api.BadRequest
	}

	rangeData, err := h.app.GetRangeMonthly.Execute(app.GetRangeMonthly{
		Website: ps.ByName("name"),
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return NewRangeDataSlice(rangeData), nil
}

func mapError(err error) api.Error {
	switch {
	case errors.Is(err, app.ErrWebsiteNotFound):
		return api.BadRequest
	case errors.Is(err, app.ErrDataNotFound):
		return api.NotFound
	default:
		return api.InternalServerError
	}
}

func getFilter(r *http.Request) (core.Filter, error) {
	q := r.URL.Query()

	category, err := getCategory(q.Get("category"))
	if err != nil {
		return core.Filter{}, errors.Wrap(err, "could not get the category")
	}

	return core.Filter{
		Category:  category,
		Uri:       q.Get("uri"),
		Status:    q.Get("status"),
		Referer:   q.Get("referer"),
		UserAgent: q.Get("userAgent"),
	}, nil
}

func getCategory(s string) (core.Category, error) {
	if s == "" {
		return core.Category{}, nil
	}

	for _, category := range core.Categories {
		if category.String() == s {
			return category, nil
		}
	}

	return core.Category{}, errors.New("unknown category")
}

func getHour(ps httprouter.Params, yearName, monthName, dayName, hourName string) (core.Hour, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return core.Hour{}, errors.Wrap(err, "could not get the year and the month")
	}

	day, err := getParamInt(ps, dayName)
	if err != nil {
		return core.Hour{}, errors.Wrap(err, "could not get the day")
	}

	hour, err := getParamInt(ps, hourName)
	if err != nil {
		return core.Hour{}, errors.Wrap(err, "could not get the hour")
	}

	return core.NewHour(year, month, day, hour)
}

func getDay(ps httprouter.Params, yearName, monthName, dayName string) (core.Day, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return core.Day{}, errors.Wrap(err, "could not get the year and the month")
	}

	day, err := getParamInt(ps, dayName)
	if err != nil {
		return core.Day{}, errors.Wrap(err, "could not get the day")
	}

	return core.NewDay(year, month, day)
}

func getMonth(ps httprouter.Params, yearName, monthName string) (core.Month, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return core.Month{}, errors.Wrap(err, "could not get the year and the month")
	}

	return core.NewMonth(year, month)
}

func getYearAndMonth(ps httprouter.Params, yearName, monthName string) (int, time.Month, error) {
	year, err := getParamInt(ps, yearName)
	if err != nil {
		return 0, 0, errors.Wrap(err, "could not get the year")
	}

	month, err := getParamInt(ps, monthName)
	if err != nil {
		return 0, 0, errors.Wrap(err, "could not get the month")
	}

	return year, time.Month(month), nil
}

func getParamInt(ps httprouter.Params, name string) (int, error) {
	return strconv.Atoi(getParamString(ps, name))
}

func getParamString(ps httprouter.Params, name string) string {
	return strings.TrimSuffix(ps.ByName(name), ".json")
}
