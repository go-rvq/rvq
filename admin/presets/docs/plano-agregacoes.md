# Plano: agregações do listing (totalizadores)

Estado: **implementado**, sem commit, com as recomendações aceitas (as
decisões no fim). A convenção está em `aggregates.md`. 2026-11-03.

## O que foi feito, e onde mudou do plano

- **Seletor em diálogo:**
  - `AdvancedSelect.vue` com a prop `dialog`: um `VXDialog` que expande e
    fecha; o "Trocar" é um ícone à direita do valor escolhido;
  - `dialogOnTop` põe o diálogo acima do menu do filtro, cujo `z-index` é
    100000;
  - o filtro `RecordPickerItem`;
  - no Go, `Dialog` / `DialogTexts`, `vx.ItemTypeRecordPicker` +
    `RecordPickerConfig`, e `helper.ModelFilterItem`;
  - o `ModelSelect` abre em diálogo por padrão (`SetSelectMode`).
- **Agregações:**
  - utils: `WithClause.DB`;
  - presets: `AggregateBuilder`, os ajudantes, `SearchResult.Aggregates` e
    os chips no portal da tabela;
  - gorm2op: a consulta única (`aggregate.go`).
- **A contagem mudou:** em vez de `COUNT(*) OVER ()`, é
  `(SELECT COUNT(*) FROM __agg_base)`. Assim ela conta os registros uma vez
  cada, mesmo com uma junção que os repete, e sai da mesma base dos totais.
- **As chaves dos filtros** têm o prefixo `f_` na URL (o `FilterDataFunc` o
  põe). `presets.FilterQueryValue(ctx, key, mod)` e `FilterKeyPrefix` os
  leem: o `When` de Movimentações e o `ModelFilterItem` usam esses dois.
- **Defeitos achados no caminho:**
  - a busca do `ModelSelect` usava a URL com `{parent_0_id}` literal; agora a
    URL vem com os pais da requisição;
  - o seed da testapp gravava `people.organizacao_id` nulo
    (`reflectutils.Set` ignora `uuid.UUID` → `model.UUID`).
- **O build do JS:** `bun run build` em `js/vuetifyx`, e o `dist` embutido no
  Go.

## O pedido

- Um totalizador, ou seja agregações, mostrado como chips acima da tabela do
  listing.
- Um componente integrado ao `ListingBuilder`, e não montado à parte por
  quem usa.
- Calculado com um WITH, pelo `thirdpart/gorm/utils.WithClauses`.
- O componente de exibição é configurável.
- Quando a lista é pedida serializada (o evento de registros, em JSON), as
  agregações nomeadas vão junto, com os seus valores.
- **A mesma consulta** da página traz os totais.
- Em Movimentações (finance), três agregações: **Total**, **Pagamentos**
  (a soma dos valores < 0) e **Recebimentos** (> 0).
- Isso vira uma convenção do rvq.
- **Filtros** no listing de Movimentações: Tipo, Projeto, Pessoa, Situação,
  Conta Contábil e Moeda. Com Tipo = pagamento, a agregação Recebimentos
  não entra na consulta (e o inverso).
- Os filtros de um modelo (Projeto, Pessoa, Conta Contábil, Moeda) usam um
  seletor **em diálogo**, como o seletor de Organização do menu, e não um
  select. O **`ModelSelect`** (o campo de Pessoa, Moeda… nos formulários)
  passa a usar o mesmo seletor no lugar do select de hoje.

## Hoje

- **A consulta da busca:** `ListingBuilder.search` monta os `SearchParams`
  (palavra-chave, `SQLConditions`, filtros, aba, página, ordem) e chama o
  `Searcher`. O padrão é o `gorm2op.Search`, que faz
  `Prepare(Search, …)`, roda os callbacks do modelo (as junções do
  Movimentacao, o escopo da organização, a lixeira) e então o `do`, com
  **duas** consultas:
  - a contagem, numa sessão marcada com `CountingConfigKey`, sem limite e
    sem ordem;
  - a página: a página, a ordem e o `Find`.
- **O resultado:** o `SearchResult{Records, TotalCount, Page, PerPage,
  OrderBys}`. O evento de registros (`records`) o devolve em `r.Data`, em
  JSON. Com `ParamMustResult` devolve só os `Records`; com
  `ParamListingEncoder`, os registros codificados.
- **A tela:** `GetTableComponents` monta a tabela (portal `DataTable`) e as
  adições (portal `DataTableAdditions`: a paginação). O `reloadList`
  (pesquisa, filtro, aba, página) atualiza os dois portais.
