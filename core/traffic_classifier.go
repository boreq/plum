package core

import (
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/boreq/plum/logging"
	"github.com/boreq/plum/parser"
)

const (
	MaliciousRequestThreshold = 5
	trafficWindow             = 14 * 24 * time.Hour

	bucketKeyFormat   = "2006-01-02"
	maxUnescapePasses = 5
)

var automatedUserAgentNames = map[string]struct{}{
	"curl":              {},
	"prometheus":        {},
	"googlebot":         {},
	"googlebot-image":   {},
	"googlebot-video":   {},
	"googlebot-news":    {},
	"googleother":       {},
	"duckduckbot":       {},
	"duckduckbot-https": {},
	"duckassistbot":     {},
	"sogou web spider":  {},
	"moooodotfarm":      {},
	"ccbot":             {},
	"python":            {},
	"python-requests":   {},
	"python-urllib":     {},
	"go-http-client":    {},
	"scrapy":            {},
	"apache-httpclient": {},
	"bun":               {},
	"domaindrift":       {},
	"axios":             {},
	"node":              {},
	"http.rb":           {},

	"facebookexternalhit": {},
	"meta-externalagent":  {},
	"catstodon":           {},
	"ditto":               {},
	"damus":               {},
	"primal":              {},
	"amethyst":            {},
	"nos":                 {},
	"netnewswire":         {},
	"networkingextension": {},
	"roku":                {},

	"turnitin":               {},
	"navcrawl":               {},
	"crusader-worker":        {},
	"mwmbl":                  {},
	"stultussearchengine":    {},
	"aiwebindex":             {},
	"synapx-datacomp-source": {},
	"t2i-data-curation":      {},
	"moving-to-bins":         {},
	"nx-zdhub":               {},
	"alittle client":         {},
	"wpmu dev broken link checker local engine": {},

	"mastodon":      {},
	"misskey":       {},
	"blogtrottr":    {},
	"onlyread":      {},
	"reeder":        {},
	"tiny tiny rss": {},
	"newsboat":      {},

	"cms-checker": {},
}

var automatedUserAgentMarkers = []string{
	"feed",
	"rss",
	"newsblur",
	"miniflux",
	"bazqux",
	"palo alto networks",
	"com.apple.webkit.networking",
	"headlesschrome",
	"google web preview",
	"mastodon",
	"marginalia",
	"terracotta",
	"amazon-quick",
	"bot",
	"crawler",
	"trawler",
	"spider",
}

var browserProducts = map[string]struct{}{
	"safari":  {},
	"firefox": {},
	"chrome":  {},
}

var maliciousUserAgentNames = map[string]struct{}{
	"bytespider":        {},
	"vuln_scanner":      {},
	"zgrab":             {},
	"masscan":           {},
	"nuclei":            {},
	"sqlmap":            {},
	"nikto":             {},
	"wpscan":            {},
	"tlm-audit-scanner": {},
}

var maliciousUris = map[string]struct{}{
	"/_ignition/execute-solution": {},
	"/actuator/env":               {},
	"/api/v1/login":               {},
	"/api/v1/auto_login":          {},
	"/debug/default/view":         {},
	"/phpmyadmin":                 {},
	"/server-status":              {},
	"/telescope/requests":         {},

	"/wp-admin/admin-ajax.php":                          {},
	"/wp-admin/install.php":                             {},
	"/wp-admin/install.php?step=1":                      {},
	"/wp-content/plugins/hellopress/wp_filemanager.php": {},
	"/wp-content/plugins/index.php":                     {},
	"/wp-json/batch/v1":                                 {},
	"/?rest_route=/batch/v1":                            {},
}

var maliciousFileExtensions = []string{
	".env",
	".conf",
	".ini",
	".yml",
	".yaml",
	".sql",
	".bak",
	".old",
	".swp",
	".pem",
	".key",
	".log",
}

var maliciousFileNames = map[string]struct{}{
	"appsettings.json":     {},
	"aws-credentials.txt":  {},
	"aws-ses.json":         {},
	"backup.zip":           {},
	"config.inc.php.dist":  {},
	"dockerfile":           {},
	"next.config.js":       {},
	"next.config.mjs":      {},
	"next.config.ts":       {},
	"nuxt.config.js":       {},
	"nuxt.config.mjs":      {},
	"nuxt.config.ts":       {},
	"config.json":          {},
	"config.php":           {},
	"gcp-credentials.json": {},
	"secrets.json":         {},

	"wlwmanifest.xml": {},
	"wp-config.php":   {},
	"wp-login.php":    {},
	"xmlrpc.php":      {},

	"abcd.php":                      {},
	"admin.php":                     {},
	"adminfuns.php":                 {},
	"adminner.php":                  {},
	"akcc.php":                      {},
	"bless.php":                     {},
	"blurbs.php":                    {},
	"bolt.php":                      {},
	"chosen.php":                    {},
	"cjfuns.php":                    {},
	"class-t.api.php":               {},
	"classwithtostring.php":         {},
	"cord.php":                      {},
	"dex.php":                       {},
	"flower.php":                    {},
	"gifclass.php":                  {},
	"php-info":                      {},
	"php-info.php":                  {},
	"php_info":                      {},
	"php_info.php":                  {},
	"phpinfo":                       {},
	"phpinfo.php":                   {},
	"phpinfo.php3":                  {},
	"shelp.php":                     {},
	"this_is_a_new_hello_world.php": {},
}

