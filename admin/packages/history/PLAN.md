# Plano — Plugin de histórico git-like para models (rvq/admin/packages/history)

## Context

Hoje `Page` e `Post` (hermon-cms) versionam só via `publish.Status`
(draft/online) + snapshots `PublishedPage`/`PublishedPost`; **não** há histórico
de edições nem contador de acessos. O objetivo maior é acelerar os templates
gadx (index.gadx, post_list pág.1) cacheando "N posts mais acessados + N mais
recentes" e todas as `LocaleMessage`. Mas "mais acessados" não tem fonte de
dados — precisa de um contador. E o contador deve ser **por revisão** do
conteúdo, não por registro: uma revisão que mudou muito o texto começa a contar
do zero.

Este plano cobre a **fundação**: um plugin genérico de histórico (estilo git)
em `rvq/admin/packages/history`, aplicável a qualquer model, ativado em `Page` e
`Post`. Ele fornece: revisões com hash de conteúdo, "tag" na publicação,
contador de acessos por revisão, diff (inclusive HTML) e revert total/parcial. O
cache dos templates (a tarefa original) vira a **fase 2**, consumindo o contador
por revisão que este plugin cria.

## Modelo de dados (uma tabela por model, não-invasiva)

Cada model versionado ganha **sua própria** tabela de revisões, nomeada
`<tabela_do_model>_revisions` — ex. `posts_revisions`, `pages_revisions`. Não
altera a PK do model (diferente de `publish.Version`, que entra na PK). Como a
tabela já é dedicada a um model, **não há campo `ModelName`**.

Uma única struct Go `Revision` é mapeada para cada tabela por
`db.Table("<tabela>_revisions")` (AutoMigrate e queries sempre com `.Table(name)`):

```
Revision struct {
  Hash        []byte  // PK — bytea (sha256, 32B) do snapshot canônico dos campos (git blob-like); dedup automático
  RecordKey   string  // PK do registro serializada — ex. "12|pt-BR" (ID + LocaleCode)   (index: RecordKey)
  Parent      []byte  // bytea — hash da revisão anterior desse registro (cadeia/DAG git-like); NULL na 1ª
  Fields      JSON    // snapshot {campo: valor} só dos campos versionados
  CreatedAt   time.Time
  CreatorID   uuid.UUID  // ID do autor (gorm type:uuid, index; FK users.id) — nunca nil: usuário logado ou AnonymousID
  Creator     string     // nome de exibição do autor, quando disponível
  // "tag" = publicação
  Tag         string  // nome da tag/versão publicada (vazio = revisão comum, não publicada)
  Published   bool    // marca a revisão que foi publicada (o "git tag")
  // contador de acessos — por revisão
  AccessCount int64
  LastAccess  time.Time
}
```

- **Nome da tabela**: `<mb.TableName()>_revisions`, derivado do model ativado
  (ex. `posts` → `posts_revisions`). Uma tabela por model, sem `ModelName`.
- **Hash de conteúdo**: `sha256(canonical(Fields))` guardado como **bytea (32
  bytes)**, PK. Dois estados idênticos → mesmo hash (sem duplicar); salvar sem
  mudança nos campos versionados não cria revisão nova. A UI/rotas exibem o hex
  curto (ex. primeiros 7 bytes) — encode/decode hex só na borda de apresentação.
- **Campos versionados**: default = campos do mode **EDIT** (`WRITE`) do
  `ModelBuilder`; ou lista explícita; ou todos. (conforme pedido.)
- **Tag na publicação**: no callback de publish, a revisão atual do registro
  recebe `Published=true` + `Tag` (nome da versão). Um ponteiro "revisão
  publicada corrente" por `RecordKey` é a última com `Published`.
- **Contador por revisão**: um request público resolve o registro para sua
  **revisão publicada corrente** e incrementa `AccessCount` dela (batch em
  memória, flush periódico — não bloqueia o request).

## Usuário Anônimo estático (pré-requisito da autoria)