- **O `WithClauses`:** cada cláusula recebe a consulta como **texto**
  (`Query string` e `Args`).
- **O Valor da movimentação:** é calculado em Go (`models.ValorMovimentacao`):
  o valor fixado ou a soma das parcelas ativas (nem canceladas nem divididas).
  Em SQL não existe.
- **O scan do gorm (a v1.31.2 é a mais recente):** uma coluna que o modelo
  não tem é lida numa variável local e descartada. Não há opção nem hook para
  recuperá-la.
- **O `ModelSelect`** (`admin/packages/helper/model_select.go`) desenha um
  `VXSelectOne`, que é o `VXAdvancedSelect` (`vx-advanced-select`,
  `js/vuetifyx/src/lib/AdvancedSelect.vue`): o rótulo, os itens escolhidos
  (os `items` de partida, o cache dos rótulos) e um `v-autocomplete` que busca
  pelo `searchItemsFunc`, ou seja o evento `presets_ListData` do listing do
  modelo com `keyword` (o `ItemsSearcher`), só os `Records`.
- **Os filtros** são o componente Vue `vx-filter`
  (`js/vuetifyx/src/lib/Filter`), um tipo por item (`SelectItem`,
  `MultipleSelectItem`, `StringItem`…), montado de um JSON (`FilterData`,
  sem funções). O `SelectItem` é uma lista de opções fixa.
- **O seletor de Organização** (`orgs/admin/menu_root.go`) é uma action do
  listing desenhada no servidor: um diálogo com a pesquisa, a lista e a
  paginação.

## A prova de conceito (feita)

`thirdpart/gorm/utils/capture.go` e `capture_test.go`:

- **`RegisterColumnCapture(db)`** troca o callback `gorm:query` por um que,
  numa consulta marcada, envolve as linhas lidas (a interface `gorm.Rows`).
  Uma consulta sem marca roda a do gorm.
- **`CaptureColumns(db, nomes…)`** marca a consulta (`InstanceSet`: só ela,
  não os preloads). Depois de cada `Scan` do gorm, o envoltório lê do destino
  o valor das colunas pedidas, que o gorm descartaria.
- **`TestCaptureColumns`**, no SQLite e no PostgreSQL (`RVQ_TEST_POSTGRES`):
  numa consulta só, com `WITH base`, `LIMIT 2` e `Preload`, vieram:
  - a página com os seus 2 registros, com o `AfterFind` executado e o
    preload carregado;
  - `__count`, a contagem de todo o filtro, por `COUNT(*) OVER ()`;
  - `__aggs`, os totais em JSON, pela subconsulta do WITH.
- **No PostgreSQL** o JSON chega como `[]byte` dentro do `interface{}`: a
  captura o copia para texto, porque o driver pode reaproveitar o buffer.

## O desenho

### A convenção

Uma agregação é declarada no listing, com nome, rótulo, expressão SQL e,
opcionalmente, formato e componente:

```go
lb := mb.Listing()
lb.AggregateBase("valor", models.ValorMovimentacaoSQL("movimentacoes")) // a coluna da base
lb.Aggregate("total").
	SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).Total }).
	Expr(presets.Sum("valor"))
lb.Aggregate("pagamentos").
	Expr(presets.SumIf("valor", "valor < 0"))
lb.Aggregate("recebimentos").
	Expr(presets.SumIf("valor", "valor > 0"))
lb.AggregatesGroupBy("moeda_id") // opcional: um valor por grupo (decisão 1)
```

- **`AggregateBase(nome, expr)`:** uma coluna calculada por registro, na
  base (o WITH). Por padrão, a base tem só a chave primária.
- **`Aggregate(nome)`** devolve um `*AggregateBuilder`:
  - `Expr(sql)`: a função de agregação sobre as colunas da base;
  - `SetI18nLabel` / `Label`;
  - `Scan(T)`: o tipo do valor (`decimal.Decimal`, `int64`…), convertido do
    valor cru do driver;
  - `Format(func(ctx, v) string)`;
  - `ComponentFunc(func(ctx, *AggregateValue) h.HTMLComponent)`: o chip
    desta agregação;
  - `Color(func(v) string)`;
  - `Hidden()`: serializada, mas sem chip;
  - `Serialize(false)`: só na tela.
- **Os ajudantes de expressão:** `presets.Sum`, `presets.SumIf`,
  `presets.Count`, `presets.CountIf`, `presets.Avg`, `presets.Min` e
  `presets.Max`. Fazem `SUM(CASE WHEN cond THEN col END)`, válido no
  PostgreSQL, no SQLite e no MySQL; o `FILTER (WHERE …)` é só do PostgreSQL
  e do SQLite.
