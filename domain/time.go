package domain

import (
	"time"

	"github.com/boreq/errors"
)

var (
	ErrInvalidMonth = errors.New("invalid month")
	ErrInvalidDay   = errors.New("invalid day")
	ErrInvalidHour  = errors.New("invalid hour")
)

type Month struct {
	t time.Time
}

func NewMonth(year int, month time.Month) (Month, error) {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month {
		return Month{}, ErrInvalidMonth
	}

	return Month{t: t}, nil
}

func (m Month) Year() int {
	return m.t.Year()
}

func (m Month) Month() time.Month {
	return m.t.Month()
}

func (m Month) Next() Month {
	return Month{t: time.Date(m.t.Year(), m.t.Month()+1, 1, 0, 0, 0, 0, m.t.Location())}
}

func (m Month) After(o Month) bool {
	return m.t.After(o.t)
}

func (m Month) StartingPoint() time.Time {
	return m.t
}

type Day struct {
	t time.Time
}

func NewDay(year int, month time.Month, day int) (Day, error) {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Day{}, ErrInvalidDay
	}

	return Day{t: t}, nil
}

func (d Day) Year() int {
	return d.t.Year()
}

func (d Day) Month() time.Month {
	return d.t.Month()
}

func (d Day) Day() int {
	return d.t.Day()
}

func (d Day) Next() Day {
	return Day{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day()+1, 0, 0, 0, 0, d.t.Location())}
}

func (d Day) After(o Day) bool {
	return d.t.After(o.t)
}

func (d Day) StartingPoint() time.Time {
	return d.t
}

type Hour struct {
	t time.Time
}

func NewHour(year int, month time.Month, day int, hour int) (Hour, error) {
	t := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day || t.Hour() != hour {
		return Hour{}, ErrInvalidHour
	}

	return Hour{t: t}, nil
}

func (h Hour) Year() int {
	return h.t.Year()
}

func (h Hour) Month() time.Month {
	return h.t.Month()
}

func (h Hour) Day() int {
	return h.t.Day()
}

func (h Hour) Hour() int {
	return h.t.Hour()
}

func (h Hour) Next() Hour {
	return Hour{t: time.Date(h.t.Year(), h.t.Month(), h.t.Day(), h.t.Hour()+1, 0, 0, 0, h.t.Location())}
}

func (h Hour) After(o Hour) bool {
	return h.t.After(o.t)
}

func (h Hour) StartingPoint() time.Time {
	return h.t
}
