package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/plum/core"
	"github.com/boreq/plum/logging"
	"github.com/boreq/plum/server/api"
	_ "github.com/boreq/plum/statik"
	"github.com/julienschmidt/httprouter"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
)

var log = logging.New("server")

type handler struct {
	repositories *core.Repositories
}

func (h *handler) Website(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	return h.repositories.Names(), nil
}

func (h *handler) Hour(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	year, err := getParamInt(ps, "year")
	if err != nil {
		return nil, api.BadRequest
	}

	month, err := getParamInt(ps, "month")
	if err != nil {
		return nil, api.BadRequest
	}

	day, err := getParamInt(ps, "day")
	if err != nil {
		return nil, api.BadRequest
	}

	hour, err := getParamInt(ps, "hour")
	if err != nil {
		return nil, api.BadRequest
	}

	if month < 1 || month > 12 {
		return nil, api.BadRequest
	}

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	data, ok := repository.RetrieveHour(year, time.Month(month), day, hour)
	if !ok {
		return nil, api.NotFound
	}
	rangeData := RangeData{
		Time: time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC),
		Data: data,
	}
	return rangeData, nil
}

func (h *handler) Day(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	year, err := getParamInt(ps, "year")
	if err != nil {
		return nil, api.BadRequest
	}

	month, err := getParamInt(ps, "month")
	if err != nil {
		return nil, api.BadRequest
	}

	day, err := getParamInt(ps, "day")
	if err != nil {
		return nil, api.BadRequest
	}

	if month < 1 || month > 12 {
		return nil, api.BadRequest
	}

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	data, ok := repository.RetrieveDay(year, time.Month(month), day)
	if !ok {
		return nil, api.NotFound
	}
	rangeData := RangeData{
		Time: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		Data: data,
	}
	return rangeData, nil
}

func (h *handler) Month(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	year, err := getParamInt(ps, "year")
	if err != nil {
		return nil, api.BadRequest
	}

	month, err := getParamInt(ps, "month")
	if err != nil {
		return nil, api.BadRequest
	}

	if month < 1 || month > 12 {
		return nil, api.BadRequest
	}

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	data, ok := repository.RetrieveMonth(year, time.Month(month))
	if !ok {
		return nil, api.NotFound
	}
	rangeData := RangeData{
		Time: time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC),
		Data: data,
	}
	return rangeData, nil
}

