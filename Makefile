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

.PHONY: log-moooodotfarm
log-moooodotfarm:
	mkdir -p _logs
	scp server:/var/log/nginx/moooodotfarm.access.log ./_logs/moooodotfarm.access.log

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf ./build
