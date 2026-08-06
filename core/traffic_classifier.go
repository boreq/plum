package core

import (
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/boreq/plum/parser"
)

const (
	MaliciousRequestThreshold = 10

	trafficWindow   = 7 * 24 * time.Hour
	bucketKeyFormat = "2006-01-02"
)

var automatedUserAgentNames = map[string]struct{}{
	"curl":              {},
	"prometheus":        {},
	"googlebot":         {},
	"googlebot-image":   {},
	"googlebot-video":   {},
	"googlebot-news":    {},
	"duckduckbot":       {},
	"duckduckbot-https": {},
	"duckassistbot":     {},
	"sogou web spider":  {},
	"moooodotfarm":      {},

	"mastodon":      {},
	"misskey":       {},
	"blogtrottr":    {},
	"onlyread":      {},
	"reeder":        {},
	"tiny tiny rss": {},
}

var automatedUserAgentMarkers = []string{
	"feed",
	"rss",
}

var scanUris = map[string]struct{}{
	"/_ignition/execute-solution": {},
	"/actuator/env":               {},
	"/api/v1/login":               {},
	"/api/v1/auto_login":          {},
	"/appsettings.json":           {},
	"/aws-ses.json":               {},
	"/aws-credentials.txt":        {},
	"/gcp-credentials.json":       {},
	"/config.json":                {},
	"/debug/default/view":         {},
	"/phpmyadmin":                 {},
	"/server-status":              {},
	"/telescope/requests":         {},

	"/wp-config.php":                                    {},
	"/blog/wp-config.php":                               {},
	"/wp-admin/install.php":                             {},
	"/wp-admin/install.php?step=1":                      {},
	"/wp-content/plugins/hellopress/wp_filemanager.php": {},
	"/wp-json/batch/v1":                                 {},
	"/?rest_route=/batch/v1":                            {},

	"/adminfuns.php":                 {},
	"/adminner.php":                  {},
	"/this_is_a_new_hello_world.php": {},
}

var scanRules = []func(r scanRequest) bool{
	matchesKnownScanUri,
	attemptsPathTraversal,
	requestsHiddenFile,
	requestsWordPressProbe,
	requestsNumberedScript,
	requestsSensitiveFile,
	containsInjection,
}

var (
	sensitiveFileExtensions = []string{
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

	wordPressIncludeDirectories = []string{
		"wp-includes",
	}

	wordPressProbeFiles = []string{
		"wlwmanifest.xml",
		"wp-config.php",
		"xmlrpc.php",
	}

	injectionMarkers = []string{
		"$(",
		"${",
		"`",
		"cmd=",
		"<",
		">",
	}
)

type requestCounters struct {
	Malicious int
	Automated int
	Other     int
}

type TrafficClassifier struct {
	ips map[string]map[string]*requestCounters
}

func NewTrafficClassifier() *TrafficClassifier {
	return &TrafficClassifier{
		ips: make(map[string]map[string]*requestCounters),
	}
}

func (c *TrafficClassifier) Classify(entry *parser.Entry) Category {
	category := classifyUserAgent(entry.UserAgent)

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

func classifyUserAgent(userAgent string) Category {
	name := strings.ToLower(UserAgentName(userAgent))

	if _, ok := automatedUserAgentNames[name]; ok {
		return CategoryAutomated
	}

	for _, marker := range automatedUserAgentMarkers {
		if strings.Contains(name, marker) {
			return CategoryAutomated
		}
	}

	return CategoryUnclassified
}

type scanRequest struct {
	Uri      string
	Path     string
	Query    string
	Segments []string
	FileName string
}

func newScanRequest(uri string) scanRequest {
	uri = strings.ToLower(uri)
	if unescaped, err := url.PathUnescape(uri); err == nil {
		uri = unescaped
	}
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

	return scanRequest{
		Uri:      uri,
		Path:     path,
		Query:    query,
		Segments: segments,
		FileName: fileName,
	}
}

func isScanRequest(uri string) bool {
	r := newScanRequest(uri)

	for _, rule := range scanRules {
		if rule(r) {
			return true
		}
	}

	return false
}

func matchesKnownScanUri(r scanRequest) bool {
	_, ok := scanUris[r.Uri]
	return ok
}

func attemptsPathTraversal(r scanRequest) bool {
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

func requestsHiddenFile(r scanRequest) bool {
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

func requestsWordPressProbe(r scanRequest) bool {
	for _, file := range wordPressProbeFiles {
		if r.FileName == file {
			return true
		}
	}

	if !strings.HasSuffix(r.FileName, ".php") {
		return false
	}

	for _, segment := range r.Segments {
		for _, directory := range wordPressIncludeDirectories {
			if segment == directory {
				return true
			}
		}
	}

	return false
}

func requestsNumberedScript(r scanRequest) bool {
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

func requestsSensitiveFile(r scanRequest) bool {
	for _, extension := range sensitiveFileExtensions {
		if strings.HasSuffix(r.FileName, extension) {
			return true
		}
	}

	return false
}

func containsInjection(r scanRequest) bool {
	for _, marker := range injectionMarkers {
		if strings.Contains(r.Uri, marker) {
			return true
		}
	}

	return false
}

func UserAgentName(userAgent string) string {
	name, _, _ := strings.Cut(userAgent, "/")
	name, _, _ = strings.Cut(name, " (")
	name = strings.TrimSpace(name)

	if !isSimpleUserAgentName(name) {
		return userAgent
	}

	return name
}

func isSimpleUserAgentName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && r != ' ' && r != '-' {
			return false
		}
	}

	return true
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
