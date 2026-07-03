BUN ?= bun

.PHONY: build
build:
	$(BUN) install
	$(BUN) run build

.PHONY: dist
dist:
	$(BUN) install
	$(BUN) run build


.PHONY: dist-clean
dist-clean:
	rm -rf ../../x/ui/vue-chart/vue-chart-js/dist

.PHONY: clean
clean:
	rm -rf node_modules

.PHONY: dev
dev:
	$(BUN) run dev
