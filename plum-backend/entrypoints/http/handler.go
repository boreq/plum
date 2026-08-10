package http

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/request"
	_ "github.com/boreq/plum/plum-backend/entrypoints/http/statik"
	"github.com/boreq/plum/plum-backend/logging"
	"github.com/boreq/rest"
	"github.com/julienschmidt/httprouter"
	"github.com/rakyll/statik/fs"
)

const (
	headerRemoteAddress = "X-Plum-Remote-Address"
	headerUri           = "X-Plum-Uri"
	headerUserAgent     = "X-Plum-User-Agent"
	headerMalicious     = "X-Plum-Malicious"
)

type Handler struct {
	app    *app.Application
	router *httprouter.Router
	log    logging.Logger
}

func NewHandler(app *app.Application) (*Handler, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, errors.Wrap(err, "could not create the statik filesystem")
	}

	h := &Handler{
		app:    app,
		router: httprouter.New(),
		log:    logging.New("entrypoints/http.Handler"),
	}

	h.router.HandlerFunc(http.MethodGet, "/api/websites", rest.Wrap(h.websites))
	h.router.HandlerFunc(http.MethodGet, "/api/malicious", rest.Wrap(h.malicious))

	// Discrete endpoints
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/hour/:year/:month/:day/:hour", rest.Wrap(h.hour))
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/day/:year/:month/:day", rest.Wrap(h.day))
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/month/:year/:month", rest.Wrap(h.month))

	// Range endpoints
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/range/hourly/:yearFrom/:monthFrom/:dayFrom/:hourFrom/:yearTo/:monthTo/:dayTo/:hourTo", rest.Wrap(h.rangeHourly))
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/range/daily/:yearFrom/:monthFrom/:dayFrom/:yearTo/:monthTo/:dayTo", rest.Wrap(h.rangeDaily))
	h.router.HandlerFunc(http.MethodGet, "/api/websites/:name/range/monthly/:yearFrom/:monthFrom/:yearTo/:monthTo", rest.Wrap(h.rangeMonthly))

	// Frontend
	h.router.NotFound = http.FileServer(statikFS)

	return h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) websites(r *http.Request) rest.RestResponse {
	websites, err := h.app.GetWebsites.Execute(app.GetWebsites{})
	if err != nil {
		return mapError(err)
	}
	return rest.NewResponse(newWebsites(websites))
}

func (h *Handler) malicious(r *http.Request) rest.RestResponse {
	remoteAddress := r.Header.Get(headerRemoteAddress)
	if remoteAddress == "" {
		h.log.Error("the remote address header is empty", "header", headerRemoteAddress)
		return rest.ErrBadRequest
	}

	uri := r.Header.Get(headerUri)
	if uri == "" {
		h.log.Error("the uri header is empty", "header", headerUri)
		return rest.ErrBadRequest
	}

	malicious, err := h.app.IsRequestMalicious.Execute(app.IsRequestMalicious{
		RemoteAddress: request.NewRemoteAddress(remoteAddress),
		Uri:           request.NewUri(uri),
		UserAgent:     request.NewUserAgent(r.Header.Get(headerUserAgent)),
	})
	if err != nil {
		h.log.Error("could not check if the request is malicious", "err", err)
		return mapError(err)
	}

	return maliciousResponse(malicious)
}

