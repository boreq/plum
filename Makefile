.PHONY: all
all: test build

.PHONY: build
build:
	mkdir -p build
	go build -o ./build/plum ./cmd/plum

.PHONY: build-race
build-race:
	mkdir -p build
	go build -race -o ./build/plum ./cmd/plum

.PHONY: frontend
frontend:
	./tools/build_frontend.sh

.PHONY: log
log:
	mkdir -p _logs
	scp server:/var/log/nginx/blender.access.log ./_logs/blender.access.log
	scp server:/var/log/nginx/moooodotfarm.access.log ./_logs/moooodotfarm.access.log
	scp server:/var/log/nginx/drive.access.log ./_logs/drive.access.log

.PHONY: fmt
fmt:
	goimports -w $(shell find . -name '*.go' -not -path './plum-frontend/*' -not -path './statik/*')

.PHONY: test
test:
	go test ./...

.PHONY: benchmark
benchmark:
ifndef LOG
	$(error pass a log to run the benchmarks against eg. `make benchmark LOG=./_logs/drive.access.log`)
endif
	go test ./core/ -run '^$$' -bench 'BenchmarkLoad|BenchmarkParse' -benchtime 1x -count=3 -benchmark-log $(abspath $(LOG))

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf ./build
