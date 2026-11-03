package gorm2op

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/go-rvq/rvq/admin/presets/data"
	gormutils "github.com/go-rvq/rvq/thirdpart/gorm/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AggregatingConfigKey marks (gorm's Set) the session of the base of the
// aggregates of a search, as CountingConfigKey marks the one of the count: a
// callback may leave out what the base needs not.
const AggregatingConfigKey = "gorm2op:aggregating"

const (
	aggBaseTable  = "__agg_base"
	aggGroupCol   = "__agg_group"
	aggCountCol   = data.AggregateCountName
	aggValuesCol  = "__aggs"
	aggGroupKey   = "__group"
	aggSubqueryAs = "__agg"
)

var captureRegistered sync.Map // *gorm.Config → true

// registerCapture puts the capture of columns (utils.RegisterColumnCapture)
// on the callbacks of db, once.
func registerCapture(db *gorm.DB) error {
	if _, ok := captureRegistered.Load(db.Config); ok {
		return nil
	}
	if err := gormutils.RegisterColumnCapture(db); err != nil {
		return err
	}
	captureRegistered.Store(db.Config, true)
	return nil
}

// aggregateSearch is the page of the search with the count and the
// aggregates of all the records of the filter, in the same query:
//
//	WITH __agg_base AS (SELECT DISTINCT <pk>, <base>, <group> FROM … <the filter> …)
//	SELECT (SELECT COUNT(*) FROM __agg_base) AS __count,
//		(SELECT <json of the aggregates> FROM (SELECT … FROM __agg_base [GROUP BY]) __agg) AS __aggs,
//		<the columns of the page> FROM … <the filter> … ORDER BY … LIMIT …
//
// Both columns are read by the capture (the model has no field for them).
// A page with no row: on the first page the filter is empty (count 0, no
// value); past the end, a query of the aggregates alone, in named rows.
func (b *DataOperatorBuilder) aggregateSearch(filtered, base, page *gorm.DB, params *data.SearchParams, find func(db *gorm.DB) (*gorm.DB, error)) (totalCount int, err error) {
	q := params.Aggregates
	if err = registerCapture(filtered); err != nil {
		return
	}

	with := gormutils.WithClauses{{Name: aggBaseTable, DB: base}}
	page = page.Clauses(with.DbClause(), prependColumns{
		"(SELECT COUNT(*) FROM " + aggBaseTable + ") AS " + aggCountCol,
		"(" + aggregatesJSON(filtered, q) + ") AS " + aggValuesCol,
	})
	page, capt := gormutils.CaptureColumns(page, aggCountCol, aggValuesCol)

	if _, err = find(page); err != nil {
		return
	}

	if raw, ok := capt.Values[aggCountCol]; ok {
		if totalCount, err = strconv.Atoi(fmt.Sprint(raw)); err != nil {
			return 0, fmt.Errorf("the count of the aggregates %q: %w", raw, err)
		}
		params.AggregateRows, err = parseAggregatesJSON(capt.Values[aggValuesCol], q)
		return
	}

	if params.Page <= 1 {
		// the filter is empty
		return 0, nil
	}
	// past the end: the aggregates alone
	rows, err := aggregateRows(filtered, base, q)
	if err != nil {
		return
	}
	for _, r := range rows {
		if r.Name == aggCountCol {
			totalCount, _ = strconv.Atoi(r.Value)
			continue
		}
		params.AggregateRows = append(params.AggregateRows, r)
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

// aggregateRows are the aggregates of the base alone, in named rows — name,
// group, value (as text) —, and the row of the count.
func aggregateRows(db, base *gorm.DB, q *data.AggregateQuery) (rows []*data.AggregateRow, err error) {
	group, groupBy := "NULL", ""
	if q.GroupBy != "" {
		group, groupBy = "CAST("+aggGroupCol+" AS TEXT)", " GROUP BY "+aggGroupCol
	}
	var parts []string
	for _, it := range q.Items {
		parts = append(parts, "SELECT "+sqlString(it.Name)+" AS name, "+group+" AS grp, CAST("+it.Expr+" AS TEXT) AS value FROM "+aggBaseTable+groupBy)
	}
	parts = append(parts, "SELECT "+sqlString(aggCountCol)+", NULL, CAST(COUNT(*) AS TEXT) FROM "+aggBaseTable)

	var found []struct {
		Name  string
		Grp   *string
		Value *string
	}
	err = db.Session(&gorm.Session{NewDB: true}).
		Raw("WITH "+aggBaseTable+" AS (?) "+strings.Join(parts, " UNION ALL "), base).
		Scan(&found).Error
	if err != nil {
		return nil, fmt.Errorf("the aggregates: %w", err)
	}
	for _, f := range found {
		r := &data.AggregateRow{Name: f.Name, HasGroup: q.GroupBy != "" && f.Name != aggCountCol}
		if f.Grp != nil {
			r.Group = *f.Grp
		}
		if f.Value == nil {
			r.Null = true
		} else {
			r.Value = *f.Value
		}
		rows = append(rows, r)
	}
	return
}

// sqlString is s as an SQL string literal.
func sqlString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// prependColumns are columns (SQL expressions with their aliases) written
// before the columns of the SELECT (its AfterNameExpression).
type prependColumns []string

func (p prependColumns) ModifyStatement(stmt *gorm.Statement) {
	c := stmt.Clauses["SELECT"]
	c.AfterNameExpression = p
	stmt.Clauses["SELECT"] = c
}

func (p prependColumns) Build(builder clause.Builder) {
	for _, col := range p {
		builder.WriteString(col)
		builder.WriteString(", ")
	}
}
