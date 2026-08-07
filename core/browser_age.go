package core

import (
	"strconv"
	"strings"
	"time"
)

const (
	maxBrowserAge            = 2 * 365 * 24 * time.Hour
	browserReleaseDataMaxAge = 400 * 24 * time.Hour
)

type browserRelease struct {
	Major int
	Date  time.Time
}

type browserReleaseHistory struct {
	Markers  []string
	Requires string
	Releases []browserRelease
}

var browserReleaseHistories = []browserReleaseHistory{
	{
		Markers:  []string{"chromium/", "crios/", "chrome/"},
		Releases: chromeReleases,
	},
	{
		Markers:  []string{"fxios/", "firefox/"},
		Releases: firefoxReleases,
	},
	{
		Markers:  []string{"version/"},
		Requires: "safari",
		Releases: safariReleases,
	},
}

func releaseDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func EstimateBrowserReleaseAge(rawUserAgent string, t time.Time) (time.Duration, bool) {
	raw := strings.ToLower(rawUserAgent)

	for _, history := range browserReleaseHistories {
		if history.Requires != "" && !strings.Contains(raw, history.Requires) {
			continue
		}

		major, ok := majorVersionAfterMarkers(raw, history.Markers)
		if !ok {
			continue
		}

		date, ok := history.ReleaseDate(major)
		if !ok {
			return 0, false
		}

		return t.UTC().Sub(date), true
	}

	return 0, false
}

func BrowserReleaseDataIsOutOfDate(t time.Time) (time.Time, bool) {
	newest := newestKnownBrowserRelease()
	return newest, t.UTC().Sub(newest) > browserReleaseDataMaxAge
}

func isOutdatedBrowser(rawUserAgent string, t time.Time) bool {
	age, ok := EstimateBrowserReleaseAge(rawUserAgent, t)
	if !ok {
		return false
	}

	return age > maxBrowserAge
}

func newestKnownBrowserRelease() time.Time {
	var newest time.Time

	for _, history := range browserReleaseHistories {
		date := history.Releases[len(history.Releases)-1].Date
		if newest.IsZero() || date.Before(newest) {
			newest = date
		}
	}

	return newest
}

func (h browserReleaseHistory) ReleaseDate(major int) (time.Time, bool) {
	if major > h.Releases[len(h.Releases)-1].Major {
		return time.Time{}, false
	}

	date := h.Releases[0].Date
	for _, release := range h.Releases {
		if release.Major > major {
			break
		}
		date = release.Date
	}

	return date, true
}

func majorVersionAfterMarkers(raw string, markers []string) (int, bool) {
	for _, marker := range markers {
		index := strings.Index(raw, marker)
		if index < 0 {
			continue
		}

		if major, ok := parseMajorVersion(raw[index+len(marker):]); ok {
			return major, true
		}
	}

	return 0, false
}

func parseMajorVersion(s string) (int, bool) {
	end := 0
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}

	if end == 0 {
		return 0, false
	}

	major, err := strconv.Atoi(s[:end])
	if err != nil || major <= 0 {
		return 0, false
	}

	return major, true
}
