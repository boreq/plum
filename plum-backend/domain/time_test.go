package domain

import (
	"testing"
	"time"
)

func TestNewMonth(t *testing.T) {
	testCases := []struct {
		Name  string
		Year  int
		Month time.Month
		Valid bool
	}{
		{Name: "valid", Year: 2026, Month: time.August, Valid: true},
		{Name: "month_too_large", Year: 2026, Month: time.Month(13), Valid: false},
		{Name: "month_zero", Year: 2026, Month: time.Month(0), Valid: false},
		{Name: "month_negative", Year: 2026, Month: time.Month(-1), Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			m, err := NewMonth(testCase.Year, testCase.Month)
			if testCase.Valid {
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				if m.Year() != testCase.Year || m.Month() != testCase.Month {
					t.Fatalf("error: %v", m)
				}
			} else if err == nil {
				t.Fatalf("expected an error, got: %v", m)
			}
		})
	}
}

func TestNext(t *testing.T) {
	month, err := NewMonth(2026, time.December)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if next := month.Next(); next.Year() != 2027 || next.Month() != time.January {
		t.Fatalf("error: %v", next)
	}

	day, err := NewDay(2024, time.February, 28)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if next := day.Next(); next.Month() != time.February || next.Day() != 29 {
		t.Fatalf("error: %v", next)
	}

	hour, err := NewHour(2026, time.August, 5, 23)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if next := hour.Next(); next.Day() != 6 || next.Hour() != 0 {
		t.Fatalf("error: %v", next)
	}
}

func TestNewDay(t *testing.T) {
	testCases := []struct {
		Name  string
		Year  int
		Month time.Month
		Day   int
		Valid bool
	}{
		{Name: "valid", Year: 2026, Month: time.August, Day: 5, Valid: true},
		{Name: "leap_day", Year: 2024, Month: time.February, Day: 29, Valid: true},
		{Name: "not_a_leap_day", Year: 2026, Month: time.February, Day: 29, Valid: false},
		{Name: "day_too_large", Year: 2026, Month: time.August, Day: 32, Valid: false},
		{Name: "day_zero", Year: 2026, Month: time.August, Day: 0, Valid: false},
		{Name: "invalid_month", Year: 2026, Month: time.Month(13), Day: 1, Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			d, err := NewDay(testCase.Year, testCase.Month, testCase.Day)
			if testCase.Valid {
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				if d.Year() != testCase.Year || d.Month() != testCase.Month || d.Day() != testCase.Day {
					t.Fatalf("error: %v", d)
				}
			} else if err == nil {
				t.Fatalf("expected an error, got: %v", d)
			}
		})
	}
}

func TestNewHour(t *testing.T) {
	testCases := []struct {
		Name  string
		Year  int
		Month time.Month
		Day   int
		Hour  int
		Valid bool
	}{
		{Name: "valid", Year: 2026, Month: time.August, Day: 5, Hour: 13, Valid: true},
		{Name: "midnight", Year: 2026, Month: time.August, Day: 5, Hour: 0, Valid: true},
		{Name: "last_hour", Year: 2026, Month: time.August, Day: 5, Hour: 23, Valid: true},
		{Name: "hour_too_large", Year: 2026, Month: time.August, Day: 5, Hour: 24, Valid: false},
		{Name: "hour_negative", Year: 2026, Month: time.August, Day: 5, Hour: -1, Valid: false},
		{Name: "invalid_day", Year: 2026, Month: time.February, Day: 29, Hour: 0, Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			h, err := NewHour(testCase.Year, testCase.Month, testCase.Day, testCase.Hour)
			if testCase.Valid {
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				if h.Year() != testCase.Year || h.Month() != testCase.Month || h.Day() != testCase.Day || h.Hour() != testCase.Hour {
					t.Fatalf("error: %v", h)
				}
			} else if err == nil {
				t.Fatalf("expected an error, got: %v", h)
			}
		})
	}
}