Hoje "anônimo" = `GetCurrentUser(r)` retornar `nil`
(`packages/user/helper.go:12`). Isso não serve para `CreatorID`/`UserID` com FK:
uma revisão (ou `activity_log`) criada sem usuário na sessão ficaria com autor
nulo/quebrado. Muda-se para um **usuário estático Anônimo, com UUID fixo e
imutável**:

- `packages/user`: constante `AnonymousID uuid.UUID` (imutável) e um valor
  estático do Anônimo. A interface `User` ganha `Anonymous() bool`
  (`GetID() == AnonymousID`) — a comparação é com o UUID estático.
- `GetCurrentUser(r)` passa a **nunca** retornar `nil`: sem usuário na sessão,
  devolve o Anônimo estático. Assim, ao **salvar um model** ou **gerar um
  activity_log**, sempre há usuário no Request, mesmo que `Anonymous()`.
- **Persistência para a FK**: a linha do Anônimo é semeada no BD no migrate (para
  `revisions.CreatorID` / `activity_logs.UserID` referenciarem algo real). É
  **não-editável e não-excluível**: no admin, `SetDeletingDisabled(true)` +
  campos read-only + filtro que a protege/esconde; (opcional) trigger/constraint
  no BD reforçando.
- **Username traduzível, sem ir ao banco**: o nome do Anônimo é resolvido **por
  request** conforme o idioma (i18n) e setado no valor estático em contexto — o
  Request **não** consulta o BD por ele (é estático). `GetName()` do Anônimo
  devolve o rótulo traduzido do request.
- Arquivos: `packages/user/api.go` (interface + `Anonymous()`), novo
  `packages/user/anonymous.go` (UUID estático, valor do Anônimo, seeding, nome
  i18n), `helper.go`/`middlewares.go` (retornar o Anônimo no lugar de `nil`),
  `presets/record_stamp.go:308` (usar o Anônimo no branch hoje `user == nil`).

Com isso a `Revision` sempre grava `CreatorID` válido (o do usuário logado ou
`AnonymousID`), e a FK `revisions.CreatorID → users.id` nunca quebra.

## O plugin `rvq/admin/packages/history`

Segue o padrão de `packages/people` (Builder que é `presets.Plugin`):

- `models/revision.go` — a struct `Revision` acima (sem `TableName` fixo; a
  tabela é resolvida por model via `db.Table(...)`).
- `admin/builder.go` — `Builder` com `Install(pb)` e `Configure(b, db,
  customize...)` → `b.Use(pb)`. Guarda a lista de models ativados para o admin.
- `admin/model.go` — API de ativação por model:
  `history.New(db).Model(mb).Fields("A","B"...).Build()` (sem `Fields` = campos
  EDIT; `.AllFields()` = todos). Isso:
  - resolve o nome da tabela `<mb.TableName()>_revisions` e faz
    `db.Table(name).AutoMigrate(&Revision{})`.
  - envolve o `Editing` do model (callback pós-save) para **capturar** os campos
    versionados, calcular o hash e gravar a `Revision` (na tabela do model) com
    `Parent` = revisão anterior — reusando `activity.DiffBuilder` para saber se
    mudou.
  - registra no callback de `publish` para marcar a tag/`Published`.
- `admin/access.go` — `Recorder`: `Hit(table, recordKey)` incrementa em
  memória; goroutine faz flush em lote (`UPDATE <table> SET access_count = …`)
  a cada N segundos. Expõe consultas por tabela: `MostAccessed(table, n)` e
  `Latest(table, n)` (a fase 2 usa).
- `admin/diff.go` — comparação em dois níveis:
  - **entre revisões (todos os campos)**: reusa `activity.DiffBuilder` →
    `[]Diff{Field,Old,Now}`; campos HTML (ex. `Body`) com diff inline via
    `github.com/sergi/go-diff/diffmatchpatch` (já em rvq/go.mod v1.3.1, hoje
    indirect). UI reusa `activity.DiffComponent`/`diffTable`.
  - **de um único campo** (`FieldDiff(table, recordKey, field, hashA, hashB)`):
    extrai `Fields[field]` de duas revisões e diffa só ele — base da navegação e
    do revert de conteúdo de um campo.
