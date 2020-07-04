// Package config holds the configuration struct.
package config

import "github.com/pkg/errors"

type Config struct {
	// HTTP server address eg. "127.0.0.1:8118".
	ServeAddress string `json:"serveAddress"`

	// Groups of logs which will be monitored.
	Websites []Website `json:"websites"`
}

func (c Config) Valid() error {
	if len(c.Websites) == 0 {
		return errors.New("no websites defined")
	}

	for _, website := range c.Websites {
		if err:= website.Valid(); err != nil {
			return errors.Wrap(err, "website definition invalid")
		}
	}

	return nil
}

type Website struct {
	// Name of this website visible in the frontend.
	Name string `json:"name"`

	// Log file to be monitored eg. "path/to/my.access.log".
	Follow string `json:"follow"`

	// Log files to be initially loaded eg. "path/to/my.access.log.*". Will be
	// expanded using filepath.Glob.
	Load []string `json:"load"`

	// Specifies the predefined format of the logs being analyzed eg.
	// "combined" or a custom log format. Please refer to PredefinedFormats
	// defined in parser/parser.go.
	LogFormat string `json:"logFormat"`

	// Trims trailing slashes from request URIs before aggregating them.
	NormalizeSlash bool `json:"-"`

	// Trims queries from request URIs before aggregating them.
	NormalizeQuery bool `json:"-"`

	// Trims "http://" and "https://" prefixes from referrer URIs before
	// aggregating them.
	StripRefererProtocol bool `json:"-"`
}

func (w Website) Valid() error {
	if w.Name == "" {
		return errors.New("blank name")
	}

	if w.Follow == "" {
		return errors.New("blank follow path")
	}

	for _,load :=  range w.Load {
		if load == "" {
			return errors.New("blank load path")
		}
	}

	if w.LogFormat == "" {
		return errors.New("blank log format")
	}

	return nil
}

// Default returns the default config.
func Default() *Config {
	conf := &Config{
		ServeAddress: "127.0.0.1:8118",
		Websites: []Website{
			{
				Name:   "example.com",
				Follow: "path/to/my.access.log",
				Load: []string{
					"path/to/my.access.log.*",
				},
				LogFormat:            "combined",
				NormalizeSlash:       true,
				NormalizeQuery:       true,
				StripRefererProtocol: true,
			},
		},
	}
	return conf
}
