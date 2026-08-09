# Plum

Plum is a self-hosted real-time web server access log analyser. It allows you
to view analytics derived from your web server logs using a web dashboard. No
analytics services integrating with your software required. Just tell Plum
where your access logs are.

I built Plum for several reasons. Due to laziness, I wanted to avoid manually
adding analytics to all my software. Due to privacy concerns, I wanted to avoid
existing analytics services, which usually means sending data to external
parties. You can find open source and self-hosted alternatives, but if they are
based on some kind of tracking on the user side of things then most people
block them anyway, which makes them ineffective.

## Installation

### Docker

The repository ships a `Dockerfile`. Mount your config at `/config.json` and
the log files at the paths it references:

    $ docker build -t plum .
    $ docker run \
        -v /var/log/nginx:/var/log/nginx:ro \
        -v $(pwd)/config.json:/config.json:ro \
        -p 9000:8118 \
        plum

Or using Docker Compose:

    services:
      plum:
        build: ./plum
        volumes:
          - /var/log/nginx:/var/log/nginx:ro
          - ./config.json:/config.json:ro
        ports:
          - "9000:8118"
        restart: always

Here `./plum` is where you checked out this repository.

The paths in the config file are the paths *inside* the container, so they have
to match the mounted log directory:

    $ cat config.json
    {
        "serveAddress": ":8118",
        "websites": [
            {
                "name": "example.com",
                "follow": "/var/log/nginx/example.access.log",
                "load": [
                    "/var/log/nginx/example.access.log.*"
                ],
                "logFormat": "combined"
            }
        ]
    }

With the compose file above the dashboard is available on the host under
http://127.0.0.1:9000.

### Manual

Plum is written in Go which means that the Go tools can be used to install the
program using the following command:

    $  go install github.com/boreq/plum/plum-backend/cmd/plum@latest

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
When using a custom format a number of elements can be used to construct it, check out [parser.go](https://github.com/boreq/plum/blob/master/plum-backend/domain/parser/parser.go).

## Development

Run the checks which run in CI, both for the backend and the frontend:

    $ make ci

### Building the frontend

The dashboard lives in `plum-frontend` and is compiled into the binary, so a
plain `go install` or `docker build` does not require Node. The generated
assets are committed to the repository in
`plum-backend/entrypoints/http/statik`.

After changing anything in `plum-frontend` regenerate them:

    $ make frontend

This installs the frontend dependencies, runs the Vite build and packs the
result into `statik.go`. The output is deterministic, so an unchanged frontend
produces an unchanged `statik.go`. Commit the regenerated file together with
the frontend changes.

