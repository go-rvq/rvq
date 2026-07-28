# RVQ (de Rebeca)

Rebeca (do hebraico, significa mulher que prende os homens com sua grande beleza, ou união, aquela que une רִבְקָה, 
transl. Rivqá). Uma homenagem à minha então filha caçula.

Este é a unificação de [rvq/web](https://github.com/go-rvq/web), [rvq/x](https://github.com/go-rvq/x) e [rvq/admin](https://github.com/go-rvq/admin).
com minhas alterações.

## Documentação

- [Índice da documentação](docs/README.md) — visão geral e navegação.
- [Referência dos pacotes do admin](admin/docs/reference.md) — presets,
  packages (perms, …) e features (activity, …).

> Os docs do RVQ em `admin/docs` são de uma versão antiga (upstream) e serão
> atualizados posteriormente; prefira os docs por pacote referenciados acima.

## Testes

- **`make check`** — o portão antes de commitar: `gofmt -s`, build, os testes Go
  que passam e as suítes de integração de UI. Pacotes quebrados de longa data
  (docs/examples do RVQ, pagebuilder, seo, media, libvips, testes de `web` e
  `x/perm`) ficam de fora de propósito.
- Go: `go test ./...` (roda tudo, inclusive o que já está quebrado).
- **Testes de integração de UI: `js/integration_tests/`** — rodam com `bun test`
  a partir do diretório da suíte (é lá que fica o `bunfig.toml` com o `preload`),
  p.ex. `cd js/integration_tests/admin/presets && bun test`. Cada suíte sobe o
  servidor Go de fixtures (SQLite em memória, semeado) e o dirige por HTTP;
  os `*.dom.test.ts` ainda montam o Vue de verdade (happy-dom) e clicam nos
  componentes.