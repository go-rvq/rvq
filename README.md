# RVQ (de Rebeca)

Rebeca (do hebraico, significa mulher que prende os homens com sua grande beleza, ou união, aquela que une רִבְקָה, 
transl. Rivqá). Uma homenagem à minha então filha caçula.

Este é a unificação de [qor5/web](https://github.com/go-rvq/web), [qor5/x](https://github.com/go-rvq/x) e [qor5/admin](https://github.com/go-rvq/admin).
com minhas alterações.

## Documentação

- [Índice da documentação](docs/README.md) — visão geral e navegação.
- [Referência dos pacotes do admin](admin/docs/reference.md) — presets,
  packages (perms, …) e features (activity, …).

> Os docs do QOR5 em `admin/docs` são de uma versão antiga (upstream) e serão
> atualizados posteriormente; prefira os docs por pacote referenciados acima.

## Testes

- Go: `go test ./...`.
- **Testes de integração de UI: `js/integration_tests/`** — rodam com `bun test`
  a partir do diretório da suíte (é lá que fica o `bunfig.toml` com o `preload`),
  p.ex. `cd js/integration_tests/admin/presets && bun test`. Cada suíte sobe o
  servidor Go de fixtures (SQLite em memória, semeado) e o dirige por HTTP;
  os `*.dom.test.ts` ainda montam o Vue de verdade (happy-dom) e clicam nos
  componentes.