var maliciousDirectories = []string{
	"wp-admin",
	"wp-content",
	"wp-includes",
	"wp-json",
}

var injectionMarkers = []string{
	"$(",
	"${",
	"`",
	"cmd=",
	"<",
	">",
}

var maliciousRules = []func(r request) bool{
	requestsMaliciousUri,
	requestsMaliciousDirectory,
	requestsMaliciousFile,
	requestsHiddenFile,
	attemptsPathTraversal,
	requestsNumberedPhpScript,
	attemptsInjection,
}

type requestCounters struct {
	Malicious int
	Automated int
	Other     int
}

type TrafficClassifier struct {
	ips map[string]map[string]*requestCounters
	log logging.Logger
}

func NewTrafficClassifier() *TrafficClassifier {
	log := logging.New("core/traffic_classifier")

	if newest, outOfDate := BrowserReleaseDataIsOutOfDate(time.Now()); outOfDate {
		log.Warn("browser release data has to be updated, outdated browsers are not detected reliably", "newestKnownRelease", newest.Format(bucketKeyFormat))
	}

	return &TrafficClassifier{
		ips: make(map[string]map[string]*requestCounters),
		log: log,
	}
}

func (c *TrafficClassifier) Classify(entry *parser.Entry) Category {
	category := classifyUserAgent(entry.UserAgent, entry.Time)

	if isScanRequest(entry.HttpRequestURI) {
		category = CategoryMalicious
	}

	c.insert(entry.RemoteAddress, entry.Time, category)

	if category == CategoryMalicious {
		return category
	}

	if c.maliciousRequests(entry.RemoteAddress, entry.Time) > MaliciousRequestThreshold {
		return CategoryMalicious
	}

	return category
}

func (c *TrafficClassifier) RemoveOldData(now time.Time) {
	for remoteAddress, buckets := range c.ips {
		removeOldBuckets(buckets, now)
		if len(buckets) == 0 {
			delete(c.ips, remoteAddress)
		}
	}
}

func (c *TrafficClassifier) insert(remoteAddress string, t time.Time, category Category) {
	buckets, ok := c.ips[remoteAddress]
	if !ok {
		buckets = make(map[string]*requestCounters)
		c.ips[remoteAddress] = buckets
	}

	key := bucketKey(t)
	counters, ok := buckets[key]
	if !ok {
		counters = &requestCounters{}
		buckets[key] = counters
	}

	switch category {
	case CategoryMalicious:
		counters.Malicious++
	case CategoryAutomated:
		counters.Automated++
	default:
		counters.Other++
	}

	removeOldBuckets(buckets, t)
}

func (c *TrafficClassifier) maliciousRequests(remoteAddress string, t time.Time) int {
	var malicious int

	for key, counters := range c.ips[remoteAddress] {
		if bucketWithinWindow(key, t) {
			malicious += counters.Malicious
		}
	}

	return malicious
}

func classifyUserAgent(rawUserAgent string, t time.Time) Category {
	raw := strings.ToLower(rawUserAgent)

	if isMissingUserAgent(raw) {
		return CategoryMalicious
	}

	products := userAgentProducts(raw)

	for _, product := range products {
		if _, ok := maliciousUserAgentNames[product]; ok {
			return CategoryMalicious
		}
	}

	for _, product := range products {
		if _, ok := automatedUserAgentNames[product]; ok {
			return CategoryAutomated
		}
	}

	for _, marker := range automatedUserAgentMarkers {
		if strings.Contains(raw, marker) {
			return CategoryAutomated
		}
	}

	if pretendsToBeBrowser(raw) && !containsBrowserProduct(products) {
		return CategoryPossiblyAutomated
	}

	if isOutdatedBrowser(raw, t) {
		return CategoryPossiblyAutomated
	}

	if isAllLowercaseUserAgent(rawUserAgent) {
		return CategoryPossiblyAutomated
	}

	if isQuotedUserAgent(rawUserAgent) {
		return CategoryPossiblyAutomated
	}

	return CategoryUnclassified
}

func isMissingUserAgent(raw string) bool {
	raw = strings.TrimSpace(raw)
	return raw == "" || raw == "-"
}

