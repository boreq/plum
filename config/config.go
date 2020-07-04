// Package config holds the configuration struct.
package config

type Config struct {
	// HTTP server address eg. "127.0.0.1:8118".
	ServeAddress         string `json:"serveAddress"`

	// Specifies the predefined format of the logs being analyzed eg.
	// "combined" or a custom log format. Please refer to PredefinedFormats
	// defined in parser/parser.go.
	LogFormat            string `json:"logFormat"`

	// Trims trailing slashes from request URIs before aggregating them.
	NormalizeSlash       bool `json:"-"`

	// Trims queries from request URIs before aggregating them.
	NormalizeQuery       bool `json:"-"`

	// Trims "http://" and "https://" prefixes from referrer URIs before
	// aggregating them.
	StripRefererProtocol bool `json:"-"`
}

// Default returns the default config.
func Default() *Config {
	conf := &Config{
		ServeAddress:         "127.0.0.1:8118",
		LogFormat:            "combined",
		NormalizeSlash:       true,
		NormalizeQuery:       true,
		StripRefererProtocol: true,
	}
	return conf
}