- **`ListingBuilder.AggregatesComponentFunc(func(ctx, []*AggregateValue)
  h.HTMLComponent)`** troca o componente inteiro, ou seja a barra de chips.

### A consulta: uma só, a da página

A página, a contagem e os totais saem da mesma consulta:

```sql
WITH base AS (
	SELECT DISTINCT movimentacoes.id, <valor> AS valor, movimentacoes.moeda_id
	FROM … <as junções, o escopo e os filtros da busca> …
)
SELECT movimentacoes.*, …,                          -- o que a página já seleciona
	COUNT(*) OVER () AS __count,                     -- a contagem do filtro
	(SELECT <json>(…) FROM (                         -- os totais do filtro
		SELECT moeda_id,
			SUM(valor) AS total,
			SUM(CASE WHEN valor < 0 THEN valor END) AS pagamentos,
			SUM(CASE WHEN valor > 0 THEN valor END) AS recebimentos
		FROM base [GROUP BY moeda_id]) g) AS __aggs
FROM movimentacoes …
WHERE movimentacoes.id IN (SELECT id FROM base)
ORDER BY … LIMIT … OFFSET …
```

- **Base:** é a consulta filtrada da busca (o `state.DB`, sem limite nem
  ordem), com `DISTINCT` pela chave. Uma junção que repete registros não os
  conta duas vezes. É o uso do `WithClauses`, pelo `WithClause.DB` (abaixo).
- **Contagem:** é uma função de janela (`COUNT(*) OVER ()`), calculada antes
  do `LIMIT`, e tira a consulta de contagem de hoje.
- **Totais:** uma subconsulta da base, num JSON (`json_build_object` no
  PostgreSQL, `json_object` no SQLite e no MySQL). O PostgreSQL a executa
  **uma vez** (`InitPlan`, visto no `EXPLAIN`), não por linha. Serve com e
  sem agrupamento.
  - Uma janela (`SUM(…) OVER (PARTITION BY moeda_id)`) não serve para
    agrupar: cada linha traria só a sua moeda, e uma moeda que não aparece
    na página sumiria dos totais.
- **Leitura:** pela captura (`CaptureColumns(db, "__count", "__aggs")`). O
  modelo continua o mesmo, com os seus hooks e os seus preloads.
- **Onde:** no `do` do `gorm2op.Search`. A sessão leva
  `AggregatingConfigKey`, como a contagem leva a `CountingConfigKey`.
- **Sem agregações declaradas:** a contagem também pode ir para a janela; é
  um ganho à parte (decisão 3).

### A página vazia

Sem linhas na página, nada é capturado:

- **Página 1 sem linhas:** o filtro está vazio. A contagem é 0 e os totais
  são zero (os `SUM` de nada são nulos, e são mostrados como 0). Não precisa
  de outra consulta.
- **Página além do fim** (página > 1 sem linhas: um link antigo, uma
  exclusão que encolheu a lista): uma consulta **só de totais**, sobre a
  mesma base e sem página, em **linhas nomeadas**:

  ```sql
  WITH base AS (…)
  SELECT 'total' AS name, moeda_id AS group_key, SUM(valor) AS value FROM base GROUP BY moeda_id
  UNION ALL
  SELECT 'pagamentos', moeda_id, SUM(CASE WHEN valor < 0 THEN valor END) FROM base GROUP BY moeda_id
  UNION ALL
  SELECT 'recebimentos', moeda_id, SUM(CASE WHEN valor > 0 THEN valor END) FROM base GROUP BY moeda_id
  ```

  O resultado é lido num `[]AggregateRow{Name, Group, Value}`, com
  `COUNT(*)` como a linha `__count`. Sem agrupamento, `group_key` é nulo.
  - Por que as linhas nomeadas não vão na consulta da página: num
    `UNION ALL` com os registros, o gorm leria as linhas de totais como
    registros do modelo, e elas passariam pelo `AfterFind` e pelo preload.
    Também seria preciso uma subconsulta entre parênteses para o `ORDER BY`
    e o `LIMIT`, que o gorm não monta.
- **Opcional** (decisão 4): numa página além do fim, ir para a última página
  que existe. Os totais voltam a sair da mesma consulta.

### O caminho dos valores

- **Pedido:** os `SearchParams` ganham `Aggregates []*AggregateBuilder`.
- **Resposta:** o operador devolve `SearchParams.AggregateValues`, porque a
  assinatura do `SearchFunc` não muda (decisão 5).
