YARN := corepack yarn

.PHONY: ci
ci: install lint typecheck build

.PHONY: install
install:
	$(YARN) install --immutable

.PHONY: dev
dev: install
	$(YARN) dev

.PHONY: lint
lint:
	$(YARN) lint

.PHONY: typecheck
typecheck:
	$(YARN) typecheck

.PHONY: build
build:
	$(YARN) build