func (h *Handler) hour(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	hour, err := getHour(ps, "year", "month", "day", "hour")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	summary, err := h.app.GetHour.Execute(app.GetHour{
		Website: websiteName,
		Hour:    hour,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(newData(summary))
}

func (h *Handler) day(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	day, err := getDay(ps, "year", "month", "day")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	summary, err := h.app.GetDay.Execute(app.GetDay{
		Website: websiteName,
		Day:     day,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(newData(summary))
}

func (h *Handler) month(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	month, err := getMonth(ps, "year", "month")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	summary, err := h.app.GetMonth.Execute(app.GetMonth{
		Website: websiteName,
		Month:   month,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(newData(summary))
}

func (h *Handler) rangeHourly(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	from, err := getHour(ps, "yearFrom", "monthFrom", "dayFrom", "hourFrom")
	if err != nil {
		return rest.ErrBadRequest
	}

	to, err := getHour(ps, "yearTo", "monthTo", "dayTo", "hourTo")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	rangeResult, err := h.app.GetRangeHourly.Execute(app.GetRangeHourly{
		Website: websiteName,
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(NewRangeResult(rangeResult))
}

func (h *Handler) rangeDaily(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	from, err := getDay(ps, "yearFrom", "monthFrom", "dayFrom")
	if err != nil {
		return rest.ErrBadRequest
	}

	to, err := getDay(ps, "yearTo", "monthTo", "dayTo")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	rangeResult, err := h.app.GetRangeDaily.Execute(app.GetRangeDaily{
		Website: websiteName,
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(NewRangeResult(rangeResult))
}

func (h *Handler) rangeMonthly(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	from, err := getMonth(ps, "yearFrom", "monthFrom")
	if err != nil {
		return rest.ErrBadRequest
	}

	to, err := getMonth(ps, "yearTo", "monthTo")
	if err != nil {
		return rest.ErrBadRequest
	}

	filter, err := getFilter(r)
	if err != nil {
		return rest.ErrBadRequest
	}

	websiteName, err := getWebsiteName(ps)
	if err != nil {
		return rest.ErrBadRequest
	}

	rangeResult, err := h.app.GetRangeMonthly.Execute(app.GetRangeMonthly{
		Website: websiteName,
		From:    from,
		To:      to,
		Filter:  filter,
	})
	if err != nil {
		return mapError(err)
	}

	return rest.NewResponse(NewRangeResult(rangeResult))
}

func getWebsiteName(ps httprouter.Params) (domain.WebsiteName, error) {
	return domain.NewWebsiteName(ps.ByName("name"))
}

func mapError(err error) rest.Error {
	switch {
	case errors.Is(err, app.ErrWebsiteNotFound):
		return rest.ErrBadRequest
	case errors.Is(err, app.ErrDataNotFound):
		return rest.ErrNotFound
	default:
		return rest.ErrInternalServerError
	}
}

func getFilter(r *http.Request) (domain.Filter, error) {
	q := r.URL.Query()

	category, err := getCategory(q.Get("category"))
	if err != nil {
		return domain.Filter{}, errors.Wrap(err, "could not get the category")
	}

	return domain.Filter{
		Category:  category,
		Uri:       getFilterValue(q, "uri", request.NewUri),
		Status:    getFilterValue(q, "status", request.NewStatus),
		Referer:   getFilterValue(q, "referer", request.NewReferer),
		UserAgent: getFilterValue(q, "userAgent", request.NewUserAgent),
	}, nil
}

func getFilterValue[T any](q url.Values, key string, newValue func(string) T) *T {
	if !q.Has(key) {
		return nil
	}

	value := newValue(q.Get(key))
	return &value
}

func getCategory(s string) (domain.Category, error) {
	if s == "" {
		return domain.Category{}, nil
	}

	for _, category := range domain.Categories {
		if category.String() == s {
			return category, nil
		}
	}

	return domain.Category{}, errors.New("unknown category")
}

func getHour(ps httprouter.Params, yearName, monthName, dayName, hourName string) (domain.Hour, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return domain.Hour{}, errors.Wrap(err, "could not get the year and the month")
	}

	day, err := getParamInt(ps, dayName)
	if err != nil {
		return domain.Hour{}, errors.Wrap(err, "could not get the day")
	}

	hour, err := getParamInt(ps, hourName)
	if err != nil {
		return domain.Hour{}, errors.Wrap(err, "could not get the hour")
	}

	return domain.NewHour(year, month, day, hour)
}

func getDay(ps httprouter.Params, yearName, monthName, dayName string) (domain.Day, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return domain.Day{}, errors.Wrap(err, "could not get the year and the month")
	}

	day, err := getParamInt(ps, dayName)
	if err != nil {
		return domain.Day{}, errors.Wrap(err, "could not get the day")
	}

	return domain.NewDay(year, month, day)
}

func getMonth(ps httprouter.Params, yearName, monthName string) (domain.Month, error) {
	year, month, err := getYearAndMonth(ps, yearName, monthName)
	if err != nil {
		return domain.Month{}, errors.Wrap(err, "could not get the year and the month")
	}

	return domain.NewMonth(year, month)
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
