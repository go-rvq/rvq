# history — histórico de revisões (git-like) para models

Plugin genérico de histórico para o admin `go-rvq/rvq`: cada save de um model
ativado grava uma **revisão** dos campos versionados, com hash de conteúdo,
autor, cadeia pai→filho, tag de publicação e contador de acessos. Fornece
comparação (diff por campo, inclusive HTML/JSON), navegação por campo e
**revert** total, por campo e parcial (por trechos), tudo pela UI do admin.

Import: `history "github.com/go-rvq/rvq/admin/packages/history/admin"`
Modelos: `github.com/go-rvq/rvq/admin/packages/history/models`.

---

## Modelo de dados

Cada model versionado ganha **sua própria** tabela `<tabela>_revisions` (ex.
`posts_revisions`) — sem coluna `ModelName`. Uma única struct Go `Revision` é
mapeada para cada tabela por `db.Table(name)`.

`models.Revision`:

| Campo | Tipo | Papel |
|-------|------|-------|
| `Hash` | `models.Hash` (bytea, PK) | sha256(recordKey + snapshot canônico). Dedup: estados iguais → mesmo hash, sem nova revisão. É a PK e o id na URL (hex, via `PrimarySlug`). |
| `RecordKey` | string (index) | PK do registro serializada (ex. `12_pt-BR`). |
| `Parent` | `models.Hash` | hash da revisão anterior do registro (cadeia/DAG); NULL na 1ª. |
| `Fields` | JSONB | snapshot `{campo: valor}` **completo** dos campos versionados. |
| `ChangedFields` | JSONB `[]string` | fields que mudaram vs o pai (todos na 1ª). Torna o histórico consultável por campo. |
| `CreatedAt` | time | quando. |
| `CreatorID` / `Creator` | uuid / string | autor (nunca nil: usuário logado ou AnonymousID). |
| `Tag` / `Published` | string / bool | a revisão publicada (o "git tag"). |
| `AccessCount` / `LastAccess` | int64 / time | contador de acessos por revisão (público). |

`models.Hash` é `[]byte` com `String()` (hex), `Parse`, e `Value`/`Scan` (bytea).

---

## Uso

### 1. Instalar o plugin (uma vez por admin)

```go
history.Configure(p, db) // p é o *presets.Builder
```

Instala o `Recorder` (contador de acessos), liga o diff do `activity_log` às
revisões e registra as mensagens i18n.

### 2. Ativar por model

```go
history.New(db).
    Model(postM).            // *presets.ModelBuilder
    HTMLFields("Body").      // diff HTML
    WholeFields("Cover").    // revert só inteiro (sem trechos)
    Build()
```

`Build()` resolve a tabela `<tabela>_revisions`, migra-a, faz backfill do
`ChangedFields` em revisões legadas, e embrulha o save (create **e** update) do
model para capturar a revisão. Também monta o NestedModel de revisões
(`/<model>/{id}/revisions`) e o gancho de publicação.

### API fluente (`*ModelHistory`)

| Método | Efeito |
|--------|--------|
| `Model(mb)` | model a versionar (obrigatório). |
| `Fields("A","B")` | campos versionados explícitos (default: os do mode EDIT). |
| `AllFields()` | versiona todos os campos do schema. |
| `ExtraFields("Seo")` | adiciona campos além do conjunto resolvido — úteis para valores editados **fora** do form (ex. um SEO gerido por nested model). |
| `HTMLFields("Body")` | trata como HTML (diff e revert parcial por bloco). |
| `WholeFields("Cover")` | só revert inteiro (nunca por trechos) — para FK/estruturados. |
| `ReferenceFields("Author")` | força tratamento de referência (ver Diff); relações são autodetectadas do schema gorm. |
| `FieldDiffHandler(field, fn)` | handler de diff próprio para um campo (ver Diff). |
| `Fetcher(fn)` | como recarregar o registro antes do snapshot (default: o fetcher de Editing do model, que aplica os preloads). |
| `Build()` | aplica tudo. |

---

## Captura

