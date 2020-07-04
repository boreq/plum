# Plum

Plum is a real-time web server access log analyser. It allows the user to
access the statistics using a web dashboard.

## Installation

Plum is written in Go which means that the Go tools can be used to install the
program using the following command:

    $  go get github.com/boreq/plum/cmd/plum

If you prefer to do this by hand clone the repository and execute the `make`
program:

    $ git clone https://github.com/boreq/plum
    $ make
    $ ls build
    plum

## Usage

Plum can track a single log file in order to produce live data as well as load
any number of past log files in order to present historical data points.

A common scenario on a `logrotate` enabled server is following a single log
file in order to produce live data as well as loading any number of past log
files in order to present historical data points.

    $ ls
    my.access.log
    my.access.log.1.gz
    my.access.log.2.gz
    my.access.log.3.gz
    my.access.log.4.gz
    
You can use the following command to generate a default config file:

    $ plum default_config > config.json

Based on that prepare a config which will load the current access log as well
as the older logs:

    $ cat config.json
    {
        "serveAddress": "127.0.0.1:8118",
        "websites": [
            {
                "name": "example.com",
                "follow": "/path/to/example.access.log",
                "load": [
                    "/path/to/example.access.log.*"
                ],
                "logFormat": "combined"
            }
        ]
    }
    
Execute the program:
    
    $ plum run config.json
    INFO starting listening                       source=server address=127.0.0.1:8118

Navigate to http://127.0.0.0:8118 to see the results.

## Configuration

### `websites.logFormat`

Specifies the log format. A custom or a predefined format can be used.

#### Predefined formats
When using a predefined format simply pass its name as the argument of this
option.

- `combined` - a default format used by NGINX


#### Custom formats
When using a custom format a number of elements can be used to construct it, check out [parser.go](https://github.com/boreq/plum/blob/master/parser/parser.go).

