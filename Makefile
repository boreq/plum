.PHONY: all
all: test build

.PHONY: dev
dev: test lint build

.PHONY: ci
ci: ci-backend ci-frontend

.PHONY: ci-backend
ci-backend:
	$(MAKE) -C plum-backend ci

.PHONY: ci-frontend
ci-frontend:
	$(MAKE) -C plum-frontend ci

.PHONY: lint
lint:
	$(MAKE) -C plum-backend lint

.PHONY: build
build:
	$(MAKE) -C plum-backend build

.PHONY: build-race
build-race:
	$(MAKE) -C plum-backend build-race

.PHONY: fmt
fmt:
	$(MAKE) -C plum-backend fmt

.PHONY: test
test:
	$(MAKE) -C plum-backend test

.PHONY: test-verbose
test-verbose:
	$(MAKE) -C plum-backend test-verbose

.PHONY: benchmark
benchmark:
	$(MAKE) -C plum-backend benchmark LOG=$(abspath $(LOG))

.PHONY: frontend
frontend:
	./tools/build_frontend.sh

.PHONY: log
log:
	mkdir -p _logs
	scp server:/var/log/nginx/blender.access.log ./_logs/blender.access.log
	scp server:/var/log/nginx/moooodotfarm.access.log ./_logs/moooodotfarm.access.log
	scp server:/var/log/nginx/drive.access.log ./_logs/drive.access.log

.PHONY: clean
clean:
	$(MAKE) -C plum-backend clean
