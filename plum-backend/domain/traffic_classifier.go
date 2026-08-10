package domain

import (
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/boreq/plum/plum-backend/domain/request"
	"github.com/boreq/plum/plum-backend/logging"
)

const (
	MaliciousRequestThreshold = 5
	TrafficWindow             = 30 * 24 * time.Hour

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
	"claude-user":       {},
	"python":            {},
	"python-requests":   {},
	"python-urllib":     {},
	"go-http-client":    {},
	"scrapy":            {},
	"apache-httpclient": {},
	"java":              {},
	"bun":               {},
	"domaindrift":       {},
	"axios":             {},
	"dart":              {},
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

	"turnitin":                 {},
	"navcrawl":                 {},
	"crusader-worker":          {},
	"mwmbl":                    {},
	"stultussearchengine":      {},
	"aiwebindex":               {},
	"synapx-datacomp-source":   {},
	"t2i-data-curation":        {},
	"moving-to-bins":           {},
	"nx-zdhub":                 {},
	"owler":                    {},
	"frontpagedomainpipeline":  {},
	"gdelt-media-org-research": {},
	"hnblogarchive":            {},
	"alittle client":           {},
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
	"scaner",
	"scanner",
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
	"~",
}

var maliciousFileNames = map[string]struct{}{
	"appsettings.json":             {},
	"appsettings.development.json": {},
	"appsettings.production.json":  {},
	"appsettings.staging.json":     {},

	"application-dev.properties": {},

	"asset-manifest.json":  {},
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
	"credentials.json":     {},
	"gcp-credentials.json": {},
	"gcp.json":             {},
	"gradle.properties":    {},
	"mail.config.js":       {},
	"manifest.json":        {},
	"secrets.json":         {},
	"webpack-stats.json":   {},

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

var maliciousRules = []func(r requestURI) bool{
	requestsMaliciousUri,
	requestsMaliciousDirectory,
	requestsMaliciousFile,
	requestsHiddenFile,
	attemptsPathTraversal,
	requestsNumberedPhpScript,
	attemptsInjection,
	obfuscatesPath,
	requestsPunctuationPrefixedSegment,
}

type TrafficClassifier struct {
	log logging.Logger
}

func NewTrafficClassifier() *TrafficClassifier {
	log := logging.New("core/traffic_classifier")

	if newest, outOfDate := BrowserReleaseDataIsOutOfDate(time.Now()); outOfDate {
		log.Warn("browser release data has to be updated, outdated browsers are not detected reliably", "newestKnownRelease", newest.Format(bucketKeyFormat))
	}

	return &TrafficClassifier{
		log: log,
	}
}

func (c *TrafficClassifier) Classify(uri request.Uri, userAgent request.UserAgent, t time.Time) Category {
	category := c.classifyUserAgent(userAgent.String(), t)

	if c.isScanRequest(uri.String()) {
		category = CategoryMalicious
	}

	return category
}

func (c *TrafficClassifier) classifyUserAgent(rawUserAgent string, t time.Time) Category {
	raw := strings.ToLower(rawUserAgent)

	if c.isMissingUserAgent(raw) {
		return CategoryMalicious
	}

	products := c.userAgentProducts(raw)

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

	if c.pretendsToBeBrowser(raw) && !c.containsBrowserProduct(products) {
		return CategoryPossiblyAutomated
	}

	if isOutdatedBrowser(raw, t) {
		return CategoryPossiblyAutomated
	}

	if c.isAllLowercaseUserAgent(rawUserAgent) {
		return CategoryPossiblyAutomated
	}

	if c.isQuotedUserAgent(rawUserAgent) {
		return CategoryPossiblyAutomated
	}

	return CategoryUnclassified
}

func (c *TrafficClassifier) isMissingUserAgent(raw string) bool {
	raw = strings.TrimSpace(raw)
	return raw == "" || raw == "-"
}

func (c *TrafficClassifier) isAllLowercaseUserAgent(rawUserAgent string) bool {
	return rawUserAgent == strings.ToLower(rawUserAgent)
}

func (c *TrafficClassifier) isQuotedUserAgent(rawUserAgent string) bool {
	for _, q := range []string{`"`, `'`, `\x22`, `\x27`} {
		if strings.HasPrefix(rawUserAgent, q) || strings.HasSuffix(rawUserAgent, q) {
			return true
		}
	}
	return false
}

func (c *TrafficClassifier) userAgentProducts(raw string) []string {
	var rv []string
	seen := make(map[string]struct{})

	add := func(s string) {
		product := c.userAgentProduct(s)
		if product == "" {
			return
		}
		if _, ok := seen[product]; ok {
			return
		}
		seen[product] = struct{}{}
		rv = append(rv, product)
	}

	for _, segment := range c.userAgentSegments(raw) {
		add(segment)

		for _, field := range strings.Fields(segment) {
			add(field)
		}
	}

	return rv
}

func (c *TrafficClassifier) userAgentSegments(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == '(' || r == ')' || r == ';' || r == ','
	})
}

func (c *TrafficClassifier) userAgentProduct(s string) string {
	name, _, _ := strings.Cut(s, "/")
	return strings.Trim(name, " +")
}

func (c *TrafficClassifier) pretendsToBeBrowser(raw string) bool {
	return strings.HasPrefix(raw, "mozilla/")
}

func (c *TrafficClassifier) containsBrowserProduct(products []string) bool {
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

type requestURI struct {
	RawUri   string
	Uri      string
	Path     string
	Query    string
	Segments []string
	FileName string
}

func newRequestURI(rawUri string) requestURI {
	uri := unescapeUri(rawUri)
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

	return requestURI{
		RawUri:   rawUri,
		Uri:      uri,
		Path:     path,
		Query:    query,
		Segments: segments,
		FileName: fileName,
	}
}

func (c *TrafficClassifier) isScanRequest(uri string) bool {
	r := newRequestURI(uri)

	for _, rule := range maliciousRules {
		if rule(r) {
			return true
		}
	}

	return false
}

func requestsMaliciousUri(r requestURI) bool {
	_, ok := maliciousUris[r.Uri]
	return ok
}

func attemptsPathTraversal(r requestURI) bool {
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

func requestsHiddenFile(r requestURI) bool {
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

func requestsMaliciousDirectory(r requestURI) bool {
	for _, segment := range r.Segments {
		for _, directory := range maliciousDirectories {
			if segment == directory {
				return true
			}
		}
	}

	return false
}

func requestsNumberedPhpScript(r requestURI) bool {
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

func requestsMaliciousFile(r requestURI) bool {
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

func requestsPunctuationPrefixedSegment(r requestURI) bool {
	for _, segment := range r.Segments {
		c := rune(segment[0])

		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			continue
		}

		if strings.ContainsRune("._-~@", c) {
			continue
		}

		return true
	}

	return false
}

func obfuscatesPath(r requestURI) bool {
	uri := strings.ToLower(r.RawUri)

	for i := 0; i < maxUnescapePasses; i++ {
		path, _, _ := strings.Cut(uri, "?")
		if encodesAlphanumericCharacters(path) {
			return true
		}

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

	return false
}

func encodesAlphanumericCharacters(path string) bool {
	for i := 0; i+2 < len(path); i++ {
		if path[i] != '%' {
			continue
		}

		b, err := strconv.ParseUint(path[i+1:i+3], 16, 8)
		if err != nil {
			continue
		}

		c := byte(b)
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			return true
		}
	}

	return false
}

func attemptsInjection(r requestURI) bool {
	for _, marker := range injectionMarkers {
		if strings.Contains(r.Uri, marker) {
			return true
		}
	}

	return false
}
