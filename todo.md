# TASK — refresh pós-save (LISTING / DETAIL / EDIT / SINGLETON)

Quando um formulário é salvo com sucesso, só o que reflete a mudança é
atualizado — sem recarregar a página.

- [x] Mecanismo: **`onSaveCallbacks`** desce pelo escopo como `closer`/`form`
      (provide no `app.ts`, no localScope do `component-by-template`, global
      property e em todo `web.Plaid()`), e cada FormHost ACRESCENTA seu hook antes
      de repassar. O save responde com
      `(onSaveCallbacks || []).forEach(f => f("<id>"))` — nada vai na request.
      `ParamPostChangeCallback` deixou de ser necessário nesses fluxos.
- [x] LISTING (página ou diálogo) → NEW → save ⇒ listing recarrega.
- [x] LISTING → item DETAIL → EDIT → save ⇒ o DETAIL se re-renderiza no próprio
      portal e o LISTING que o abriu recarrega (só ele: é o hook que estava na
      lista daquele nível).
- [x] Regra 1: o DETAIL não depende de quem o abriu — ele se recarrega sozinho,
      então funciona aberto por qualquer coisa.
- [x] Regra 2: sem DETAIL, a linha abre o EDIT e o listing recarrega igual.
- [x] Regra 3: SINGLETON com DETAIL → o save do EDIT re-renderiza o DETAIL; sem
      DETAIL só existe a página de EDIT, que se re-renderiza.
- [x] Regra 4: nada recarrega a página. O corpo do DETAIL em página vive no portal
      `presets_detail_page`; `actions.ReloadDetail` re-renderiza esse portal e
      devolve `pageTitle`, que o cliente escreve no `<title>`.
- [x] Testes: `refresh.test.ts` (8, contrato) e `refresh.dom.test.ts` (4, UI real
      com clique/preenchimento/save) + fixture `integration/refresh_app.go`
      (Article com detailing, Tag sem detailing, Site singleton+detailing, Pref
      singleton só edit). Suíte total: 49 pass.
- [x] Documentado em `admin/presets/docs/post-save-refresh.md` com diagramas
      (sequência para NEW e para o DETAIL em página, estado para DETAIL→EDIT e
      para os singletons).
- [x] `make check` criado (gofmt + build + testes Go que passam + suítes bun) e
      registrado no CLAUDE.md/README como portão antes do commit.

# FORM HOST (evolução do EditForm/NewForm)

- [x] Conceito do usuário generalizado em `presets.FormHost` (`form_host.go`): a
      PÁGINA hospeda o form — declara uma var de estado que É o closer do overlay
      e guarda o bloco do form com `v-if='<ref>?.show'`. Ligar monta o bloco, cujo
      run-script carrega o form direto no portal do host (UM request, sem o
      round-trip do wrapper); desligar destrói o form e seu escopo por completo;
      qualquer botão pode reabrir (form fresco).
- [x] O estado vive em `vars` (não em slot scope): na PÁGINA o botão vai para o
      app-bar via layout, FORA do componente do host — com slot scope dava
      `TypeError: editing is undefined`. Vars são alcançáveis de qualquer ponto.
      Nomes prefixados: `$presetsEditing`, `$presetsCreating`,
      `$presetsItemEditing`, `$presetsItemDetailing` (sem colisão com a app).
- [x] Aplicado a: detailing (evento E página/defaultPageFunc), listing (New) e
      ITENS do listing — UM host por ação para toda a listagem (não um por linha):
      a linha faz `vars.$presetsItemDetailing.id = "<id>"; ….show = true` e o host
      carrega aquele registro. Hosts publicados no ctx (`WithItemFormHosts`) e as
      linhas/row-menu os consultam, com fallback ao evento self-hosted fora de um
      listing.
- [x] `ParamCloserProvided` + `closerScope` no **Dialog e no Drawer**: quando o
      host provê o closer, o overlay NÃO cria um filho — fechar/salvar desliga o
      host e destrói o form. (O Drawer importa: na página o edit abre em
      RightDrawer.)
- [x] `formScope` (EditForm/NewForm) reimplementado sobre o FormHost (Show(true)),
      para os call sites que não podem hospedar (row menu fora de listing,
      publish, l10n, model_select, pagebuilder).
- [x] Fix: o bloco do form usa `go-plaid-scope` (FormInit), não um user-component
      aninhado — `<template v-slot>` dentro de outro NÃO é compilado pelo parser
      de runtime e o conteúdo renderizava inerte (os campos sumiam).
- [x] Testes (bun, 27 pass): `formhost.dom.test.ts` (EDIT/NEW abrir-destruir-
      reabrir, 1 request por abertura), `listing.dom.test.ts` + `listing.test.ts`
      (hosts por item), `detailing.test.ts` (evento E página hospedam igual),
      `editform.test.ts` atualizado.
