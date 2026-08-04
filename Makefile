all: test build

build:
	mkdir -p build
	go build -o ./build/plum ./cmd/plum

build-race:
	mkdir -p build
	go build -race -o ./build/plum ./cmd/plum

frontend:
	./tools/build_frontend.sh

test:
	go test ./...

test-verbose:
	go test -v ./...

clean:
	rm -rf ./build

.PHONY: all build frontend build-race test test-verbose clean
