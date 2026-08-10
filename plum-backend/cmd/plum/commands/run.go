package commands

import (
	"encoding/json"
	"os"
	"runtime"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/boreq/plum/plum-backend/adapters"
	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/config"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
	"github.com/boreq/plum/plum-backend/entrypoints/http"
	"github.com/boreq/plum/plum-backend/entrypoints/logs"
	"github.com/boreq/plum/plum-backend/entrypoints/timers"
	"github.com/dustin/go-humanize"
)

const removeOldDataEvery = time.Hour

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

	whitelist, err := domain.NewWhitelist(conf.Whitelist)
	if err != nil {
		return errors.Wrap(err, "could not create the whitelist")
	}

	errC := make(chan error)

	repositories := adapters.NewRepositories()
	maliciousAddresses := adapters.NewMaliciousAddresses(whitelist)

	go logMemoryStats()

	application := app.New(repositories, maliciousAddresses, whitelist)

	for i := range conf.Websites {
		website := conf.Websites[i]

		websiteName, err := domain.NewWebsiteName(website.Name)
		if err != nil {
			return errors.Wrap(err, "could not create the website name")
		}

		p, err := parser.NewParser(getLogFormat(website.LogFormat))
		if err != nil {
			return err
		}

		if err := repositories.Add(websiteName, adapters.NewRepository(website, maliciousAddresses)); err != nil {
			return errors.Wrap(err, "could not add a repository")
		}

		logReader := adapters.NewLogReader()

		entrypoint := logs.New(
			websiteName,
			p,
			website.Load,
			website.Follow,
			logReader,
			application.AddRequest,
		)

		go printStats(website.Name, logReader)

		go func() {
			errC <- entrypoint.Run()
		}()
	}

	removeOldData := timers.NewRemoveOldData(application.RemoveOldData, removeOldDataEvery)

	go func() {
		errC <- removeOldData.Run()
	}()

	handler, err := http.NewHandler(application)
	if err != nil {
		return errors.Wrap(err, "could not create the http handler")
	}

	server := http.NewServer(handler)

	go func() {
		errC <- server.Serve(conf.ServeAddress)
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

func printStats(websiteName string, logReader *adapters.LogReader) {
	lastLines, _ := logReader.GetStats()
	duration := 5 * time.Second
	for range time.Tick(duration) {
		lines, _ := logReader.GetStats()
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
