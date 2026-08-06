package core

import (
	"testing"
	"time"

	"github.com/boreq/plum/parser"
)

func classifierEntry(remoteAddress, userAgent, uri string, t time.Time) *parser.Entry {
	return &parser.Entry{
		RemoteAddress:  remoteAddress,
		UserAgent:      userAgent,
		HttpRequestURI: uri,
		Time:           t,
	}
}

func TestClassifyUsesTheUserAgent(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	if category := c.Classify(classifierEntry("1.1.1.1", "curl/7.64.0", "/", now)); category != CategoryAutomated {
		t.Errorf("got %q, want %q", category, CategoryAutomated)
	}

	if category := c.Classify(classifierEntry("2.2.2.2", "Mozilla/5.0", "/", now)); category != CategoryUnclassified {
		t.Errorf("got %q, want %q", category, CategoryUnclassified)
	}
}

func TestClassifyDetectsScans(t *testing.T) {
	testCases := []struct {
		Uri      string
		Category Category
	}{
		{Uri: "/.env", Category: CategoryMalicious},
		{Uri: "/.git/config", Category: CategoryMalicious},
		{Uri: "/", Category: CategoryUnclassified},
		{Uri: "/index.php", Category: CategoryUnclassified},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Uri, func(t *testing.T) {
			c := NewTrafficClassifier()

			category := c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", testCase.Uri, time.Now().UTC()))
			if category != testCase.Category {
				t.Errorf("got %q, want %q", category, testCase.Category)
			}
		})
	}
}

