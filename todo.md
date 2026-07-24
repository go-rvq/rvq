# Refactors

- [ ] Mover as funções helper de gorm para `RVQ/thirdpart/gorm/utils` (um pacote
      reutilizável, desacoplado de `admin/packages/helper`). Inclui as que criei
      até agora:
  - `helper.foreignKeyFieldsOf(db, model, field)` — nomes dos campos de FK
    (belongs-to) em ordem de PK do relacionado.
  - `helper.hasManyParentFilter(rel)` / `helper.m2mParentFilter(rel)` — WHERE de
    filtro do pai (has-many / m2m), com suporte a PK composta.
  - `helper.parentFilterArgs(id, fieldNames)` — extrai os valores da PK do pai
    para os placeholders do filtro.
  - `gorm2op` (assoc.go): `primaryFieldNames(schema)`, `pkAllZero(item, pkNames)`,
    `pkMapKey(item, pkNames)`, `normalizeKey(v)` — casamento por PK
    (possivelmente composta) na reconciliação de has-many.
      Objetivo: uma única fonte de verdade para leitura de schema/relacionamento
      gorm (PK/FK, composta, qualquer tipo), reutilizada por `helper`, `gorm2op`
      e demais consumidores. Atualizar os call sites e os testes.

# Nao prioritárias
- [ ] crie o componente /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/js/vuetify/src/components/JSONInput.vsx com API em /mnt/MPS-WORK/work/.goenv/ipc-vicosa/src/github.com/go-rvq/rvq/x/ui/vuetify/json-input.go
  (como number-input.go). Se este componente for `:readonly=true`, ele renderiza o conteudo em codemirror-json. ele adiciona em nested para array e object. quado seu valor é vazio, ele tem flag para iniciar com `{}` ou `[]`, mas o atributo `json-type=ARRAY|MAP` tiver definido, ele já inicia pelo type e não permite o usuario alterar o tipo do objeto raiz.
  crie exemplos de uso e documente o componente. o componente tambem tem um icone que permite alternar para o modo readonly (apenas quando nao tiver sido iniciado em readonly).