- **Resultado:** o `SearchResult` ganha `Aggregates`.
- **Quando:** só se o listing tiver agregações, e não com `ParamMustResult`
  (que já devolve só os registros). O `?aggregates=0` as desliga numa
  chamada.
- **Outro `DataOperator`** (não gorm2op): uma interface opcional
  `AggregatingDataOperator`. Sem ela, nenhuma agregação, e um aviso no log
  na primeira vez.
- **A captura:** o `RegisterColumnCapture` é chamado uma vez, no banco que o
  gorm2op recebe (`gorm2op.DataOperator(db)`).

### O componente

- **Onde fica:** no portal `DataTable`, acima da tabela. O `reloadList` já o
  atualiza: uma pesquisa, um filtro ou uma troca de aba refazem os totais.
  Não fica no `DataTableAdditions`, que é a barra de baixo.
- **O padrão:** uma linha de `VChip` (rótulo e valor formatado, a cor do
  `Color`), `d-flex flex-wrap ga-2 mb-2`. Com `AggregatesGroupBy`, um chip
  por grupo (rótulo do grupo e valor).
- **No diálogo** (`UpdateListingDialog`): igual.
- **Sem registros:** nenhum chip. O aviso de lista vazia fica no lugar da
  tabela.

### A serialização

```json
{
  "Records": [...],
  "TotalCount": 120,
  "Aggregates": {
    "total":        {"value": "1250.00"},
    "pagamentos":   {"value": "-3400.00"},
    "recebimentos": {"value": "4650.00"}
  }
}
```

- **O formato:** os valores crus (o `decimal` como string, sem formatar), por
  nome. Com `GroupBy`, `{"groups": [{"key": "<id>", "label": "USD", "value":
  "…"}]}`.
- **Fora:** as agregações com `Serialize(false)`.
- **O `ParamListingEncoder`:** as agregações vão junto, já que o encoder só
  muda os registros.
- **O `ParamMustResult`:** não vão (só os registros, como hoje).
- **O `http_api` do cli:** passa a mostrá-las, porque usa o mesmo evento.

### `thirdpart/gorm/utils`

- **`WithClause.DB *gorm.DB`:** a consulta como subconsulta do gorm, com os
  seus valores, em vez do texto. A base vem da busca sem virar texto à mão.
  Com `DB`, `Query` e `Args` ficam vazios.
- **A captura** sai da prova de conceito com o que falta:
  - os tipos (o `Scan` do `AggregateBuilder` converte o valor cru);
  - a documentação (README, `docs/capture.md`, um exemplo testável);
  - o teste no PostgreSQL pelo `RVQ_TEST_POSTGRES`.

### O seletor em diálogo

- **No JS**, o `AdvancedSelect.vue` ganha o modo `dialog`. Os itens
  escolhidos, os `items` de partida (o cache dos rótulos), o `item-value`,
  o `item-text`, o `many` e as mensagens de erro continuam os mesmos.
  - **Abrir:** no lugar do `v-autocomplete`, um chip "Selecionar" abre um
    `v-dialog` com a pesquisa (com debounce), a lista e a paginação, como o
    da Organização.
  - **Buscar:** a busca é a `searchPageFunc(keyword, page)`, que devolve o
    `SearchResult` inteiro: `Records`, `TotalCount` e `PerPage`, para as
    páginas.
  - **Escolher:** sem `many`, escolher fecha o diálogo; com `many`, marca e
    desmarca, e o diálogo fica aberto.
- **No Go**, `VXAdvancedSelectBuilder.Dialog(true)`. Com um `ItemsSearcher`,
  o `Write` gera também a `:search-page-func` (`keyword` e `page`).
- **No `ModelSelect`**, o diálogo passa a ser o padrão.
  `SetSelectMode(helper.SelectAutocomplete)` volta ao `v-autocomplete`. O
  `ModelSelectorConfiguror` continua recebendo o mesmo
  `*vx.VXAdvancedSelectBuilder`.
- **No filtro**, um tipo novo, `ItemTypeRecordPicker` (`RecordPickerItem`):
  - é o mesmo `AdvancedSelect` em diálogo, sem `many` ou com `many`
    (`valuesAre`);
  - a busca vem do `Config` (`RecordPickerConfig{URL, ItemText, ItemValue,
    Many, Items}`): um `fetch` da URL do listing do modelo com
    `__execute_event__=presets_ListData`, `keyword` e `page`;
  - os `Items` são os escolhidos, para os rótulos na primeira vez.
- **`helper.ModelFilterItem(mb, key, column)`** monta o item de filtro de um
  modelo: a URL do seu listing, o encoder do título e a condição
  `column = ?` (com `many`, `column IN ?`).

