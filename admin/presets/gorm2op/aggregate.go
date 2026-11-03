package gorm2op

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-rvq/rvq/admin/presets/data"
	"gorm.io/gorm"
)

// AggregatingConfigKey marks (gorm's Set) the session of the base of the
// aggregates of a search — made from the one of the count (CountingConfigKey,
// kept) —: a callback may leave out what the base needs not.
const AggregatingConfigKey = "gorm2op:aggregating"

const (
	aggBaseTable  = "__agg_base"
	aggGroupCol   = "__agg_group"
	aggGroupKey   = "__group"
	aggSubqueryAs = "__agg"
)

// countWithAggregates is the select of the count of a search that asks for
// aggregates — the count and the aggregates of all the records of the
// filter, always one row:
//
//	WITH __agg_base AS (SELECT DISTINCT <pk>, <base>, <group> FROM … <the filter> …)
//	SELECT (SELECT COUNT(*) FROM __agg_base),
//		(SELECT <json of the aggregates> FROM (SELECT … FROM __agg_base [GROUP BY]) __agg)
//
// The count is of the records (DISTINCT by the primary key: a join that
// repeats a record counts it once). The select of the records of the page
// stays as it is.
func countWithAggregates(counting *gorm.DB, obj any, q *data.AggregateQuery) (count int64, rows []*data.AggregateRow, err error) {
	base, err := aggregateBase(counting, obj, q)
	if err != nil {
		return
	}
	var aggs *string
	err = counting.Session(&gorm.Session{NewDB: true}).
		Raw("WITH "+aggBaseTable+" AS (?) SELECT (SELECT COUNT(*) FROM "+aggBaseTable+"), ("+aggregatesJSON(counting, q)+")", base).
		Row().Scan(&count, &aggs)
	if err != nil {
		return 0, nil, fmt.Errorf("the count and the aggregates: %w", err)
	}
	// an empty filter has no value (not a null of each)
	if aggs != nil && count > 0 {
		rows, err = parseAggregatesJSON(*aggs, q)
	}
	return
}

// aggregateBase is the base of the aggregates: the filtered query, its
// records once each (DISTINCT by the primary key — a join that repeats them
// counts them once), the columns of the base and the group.
func aggregateBase(filtered *gorm.DB, obj any, q *data.AggregateQuery) (*gorm.DB, error) {
	stmt := &gorm.Statement{DB: filtered}
	if err := stmt.Parse(obj); err != nil {
		return nil, fmt.Errorf("the schema of the aggregates: %w", err)
	}
	s := stmt.Schema
	if len(s.PrimaryFieldDBNames) == 0 {
		return nil, fmt.Errorf("the aggregates of %s: no primary key", s.Table)
	}
	cols := make([]string, 0, len(s.PrimaryFieldDBNames)+len(q.Base)+1)
	for _, pk := range s.PrimaryFieldDBNames {
		cols = append(cols, stmt.Quote(s.Table)+"."+stmt.Quote(pk))
	}
	for _, c := range q.Base {
		cols = append(cols, c.Expr+" AS "+stmt.Quote(c.Name))
	}
	if q.GroupBy != "" {
		cols = append(cols, q.GroupBy+" AS "+aggGroupCol)
	}
	base := filtered.Session(&gorm.Session{}).
		Set(AggregatingConfigKey, true).
		Select("DISTINCT " + strings.Join(cols, ", "))
	// the base has no record to load the associations of
	base.Statement.Preloads = map[string][]any{}
	return base, nil
}

// aggregatesJSON is the SELECT of the aggregates of the base as a JSON array
// of rows (one, or one per group): [{"__group": g, "<name>": value, …}].
func aggregatesJSON(db *gorm.DB, q *data.AggregateQuery) string {
	obj, arr := "json_object", "json_group_array"
	switch db.Dialector.Name() {
	case "postgres":
		obj, arr = "json_build_object", "json_agg"
	case "mysql":
		obj, arr = "JSON_OBJECT", "JSON_ARRAYAGG"
	}

	var (
		pairs []string
		sel   []string
	)
	if q.GroupBy != "" {
		pairs = append(pairs, sqlString(aggGroupKey), aggGroupCol)
		sel = append(sel, aggGroupCol)
	}
	for _, it := range q.Items {
		alias := db.Statement.Quote(it.Name)
		pairs = append(pairs, sqlString(it.Name), aggSubqueryAs+"."+alias)
		sel = append(sel, it.Expr+" AS "+alias)
	}
	inner := "SELECT " + strings.Join(sel, ", ") + " FROM " + aggBaseTable
	if q.GroupBy != "" {
		inner += " GROUP BY " + aggGroupCol
	}
	return "SELECT " + arr + "(" + obj + "(" + strings.Join(pairs, ", ") + ")) FROM (" + inner + ") " + aggSubqueryAs
}

// parseAggregatesJSON reads the rows of aggregatesJSON.
func parseAggregatesJSON(raw any, q *data.AggregateQuery) (rows []*data.AggregateRow, err error) {
	if raw == nil {
		return nil, nil
	}
	var s string
	switch v := raw.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprint(v)
	}
	if s == "" {
		return nil, nil
	}
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var list []map[string]any
	if err = dec.Decode(&list); err != nil {
		return nil, fmt.Errorf("the aggregates %q: %w", s, err)
	}
	for _, m := range list {
		group, hasGroup := "", q.GroupBy != ""
		if hasGroup && m[aggGroupKey] != nil {
			group = fmt.Sprint(m[aggGroupKey])
		}
		for _, it := range q.Items {
			r := &data.AggregateRow{Name: it.Name, Group: group, HasGroup: hasGroup}
			if v := m[it.Name]; v == nil {
				r.Null = true
			} else {
				r.Value = fmt.Sprint(v)
			}
			rows = append(rows, r)
		}
	}
	return
}

// sqlString is s as an SQL string literal.
func sqlString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