func TestIsScanRequest(t *testing.T) {
	testCases := []struct {
		Uri  string
		Scan bool
	}{
		{Uri: "/.env", Scan: true},
		{Uri: "/app/.env", Scan: true},
		{Uri: "/.env.dev", Scan: true},
		{Uri: "/.env.development", Scan: true},
		{Uri: "/.env.staging", Scan: true},
		{Uri: "/.env.docker", Scan: true},
		{Uri: "/.env.local", Scan: true},
		{Uri: "/.env.production", Scan: true},
		{Uri: "/laravel/.env", Scan: true},
		{Uri: "/config/.env", Scan: true},
		{Uri: "/environment.env", Scan: true},
		{Uri: "/config/database.env", Scan: true},
		{Uri: "/aws/%2eenv", Scan: true},
		{Uri: "/%2eenv/", Scan: true},

		{Uri: "/.git/", Scan: true},
		{Uri: "/.git/config", Scan: true},
		{Uri: "/.git/HEAD", Scan: true},
		{Uri: "/api/.git/config", Scan: true},
		{Uri: "/app/.git/config", Scan: true},
		{Uri: "/%2egit/", Scan: true},

		{Uri: "/.aws", Scan: true},
		{Uri: "/.aws/config", Scan: true},
		{Uri: "/.AWS_/credentials", Scan: true},
		{Uri: "/%2eaws/", Scan: true},
		{Uri: "/.ssh", Scan: true},
		{Uri: "/.htpasswd", Scan: true},
		{Uri: "/.htaccess", Scan: true},
		{Uri: "/.circleci/configs/development.yml", Scan: true},
		{Uri: "/.anthropic/config.json", Scan: true},
		{Uri: "/.aider.conf.yml", Scan: true},
		{Uri: "/.DS_Store", Scan: true},

		{Uri: "/gcp-credentials.json", Scan: true},
		{Uri: "/aws-credentials%2etxt", Scan: true},
		{Uri: "/aws-secret%2eyaml", Scan: true},
		{Uri: "/aws-ses%2ejson", Scan: true},
		{Uri: "/metrics%2econf", Scan: true},
		{Uri: "/backup.sql", Scan: true},
		{Uri: "/database.yml", Scan: true},
		{Uri: "/id_rsa.key", Scan: true},
		{Uri: "/storage/logs/laravel.log", Scan: true},
		{Uri: "/appsettings.json", Scan: true},
		{Uri: "/config.json", Scan: true},

		{Uri: "/wp-config.php", Scan: true},
		{Uri: "/blog/wp-config.php", Scan: true},
		{Uri: "/xmlrpc.php", Scan: true},
		{Uri: "/wp-admin/install.php?step=1", Scan: true},
		{Uri: "/wp-json/batch/v1", Scan: true},
		{Uri: "/wp-includes/interactivity-api/index.php", Scan: true},
		{Uri: "/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "/blog/wp-includes/registration.php", Scan: true},
		{Uri: "/?rest_route=%2Fbatch%2Fv1", Scan: true},
		{Uri: "/wp-content/plugins/hellopress/wp_filemanager.php", Scan: true},
		{Uri: "//2018/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "//2019/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "//blog/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "//shop/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "//cms/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "/2024/wp-includes/wlwmanifest.xml", Scan: true},

		{Uri: "/1.php", Scan: true},
		{Uri: "/8.php", Scan: true},
		{Uri: "/9.php", Scan: true},
		{Uri: "/222.php", Scan: true},
		{Uri: "//admin.php", Scan: true},
		{Uri: "///admin.php", Scan: true},
		{Uri: "//aa.php", Scan: true},
		{Uri: "//av.php", Scan: true},
		{Uri: "//index.html", Scan: true},
		{Uri: "/../../etc/passwd", Scan: true},
		{Uri: "/static/../../secrets", Scan: true},

		{Uri: "/adminfuns.php", Scan: true},
		{Uri: "///adminfuns.php", Scan: true},
		{Uri: "/adminner.php", Scan: true},
		{Uri: "/this_is_a_new_hello_world.php", Scan: true},

		// Plausible names on a real site, scanners hit them too but the cost
		// of flagging every visitor of such a page is higher.
		{Uri: "/admin.php", Scan: false},
		{Uri: "/config.php", Scan: false},
		{Uri: "/media.php", Scan: false},
		{Uri: "/ops.php", Scan: false},
		{Uri: "/aa.php", Scan: false},

		{Uri: "/_ignition/execute-solution", Scan: true},
		{Uri: "/actuator/env", Scan: true},
		{Uri: "/api/v1/login", Scan: true},
		{Uri: "/api/v1/auto_login", Scan: true},
		{Uri: "/debug/default/view", Scan: true},
		{Uri: "/telescope/requests", Scan: true},
		{Uri: "/server-status", Scan: true},
		{Uri: "/phpmyadmin", Scan: true},

		{Uri: "/$(pwd)/%2egit/config", Scan: true},
		{Uri: "/?cmd=%3B+echo+GSCAN_CMDI+%23", Scan: true},
		{Uri: "/?exec=`id`", Scan: true},
		{Uri: "/${jndi:ldap://x}", Scan: true},
		{Uri: "/projects/$(pwd)/index.html", Scan: true},
		{Uri: "/images/logo.png?size=$(whoami)", Scan: true},
		{Uri: "/%24%28pwd%29/config", Scan: true},
		{Uri: "/api/%60id%60", Scan: true},
		{Uri: "/%24%7bjndi:ldap://x%7d", Scan: true},
		{Uri: "/thou%3C/p%3E%3Ca%20class=", Scan: true},

		{Uri: "/", Scan: false},
		{Uri: "/index.html", Scan: false},
		{Uri: "/projects/2014/12/01/botnet", Scan: false},
		{Uri: "/posts/how-to-use-.env-files", Scan: false},
		{Uri: "/static/style.css", Scan: false},
		{Uri: "/feed.xml", Scan: false},
		{Uri: "/images/logo.png", Scan: false},
		{Uri: "/2014/12/01/some-post", Scan: false},
		{Uri: "/.well-known/security.txt", Scan: false},
		{Uri: "/.well-known/acme-challenge/token", Scan: false},

		// A site may be built with any stack, scripts alone are not a scan.
		{Uri: "/index.php", Scan: false},
		{Uri: "/blog/post.php?id=5", Scan: false},
		{Uri: "/shop/checkout.aspx", Scan: false},
		{Uri: "/app/handler.jsp", Scan: false},
		{Uri: "/wp-login.php", Scan: false},
		{Uri: "/wp-content/themes/theme/style.css", Scan: false},
		{Uri: "/wp-includes/js/jquery/jquery.min.js", Scan: false},
		{Uri: "/wp-includes/css/dist/block-library/style.min.css", Scan: false},
		{Uri: "/wp-includes/js/wp-emoji-release.min.js", Scan: false},
		{Uri: "/wp-json/wp/v2/posts", Scan: false},
		{Uri: "/blog/password-managers", Scan: false},
		{Uri: "/articles/my-secret-recipe", Scan: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Uri, func(t *testing.T) {
			if scan := isScanRequest(testCase.Uri); scan != testCase.Scan {
				t.Errorf("got %v, want %v", scan, testCase.Scan)
			}
		})
	}
}

func TestClassifyMarksMaliciousAddresses(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	for i := 0; i <= MaliciousRequestThreshold; i++ {
		c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/.env", now.Add(-time.Duration(i)*time.Hour)))
	}

	if category := c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/", now)); category != CategoryMalicious {
		t.Errorf("got %q, want %q", category, CategoryMalicious)
	}

	if category := c.Classify(classifierEntry("2.2.2.2", "Mozilla/5.0", "/", now)); category != CategoryUnclassified {
		t.Errorf("got %q, want %q", category, CategoryUnclassified)
	}
}

func TestClassifyIgnoresRequestsOutsideOfTheWindow(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()
	old := now.Add(-trafficWindow).Add(-24 * time.Hour)

	for i := 0; i <= MaliciousRequestThreshold; i++ {
		c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/.env", old))
	}

	if category := c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/", now)); category != CategoryUnclassified {
		t.Errorf("got %q, want %q", category, CategoryUnclassified)
	}
}