- `admin/field_history.go` — histórico de **um campo específico**:
  `FieldHistory(table, recordKey, field)` percorre a cadeia `Parent` e devolve
  só as revisões em que aquele campo mudou (ex. todas as alterações de
  `Post.Body`), para navegar/comparar versão a versão daquele campo.
- `admin/revert.go` — aplicar revisão, sempre gravando uma **nova** revisão com
  `Parent` = atual (estilo `git revert`, sem reescrever histórico).
  Granularidades:
  - **total (registro)**: restaura todos os campos versionados da revisão.
  - **por campo (subset)**: restaura só os campos selecionados do registro.
  - **campo inteiro (conteúdo total)**: restaura todo o conteúdo de um campo
    específico a partir de uma revisão dele (ex. voltar `Post.Body` inteiro para
    uma versão anterior). **Sempre disponível**, inclusive para campos que não
    aceitam parcial.
  - **conteúdo parcial (hunks, dentro de um campo)**: restaura só **parte** do
    conteúdo de um campo (ex. alguns trechos de `Post.Body`) — computa os hunks
    entre o valor atual e o da revisão-alvo com `diffmatchpatch`
    (`PatchMake`/`PatchApply`), aplica só os hunks selecionados (como
    `git checkout -p`). **Só para campos marcados como aceitando parcial**
    (ver abaixo).
  - **Quais campos aceitam parcial**: vem de uma struct tag no Model
    (`history:"nopartial"` no campo) **ou** de override na config do field
    (`.WholeFields("Cover", "Config", …)` na API fluente). Default:
    FK/relacionais/estruturados (ex. `ConfigID`, `Cover`, associações,
    JSON/slice) são **whole-only** — parcial corromperia o valor; texto/HTML
    simples (ex. `Body`, `Summary`, `Title`) aceitam **parcial**. O
    `admin/model.go` resolve isso por campo ao registrar (tag → tipo → override),
    e a UI só oferece seleção de hunks quando o campo aceita parcial.
- `admin/compo.go` + `messages/` — no Detailing do model:
  - aba **"Histórico"**: lista de revisões (hash curto, autor, data,
    tag/publicada, acessos), comparar (diff de todos os campos) e revert
    total / parcial-por-campo. Ref. visual: `publish/version_compo.go`.
  - visão **por campo** (a partir do campo no Detailing, ex. `Body`): linha do
    tempo das revisões daquele campo (`FieldHistory`), comparação entre duas
    (`FieldDiff`), revert do **campo inteiro** (sempre) e — quando o campo aceita
    parcial — revert **parcial-de-conteúdo** por seleção de hunks.

## Ligação com o activity_log (sem duplicar o diff)

Quando o `activity` também está ativo no model, o `activity_log` **aponta para a
revisão** em vez de duplicar o conteúdo do diff:

- `activity.ActivityLog` ganha `RevisionHash []byte` (+ `RevisionTable string`
  para saber qual `<tabela>_revisions`) e os acessores na interface.
- O plugin de histórico é o único a criar a revisão no save. Quando presente,
  ele registra a revisão e faz o log de edição **referenciar** essa revisão
  (`RevisionHash`/`RevisionTable` setados) e **não** grava `ModelDiffs` (fica
  vazio). Para não gerar diff duas vezes, o histórico desliga o diff-log do
  activity para aquele model (o activity já tem flags `skip`) e cria/anota o log
  apontando para a revisão.
- Na renderização (log view), quando o log tem `RevisionHash`, o diff é
  **derivado da revisão** (parent → revisão) via `admin/diff.go`, em vez de ler
  `ModelDiffs`. Assim o conteúdo mora só em `<tabela>_revisions`.

## Reuso (não reimplementar)

- `admin/activity/diff.go` — `NewDiffBuilder(mb).Diff(old,now) []Diff{Field,Old,Now}`
  com handlers de tipo (time.Time, MediaBox). Base do diff de campos e da
  detecção de mudança.
