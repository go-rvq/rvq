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
| `Origin` | `models.Origin` (colunas `origin_*`) | de onde o autor fez a revisão: IP, navegador e — quando a aplicação localiza endereços (`SetOriginFunc`) — país, região, cidade e coordenadas aproximadas. Vazio no seed. |
| `Event` | string | o que a revisão registra além de uma edição: `deleted` (`EventDeleted`), a exclusão do registro — o registro como estava, quem excluiu e de onde —; vazio numa edição. |
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

`Build()`, na ordem, faz:
1. resolve a tabela `<tabela>_revisions` e a migra (`AutoMigrate`);
2. **backfill** do `ChangedFields` em revisões legadas (anteriores à coluna);
3. **seed inicial** — para cada registro do model sem revisão, cria a primeira
   como espelho do registro atual (ver "Seed inicial");
4. embrulha o save (create **e** update) do model para capturar a revisão;
5. monta o NestedModel de revisões (`/<model>/{id}/revisions`) e o gancho de
   publicação.

Os passos 2 e 3 são idempotentes e best-effort (uma falha é logada, não bloqueia
o boot), então rodam a cada boot sem efeito depois de aplicados.

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
| `FieldContentHandler(field, lang, codec)` | codec de texto para um campo estruturado (ver Revert parcial de campo estruturado): faz um map JSON participar do revert parcial, editado como texto (ex. YAML). Registrar um codec habilita `AcceptsPartial` no campo. |
| `Fetcher(fn)` | como recarregar o registro antes do snapshot (default: o fetcher de Editing do model, que aplica os preloads). |
| `Build()` | aplica tudo. |
| `Record(obj, creatorID, creator, r)` | grava a revisão de um registro alterado **fora** dos forms do admin (um site, um job, um comando), com o autor e a origem (a requisição `r`; nil: nenhuma); `obj` carregado inteiro. Nada quando não mudou (dedup). |

`Build()` devolve o `*ModelHistory` — encadeie `FieldDiffHandler`/`FieldContentHandler` **depois** de `Build()`.

---

## Captura

A cada save (create ou update), `capture`:
1. recarrega o registro pelo `Fetcher` (preloads) — evita associações nulas no snapshot;
2. serializa os campos versionados → snapshot + hash `sha256(recordKey + canonical(snapshot))`;
3. se já existe revisão com aquele (hash, recordKey), não faz nada (**dedup**);
4. calcula `ChangedFields` vs o pai e grava a `Revision` com `Parent` = revisão atual;
5. anota no request (`activity.WithRevisionRef`) para o `activity_log` referenciar a revisão em vez de duplicar o diff.

A criação e a atualização são capturadas (o admin usa `Saver` e `Creator`
separados; ambos são embrulhados), então criar um registro já grava sua 1ª
revisão.

### Fora do admin: `Record`

Um registro que muda por outro caminho que não os forms do admin — o site que
deixa o autor editar o próprio comentário, um job — grava a sua revisão com
`Record`, depois de salvar, com o registro recarregado e o autor:

```go
mh := history.New(db).Model(commentM).Fields("Body").Build()

// no handler do site, depois de salvar
var c PostComment
db.First(&c, "id = ?", id)
mh.Record(&c, user.ID, user.Name, r) // r: a requisição, a origem da revisão
```

A revisão é a mesma de um save do admin: entra na cadeia, é deduplicada e
aparece no histórico do registro (`/<model>/{id}/revisions`).

### A origem de uma revisão

Cada revisão guarda de onde o autor a fez (`Revision.Origin`): o endereço — atrás
de um proxy reverso, o primeiro do `X-Forwarded-For` (ou `X-Real-Ip`) — e o
navegador, tirados da requisição do save (ou da passada a `Record`). Onde fica o
endereço, o history não sabe: a aplicação que localiza endereços (uma base GeoIP)
diz, trocando a função que faz a origem:

```go
history.SetOriginFunc(func(r *http.Request) histmodels.Origin {
    o := history.DefaultOrigin(r) // o endereço e o navegador
    loc := geo.Lookup(o.IP)       // da aplicação
    o.Country, o.CountryName, o.Region, o.City = loc.Country, loc.CountryName, loc.Region, loc.City
    return o
})
```

