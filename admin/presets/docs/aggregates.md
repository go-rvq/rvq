# Aggregates of a listing (totals)

Totals of **all** the records of the filter of a listing (the search, the
filters, the tab, the organization of the page), not only those of the page.
They are shown as chips above the table, and they go with the records when
the listing is asked serialized (the `presets_ListData` event).

```go
lb := mb.Listing()
lb.AggregateBase("valor", "<an SQL expression of a record>") // a column of the base
lb.Aggregate("total").
	SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).Total }).
	Expr(presets.Sum("valor"))
lb.Aggregate("negatives").
	Expr(presets.SumIf("valor", "valor < 0")).
	Color(func(*presets.AggregateValue) string { return "error" })
lb.Aggregate("positives").
	Expr(presets.SumIf("valor", "valor > 0")).
	// out of the query when the filter makes it useless
	When(func(ctx *web.EventContext) bool { return presets.FilterQueryValue(ctx, "tipo") != "negative" })
lb.AggregatesGroupBy("t.currency_id", func(ctx *web.EventContext, key string) string { return codeOf(key) })
```

## The declaration

- **`AggregateBase(name, expr)`:** a column computed for each record. The base
  is the filtered query, one row per record (`DISTINCT` by the primary key),
  plus these columns.
- **`Aggregate(name)`:** an aggregate, made the first time. Its name is the key
  of its value when serialized.
  - `Expr(sql)`: an aggregate function on the columns of the base. The
    helpers are `Sum`, `SumIf`, `Count`, `CountIf`, `Avg`, `Min` and `Max`.
    `SumIf` writes `SUM(CASE WHEN cond THEN col END)`, which works in
    PostgreSQL, SQLite and MySQL.
  - `SetI18nLabel` / `Label`: the label of the chip.
  - `When(func(ctx) bool)`: asked only when it returns true. An aggregate
    left out does not go into the query.
  - `Scan(func(raw string) (any, error))`: the value from the text of the
    database (a decimal, an int).
  - `Format(func(ctx, v) string)`: the text of the chip. Without it, the text
    of the database (a null shows as 0).
  - `Color(func(v) string)`: a color of the theme for the chip.
  - `ComponentFunc`: draws the value in place of the chip.
  - `Hidden()`: serialized, with no chip.
  - `Serialize(false)`: only on the chip.
- **`AggregatesGroupBy(expr, label)`:** a value of each aggregate for each
  value of `expr` (a column), labeled by `label`. Use it for one total per
  currency, for example, where a sum across currencies means nothing.
- **`AggregatesComponentFunc`:** draws the whole bar in place of the chips.
- **`?aggregates=0`** leaves the aggregates out of a request.
  `ParamMustResult` (the records alone) leaves them out too.

## The query: one, the one of the page

The gorm data operator (`gorm2op`) computes the page, the count and the
aggregates in the same query:

```sql
WITH __agg_base AS (
	SELECT DISTINCT t.id, <base> AS valor, <group> AS __agg_group FROM … <the filter> …
)
SELECT (SELECT COUNT(*) FROM __agg_base) AS __count,
	(SELECT json_agg(json_build_object('total', __agg."total", …))
		FROM (SELECT SUM(valor) AS "total", … FROM __agg_base [GROUP BY __agg_group]) __agg) AS __aggs,
	<the columns of the page>
FROM … <the filter> … ORDER BY … LIMIT … OFFSET …
```

- **The base** is the filtered query of the search, before the page, put in
  the WITH (`gormutils.WithClause.DB`). A join that repeats a record counts it
  once.
- **`__count` and `__aggs`** are columns the model has no field for. gorm
  discards them, so the capture of columns (`gormutils.CaptureColumns`) reads
  them. The model keeps its hooks (`AfterFind`) and its preloads. PostgreSQL
  runs each subquery once (an `InitPlan`), not once per row.
- **The JSON** is `json_build_object` / `json_agg` in PostgreSQL, `json_object`
  / `json_group_array` in SQLite and `JSON_OBJECT` / `JSON_ARRAYAGG` in MySQL.
- **A page with no row:**
  - on the first page, the filter is empty: the count is 0 and there is no
    value, so no other query runs;
  - past the end, a query of the aggregates alone runs, in named rows
    (`name`, group, value).
- **Without aggregates** the search is the same as before: a count, then the
  page.

The values come back in `SearchParams.AggregateRows`, and the listing turns
them into `SearchResult.Aggregates`. Another data operator has no aggregates.

## On the page and serialized

- **The chips** go above the table, in the portal of the table, so a search,
  a filter, a tab or a page redoes them. There are none when the filter is
  empty.
- **Serialized** (the `presets_ListData` event), the values are raw, by name:

  ```json
  {"Records": [...], "TotalCount": 120,
   "Aggregates": {"total": {"value": "1250.00"},
                  "negatives": {"groups": [{"key": "<id>", "label": "USD", "value": "-3400.00"}]}}}
  ```

  A null is no value. The aggregates with `Serialize(false)` are left out.

## Tests

- `listing_aggregates_test.go`: the query of the request (`When`,
  `?aggregates=0`), the values and the serialization.
- `gorm2op/aggregate_test.go` (SQLite):
  - the totals of the whole filter, and a record repeated by a join counted
    once;
  - by group, including a group out of the page;
  - an empty filter;
  - past the end;
  - without aggregates.
- `thirdpart/gorm/utils/capture_test.go`: the capture, in SQLite and in
  PostgreSQL (`RVQ_TEST_POSTGRES`).