- [x] FIX (colisão em ANINHAMENTO): o estado saiu de `vars` (global) e passou a ser
      VARIÁVEL DE ESCOPO. Em `vars`, cada listing/detailing vivo declarava o MESMO
      nome e guardava seu bloco com `v-if` nele: abrir uma linha do listing de
      dentro do dialog montava TAMBÉM o bloco do listing de baixo — 2 requests e
      2 overlays (medido: `Received: 2`). Regra do usuário: em `*PageFunc` o botão
      vai para a app bar (fora do host), então lá o host embrulha o BOTÃO
      (`WrapsOpener`, botão + bloco guardado juntos); nos eventos o botão já está
      dentro do host, que embrulha o conteúdo.
- [x] `FormHosts(children, hosts...)`: os N hosts de um listing viram UM
      user-component com N vars — hosts aninhados renderizariam
      `<template v-slot>` dentro de outro, que o parser de runtime não compila.
      Portais são compilados à parte, por isso dialog dentro de página funciona.
- [x] Consumidores não podem mais hardcodar o nome: `GetItemFormHosts` ganhou
      `New`/`OpenNewExpr` e há `WithDetailingEditHost`/`GetDetailingEditHost` +
      `ShowExprOr`.
- [x] Testes novos (37 no total): `nested_hosts.{test,dom}.test.ts` (dois listings
      vivos; falham com o código antigo) e `nestedmodels.dom.test.ts` + fixture
      `integration/nestedmodels_app.go` — 2 níveis de `helper.NewNestedSliceBuilder`
      (Survey → Places → Products), navegando pelo menu "..." da linha, com os 4
      diálogos EMPILHADOS. Stub de `visualViewport` no `setup.ts` para o v-menu.
- [ ] PENDENTE: validar no NAVEGADOR as demais telas (drawer/teleport reais).

# Base — refactor do form scope (commitado em 4d0de132 / e24e0ee4)

Estado (validar no navegador antes de publicar):

- [x] Novas ações `actions.EditForm`/`actions.NewForm` + handler `formScope`
      (`admin/presets/editing_form.go`): responde no `ParamTargetPortal` com
      `CloserScope(UserComponent(Portal(inner), RunScript(onclick)).Scope("form",
      Var("{$parent: form}")), true)` e `v-if='closer.show'`. O closer liga via
      setup (show=true → dispara o plaid interno de Edit/New) e destrói tudo em
      show=false. O plaid interno recebe `.scope({closer: closer})`.
- [x] Removidos (obsoletos): `Form.ScopeDisabled`, `ParamEditFormUnscoped`,
      `EditFormUnscoped`/`GetEditFormUnscoped`/`ctxEditFormUnscoped`, os 3
      `Set(ParamEditFormUnscoped)` no list editor, o `FormInit`/`Wrap` do
      `RespondToPortal` e do `formNew`. `Form.Component()` overlay não cria mais o
      escopo `form` (vem do wrapper; cria só locals/vars); página continua criando.
- [x] Call sites migrados p/ EditForm/NewForm: detailing, listing (Edit+New),
      row menu, publish, l10n, model_select, pagebuilder.
- [x] Fixture de teste movida p/ `admin/presets/integration` (`Product`, `Item`,
      `NewDB`, `NewApp`, `NewSeededHandler`) + `ServeEnv`/`ServeEnvMain`. O
      `admin/presets/tests/listeditor` importa a fixture (dot-import). Fetcher/
      listing dão `Preload("Items")` (Fetch+Search).
- [x] Fix: `EditingBuilder.FetchAndUnmarshal`/`UnmarshalForm` agora usam o
      FormData (`r.Form`) quando o request NÃO é multipart (urlencoded/query) —
      `ParseMultipartForm` faz fallback p/ `ParseForm` em `http.ErrNotMultipart`.
      Antes, um submit não-multipart era ignorado (campos vazios perdidos).
- [x] Testes de integração bun (`js/integration_tests/admin/presets/`), TUDO em
      bun (sobem o servidor Go com sqlite em memória, dirigem via HTTP e validam a
      variável reativa `form` lida dos v-assign):
      - `editform.test.ts` (4): EditForm→wrapper/closer/scope, follow do run-script
        → Edit dialog, re-render com erro, save com `closer.show=false`.
      - `nested.test.ts` (5): list editor aninhado em 4 níveis (L0→L1→L2→L3) —
        abertura com o `form` de todos os níveis, serialização do `form` completo,
        falha de validação no nível mais profundo + re-render preservando os
        demais, inclusão aninhada (`__new`), exclusão profunda (re-seed deleted) e
        exclusão de 1º nível persistida.
      Fixtures Go em `admin/presets/integration` (`NewApp`/`NewSeededHandler`,
      `NewNestedApp`/`NewNestedSeededHandler`), servidor seleciona via env `APP`.
      9 pass.
