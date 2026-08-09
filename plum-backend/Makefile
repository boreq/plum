.PHONY: all
all: test build

.PHONY: dev
dev: test lint build

.PHONY: ci
ci: tools dependencies tidy fmt check-repository-unchanged test lint build

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
	go install golang.org/x/tools/cmd/goimports@v0.45.0

.PHONY: dependencies
dependencies:
	go get ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: check-repository-unchanged
check-repository-unchanged:
	../tools/check_repository_unchanged.sh

.PHONY: build
build:
	mkdir -p build
	go build -o ./build/plum ./cmd/plum

.PHONY: build-race
build-race:
	mkdir -p build
	go build -race -o ./build/plum ./cmd/plum

.PHONY: fmt
fmt:
	goimports -w $(shell find . -name '*.go' -not -path './statik/*')

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: benchmark
benchmark:
ifndef LOG
	$(error pass a log to run the benchmarks against eg. `make benchmark LOG=../_logs/drive.access.log`)
endif
	go test ./domain/ ./adapters/ -run '^$$' -bench 'BenchmarkLoad|BenchmarkParse' -benchtime 1x -count=3 -benchmark-log $(abspath $(LOG))

.PHONY: clean
clean:
	rm -rf ./build
