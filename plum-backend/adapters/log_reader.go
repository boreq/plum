package adapters

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/boreq/errors"
	"github.com/boreq/plum/plum-backend/logging"
	"github.com/nxadm/tail"
)

type LineHandler interface {
	Handle(line string) error
}

type LogReader struct {
	log logging.Logger

	lines      int
	errors     int
	statsMutex sync.Mutex
}

func NewLogReader() *LogReader {
	return &LogReader{
		log: logging.New("adapters/log_reader"),
	}
}

func (r *LogReader) Load(globs []string, handler LineHandler) error {
	for _, glob := range globs {
		paths, err := filepath.Glob(glob)
		if err != nil {
			return errors.Wrapf(err, "could not process a glob pattern '%s'", glob)
		}

		for _, path := range paths {
			if err := r.loadFile(path, handler); err != nil {
				return errors.Wrapf(err, "could not load '%s'", path)
			}
		}
	}

	return nil
}

func (r *LogReader) Follow(path string, handler LineHandler) error {
	config := tail.Config{Follow: true, ReOpen: true}
	ta, err := tail.TailFile(path, config)
	if err != nil {
		return err
	}
	return r.processTail(ta, handler)
}

func (r *LogReader) GetStats() (lines int, errors int) {
	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()
	lines = r.lines
	errors = r.errors
	return
}

func (r *LogReader) loadFile(path string, handler LineHandler) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return r.processFile(f, handler)
}

func (r *LogReader) processFile(f *os.File, handler LineHandler) error {
	reader, err := r.getReader(f)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		r.addLine()
		if err := handler.Handle(scanner.Text()); err != nil {
			r.addError()
			r.log.Error("error processing a line", "err", err, "line", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (r *LogReader) getReader(f *os.File) (io.Reader, error) {
	if strings.HasSuffix(f.Name(), ".gz") {
		reader, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	return f, nil
}

func (r *LogReader) processTail(ta *tail.Tail, handler LineHandler) error {
	for line := range ta.Lines {
		r.addLine()
		if err := handler.Handle(line.Text); err != nil {
			r.addError()
			r.log.Error("error processing a line", "err", err, "line", line.Text)
		}
	}
	return nil
}

func (r *LogReader) addLine() {
	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()
	r.lines++
}

func (r *LogReader) addError() {
	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()
	r.errors++
}
