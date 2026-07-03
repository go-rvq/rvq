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
	rm -rf ../../x/ui/ckeditor5/ckeditor5-js/dist

.PHONY: clean
clean:
	rm -rf node_modules

.PHONY: dev
dev:
	$(BUN) run dev