func isAllLowercaseUserAgent(rawUserAgent string) bool {
	return rawUserAgent == strings.ToLower(rawUserAgent)
}

func isQuotedUserAgent(rawUserAgent string) bool {
	for _, q := range []string{`"`, `'`, `\x22`, `\x27`} {
		if strings.HasPrefix(rawUserAgent, q) || strings.HasSuffix(rawUserAgent, q) {
			return true
		}
	}
	return false
}

func userAgentProducts(raw string) []string {
	var rv []string
	seen := make(map[string]struct{})

	add := func(s string) {
		product := userAgentProduct(s)
		if product == "" {
			return
		}
		if _, ok := seen[product]; ok {
			return
		}
		seen[product] = struct{}{}
		rv = append(rv, product)
	}

	for _, segment := range userAgentSegments(raw) {
		add(segment)

		for _, field := range strings.Fields(segment) {
			add(field)
		}
	}

	return rv
}

func userAgentSegments(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == '(' || r == ')' || r == ';' || r == ','
	})
}

func userAgentProduct(s string) string {
	name, _, _ := strings.Cut(s, "/")
	return strings.Trim(name, " +")
}

func pretendsToBeBrowser(raw string) bool {
	return strings.HasPrefix(raw, "mozilla/")
}

func containsBrowserProduct(products []string) bool {
	for _, product := range products {
		if _, ok := browserProducts[product]; ok {
			return true
		}
	}

	return false
}

func unescapeUri(uri string) string {
	uri = strings.ToLower(uri)

	for i := 0; i < maxUnescapePasses; i++ {
		unescaped, err := url.PathUnescape(uri)
		if err != nil {
			break
		}

		unescaped = strings.ToLower(unescaped)
		if unescaped == uri {
			break
		}

		uri = unescaped
	}

	return uri
}

type request struct {
	Uri      string
	Path     string
	Query    string
	Segments []string
	FileName string
}

func newRequest(uri string) request {
	uri = unescapeUri(uri)
	path, query, _ := strings.Cut(uri, "?")

	var segments []string
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	var fileName string
	if len(segments) > 0 {
		fileName = segments[len(segments)-1]
	}

	return request{
		Uri:      uri,
		Path:     path,
		Query:    query,
		Segments: segments,
		FileName: fileName,
	}
}

func isScanRequest(uri string) bool {
	r := newRequest(uri)

	for _, rule := range maliciousRules {
		if rule(r) {
			return true
		}
	}

	return false
}

func requestsMaliciousUri(r request) bool {
	_, ok := maliciousUris[r.Uri]
	return ok
}

func attemptsPathTraversal(r request) bool {
	if strings.HasPrefix(r.Uri, "//") {
		return true
	}

	for _, segment := range r.Segments {
		if segment == ".." {
			return true
		}
	}

	return false
}

func requestsHiddenFile(r request) bool {
	for _, segment := range r.Segments {
		if !strings.HasPrefix(segment, ".") {
			continue
		}

		if segment == "." || segment == ".." || strings.HasPrefix(segment, ".well-known") {
			continue
		}

		return true
	}

	return false
}

func requestsMaliciousDirectory(r request) bool {
	for _, segment := range r.Segments {
		for _, directory := range maliciousDirectories {
			if segment == directory {
				return true
			}
		}
	}

	return false
}

func requestsNumberedPhpScript(r request) bool {
	name, extension, found := strings.Cut(r.FileName, ".")
	if !found || extension != "php" || name == "" {
		return false
	}

	for _, c := range name {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

func requestsMaliciousFile(r request) bool {
	if _, ok := maliciousFileNames[r.FileName]; ok {
		return true
	}

	for _, extension := range maliciousFileExtensions {
		if strings.HasSuffix(r.FileName, extension) {
			return true
		}
	}

	return false
}

func attemptsInjection(r request) bool {
	for _, marker := range injectionMarkers {
		if strings.Contains(r.Uri, marker) {
			return true
		}
	}

	return false
}

func bucketKey(t time.Time) string {
	return t.UTC().Format(bucketKeyFormat)
}

func bucketWithinWindow(key string, t time.Time) bool {
	bucket, err := parseBucketKey(key)
	if err != nil {
		return false
	}
	return !bucket.Before(windowStart(t)) && !bucket.After(t.UTC())
}

func removeOldBuckets(buckets map[string]*requestCounters, t time.Time) {
	for key := range buckets {
		bucket, err := parseBucketKey(key)
		if err != nil || bucket.Before(windowStart(t)) {
			delete(buckets, key)
		}
	}
}

func parseBucketKey(key string) (time.Time, error) {
	return time.ParseInLocation(bucketKeyFormat, key, time.UTC)
}

func windowStart(t time.Time) time.Time {
	return t.UTC().Add(-trafficWindow)
}
