# Aggregates of a listing (totals)

Totals of **all** the records of the filter of a listing (the search, the
filters, the tab, the organization of the page), not only those of the page.
They are shown as chips above the table, and they go with the records when
the listing is asked serialized (the `presets_ListData` event). The data
operator computes them in the select of the count.

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

## The query: the select of the count

A search makes two selects: the count of the filter, then the records of the
page. The aggregates go into the **select of the count**, the one that exists
only to aggregate. The select of the records of the page stays as it is.

```sql
WITH __agg_base AS (
	SELECT DISTINCT t.id, <base> AS valor, <group> AS __agg_group FROM … <the filter> …
)
SELECT (SELECT COUNT(*) FROM __agg_base),
	(SELECT json_agg(json_build_object('total', __agg."total", …))
		FROM (SELECT SUM(valor) AS "total", … FROM __agg_base [GROUP BY __agg_group]) __agg)
```

- **The base** is the filtered query of the search, the one of the count (no
  page, no order), put in the WITH (`gormutils.WithClause.DB`). The count is
  of the records: a join that repeats a record counts it once.
- **Always one row:** the count and the aggregates, whatever the page. On an
  empty filter the count is 0 and there is no value. With the count 0, the
  select of the records does not run, as before.
- **The JSON** is `json_build_object` / `json_agg` in PostgreSQL, `json_object`
  / `json_group_array` in SQLite and `JSON_OBJECT` / `JSON_ARRAYAGG` in MySQL:
  one row of totals, or one per group.
- **Without aggregates** the count is the same as before (gorm's `Count`).

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
  - the aggregates in the select of the count, and the select of the records
    of the page as it was;
  - an empty filter;
  - past the end;
  - without aggregates.
