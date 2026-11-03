package gorm2op

import (
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets/data"
	"github.com/go-rvq/rvq/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type aggCur struct {
	ID   uint `gorm:"primarykey"`
	Code string
}

type aggTag struct {
	ID    uint `gorm:"primarykey"`
	MovID uint
	Name  string
}

type aggMov struct {
	ID    uint `gorm:"primarykey"`
	Value int
	CurID uint
	Cur   *aggCur
	found bool
}

func (m *aggMov) AfterFind(*gorm.DB) error { m.found = true; return nil }

// 5 records: USD -30, -20, 10; BRL 40, 50 — the first with two tags (a join
// of the tags repeats it)
func newAggOperator(t *testing.T, joinTags bool) (*DataOperatorBuilder, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&aggCur{}, &aggMov{}, &aggTag{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]*aggCur{{ID: 1, Code: "USD"}, {ID: 2, Code: "BRL"}})
	for i, v := range []struct{ value, cur int }{{-30, 1}, {-20, 1}, {10, 1}, {40, 2}, {50, 2}} {
		db.Create(&aggMov{ID: uint(i + 1), Value: v.value, CurID: uint(v.cur)})
	}
	db.Create(&[]*aggTag{{MovID: 1, Name: "a"}, {MovID: 1, Name: "b"}})

	op := DataOperator(db).WrapPrepare(func(old Preparer) Preparer {
		return func(db *gorm.DB, mode Mode, obj interface{}, id model.ID, params *data.SearchParams, ctx *web.EventContext) *gorm.DB {
			db = old(db, mode, obj, id, params, ctx).Joins("Cur")
			if joinTags {
				db = db.Joins("LEFT JOIN agg_tags ON agg_tags.mov_id = agg_movs.id")
			}
			return db
		}
	})
	return op, db
}

func aggParams(page int64) *data.SearchParams {
	return &data.SearchParams{
		PerPage: 2, Page: page, OrderBy: "agg_movs.id",
		Aggregates: &data.AggregateQuery{
			Base: []data.NamedExpr{{Name: "valor", Expr: "agg_movs.value"}},
			Items: []data.NamedExpr{
				{Name: "total", Expr: "SUM(valor)"},
				{Name: "neg", Expr: "SUM(CASE WHEN valor < 0 THEN valor END)"},
				{Name: "pos", Expr: "SUM(CASE WHEN valor > 0 THEN valor END)"},
			},
		},
	}
}

func aggValues(rows []*data.AggregateRow) map[string]string {
	m := map[string]string{}
	for _, r := range rows {
		k := r.Name
		if r.HasGroup {
			k = r.Group + "/" + r.Name
		}
		if r.Null {
			m[k] = "null"
		} else {
			m[k] = r.Value
		}
	}
	return m
}

func TestSearchAggregates(t *testing.T) {
	ctx := ctxWithForm(url.Values{})

	t.Run("the page and the totals of the whole filter, a repeated record once", func(t *testing.T) {
		op, _ := newAggOperator(t, true)
		params := aggParams(1)
		var movs []*aggMov
		r, total, err := op.Search(&movs, params, ctx)
		if err != nil {
			t.Fatal(err)
		}
		page := r.([]*aggMov)
		if total != 5 {
			t.Fatalf("total %d, want 5 (the records, not the rows of the join)", total)
		}
		if len(page) != 2 || !page[0].found || page[0].Cur == nil || page[0].Cur.Code != "USD" {
			t.Fatalf("the page: %+v", page)
		}
		got := aggValues(params.AggregateRows)
		if got["total"] != "50" || got["neg"] != "-50" || got["pos"] != "100" {
			t.Fatalf("aggregates %v", got)
		}
	})

	t.Run("grouped: a value of each group, a group out of the page too", func(t *testing.T) {
		op, _ := newAggOperator(t, false)
		params := aggParams(1)
		params.Aggregates.GroupBy = "agg_movs.cur_id"
		var movs []*aggMov
		if _, _, err := op.Search(&movs, params, ctx); err != nil {
			t.Fatal(err)
		}
		got := aggValues(params.AggregateRows)
		if got["1/total"] != "-40" || got["1/pos"] != "10" || got["2/total"] != "90" || got["2/neg"] != "null" {
			t.Fatalf("aggregates %v", got)
		}
	})

	t.Run("an empty filter: no value, count 0", func(t *testing.T) {
		op, _ := newAggOperator(t, false)
		params := aggParams(1)
		params.SQLConditions = data.SQLConditions{{Query: "agg_movs.value > ?", Args: []any{1000}}}
		var movs []*aggMov
		r, total, err := op.Search(&movs, params, ctx)
		if err != nil || total != 0 || len(r.([]*aggMov)) != 0 || len(params.AggregateRows) != 0 {
			t.Fatalf("%v %d %v %v", err, total, r, params.AggregateRows)
		}
	})

	t.Run("past the end: the totals by the query of them alone", func(t *testing.T) {
		op, _ := newAggOperator(t, false)
		params := aggParams(9)
		var movs []*aggMov
		r, total, err := op.Search(&movs, params, ctx)
		if err != nil {
			t.Fatal(err)
		}
		if total != 5 || len(r.([]*aggMov)) != 0 {
			t.Fatalf("total %d, records %v", total, r)
		}
		got := aggValues(params.AggregateRows)
		if got["total"] != "50" || got["neg"] != "-50" || got["pos"] != "100" {
			t.Fatalf("aggregates %v", got)
		}
	})

	t.Run("no aggregates: the search as before", func(t *testing.T) {
		op, _ := newAggOperator(t, false)
		params := aggParams(1)
		params.Aggregates = nil
		var movs []*aggMov
		r, total, err := op.Search(&movs, params, ctx)
		if err != nil || total != 5 || len(r.([]*aggMov)) != 2 || params.AggregateRows != nil {
			t.Fatalf("%v %d %v", err, total, params.AggregateRows)
		}
	})
}