- `admin/activity` autor: `UserID uuid.UUID`/`Creator string` obtidos do contexto
  via `CreatorContextKey`/`getCreatorFromContext` (activity.go). O plugin captura
  `CreatorID`/`Creator` da mesma fonte ao gravar a revisão — sempre via
  `GetCurrentUser(r)`, que passa a devolver o Anônimo estático em vez de `nil`
  (ver "Usuário Anônimo estático"). (Alternativa: `presets.RecordUser`/
  `RecordUserFinder` em `record_stamp.go`.)
- `admin/publish` — callbacks `WithPublishCallback`/`WithUnpublishCallback`
  (já usados em `hermon-cms/admin/page.go`) para disparar a tag na publicação.
- `admin/presets` — `mb.NewFieldsBuilder(WRITE, …)` / `FieldBuilders.HasMode(WRITE)`
  para enumerar os campos EDIT (default). Padrão do builder fluente igual ao
  `presets/fields/tiptap` (`New().Model(mb).Fields(...).Build(mode)`).
- Padrão de plugin: `packages/people/admin/{builder,admin}.go` (`Install`,
  `Configure`, `b.Use`, `GetModelByID`).

## Integração em hermon-cms

Em `admin/config.go`, depois de `configPage`/`configPost` (que devolvem os
`*presets.ModelBuilder`):

```go
history.Configure(b, db) // instala o plugin
for _, m := range []*presets.ModelBuilder{pageM, postM} {
    history.New(db).Model(m).Build() // cria/migra <tabela>_revisions; campos EDIT
}
```

Isso cria `pages_revisions` e `posts_revisions`. No `site`: no render público de
página/post, chamar `recorder.Hit("posts", recordKey)` (ou `"pages"`) — 1 linha
no handler, sem query extra.

## Fase 2 (consome esta fundação — fora deste plano, só o gancho)

- Generalizar o pg_notify/LISTEN de `hermon-cms/server/notify.go` para payload
  `{table, pk, op}` + registro de callbacks por model (hoje é só
  `PageChange{Op,ID,Locale}` no canal `hermon_pages`).
- Caches em memória alimentados por esses callbacks + pelo `Recorder`:
  N posts mais acessados (via `AccessCount` da revisão publicada em
  `posts_revisions`) + N mais recentes (sem duplicar), N do `.env` (default
  100); todas as `LocaleMessage`.
- index.gadx e post_list pág.1 leem do cache em vez de `DB(models.PublishedPost)(…)`.

## Decisões/def­aults assumidos (ajustáveis)

- Uma tabela por model `<tabela>_revisions` (sem `ModelName`), não a
  PK-versioning do `publish` (invasiva) — plugin plugável, "ID/hash mapeado" sem
  mexer na PK, e cada model isolado na sua tabela.
- Snapshot completo dos campos por revisão (não delta) — diff é calculado sob
  demanda; simples e robusto para revert.
- Contador incrementa a **revisão publicada corrente** do registro (é o que o
  público vê); flush em lote assíncrono.
- Revert cria nova revisão (não reescreve histórico), como `git revert`.

## Verificação

- `go build ./...` em rvq e hermon-cms; `gofmt`/`go vet`.
- Testes de unidade no plugin: nome da tabela `<tabela>_revisions`; hashing
  determinístico + dedup; cadeia `Parent`; captura só dos campos EDIT vs
  `Fields(...)` vs `AllFields()`; `Recorder` acumula e faz flush; diff de campos
  e HTML; `FieldHistory`/`FieldDiff` de um campo; resolução do "aceita parcial"
  (tag `history:"nopartial"` / tipo FK-relacional / override `.WholeFields`);
  revert total, parcial-por-campo, campo-inteiro e parcial-de-conteúdo (hunks)
  geram a revisão esperada; campo whole-only recusa hunks. Referência:
  `activity/diff_test.go`, `publish/publish_test.go`.
- Manual no admin: editar um Post → nova revisão em `posts_revisions`; publicar →
  vira tag/Published; acessar a página pública → `AccessCount` sobe; comparar
  duas revisões (diff HTML); navegar o histórico só do `Body`, comparar duas
  versões dele e reverter apenas alguns trechos (hunks) → nova revisão só com
  aquela parte do `Body` revertida.
