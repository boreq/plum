package adapters_test

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/boreq/plum/plum-backend/adapters"
	"github.com/boreq/plum/plum-backend/config"
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
)

var benchmarkLogPath = flag.String("benchmark-log", "", "path to an access log in the combined format which the benchmarks are run against")

const benchmarkLines = 500000

type insertingHandler struct {
	parser     *parser.Parser
	repository *domain.Repository
}

func (h insertingHandler) Handle(line string) error {
	entry, err := h.parser.Parse(line)
	if err != nil {
		return err
	}
	return h.repository.Insert(entry)
}

func BenchmarkLoadOldEntries(b *testing.B) {
	benchmarkLoad(b, false)
}

func BenchmarkLoadRecentEntries(b *testing.B) {
	benchmarkLoad(b, true)
}

func benchmarkLoad(b *testing.B, fromEnd bool) {
	path, lines := benchmarkLog(b, fromEnd)

	p, err := parser.NewParser(parser.PredefinedFormats["combined"])
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler := insertingHandler{
			parser:     p,
			repository: domain.NewRepository(config.Website{}),
		}

		if err := adapters.NewLogReader().Load([]string{path}, handler); err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	b.ReportMetric(float64(lines)*float64(b.N)/b.Elapsed().Seconds(), "lines/s")
}

func benchmarkLog(b *testing.B, fromEnd bool) (string, int) {
	b.Helper()

	if *benchmarkLogPath == "" {
		b.Skip("pass -benchmark-log with a path to an access log")
	}

	source, err := os.Open(*benchmarkLogPath)
	if err != nil {
		b.Fatal(err)
	}
	defer source.Close()

	path := filepath.Join(b.TempDir(), "access.log")

	target, err := os.Create(path)
	if err != nil {
		b.Fatal(err)
	}
	defer target.Close()

	writer := bufio.NewWriter(target)
	defer func() {
		if err := writer.Flush(); err != nil {
			b.Error(err)
		}
	}()

	scanner := bufio.NewScanner(source)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	if fromEnd {
		tail := make([]string, 0, benchmarkLines)
		for scanner.Scan() {
			if len(tail) == benchmarkLines {
				tail = tail[1:]
			}
			tail = append(tail, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			b.Fatal(err)
		}

		for _, line := range tail {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				b.Fatal(err)
			}
		}
		return path, len(tail)
	}

	var lines int
	for scanner.Scan() && lines < benchmarkLines {
		if _, err := writer.WriteString(scanner.Text() + "\n"); err != nil {
			b.Fatal(err)
		}
		lines++
	}
	if err := scanner.Err(); err != nil {
		b.Fatal(err)
	}

	return path, lines
}
