# RVQ (from Rebeca)

Rebeca (from Hebrew רִבְקָה, transl. Rivqá — she who ties men with her great
beauty, or union, she who unites). Named after my youngest daughter.

This unifies [rvq/web](https://github.com/go-rvq/web),
[rvq/x](https://github.com/go-rvq/x) and
[rvq/admin](https://github.com/go-rvq/admin), with my own changes.

## Documentation

- [Documentation index](docs/README.md) — overview and navigation.
- [Admin package reference](admin/docs/reference.md) — presets, packages
  (perms, …) and features (activity, …).
- [Developing a package](developer.md) — the naming a package must not take
  from the application: table prefix, registration id, menu key and the
  migration a rename owes the database.

> The RVQ docs under `admin/docs` are from an old (upstream) version and will be
> updated later; prefer the per-package docs linked above.

## Code Conventions

- **Everything written for developers is in English** — code comments, doc
  comments, `docs/**.md`, `README.md`, test names and test comments.
- **What the end user reads is translated.** Every i18n module ships
  `Messages_en_US` and `Messages_pt_BR`, both registered on the `i18n.Builder`:

  ```go
  i18nB.RegisterForModule(language.English, I18nSomethingKey, Messages_en_US).
      RegisterForModule(language.BrazilianPortuguese, I18nSomethingKey, Messages_pt_BR)
  ```

  Other languages (zh-CN, ja-JP) come from upstream: keep the ones that are
  there, do not add new ones.

  A language's module messages are stored **whole**, with no per-field fallback —
  a partial translation renders as blanks, not as English. When adding a field to
  a `Messages` struct, fill it in for every language that module already has.
- **Idiomatic Go**, with errors wrapped as `fmt.Errorf("...: %w", err)`.

## Tests

- **`make check`** — the gate before committing: `gofmt -s`, build, the Go tests
  that pass, and the UI integration suites. Long-broken packages (RVQ
  docs/examples, pagebuilder, seo, media, libvips, the `web` and `x/perm` tests)
  are deliberately left out.
- Go: `go test ./...` (runs everything, including what is already broken).
- **UI integration tests: `js/integration_tests/`** — run them with `bun test`
  from the suite's own directory (that is where `bunfig.toml` holds the
  `preload`), e.g. `cd js/integration_tests/admin/presets && bun test`. Each
  suite boots the Go fixture server (in-memory SQLite, seeded) and drives it over
  HTTP; the `*.dom.test.ts` ones also mount the real Vue app (happy-dom) and
  click through the components.
