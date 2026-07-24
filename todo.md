# Em andamento — refactor do form scope (EditForm/NewForm)

Estado (NÃO commitado ainda; validar no navegador antes do commit):

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
- [ ] `pagebuilder` tem erro de build PRÉ-EXISTENTE (assinatura
      `defaultTemplateInstall`/`categoryInstall`), não relacionado a este refactor.

# Nao prioritárias
- [ ] crie o componente /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/js/vuetify/src/components/JSONInput.vsx com API em /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/x/ui/vuetify/json-input.go
  (como number-input.go). Se este componente for `:readonly=true`, ele renderiza o conteudo em codemirror-json. ele adiciona em nested para array e object. quado seu valor é vazio, ele tem flag para iniciar com `{}` ou `[]`, mas o atributo `json-type=ARRAY|MAP` tiver definido, ele já inicia pelo type e não permite o usuario alterar o tipo do objeto raiz.
  crie exemplos de uso e documente o componente. o componente tambem tem um icone que permite alternar para o modo readonly (apenas quando nao tiver sido iniciado em readonly).