func TestClassifyKeepsAddressesBelowTheThresholdUnclassified(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	for i := 0; i < MaliciousRequestThreshold; i++ {
		c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/.env", now))
	}

	if category := c.Classify(classifierEntry("1.1.1.1", "Mozilla/5.0", "/", now)); category != CategoryUnclassified {
		t.Errorf("got %q, want %q", category, CategoryUnclassified)
	}
}

func TestRemoveOldDataDiscardsAddresses(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	c.Classify(classifierEntry("1.1.1.1", "curl/7.64.0", "/", now))

	c.RemoveOldData(now)
	if len(c.ips) != 1 {
		t.Fatalf("got %v", c.ips)
	}

	c.RemoveOldData(now.Add(trafficWindow).Add(24 * time.Hour))
	if len(c.ips) != 0 {
		t.Fatalf("got %v", c.ips)
	}
}

func TestUserAgentName(t *testing.T) {
	testCases := []struct {
		UserAgent string
		Name      string
	}{
		{
			UserAgent: "curl/7.64.0",
			Name:      "curl",
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Name:      "Mozilla",
		},
		{
			UserAgent: "moooodotfarm (+https://moooo.farm)",
			Name:      "moooodotfarm",
		},
		{
			UserAgent: "Turnitin (https://turnitin.com/robot/crawlerinfo.html)",
			Name:      "Turnitin",
		},
		{
			UserAgent: "OnlyRead",
			Name:      "OnlyRead",
		},
		{
			UserAgent: "python-requests/2.32.3",
			Name:      "python-requests",
		},
		{
			UserAgent: "Feedbin feed-id:1602368 - 4 subscribers",
			Name:      "Feedbin feed-id:1602368 - 4 subscribers",
		},
		{
			UserAgent: "http.rb/5.1.0",
			Name:      "http.rb/5.1.0",
		},
		{
			UserAgent: "com.apple.webkit.networking",
			Name:      "com.apple.webkit.networking",
		},
		{
			UserAgent: "vuln_scanner",
			Name:      "vuln_scanner",
		},
		{
			UserAgent: "vuln_scanner/3.1.0 (CVE-2026-4020)",
			Name:      "vuln_scanner",
		},
		{
			UserAgent: "crusader-worker/2.1",
			Name:      "crusader-worker",
		},
		{
			UserAgent: "",
			Name:      "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			if name := UserAgentName(testCase.UserAgent); name != testCase.Name {
				t.Errorf("got %q, want %q", name, testCase.Name)
			}
		})
	}
}

func TestClassifyUserAgentName(t *testing.T) {
	testCases := []struct {
		UserAgent string
		Category  Category
	}{
		{
			UserAgent: "curl/7.64.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "curl",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Prometheus/2.48.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "Googlebot-Image/1.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "DuckDuckBot/1.1; (+http://duckduckgo.com/duckduckbot.html)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "console-feedreader/1.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "vuln_scanner/3.1.0 (CVE-2026-4020)",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "sqlmap/1.9#stable (https://sqlmap.org)",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "Mastodon/4.4.2 (+https://mastodon.social/)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Feedbin feed-id:1602368 - 4 subscribers",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "NewsBlur Feed Fetcher - 2 subscribers - https://www.newsblur.com",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "FreshRSS/1.24.3 (Linux; https://freshrss.org)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "CommaFeed/4.6.0 (https://github.com/Athou/commafeed)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "ScourRssBot/1.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Misskey/2025.4.1",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Blogtrottr/2.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "-",
			Category:  CategoryMalicious,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			if category := classifyUserAgent(testCase.UserAgent); category != testCase.Category {
				t.Errorf("got %q, want %q", category, testCase.Category)
			}
		})
	}
}
