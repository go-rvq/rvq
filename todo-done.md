# Concluídas

## Associação / list editor (2026-07)

- [x] List editor: chave de form estável `__index` (desacoplada do índice do
      slice), classificação por flags `__new`/`__deleted` (sem checar ID),
      helper `presets.PartitionListEditorItems` (retorna struct ou `nil` se não
      iniciado), e reconciliação has-many em `gorm2op.SaveHasManyAssociation`.
- [x] List editor: não validar linhas marcadas `__deleted`
      (`FieldWalkHandleOptions.SkipListEditorDeleted`); linhas novas removidas
      somem (sem cache/vazamento entre formulários); o evento de adicionar marca
      `__new` na fonte que o render lê.
- [x] `gorm2op.SaveHasManyAssociation`: suporte a PK composta
      (`primaryFieldNames`/`pkAllZero`/`pkMapKey`), com guard para linha sem PK.
- [x] `helper.NestedSlice` (m2m) e `helper.ModelSelectorBuilder` (belongs-to):
      chaves de qualquer tipo (UUID) e compostas; link/unlink via gorm Association
      API; filtro do pai via `References` (`hasManyParentFilter`/`m2mParentFilter`).
- [x] Documentar `admin/packages/helper` (`doc.go`, `README.md`, `docs/keys.md`).
- [x] Refactor: consolidar as funções de leitura de schema/relacionamento gorm
      (PK/FK, composta, qualquer tipo) em `thirdpart/gorm/utils` como fonte única
      — `ForeignKeyFields`, `HasManyParentFilter`, `M2MParentFilter`,
      `ParentFilterArgs`, `PrimaryFieldNames`, `PKAllZero`, `PKMapKey`,
      `NormalizeKey`. Atualizados os call sites (`gorm2op`, `helper`) e os testes.
      Mesclado o pacote `gormutils` da app (RawColumn/SetRawColumn, AssociationDB,
      WithClauses) em `thirdpart/gorm/utils` e removido da app. Documentado em
      `thirdpart/gorm/utils/{README.md,docs/keys.md}`.
- [x] Testes: `admin/presets/tests/listeditor/` (UI + persistência + `__index` +
      helper + PK composta + m2m) e `admin/packages/helper/*_test.go` (m2m
      UUID/composta, FK composta belongs-to).