### Movimentações (finance)

- **`models.ValorMovimentacaoSQL(tabela)`:** o Valor em SQL, a mesma regra do
  `ValorMovimentacao`:

  ```sql
  CASE WHEN m.valor_fixado THEN m.valor_fixado_valor
       ELSE COALESCE((SELECT SUM(p.valor) FROM fin_movimentacoes_parcelas p
                      WHERE p.movimentacao_id = m.id AND p.deleted_at IS NULL
                        AND p.situacao NOT IN ('cancelada', 'dividida')), 0) END
  ```

  Um teste garante que as duas regras não se separam: para os registros de
  um banco de teste, o SQL e o Go dão o mesmo valor.
- **As agregações:**
  - Total (`Sum`), colorido pelo sinal;
  - Pagamentos (`SumIf valor < 0`), em bordô;
  - Recebimentos (`SumIf valor > 0`), em azul.

  As cores são as da coluna Valor (`valorTextoClasse`).
- **O formato:** pela moeda (`Moeda.Formatar`).
- **O agrupamento:** por `moeda_id`, um chip por moeda (decisão 1).
- **Uma agregação que não cabe no filtro não entra na consulta:**
  `AggregateBuilder.When(func(ctx) bool)`. Com Tipo = Pagamento, os
  Recebimentos ficam de fora (e o inverso).
- **Os filtros:**
  - Tipo (Pagamento ou Recebimento): pelo sinal do `ValorMovimentacaoSQL`;
  - Situação: as situações da movimentação;
  - Projeto, Pessoa, Conta Contábil e Moeda: `helper.ModelFilterItem`, o
    seletor em diálogo.
- **As mensagens:** `AggTotal`, `AggPagamentos`, `AggRecebimentos` (pt e
  en).

## Testes

- **utils:**
  - a captura no SQLite e no PostgreSQL (feita na prova de conceito);
  - `WithClause.DB`, com um exemplo e o SQL com os valores.
- **presets:**
  - as expressões dos ajudantes;
  - a consulta (dry run, SQLite), com o WITH, o `DISTINCT`, o `__count` e o
    `__aggs`, com e sem `GROUP BY`;
  - o mesmo filtro da busca: palavra-chave, `SQLConditions` e aba;
  - o `SearchResult` serializado com e sem `Serialize(false)`;
  - sem agregações, nenhuma coluna a mais.
- **gorm2op**, num banco de teste:
  - os totais de todos os registros do filtro, não só os da página;
  - uma junção que repete registros não os conta duas vezes;
  - a página 1 vazia sem consulta extra;
  - a página além do fim com os totais pelas linhas nomeadas;
  - uma moeda ausente da página nos totais.
- **finance:**
  - `ValorMovimentacaoSQL` igual ao `ValorMovimentacao`;
  - os três totais no listing e no JSON do evento de registros;
  - por moeda.
- **No Chrome (testapp):** os chips acima da tabela; uma pesquisa, a aba
  Lixeira e a troca de página os refazem.

## Ordem

0. O seletor em diálogo: o JS (o `AdvancedSelect` e o filtro
   `RecordPickerItem`, depois o build do `dist`), o `Dialog` do Go, o
   `ModelSelect` com o diálogo por padrão e o `helper.ModelFilterItem`.
1. `utils`: `WithClause.DB` (com o exemplo); a captura acabada (tipos,
   documentação).
2. presets: `AggregateBuilder`, os ajudantes, `SearchParams` e
   `SearchResult`, e a interface do operador.
3. gorm2op:
   - a consulta única no `Search`: a base, o `__count` e o `__aggs`, lidos
     pela captura;
   - a regra da página vazia e as linhas nomeadas.
4. presets: o componente padrão, o portal e a serialização.
5. Documentação: `admin/presets/docs/aggregates.md` (a convenção, com
   exemplos testáveis) e o README do presets.
6. finance: `ValorMovimentacaoSQL`, as três agregações, o agrupamento por
   moeda, os filtros e as mensagens.

## Decisões (tomadas: as recomendações)

1. **Mais de uma moeda:** agrupar por moeda
   (`AggregatesGroupBy("moeda_id")`), um chip por moeda.
2. **O rótulo do grupo:** uma função do rvq
   (`GroupLabel(func(ctx, key) string)`).
3. **A contagem na janela:** só nos listings com agregações, no começo.
4. **A página além do fim:** a consulta só de totais, em linhas nomeadas.
5. **Onde voltam os valores:** em `SearchParams.AggregateValues`; a
   assinatura do `SearchFunc` não muda.
6. **A serialização:** os valores crus, por nome.
