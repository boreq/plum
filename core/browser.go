package core

type Browser struct {
	id   string
	name string
}

var (
	BrowserChrome           = &Browser{"chrome", "Chrome"}
	BrowserChromium         = &Browser{"chromium", "Chromium"}
	BrowserFirefox          = &Browser{"firefox", "Firefox"}
	BrowserSafari           = &Browser{"safari", "Safari"}
	BrowserEdge             = &Browser{"edge", "Edge"}
	BrowserOpera            = &Browser{"opera", "Opera"}
	BrowserBrave            = &Browser{"brave", "Brave"}
	BrowserVivaldi          = &Browser{"vivaldi", "Vivaldi"}
	BrowserDuckDuckGo       = &Browser{"duckduckgo", "DuckDuckGo"}
	BrowserSamsungInternet  = &Browser{"samsung-internet", "Samsung Internet"}
	BrowserYandex           = &Browser{"yandex-browser", "Yandex Browser"}
	BrowserInternetExplorer = &Browser{"internet-explorer", "Internet Explorer"}
)

func (b *Browser) Name() string {
	return b.name
}

func (b *Browser) String() string {
	return b.id
}