A listagem das revisões mostra a origem ("Viçosa, Minas Gerais, Brazil
(200.1.2.3)", o navegador ao passar o mouse).

A origem é do pacote `admin/origin` (`origin.Origin`, `origin.SetFunc`,
`origin.Of`): `SetOriginFunc` troca a função de todos que a usam — o history e
a lixeira (`perms`), que guarda de onde um registro foi excluído.

### Exclusões

A exclusão de um registro pelo admin (o `Delete` do data operator gorm2op) é
uma revisão própria — `Event` `deleted`, nunca deduplicada —: o registro como
estava antes, quem excluiu e de onde. A listagem das revisões a marca
("Exclusão"). Fora do admin — um site, um job —, `RecordDeletion(obj,
creatorID, creator, r)` grava a mesma revisão, com `obj` carregado inteiro.

---

## Seed inicial (BD já populado)

Ao ativar o history sobre um banco **que já tem registros** (anteriores à tabela
de revisões), esses registros não teriam baseline para comparar/reverter. No
boot, `SeedInitialRevisions()` (chamada por `Build`) cria, para cada registro
**sem nenhuma revisão**, a primeira como **espelho do registro atual** (root),
carregando as associações pelo `Fetcher` para um snapshot fiel. É idempotente
(registros que já têm revisão são pulados) e best-effort por registro. O autor
fica vazio (seed de sistema).

---

## Diff

O diff do Model é a **soma** dos diffs dos campos alterados (um painel por
campo). Cada campo é diffado por um **handler**; o padrão é escolhido pelo valor,
e pode ser sobrescrito com `FieldDiffHandler(field, fn)`.

Handlers reutilizáveis (assinatura `FieldDiffFunc`):

| Handler | Quando (default) | O que faz |
|---------|------------------|-----------|
| `PrismDiff(lang)` | HTML, JSON (struct/slice/map), texto puro e templates gad | **default de todo campo com forma de texto**. Render lado a lado com Prism (`json`, `html`/`markup`, `yaml`, `gad`/`gadt`/`gadx`, `css`, …): números de linha, linhas add/removidas marcadas e — estilo GoLand/IntelliJ — o **conteúdo alterado destacado dentro de cada linha** (char-a-char). `JSONDiff` = `PrismDiff("json")`. |
| `BoolDiff` | campo bool | ícones antes/depois. |
| `ReferenceDiff` | FK/relação (autodetectada) e M2M | compara por valor exato; belongs-to mostra o registro renderizado + id (chip); to-many lista os itens marcando adicionados/removidos. |

`PrismDiffComponents(lang, oldTxt, newTxt) (oldC, newC)` é o núcleo reutilizável:
recebe os dois textos já prontos (ex. YAML) e devolve os dois lados. Útil num
`FieldDiffHandler` para mostrar um campo com outra linguagem — ex. um map JSON
como YAML:

```go
mh.FieldDiffHandler("Data", func(in *history.FieldDiffInput) (o, n, m h.HTMLComponent, ok bool) {
    o, n = history.PrismDiffComponents("yaml", jsonToYAML(in.OldSnap["Data"]), jsonToYAML(in.NewSnap["Data"]))
    return o, n, nil, true
})
```

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
  parcial. Fields em `WholeFields` ou reference/relação recusam parcial
  (`AcceptsPartial(field) == false`) — **exceto** quando têm um codec registrado
  (ver abaixo), que sempre aceita.

Na UI, o diff por trechos é **clicável**:
- **HTML** (ex. `Post.Body`): `vx-diff-hunks`, hunks por bloco (`<p>…</p>`).
- **texto / campos com codec**: o painel **`vx-code`** (Prism) — cada trecho
  alterado é um segmento do slot `changed`, clicável, com ✓ + super-highlight ao
  selecionar. Botão "marcar todas" e contador selecionados/total.

O `vx-code` expõe o slot `changed` (`ChangedSlot`) com escopo
`{ text, line, start, end, hunk, kind, added }`, o que permite tornar cada trecho
interativo (é assim que o revert parcial monta os toggles de hunk).

### Revert parcial de um campo estruturado (map JSON como YAML)

Um campo cujo valor é um **map JSON** não é string, então por padrão é whole-only.
Um **`FieldContentCodec`** o adapta ao fluxo de texto do revert parcial:

- `ToText(rawJSON)` → texto editável (ex. YAML);
- `Apply(obj, text)` → escreve o campo de volta a partir do texto (ex. parse do
  YAML no map).

Registrado com `FieldContentHandler(field, "yaml", codec)`, o campo passa a
aceitar parcial: é diffado/editado como YAML (Prism + realce intra-linha), e ao
aplicar os hunks selecionados o texto é reparseado para o map e salvo (nova
revisão). Todo o fluxo (painel, preview e apply) fica codec-aware.

O exemplo completo (`history.ExampleYAMLConfigHistory`, em `admin/example.go`)
ativa history + diff YAML + codec num model singleton de configuração — está na
documentação docgo (grupo *Building Admin → History*, via o snippet
`HistoryYAMLConfigExample`).

### Singleton

Um model singleton (sem `{id}` na rota) é suportado: o comparador resolve a
`record_key` pelo próprio model (`GetSingleton()` + fetcher → `MustRecordID`),
não pelo path.

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
- `admin/model.go` — API fluente, captura, backfill e seed inicial.
- `admin/builder.go` — plugin (`Configure`/`Install`), Recorder, linkage.
- `admin/diff.go`, `admin/field_diff.go`, `admin/compare.go` — diff e comparação
  (`PrismDiff`/`PrismDiffComponents`, realce de linha e intra-linha).
- `admin/field_history.go` — histórico por campo.
- `admin/revert.go`, `admin/hunks.go`, `admin/hunks_code.go`,
  `admin/revert_confirm.go` — revert (whole/fields/field/parcial). `hunks.go` é o
  parcial HTML (`vx-diff-hunks`); `hunks_code.go` é o parcial em `vx-code` (texto
  e campos com codec) + os codecs de conteúdo.
- `admin/nested.go` — NestedModel de revisões, listing (coluna "Alterações",
  filtro de campos em árvore), detail por revisão, resolução de `record_key`
  (inclui singleton).
- `admin/access.go` — Recorder. `admin/publish.go` — tag na publicação.
- `admin/example.go` — exemplo documentado (YAML config), extraído por snippetgo.
- `admin/messages.go` — i18n (en, pt-BR).

Ver `PLAN.md` para o desenho original e as decisões.
```
