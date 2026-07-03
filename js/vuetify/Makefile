BUN ?= bun

.PHONY: dist
dist:
	$(BUN) install
	$(BUN) run build

.PHONY: dist-clean
dist-clean:
	rm -rf ../../x/ui/vuetify/vuetifyjs/dist

.PHONY: clean
clean:
	rm -rf node_modules

.PHONY: dev
dev:
	$(BUN) run dev