A cada save (create ou update), `capture`:
1. recarrega o registro pelo `Fetcher` (preloads) — evita associações nulas no snapshot;
2. serializa os campos versionados → snapshot + hash `sha256(recordKey + canonical(snapshot))`;
3. se já existe revisão com aquele (hash, recordKey), não faz nada (**dedup**);
4. calcula `ChangedFields` vs o pai e grava a `Revision` com `Parent` = revisão atual;
5. anota no request (`activity.WithRevisionRef`) para o `activity_log` referenciar a revisão em vez de duplicar o diff.

---

## Diff

O diff do Model é a **soma** dos diffs dos campos alterados (um painel por
campo). Cada campo é diffado por um **handler**; o padrão é escolhido pelo valor,
e pode ser sobrescrito com `FieldDiffHandler(field, fn)`.

Handlers reutilizáveis (assinatura `FieldDiffFunc`):

| Handler | Quando (default) | O que faz |
|---------|------------------|-----------|
| `HTMLValueDiff` | `HTMLFields` | diff HTML do valor (top-level via DetailingBuilder). |
| `BoolDiff` | campo bool | ícones antes/depois. |
| `PrismDiff(lang)` | JSON (struct/slice/map) e texto puro | render com Prism (números de linha + realce), linhas add/removidas marcadas. `JSONDiff` = `PrismDiff("json")`. |
| `ReferenceDiff` | FK/relação (autodetectada) e M2M | compara por valor exato; belongs-to mostra o registro renderizado + id (chip); to-many lista os itens marcando adicionados/removidos. |
| `RenderedDiff` | default | diff do valor renderizado pelo DetailingBuilder. |

Consultas programáticas:
- `Diff(recordKey, aHash, bHash) []FieldChange` — campos que diferem.
- `FieldDiff(recordKey, field, aHash, bHash) (old, now)` — um campo.
- `FieldHistory(recordKey, field) []FieldRevision` — timeline de um campo (as revisões que o alteraram).
- `Chain(recordKey)` / `Revision(recordKey, hash)` — cadeia e leitura.

---

## Revert

Todo revert cria uma **nova** revisão (estilo `git revert`, sem reescrever
histórico) e salva pelo pipeline de edição (logga em activity + history). Pela
UI passa por uma **confirmação** que agrupa por campo e mostra, em tabela, as
alterações a reverter (+ adição / − remoção).

- `RevertRecord(obj, hash, ctx)` — todos os campos versionados.
- `RevertFields(obj, hash, fields, ctx)` — um subconjunto.
- `RevertField(obj, hash, field, ctx)` — um campo inteiro (sempre disponível).
- `RevertFieldPartial` / hunks — só trechos selecionados de um campo que aceita
  parcial (texto char-a-char; HTML por bloco). Fields em `WholeFields` ou
  reference/relação recusam parcial (`AcceptsPartial(field) == false`).

Na UI, o diff por trechos é **clicável**: cada região destacada alterna
seleção (super-highlight + ✓), com botão "marcar todas" e contador
selecionados/total.

---

## Contador de acessos

`Configure` cria um `Recorder` (contador em memória, flush em lote). No handler
público, conte um acesso da revisão publicada corrente:

```go
history.DefaultRecorder().Hit("posts_revisions", postID+"_"+locale)
```

`RevisionTableFor(db, &models.Post{})` dá o nome da tabela.

---

## Ligação com activity_log

Quando o `activity` está ativo, o log de edição **referencia** a revisão
(`RevisionTable`/`RevisionHash`) em vez de duplicar o diff; a view do log deriva
o diff da revisão (via `activity.RevisionDiffFunc`, setado por `Configure`).

---

## Arquivos

- `models/` — `Revision`, `Hash`.
- `admin/model.go` — API fluente, captura, backfill.
- `admin/builder.go` — plugin (`Configure`/`Install`), Recorder, linkage.
- `admin/diff.go`, `admin/field_diff.go`, `admin/compare.go` — diff e comparação.
- `admin/field_history.go` — histórico por campo.
- `admin/revert.go`, `admin/hunks.go`, `admin/revert_confirm.go` — revert.
- `admin/nested.go` — NestedModel de revisões, listing (coluna "Alterações",
  filtro de campos em árvore), detail por revisão.
- `admin/access.go` — Recorder. `admin/publish.go` — tag na publicação.
- `admin/messages.go` — i18n (en, pt-BR).

Ver `PLAN.md` para o desenho original e as decisões.
```