func (h *handler) RangeHourly(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	yearFrom, err := getParamInt(ps, "yearFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	monthFrom, err := getParamInt(ps, "monthFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	dayFrom, err := getParamInt(ps, "dayFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	hourFrom, err := getParamInt(ps, "hourFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	yearTo, err := getParamInt(ps, "yearTo")
	if err != nil {
		return nil, api.BadRequest
	}

	monthTo, err := getParamInt(ps, "monthTo")
	if err != nil {
		return nil, api.BadRequest
	}

	dayTo, err := getParamInt(ps, "dayTo")
	if err != nil {
		return nil, api.BadRequest
	}

	hourTo, err := getParamInt(ps, "hourTo")
	if err != nil {
		return nil, api.BadRequest
	}

	if monthFrom < 1 || monthFrom > 12 || monthTo < 1 || monthTo > 12 {
		return nil, api.BadRequest
	}

	from := time.Date(yearFrom, time.Month(monthFrom), dayFrom, hourFrom, 0, 0, 0, time.UTC)
	to := time.Date(yearTo, time.Month(monthTo), dayTo, hourTo, 0, 0, 0, time.UTC)

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	var response []RangeData
	for t := from; !t.After(to); t = t.Add(time.Hour) {
		data, ok := repository.RetrieveHour(t.Year(), t.Month(), t.Day(), t.Hour())
		if !ok {
			return nil, api.InternalServerError
		}
		rangeData := RangeData{
			Time: t,
			Data: data,
		}
		response = append(response, rangeData)

	}
	return response, nil
}

func (h *handler) RangeDaily(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	yearFrom, err := getParamInt(ps, "yearFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	monthFrom, err := getParamInt(ps, "monthFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	dayFrom, err := getParamInt(ps, "dayFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	yearTo, err := getParamInt(ps, "yearTo")
	if err != nil {
		return nil, api.BadRequest
	}

	monthTo, err := getParamInt(ps, "monthTo")
	if err != nil {
		return nil, api.BadRequest
	}

	dayTo, err := getParamInt(ps, "dayTo")
	if err != nil {
		return nil, api.BadRequest
	}

	if monthFrom < 1 || monthFrom > 12 || monthTo < 1 || monthTo > 12 {
		return nil, api.BadRequest
	}

	from := time.Date(yearFrom, time.Month(monthFrom), dayFrom, 0, 0, 0, 0, time.UTC)
	to := time.Date(yearTo, time.Month(monthTo), dayTo, 0, 0, 0, 0, time.UTC)

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	var response []RangeData
	for t := from; !t.After(to); t = t.AddDate(0, 0, 1) {
		data, ok := repository.RetrieveDay(t.Year(), t.Month(), t.Day())
		if !ok {
			return nil, api.InternalServerError
		}
		rangeData := RangeData{
			Time: t,
			Data: data,
		}
		response = append(response, rangeData)

	}
	return response, nil
}

func (h *handler) RangeMonthly(r *http.Request, ps httprouter.Params) (interface{}, api.Error) {
	yearFrom, err := getParamInt(ps, "yearFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	monthFrom, err := getParamInt(ps, "monthFrom")
	if err != nil {
		return nil, api.BadRequest
	}

	yearTo, err := getParamInt(ps, "yearTo")
	if err != nil {
		return nil, api.BadRequest
	}

	monthTo, err := getParamInt(ps, "monthTo")
	if err != nil {
		return nil, api.BadRequest
	}

	if monthFrom < 1 || monthFrom > 12 || monthTo < 1 || monthTo > 12 {
		return nil, api.BadRequest
	}

	from := time.Date(yearFrom, time.Month(monthFrom), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(yearTo, time.Month(monthTo), 1, 0, 0, 0, 0, time.UTC)

	name := ps.ByName("name")
	repository, ok := h.repositories.Get(name)
	if !ok {
		return nil, api.BadRequest
	}

	var response []RangeData
	for t := from; !t.After(to); t = t.AddDate(0, 1, 0) {
		data, ok := repository.RetrieveMonth(t.Year(), t.Month())
		if !ok {
			return nil, api.InternalServerError
		}
		rangeData := RangeData{
			Time: t,
			Data: data,
		}
		response = append(response, rangeData)

	}
	return response, nil
}

func getParamInt(ps httprouter.Params, name string) (int, error) {
	return strconv.Atoi(getParamString(ps, name))
}

func getParamString(ps httprouter.Params, name string) string {
	return strings.TrimSuffix(ps.ByName(name), ".json")
}

type RangeData struct {
	Time time.Time  `json:"time"`
	Data *core.Data `json:"data"`
}

func Serve(repositories *core.Repositories, address string) error {
	handler, err := newHandler(repositories)
	if err != nil {
		return err
	}

	// Add CORS middleware
	handler = cors.AllowAll().Handler(handler)

	// Add GZIP middleware
	handler = gziphandler.GzipHandler(handler)

	log.Info("starting listening", "address", address)
	return http.ListenAndServe(address, handler)
}

func newHandler(repositories *core.Repositories) (http.Handler, error) {
	h := &handler{
		repositories: repositories,
	}

	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	router := httprouter.New()

	// List websites
	router.GET("/api/websites", api.Wrap(h.Website))

	// Discrete endpoints
	router.GET("/api/websites/:name/hour/:year/:month/:day/:hour", api.Wrap(h.Hour))
	router.GET("/api/websites/:name/day/:year/:month/:day", api.Wrap(h.Day))
	router.GET("/api/websites/:name/month/:year/:month", api.Wrap(h.Month))

	// Range endpoints
	router.GET("/api/websites/:name/range/hourly/:yearFrom/:monthFrom/:dayFrom/:hourFrom/:yearTo/:monthTo/:dayTo/:hourTo", api.Wrap(h.RangeHourly))
	router.GET("/api/websites/:name/range/daily/:yearFrom/:monthFrom/:dayFrom/:yearTo/:monthTo/:dayTo", api.Wrap(h.RangeDaily))
	router.GET("/api/websites/:name/range/monthly/:yearFrom/:monthFrom/:yearTo/:monthTo", api.Wrap(h.RangeMonthly))

	// Frontend
	router.NotFound = http.FileServer(statikFS)

	return router, nil
}
