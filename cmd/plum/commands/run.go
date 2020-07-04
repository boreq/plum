package commands

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/boreq/guinea"
	"github.com/boreq/plum/config"
	"github.com/boreq/plum/core"
	"github.com/boreq/plum/parser"
	"github.com/boreq/plum/server"
	"github.com/dustin/go-humanize"
)

var runCmd = guinea.Command{
	Run: runRun,
	Arguments: []guinea.Argument{
		{
			Name:        "config",
			Optional:    false,
			Multiple:    false,
			Description: "Config file to be used",
		},
	},
	ShortDescription: "loads and follows log files",
}

func runRun(c guinea.Context) error {
	conf, err := loadConfig(c.Arguments[0])
	if err != nil {
		return errors.Wrap(err, "could not load the configuration")
	}

	errC := make(chan error)

	repositories := core.NewRepositories()

	go logMemoryStats()

	for i := range conf.Websites {
		website := conf.Websites[i]

		p, err := parser.NewParser(getLogFormat(website.LogFormat))
		if err != nil {
			return err
		}

		r := core.NewRepository(website)

		tracker := core.NewTracker(p, r)

		go printStats(website.Name, tracker)

		// Load the specified files
		for _, glob := range website.Load {
			paths, err := filepath.Glob(glob)
			if err != nil {
				return errors.Wrapf(err,  "could not process a glob pattern '%s", glob)
			}

			for _, path := range paths {
				if err := tracker.Load(path); err != nil {
					return err
				}
			}
		}

		// Track the specified file
		go func() {
			errC <- tracker.Follow(website.Follow)
		}()

		if err := repositories.Add(website.Name, r); err != nil {
			return errors.Wrap(err, "could not add a repository")
		}
	}

	go func() {
		errC <- server.Serve(repositories, conf.ServeAddress)
	}()

	return <-errC
}

// getLogFormat tries to find and return a predefined format with the provided
// name or otherwise returns the provided format unaltered assuming that it is
// a format string.
func getLogFormat(format string) string {
	predefinedFormat, ok := parser.PredefinedFormats[format]
	if ok {
		return predefinedFormat
	}
	return format
}

func printStats(websiteName string, tracker *core.Tracker) {
	lastLines, _ := tracker.GetStats()
	duration := 5 * time.Second
	for range time.Tick(duration) {
		lines, _ := tracker.GetStats()
		linesPerSecond := float64(lines-lastLines) / duration.Seconds()
		log.Debug("data statistics", "totalLines", lines, "linesPerSecond", linesPerSecond, "website", websiteName)
		lastLines = lines
	}
}

func logMemoryStats() {
	duration := 10 * time.Second
	for range time.Tick(duration) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		alloc := humanize.Bytes(m.Alloc)
		totalAlloc := humanize.Bytes(m.TotalAlloc)
		sys := humanize.Bytes(m.Sys)
		numGC := m.NumGC

		log.Debug("memory statistics", "alloc", alloc, "totalAlloc", totalAlloc, "sys", sys, "numGC", numGC)
	}
}

func loadConfig(path string) (*config.Config, error) {
	conf := config.Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open the config file")
	}

	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal the config")
	}

	if err := conf.Valid(); err != nil {
		return nil, errors.Wrap(err, "invalid config")
	}

	return conf, nil
}
