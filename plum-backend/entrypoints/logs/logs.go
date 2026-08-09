package logs

import (
	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/adapters"
	"github.com/boreq/plum/plum-backend/app"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
	"github.com/boreq/plum/plum-backend/logging"
)

type AddRequestHandler interface {
	Execute(cmd app.AddRequest) error
}

type LogReader interface {
	Load(globs []string, handler adapters.LineHandler) error
	Follow(path string, handler adapters.LineHandler) error
}

type Logs struct {
	website    domain.WebsiteName
	parser     *parser.Parser
	loadGlobs  []string
	followPath string
	logReader  LogReader
	addRequest AddRequestHandler
	log        logging.Logger
}

func New(
	website domain.WebsiteName,
	parser *parser.Parser,
	loadGlobs []string,
	followPath string,
	logReader LogReader,
	addRequest AddRequestHandler,
) *Logs {
	return &Logs{
		website:    website,
		parser:     parser,
		loadGlobs:  loadGlobs,
		followPath: followPath,
		logReader:  logReader,
		addRequest: addRequest,
		log:        logging.New("entrypoints/logs.Logs"),
	}
}

func (l *Logs) Run() error {
	if err := l.logReader.Load(l.loadGlobs, l); err != nil {
		return errors.Wrap(err, "could not load the log files")
	}

	return l.logReader.Follow(l.followPath, l)
}

func (l *Logs) Handle(line string) error {
	entry, err := l.parser.Parse(line)
	if err != nil {
		return errors.Wrap(err, "could not parse a line")
	}

	return l.addRequest.Execute(app.AddRequest{
		Website: l.website,
		Entry:   entry,
	})
}
