package core

import (
	"testing"
	"time"
)

func TestEstimateBrowserReleaseAge(t *testing.T) {
	testCases := []struct {
		Name      string
		UserAgent string
		Release   string
		Ok        bool
	}{
		{
			Name:      "chrome",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			Release:   "2024-10-15",
			Ok:        true,
		},
		{
			Name:      "recent_chrome",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36",
			Release:   "2026-07-28",
			Ok:        true,
		},
		{
			Name:      "old_chrome",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
			Release:   "2019-01-29",
			Ok:        true,
		},
		{
			Name:      "chrome_on_ios",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.6099.119 Mobile/15E148 Safari/604.1",
			Release:   "2023-12-05",
			Ok:        true,
		},
		{
			Name:      "chromium",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chromium/130.0.0.0 Chrome/130.0.0.0 Safari/537.36",
			Release:   "2024-10-15",
			Ok:        true,
		},
		{
			Name:      "release_which_was_never_shipped",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.0.0 Safari/537.36",
			Release:   "2020-04-07",
			Ok:        true,
		},
		{
			Name:      "chrome_based_browser",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 OPR/101.0.0.0",
			Release:   "2023-07-18",
			Ok:        true,
		},
		{
			Name:      "android_webview",
			UserAgent: "Mozilla/5.0 (Linux; Android 8.1.0; M1813 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.0.0 Mobile Safari/537.36",
			Release:   "2024-10-15",
			Ok:        true,
		},
		{
			Name:      "firefox",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0",
			Release:   "2024-09-03",
			Ok:        true,
		},
		{
			Name:      "old_firefox",
			UserAgent: "Mozilla/5.0 (Maemo; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1",
			Release:   "2012-01-31",
			Ok:        true,
		},
		{
			Name:      "firefox_on_ios",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/130.0 Mobile/15E148 Safari/605.1.15",
			Release:   "2024-09-03",
			Ok:        true,
		},
		{
			Name:      "safari",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			Release:   "2023-09-18",
			Ok:        true,
		},
		{
			Name:      "safari_on_ios",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.6 Mobile/15E148 Safari/604.1",
			Release:   "2024-09-16",
			Ok:        true,
		},
		{
			Name:      "safari_without_a_version",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Safari/604.1",
			Ok:        false,
		},
		{
			Name:      "chrome_newer_than_the_known_releases",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/900.0.0.0 Safari/537.36",
			Ok:        false,
		},
		{
			Name:      "firefox_newer_than_the_known_releases",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:900.0) Gecko/20100101 Firefox/900.0",
			Ok:        false,
		},
		{
			Name:      "safari_newer_than_the_known_releases",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/90.0 Safari/605.1.15",
			Ok:        false,
		},
		{
			Name:      "not_a_browser",
			UserAgent: "curl/8.5.0",
			Ok:        false,
		},
		{
			Name:      "no_version",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/ Safari/537.36",
			Ok:        false,
		},
		{
			Name:      "empty_user_agent",
			UserAgent: "",
			Ok:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			age, ok := EstimateBrowserReleaseAge(testCase.UserAgent, classifyReferenceTime)
			if ok != testCase.Ok {
				t.Fatalf("got %v, want %v", ok, testCase.Ok)
			}

			if !ok {
				return
			}

			if release := classifyReferenceTime.Add(-age).Format("2006-01-02"); release != testCase.Release {
				t.Errorf("got %q, want %q", release, testCase.Release)
			}
		})
	}
}

func TestEstimateBrowserReleaseAgeDependsOnTheGivenTime(t *testing.T) {
	const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"

	age, ok := EstimateBrowserReleaseAge(userAgent, releaseDate(2024, time.October, 22))
	if !ok {
		t.Fatal("the browser was not recognized")
	}

	if age != 7*24*time.Hour {
		t.Errorf("got %v, want %v", age, 7*24*time.Hour)
	}

	age, ok = EstimateBrowserReleaseAge(userAgent, releaseDate(2026, time.October, 15))
	if !ok {
		t.Fatal("the browser was not recognized")
	}

	if age != 730*24*time.Hour {
		t.Errorf("got %v, want %v", age, 730*24*time.Hour)
	}
}

func TestClassifyDetectsOutdatedBrowsers(t *testing.T) {
	testCases := []struct {
		Name      string
		UserAgent string
		Category  Category
	}{
		{
			Name:      "current_chrome",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "chrome_released_a_year_ago",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "chrome_released_three_years_ago",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
			Category:  CategoryPossiblyAutomated,
		},
		{
			Name:      "current_firefox",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:152.0) Gecko/20100101 Firefox/152.0",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "firefox_released_four_years_ago",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			Name:      "current_safari",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/19.0 Safari/605.1.15",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "safari_released_five_years_ago",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
			Category:  CategoryPossiblyAutomated,
		},
		{
			Name:      "unknown_browser_version",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Safari/604.1",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "newer_than_the_known_releases",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/900.0.0.0 Safari/537.36",
			Category:  CategoryUnclassified,
		},
		{
			Name:      "outdated_but_recognized_as_automated",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/112.0.0.0 Safari/537.36",
			Category:  CategoryAutomated,
		},
		{
			Name:      "outdated_but_recognized_as_malicious",
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 sqlmap/1.9",
			Category:  CategoryMalicious,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if category := classifyUserAgent(testCase.UserAgent, classifyReferenceTime); category != testCase.Category {
				t.Errorf("got %q, want %q", category, testCase.Category)
			}
		})
	}
}

func TestBrowserReleaseDataIsOutOfDate(t *testing.T) {
	newest := newestKnownBrowserRelease()

	if _, outOfDate := BrowserReleaseDataIsOutOfDate(newest.Add(browserReleaseDataMaxAge)); outOfDate {
		t.Error("the data should still be considered up to date")
	}

	returnedNewest, outOfDate := BrowserReleaseDataIsOutOfDate(newest.Add(browserReleaseDataMaxAge).Add(24 * time.Hour))
	if !outOfDate {
		t.Error("the data should be considered out of date")
	}

	if !returnedNewest.Equal(newest) {
		t.Errorf("got %v, want %v", returnedNewest, newest)
	}
}

func TestNewestKnownBrowserReleaseIsTheOldestOfTheNewestReleases(t *testing.T) {
	newest := newestKnownBrowserRelease()

	for _, history := range browserReleaseHistories {
		if last := history.Releases[len(history.Releases)-1].Date; last.Before(newest) {
			t.Errorf("got %v, want at most %v", newest, last)
		}
	}
}

func TestBrowserReleasesAreSortedAndPlausible(t *testing.T) {
	for _, history := range browserReleaseHistories {
		if len(history.Releases) < 2 {
			t.Fatalf("not enough releases: %v", history.Markers)
		}

		previous := history.Releases[0]
		for _, release := range history.Releases[1:] {
			if release.Major <= previous.Major {
				t.Errorf("%v: major %d is not after %d", history.Markers, release.Major, previous.Major)
			}

			if !release.Date.After(previous.Date) {
				t.Errorf("%v: release %d is not after release %d", history.Markers, release.Major, previous.Major)
			}

			previous = release
		}
	}
}
