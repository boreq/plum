package domain

import (
	"testing"
	"time"

	"github.com/boreq/plum/plum-backend/domain/request"
)

func classify(c *TrafficClassifier, userAgent, uri string, t time.Time) Category {
	return c.Classify(request.NewUri(uri), request.NewUserAgent(userAgent), t)
}

func TestClassifyUsesTheUserAgent(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	if category := classify(c, "curl/7.64.0", "/", now); category != CategoryAutomated {
		t.Errorf("got %q, want %q", category, CategoryAutomated)
	}

	if category := classify(c, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", "/", now); category != CategoryUnclassified {
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

			category := classify(c, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36", testCase.Uri, time.Now().UTC())
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
		{Uri: "/gcp.json", Scan: true},
		{Uri: "/mail.config.js", Scan: true},
		{Uri: "/aws-credentials%2etxt", Scan: true},
		{Uri: "/aws-secret%2eyaml", Scan: true},
		{Uri: "/aws-ses%2ejson", Scan: true},
		{Uri: "/metrics%2econf", Scan: true},
		{Uri: "/backup.sql", Scan: true},
		{Uri: "/database.yml", Scan: true},
		{Uri: "/gradle.properties", Scan: true},
		{Uri: "/app/gradle%2eproperties", Scan: true},
		{Uri: "/.gradle/gradle.properties", Scan: true},
		{Uri: "/id_rsa.key", Scan: true},
		{Uri: "/storage/logs/laravel.log", Scan: true},
		{Uri: "/appsettings.json", Scan: true},
		{Uri: "/api/appsettings.json", Scan: true},
		{Uri: "/appsettings.Production.json", Scan: true},
		{Uri: "/appsettings.Development.json", Scan: true},
		{Uri: "/application-dev.properties", Scan: true},
		{Uri: "/app/credentials.json", Scan: true},
		{Uri: "/config.json", Scan: true},
		{Uri: "/static/config.json", Scan: true},
		{Uri: "/config.inc.php.dist", Scan: true},
		{Uri: "/Dockerfile", Scan: true},
		{Uri: "/app/next%2econfig%2ejs", Scan: true},
		{Uri: "/app/next%2econfig%2emjs", Scan: true},
		{Uri: "/app/nuxt%2econfig%2ejs", Scan: true},
		{Uri: "/next.config.ts", Scan: true},
		{Uri: "/%252eenv%252ecredentials", Scan: true},
		{Uri: "/%2eenv", Scan: true},
		{Uri: "/%63%6f%6e%66%69%67.%70%68%70%7e", Scan: true},
		{Uri: "/%29Dotfile", Scan: true},
		{Uri: "/)dotfile", Scan: true},
		{Uri: "/'; DROP TABLE", Scan: true},
		{Uri: "/config.php~", Scan: true},
		{Uri: "/index.html~", Scan: true},
		{Uri: "/%73%65%63%72%65%74", Scan: true},
		{Uri: "/%2563%2566", Scan: true},
		{Uri: "/%25252e%25252e/%25252e%25252e/etc/passwd", Scan: true},
		{Uri: "/phpmyadmin/config.inc.php.dist", Scan: true},
		{Uri: "/deploy/aws/gcp-credentials.json", Scan: true},
		{Uri: "/deploy/aws-credentials.txt", Scan: true},
		{Uri: "/deploy/aws-ses.json", Scan: true},
		{Uri: "/backup.zip", Scan: true},
		{Uri: "/secrets.json", Scan: true},
		{Uri: "/config/secrets.json", Scan: true},
		{Uri: "/webpack-stats.json", Scan: true},
		{Uri: "/asset-manifest.json", Scan: true},
		{Uri: "/manifest.json", Scan: true},
		{Uri: "/static/asset-manifest.json", Scan: true},
		{Uri: "/static/webpack-stats.json", Scan: true},

		{Uri: "/wp-config.php", Scan: true},
		{Uri: "/blog/wp-config.php", Scan: true},
		{Uri: "/some/nested/directory/wp-config.php", Scan: true},
		{Uri: "/xmlrpc.php", Scan: true},
		{Uri: "/blog/xmlrpc.php", Scan: true},
		{Uri: "/wlwmanifest.xml", Scan: true},
		{Uri: "/blog/wlwmanifest.xml", Scan: true},
		{Uri: "/wp-admin/install.php?step=1", Scan: true},
		{Uri: "/wp-admin/admin-ajax.php", Scan: true},
		{Uri: "/wp-json/batch/v1", Scan: true},
		{Uri: "/wp-includes/interactivity-api/index.php", Scan: true},
		{Uri: "/wp-includes/wlwmanifest.xml", Scan: true},
		{Uri: "/blog/wp-includes/registration.php", Scan: true},
		{Uri: "/wp-includes/js/jquery/jquery.min.js", Scan: true},
		{Uri: "/wp-includes/css/dist/block-library/style.min.css", Scan: true},
		{Uri: "/wp-includes/js/wp-emoji-release.min.js", Scan: true},
		{Uri: "/wp-content/themes/theme/style.css", Scan: true},
		{Uri: "/wp-content/uploads/2024/photo.jpg", Scan: true},
		{Uri: "/wp-admin/", Scan: true},
		{Uri: "/wp-json/wp/v2/posts", Scan: true},
		{Uri: "/wp-json/wp/v2/users", Scan: true},
		{Uri: "/?rest_route=%2Fbatch%2Fv1", Scan: true},
		{Uri: "/wp-content/plugins/hellopress/wp_filemanager.php", Scan: true},
		{Uri: "/wp-content/plugins/index.php", Scan: true},
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
		{Uri: "/site/adminfuns.php", Scan: true},
		{Uri: "/adminner.php", Scan: true},
		{Uri: "/admin/adminner.php", Scan: true},
		{Uri: "/classwithtostring.php", Scan: true},
		{Uri: "/src/classwithtostring.php", Scan: true},
		{Uri: "/this_is_a_new_hello_world.php", Scan: true},
		{Uri: "/tmp/this_is_a_new_hello_world.php", Scan: true},
		{Uri: "/abcd.php", Scan: true},
		{Uri: "/akcc.php", Scan: true},
		{Uri: "/bless.php", Scan: true},
		{Uri: "/blurbs.php", Scan: true},
		{Uri: "/bolt.php", Scan: true},
		{Uri: "/chosen.php", Scan: true},
		{Uri: "/cjfuns.php", Scan: true},
		{Uri: "/class-t.api.php", Scan: true},
		{Uri: "/cord.php", Scan: true},
		{Uri: "/dex.php", Scan: true},
		{Uri: "/flower.php", Scan: true},
		{Uri: "/gifclass.php", Scan: true},
		{Uri: "/shelp.php", Scan: true},
		{Uri: "/config.php", Scan: true},
		{Uri: "/app/config.php", Scan: true},
		{Uri: "/phpinfo.php", Scan: true},
		{Uri: "/test/phpinfo.php", Scan: true},
		{Uri: "/phpinfo", Scan: true},
		{Uri: "/test/phpinfo", Scan: true},
		{Uri: "/php_info.php", Scan: true},
		{Uri: "/test/php_info.php", Scan: true},
		{Uri: "/php_info", Scan: true},
		{Uri: "/test/php_info", Scan: true},
		{Uri: "/php-info.php", Scan: true},
		{Uri: "/test/php-info.php", Scan: true},
		{Uri: "/php-info", Scan: true},
		{Uri: "/test/php-info", Scan: true},
		{Uri: "/phpinfo.php3", Scan: true},
		{Uri: "/test/phpinfo.php3", Scan: true},
		{Uri: "/admin.php", Scan: true},
		{Uri: "/panel/admin.php", Scan: true},
		{Uri: "/wp-login.php", Scan: true},
		{Uri: "/blog/wp-login.php", Scan: true},

		{Uri: "/posts/wp-config", Scan: false},
		{Uri: "/posts/wp-config-php-explained", Scan: false},
		{Uri: "/blog/xmlrpc-php-security", Scan: false},
		{Uri: "/app/myconfig.php", Scan: false},
		{Uri: "/app/config-loader.php", Scan: false},
		{Uri: "/docs/phpinfo-tutorial", Scan: false},
		{Uri: "/docs/using_php_info", Scan: false},
		{Uri: "/static/configuration.js", Scan: false},

		// Plausible names on a real site, scanners hit them too but the cost
		// of flagging every visitor of such a page is higher.
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
		{Uri: "/blog/password-managers", Scan: false},
		{Uri: "/articles/my-secret-recipe", Scan: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Uri, func(t *testing.T) {
			if scan := NewTrafficClassifier().isScanRequest(testCase.Uri); scan != testCase.Scan {
				t.Errorf("got %v, want %v", scan, testCase.Scan)
			}
		})
	}
}

const classifierBrowserUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"

func TestClassifyDoesNotMarkMaliciousAddresses(t *testing.T) {
	c := NewTrafficClassifier()
	now := time.Now().UTC()

	for i := 0; i <= MaliciousRequestThreshold; i++ {
		classify(c, classifierBrowserUserAgent, "/.env", now)
	}

	if category := classify(c, classifierBrowserUserAgent, "/", now); category != CategoryUnclassified {
		t.Errorf("got %q, want %q", category, CategoryUnclassified)
	}
}

var classifyReferenceTime = time.Date(2026, time.August, 7, 0, 0, 0, 0, time.UTC)

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
			Category:  CategoryPossiblyAutomated,
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
			UserAgent: "FrontPageDomainPipeline/1.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; SomeScaner/1.0)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; SomeScanner/1.0)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "OWLer",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Dart/3.11 (dart:io)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Java/1.8.0_292",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "GDELT-media-org-research/1.0 (academic)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "HNBlogArchive/0.1 (+https://news.ycombinator.com/)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "console-feedreader/1.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "CCBot/2.0 (https://commoncrawl.org/faq/)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Claude-User (claude-code/2.1.219; +https://support.anthropic.com/)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "python-requests/2.32.3",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Python/3.11 aiohttp/3.9.5",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Go-http-client/2.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Scrapy/2.11.0 (+https://scrapy.org)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Apache-HttpClient/5.3 (Java/17.0.10)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Bun/1.1.29",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "DomainDrift/0.3 (Internet Telemetry; +https://domaindrift.io)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "com.apple.WebKit.Networking/21624.2.5.11.8 Network/5812.121.1 macOS/26.5.2",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "NetworkingExtension/1 Network/5812.121.1 iOS/26.5",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Roku/DVP-14.5 (14.5.0.4126-46)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "http.rb/5.1.1 (Mastodon/4.2.17; +https://mastodon.social/)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Python-urllib/3.11",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Hello from Palo Alto Networks, find out more about our scans in https://docs-cortex.paloaltonetworks.com/r/1/Cortex-Xpanse/Scanning-activity",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "NewsBlur Page Fetcher - 2 subscribers - https://www.newsblur.com/site/7649476/thoughts",
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
			UserAgent: "Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com)",
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
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "-",
			Category:  CategoryMalicious,
		},

		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36 Edg/151.0.0.0",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 OPR/101.0.0.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 5.1; rv:38.0) Gecko/20100101 Firefox/38.0 SeaMonkey/2.35",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Maemo; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36 TungstenBrowser/2.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 9; HRY-LX1 Build/HONORHRY-L21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.93 Mobile Safari/537.36 YaApp_Android/10.20 YaSearchBrowser/10.20",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 8.1.0; M1813 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045018 Mobile Safari/537.36 MMWEBID/5434 MicroMessenger/7.0.9.1560(0x27000935) Process/tools NetType/4G Language/zh_CN ABI/arm64",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 7.1.2; Redmi 5 Build/N2G47H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/63.0.3239.83 Mobile Safari/537.36 T7/11.3 baiduboxapp/11.3.6.11 (Baidu; P1 7.1.2)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; WOW64; rv:41.0) Gecko/20100101 Firefox/140.0.2 (x64 de)",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Safari",
			Category:  CategoryUnclassified,
		},
		{
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.6099.119 Mobile/15E148 Safari/604.1",
			Category:  CategoryPossiblyAutomated,
		},

		{
			UserAgent: "Mozilla/5.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; HTML Proofer/5.0.10; +https://github.com/gjtorikian/html-proofer)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; CensysInspect/1.1; +https://about.censys.io/)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; SiteInspector/1.0)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; NetcraftSurveyAgent/1.0; +info@netcraft.com)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Android) Nextcloud-android/3.29.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Linux) mirall/3.13.0 (Nextcloud, ubuntu-6.8.0-45-generic ClientArchitecture: x86_64)",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Thunderbird/128.4.3",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (compatible; MSIE 8.0; Windows 98; Trident/5.1)",
			Category:  CategoryPossiblyAutomated,
		},

		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36 curl/8.5.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.7871.186 Mobile Safari/537.36 (compatible; GoogleOther)",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/138.0.7204.23 Safari/537.36",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36 sqlmap/1.9",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "curl/8.5.0 nuclei/3.1.0",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "TLM-Audit-Scanner/1.0",
			Category:  CategoryMalicious,
		},
		{
			UserAgent: "Newsboat/2.40.0",
			Category:  CategoryAutomated,
		},
		{
			UserAgent: "mozilla/5.0 (windows nt 10.0; win64; x64) applewebkit/537.36 (khtml, like gecko) chrome/128.0.6613.138 safari/537.36",
			Category:  CategoryPossiblyAutomated,
		},

		// Trap: an all-lowercase UA from a legitimate custom client is
		// incorrectly flagged. The mixed-case version is unclassified.
		{
			UserAgent: "someuseragent/1.0.0",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "SomeUserAgent/1.0.0",
			Category:  CategoryUnclassified,
		},

		// Quoted user agents are sent by misconfigured/automated tools.
		// A quote at either end is enough to trigger detection.
		{
			UserAgent: "\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\"",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\"",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36'",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36'",
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `\x22Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\x22`,
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `\x22Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36`,
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\x22`,
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `\x27Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\x27`,
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `\x27Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36`,
			Category:  CategoryPossiblyAutomated,
		},
		{
			UserAgent: `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36\x27`,
			Category:  CategoryPossiblyAutomated,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.UserAgent, func(t *testing.T) {
			if category := NewTrafficClassifier().classifyUserAgent(testCase.UserAgent, classifyReferenceTime); category != testCase.Category {
				t.Errorf("got %q, want %q", category, testCase.Category)
			}
		})
	}
}

func TestUnescapeUri(t *testing.T) {
	testCases := []struct {
		Name     string
		Uri      string
		Expected string
	}{
		{
			Name:     "uri without escape sequences",
			Uri:      "/blog/post",
			Expected: "/blog/post",
		},
		{
			Name:     "uri is lowercased",
			Uri:      "/Blog/Post",
			Expected: "/blog/post",
		},
		{
			Name:     "single encoding",
			Uri:      "/%2eenv",
			Expected: "/.env",
		},
		{
			Name:     "uppercase escape sequence",
			Uri:      "/%2Eenv",
			Expected: "/.env",
		},
		{
			Name:     "double encoding",
			Uri:      "/%252eenv%252ecredentials",
			Expected: "/.env.credentials",
		},
		{
			Name:     "triple encoding",
			Uri:      "/%25252e%25252e/etc/passwd",
			Expected: "/../etc/passwd",
		},
		{
			Name:     "encoded query string",
			Uri:      "/?rest_route=%2Fbatch%2Fv1",
			Expected: "/?rest_route=/batch/v1",
		},
		{
			Name:     "invalid escape sequence is left alone",
			Uri:      "/50%-off",
			Expected: "/50%-off",
		},
		{
			Name:     "decoding stops at an invalid escape sequence",
			Uri:      "/%2550%25-off",
			Expected: "/%50%-off",
		},
		{
			Name:     "plus is not treated as a space",
			Uri:      "/a+b",
			Expected: "/a+b",
		},
		{
			Name:     "decoding is bounded",
			Uri:      "/%25252525252e",
			Expected: "/%2e",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if uri := unescapeUri(testCase.Uri); uri != testCase.Expected {
				t.Errorf("got %q, want %q", uri, testCase.Expected)
			}
		})
	}
}

func TestUserAgentProducts(t *testing.T) {
	testCases := []struct {
		Name      string
		UserAgent string
		Products  []string
	}{
		{
			Name:      "product with a version",
			UserAgent: "curl/7.64.0",
			Products:  []string{"curl"},
		},
		{
			Name:      "product without a version",
			UserAgent: "onlyread",
			Products:  []string{"onlyread"},
		},
		{
			Name:      "every product is returned",
			UserAgent: "python/3.11 aiohttp/3.9.5",
			Products:  []string{"python", "aiohttp"},
		},
		{
			Name:      "products named in a comment",
			UserAgent: "mozilla/5.0 (compatible; googlebot/2.1; +http://www.google.com/bot.html)",
			Products:  []string{"mozilla", "compatible", "googlebot", "http:"},
		},
		{
			Name:      "product named using multiple words",
			UserAgent: "sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
			Products:  []string{"sogou web spider", "sogou", "web", "spider", "http:"},
		},
		{
			Name:      "product appended behind a browser",
			UserAgent: "mozilla/5.0 (x11; linux x86_64) applewebkit/537.36 (khtml, like gecko) chrome/151.0.0.0 safari/537.36 curl/8.5.0",
			Products:  []string{"mozilla", "x11", "linux x86_64", "linux", "x86_64", "applewebkit", "khtml", "like gecko", "like", "gecko", "chrome", "safari", "curl"},
		},
		{
			Name:      "empty user agent",
			UserAgent: "",
			Products:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			products := NewTrafficClassifier().userAgentProducts(testCase.UserAgent)

			if len(products) != len(testCase.Products) {
				t.Fatalf("got %q, want %q", products, testCase.Products)
			}

			for i := range products {
				if products[i] != testCase.Products[i] {
					t.Fatalf("got %q, want %q", products, testCase.Products)
				}
			}
		})
	}
}