- [x] Teste com o Vue MONTADO no DOM a partir do clique, TUDO em bun+happy-dom:
      `editform.dom.test.ts` (1) e `nested.dom.test.ts` (3). O problema era o bun
      não ter loader SFC real (o `.vue` renderizava como path) — resolvido com um
      Bun plugin em `setup.ts` que compila os SFC com `@vue/compiler-sfc`; +
      happy-dom (GlobalRegistrator, disableSameOriginPolicy, ResizeObserver stub)
      + paths `@/` no `corejs/tsconfig.json`. Um `vx-dialog` stub renderiza os
      campos. Provam o fluxo reativo do clique (closer/run-script/dialog) e leem
      a variável `form` montada.
- [x] Fix (list editor com PK composta / sem campo "ID"): `ToComponentForEach`
      (render) e `setWithChildFromObjs` (unmarshal/validação) usavam o ModelInfo
      do PARENT para os itens — só funcionava por acaso quando o item tinha "ID",
      e panicava (`no such field: ID`) para PK composta. Agora ambos derivam o
      ModelInfo do item via `f.nested.Model().Info().ChildOf(...)`.
- [x] Testes m2m/PK-composta multinível (`composite.test.ts`, 5): Cart → Items
      (PK CartID,Sku) → Notes (PK CartID,Sku,Seq) — abertura com o `form` de cada
      nível composto, validação profunda, inclusão (`__new`), exclusão profunda
      (re-seed deleted) e exclusão de 1º nível persistida. Fixture
      `admin/presets/integration/composite_app.go` (env APP=composite).
      Suite bun total: 18 pass (HTTP + DOM + composite).
- [ ] PENDENTE: validar no NAVEGADOR (dialog/drawer teleport + closer) e commitar.
- [x] `pagebuilder` volta a compilar (`go build ./admin/pagebuilder/...` ok) — o
      erro de assinatura `defaultTemplateInstall`/`categoryInstall` já não existe.

# Porte

- [x] `admin/presets/integration` portado e verde (2026-07-25); entrou no
      `make check`. Defeitos reais que ele revelou: `setToObjNilOrDelete` sem
      chamador (deleções/ordenação do list-editor não eram aplicadas no save) e
      `SectionBuilder.DefaultElementUnmarshal` zerando campos que a section não
      edita.
- [ ] `admin/docs/docsrc/content/*` (4 pacotes) — **decisão pendente**: as
      páginas usam `theplant/docgo`, que fala `theplant/htmlgo`, enquanto os
      helpers do repo devolvem `go-rvq/htmlgo`. Não há como portar sem escolher
      entre (a) fork `go-rvq/docgo` (padrão já usado para htmlgo/web/x) ou
      (b) copiar o núcleo do docgo (~700 linhas, MIT) para `admin/docs/docgo`.
      Há um `//replace github.com/theplant/docgo => ../../docgo/` comentado no
      go.mod, sem checkout local correspondente.

# Avaliar Manualmente

- [ ] **Título do painel no overlay `Content`.** `OverlayMode.Overlayed()` é só
      dialog/drawer, então um form pedido com `overlay=Content` usa
      `BuildPage()`: o título vai para o `PageTitle` da página, que não existe
      numa resposta de portal. O painel direito do pagebuilder chega sem título.
      Decidir se o `Content` deve carregá-lo e em que forma (card com toolbar,
      ou container com cabeçalho) — é desenho de UI, não porte.
      Casos que cobram isso: `examples_admin` →
      `TestPageBuilderCampaign/{Campaign_My_Contents,CampaignContents_edit}`,
      que esperam "Editing MyContent 1" / "Editing CampaignContent 1" no portal.
      O conteúdo do form já chega correto (inclusive o autosave com
      `@change-debounced`); falta só o título.

# Nao prioritárias
- [ ] crie o componente /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/js/vuetify/src/components/JSONInput.vsx com API em /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/x/ui/vuetify/json-input.go
  (como number-input.go). Se este componente for `:readonly=true`, ele renderiza o conteudo em codemirror-json. ele adiciona em nested para array e object. quado seu valor é vazio, ele tem flag para iniciar com `{}` ou `[]`, mas o atributo `json-type=ARRAY|MAP` tiver definido, ele já inicia pelo type e não permite o usuario alterar o tipo do objeto raiz.
  crie exemplos de uso e documente o componente. o componente tambem tem um icone que permite alternar para o modo readonly (apenas quando nao tiver sido iniciado em readonly).