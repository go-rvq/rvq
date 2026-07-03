BUN ?= bun

.PHONY: dev
dev:
	$(BUN) run dev

.PHONY: install
install:
	$(BUN) install

.PHONY: build
build:
	$(BUN) install
	$(BUN) run build

.PHONY: dist-clean
dist-clean:
	rm -rf dist

.PHONY: clean
clean:
	rm -rf node_modules
