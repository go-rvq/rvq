# Verificação do que está saudável hoje. Rode `make check` ANTES DE COMMITAR.
#
# O repositório tem pacotes quebrados de longa data — o que sobrou das
# estruturas antigas `github.com/go-rvq/{admin,web,x}`, forks do qor5
# (docs/examples do admin, media, os que exigem libvips). Eles NÃO entram aqui:
# `check` cobre o que passa, para que uma falha signifique de fato uma
# regressão sua.

GO ?= go
BUN ?= bun

# Pacotes Go com testes que passam.
GO_TEST_PKGS := \
	./web/... \
	./x/exchange \
	./x/i18n \
	./x/perm \
	./x/packages/db-tools \
	./admin/presets \
	./admin/presets/tests/... \
	./admin/presets/gorm2op/... \
	./admin/pagebuilder \
	./admin/publish \
	./admin/example/integration \
	./admin/seo \
	./admin/packages/... \
	./thirdpart/...

# Pacotes que precisam apenas compilar: sem testes próprios, ou com testes ainda
# presos à API antiga do fork (admin/presets/integration), mais as fixtures dos
# testes de integração.
GO_BUILD_PKGS := \
	./web/... \
	./admin/l10n/... \
	./admin/login/... \
	./admin/microsite/... \
	./admin/docs/docsrc/examples/examples_presets \
	./admin/pagebuilder/... \
	./admin/seo/... \
	./admin/presets \
	./admin/presets/examples \
	./admin/presets/integration \
	./admin/packages/... \
	./thirdpart/... \
	./x/... \
	./js/integration_tests/...

# Suítes de integração de UI (bun). Cada uma roda no seu diretório, onde está o
# bunfig.toml com o preload que os testes de DOM precisam.
BUN_TEST_DIRS := js/integration_tests/admin/presets

.PHONY: check check-fmt check-go check-ui

check: check-fmt check-go check-ui
	@echo "✓ check ok"

check-fmt:
	@echo "== gofmt =="
	@out=$$(gofmt -s -l admin web thirdpart x 2>/dev/null); \
	if [ -n "$$out" ]; then echo "arquivos não formatados:"; echo "$$out"; exit 1; fi

check-go:
	@echo "== go build =="
	$(GO) build $(GO_BUILD_PKGS)
	@echo "== go test =="
	$(GO) test $(GO_TEST_PKGS)

check-ui:
	@echo "== bun test (integração de UI) =="
	@for d in $(BUN_TEST_DIRS); do \
		echo "-- $$d"; \
		(cd $$d && $(BUN) test) || exit 1; \
	done
