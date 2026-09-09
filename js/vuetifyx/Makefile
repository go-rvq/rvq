BUN ?= bun

.PHONY: dev
dev:
	$(BUN) run dev

.PHONY: install
install:
	$(BUN) install

.PHONY: deps
deps:
	BUN=$(BUN) bash scripts/prepare-deps.sh

.PHONY: build
build:
	$(BUN) install
	$(MAKE) deps
	$(BUN) run build

.PHONY: dist-clean
dist-clean:
	rm -rf dist

.PHONY: clean
clean:
	rm -rf node_modules
