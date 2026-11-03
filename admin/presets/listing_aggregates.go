package presets

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/data"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The aggregates of a listing: totals of all the records of its filter — not
// only those of the page —, shown as chips above the table and sent with the
// records when the listing is asked serialized (SearchResult.Aggregates). The
// data operator computes them in the select of the count (gorm2op: a WITH of
// the filter, see docs/aggregates.md); the select of the records stays as it
// is.
//
//	lb.AggregateBase("valor", "<an expression of a record>")
//	lb.Aggregate("total").Expr(presets.Sum("valor"))
//	lb.Aggregate("negatives").Expr(presets.SumIf("valor", "valor < 0"))

// ParamAggregates is the query parameter that, "0", leaves the aggregates
// out of a search.
const ParamAggregates = "aggregates"

type listingAggregates struct {
	base       []data.NamedExpr
	items      []*AggregateBuilder
	groupBy    string
	groupLabel func(ctx *web.EventContext, key string) string
	compFunc   func(ctx *web.EventContext, values AggregateValues) h.HTMLComponent
}

// AggregateBuilder is an aggregate of a listing: its name (the key of its
// value when serialized), its label, its SQL and how it is shown.
type AggregateBuilder struct {
	NameLabel
	expr        string
	when        func(ctx *web.EventContext) bool
	scan        func(raw string) (any, error)
	format      func(ctx *web.EventContext, v *AggregateValue) string
	color       func(v *AggregateValue) string
	compFunc    func(ctx *web.EventContext, v *AggregateValue) h.HTMLComponent
	hidden      bool
	noSerialize bool
}

// Expr sets the SQL of the aggregate: an aggregate function on the columns of
// the base (Sum, SumIf, Count…).
func (b *AggregateBuilder) Expr(sql string) *AggregateBuilder {
	b.expr = sql
	return b
}

// SetI18nLabel sets the label of the aggregate in the language of the
// request.
func (b *AggregateBuilder) SetI18nLabel(f func(ctx context.Context) string) *AggregateBuilder {
	b.NameLabel.SetI18nLabel(f)
	return b
}

// Label sets the label of the aggregate.
func (b *AggregateBuilder) Label(v string) *AggregateBuilder {
	b.NameLabel.SetLabel(v)
	return b
}

// When says the aggregate is asked at a request (none: always): one the
// filter makes useless — the positives of a filter of the negatives, say —
// stays out of the query.
func (b *AggregateBuilder) When(f func(ctx *web.EventContext) bool) *AggregateBuilder {
	b.when = f
	return b
}

// Scan converts the value of the database (its text, a number) to the value
// of the aggregate (a decimal, an int…). None: the text.
func (b *AggregateBuilder) Scan(f func(raw string) (any, error)) *AggregateBuilder {
	b.scan = f
	return b
}

// Format is the text of the value on the chip. None: the text of the
// database (a null: 0).
func (b *AggregateBuilder) Format(f func(ctx *web.EventContext, v *AggregateValue) string) *AggregateBuilder {
	b.format = f
	return b
}

// Color is the color of the chip of a value (a color of the theme: error,
// primary…). None: the default.
func (b *AggregateBuilder) Color(f func(v *AggregateValue) string) *AggregateBuilder {
	b.color = f
	return b
}

// ComponentFunc draws the value in place of the default chip.
func (b *AggregateBuilder) ComponentFunc(f func(ctx *web.EventContext, v *AggregateValue) h.HTMLComponent) *AggregateBuilder {
	b.compFunc = f
	return b
}

// requestLabel is the label in the language of ctx: the i18n one, the one
// set, or the name.
func (b *AggregateBuilder) requestLabel(ctx context.Context) string {
	if f := b.I18nLabel(); f != nil {
		if l := f(ctx); l != "" {
			return l
		}
	}
	if l := b.NameLabel.Label(); l != "" {
		return l
	}
	return b.Name()
}

// Hidden: serialized, with no chip.
func (b *AggregateBuilder) Hidden() *AggregateBuilder {
	b.hidden = true
	return b
}

// Serialize says the value goes with the serialized records (the default);
// false: only on the chips.
func (b *AggregateBuilder) Serialize(v bool) *AggregateBuilder {
	b.noSerialize = !v
	return b
}

// AggregateValue is a value of an aggregate (of a group, when the aggregates
// are grouped).
type AggregateValue struct {
	Name       string
	Label      string
	Group      string
	GroupLabel string
	// Raw is the text of the database (a number); Null, no value (nothing to
	// aggregate).
	Raw   string
	Null  bool
	Value any
	Text  string
	Color string

	builder *AggregateBuilder
}

// AggregateValues are the values of the aggregates of a search, in the order
// they were declared (by group, when grouped).
type AggregateValues []*AggregateValue

// MarshalJSON is the values by name — {"total": {"value": "12.50"}}, or,
// grouped, {"total": {"groups": [{"key": "…", "label": "…", "value": "…"}]}} —,
// the values as the database writes them (null: no value); those with
// Serialize(false) left out.
func (vs AggregateValues) MarshalJSON() ([]byte, error) {
	type group struct {
		Key   string  `json:"key"`
		Label string  `json:"label,omitempty"`
		Value *string `json:"value"`
	}
	type agg struct {
		Value  *string  `json:"value,omitempty"`
		Groups []*group `json:"groups,omitempty"`
	}
	m := map[string]*agg{}
	for _, v := range vs {
		if v.builder != nil && v.builder.noSerialize {
			continue
		}
		var raw *string
		if !v.Null {
			s := v.Raw
			raw = &s
		}
		a := m[v.Name]
		if a == nil {
			a = &agg{}
			m[v.Name] = a
		}
		if v.Group != "" || v.GroupLabel != "" {
			a.Groups = append(a.Groups, &group{Key: v.Group, Label: v.GroupLabel, Value: raw})
		} else {
			a.Value = raw
		}
	}
	return json.Marshal(m)
}

func (b *ListingBuilder) aggs() *listingAggregates {
	if b.aggregates == nil {
		b.aggregates = &listingAggregates{}
	}
	return b.aggregates
}

// AggregateBase adds a column to the base of the aggregates: an SQL
// expression of a record (its columns, a subquery), read by the aggregates
// by its name.
func (b *ListingBuilder) AggregateBase(name, expr string) *ListingBuilder {
	a := b.aggs()
	for _, c := range a.base {
		if c.Name == name {
			panic("presets: the base of the aggregates has " + name + " already")
		}
	}
	a.base = append(a.base, data.NamedExpr{Name: name, Expr: expr})
	return b
}

// Aggregate is the aggregate name of the listing, made the first time.
func (b *ListingBuilder) Aggregate(name string) *AggregateBuilder {
	a := b.aggs()
	for _, it := range a.items {
		if it.name == name {
			return it
		}
	}
	it := &AggregateBuilder{}
	it.SetName(name)
	a.items = append(a.items, it)
	return it
}

// GetAggregates are the aggregates of the listing.
func (b *ListingBuilder) GetAggregates() []*AggregateBuilder {
	if b.aggregates == nil {
		return nil
	}
	return b.aggregates.items
}

// AggregatesGroupBy makes a value of each aggregate for each value of expr (a
// column of a record), labeled by label (none: the value itself).
func (b *ListingBuilder) AggregatesGroupBy(expr string, label func(ctx *web.EventContext, key string) string) *ListingBuilder {
	a := b.aggs()
	a.groupBy, a.groupLabel = expr, label
	return b
}

// AggregatesComponentFunc draws the values in place of the default bar of
// chips.
func (b *ListingBuilder) AggregatesComponentFunc(f func(ctx *web.EventContext, values AggregateValues) h.HTMLComponent) *ListingBuilder {
	b.aggs().compFunc = f
	return b
}

// aggregateQuery is what a search of the request asks for; nil: no
// aggregate (none declared, none at the request, or ?aggregates=0).
func (b *ListingBuilder) aggregateQuery(ctx *web.EventContext, qs url.Values) *data.AggregateQuery {
	a := b.aggregates
	if a == nil || len(a.items) == 0 || qs.Get(ParamAggregates) == "0" {
		return nil
	}
	q := &data.AggregateQuery{Base: a.base, GroupBy: a.groupBy}
	for _, it := range a.items {
		if it.expr == "" || (it.when != nil && !it.when(ctx)) {
			continue
		}
		q.Items = append(q.Items, data.NamedExpr{Name: it.name, Expr: it.expr})
	}
	if len(q.Items) == 0 {
		return nil
	}
	return q
}

// aggregateValues are the values of the rows the operator found, in the
// order of the aggregates (group by group, in the order the groups came).
func (b *ListingBuilder) aggregateValues(ctx *web.EventContext, rows []*data.AggregateRow) (vs AggregateValues) {
	a := b.aggregates
	if a == nil || len(rows) == 0 {
		return nil
	}
	var groups []string
	seen := map[string]bool{}
	byKey := map[string]*data.AggregateRow{}
	for _, r := range rows {
		if !seen[r.Group] {
			seen[r.Group] = true
			groups = append(groups, r.Group)
		}
		byKey[r.Group+"\x00"+r.Name] = r
	}
	for _, g := range groups {
		groupLabel := ""
		if a.groupBy != "" {
			groupLabel = g
			if a.groupLabel != nil {
				groupLabel = a.groupLabel(ctx, g)
			}
		}
		for _, it := range a.items {
			r := byKey[g+"\x00"+it.name]
			if r == nil {
				continue
			}
			v := &AggregateValue{Name: it.name, Label: it.requestLabel(ctx.Context()), Group: g,
				GroupLabel: groupLabel, Raw: r.Value, Null: r.Null, builder: it}
			if !v.Null {
				v.Value = v.Raw
				if it.scan != nil {
					if x, err := it.scan(v.Raw); err == nil {
						v.Value = x
					}
				}
			}
			switch {
			case it.format != nil:
				v.Text = it.format(ctx, v)
			case v.Null:
				v.Text = "0"
			default:
				v.Text = v.Raw
			}
			if it.color != nil {
				v.Color = it.color(v)
			}
			vs = append(vs, v)
		}
	}
	return
}

// aggregatesComponent is the bar of the chips of the values (the listing's
// AggregatesComponentFunc, when set); nil with no value.
func (b *ListingBuilder) aggregatesComponent(ctx *web.EventContext, vs AggregateValues) h.HTMLComponent {
	if len(vs) == 0 {
		return nil
	}
	if b.aggregates.compFunc != nil {
		return b.aggregates.compFunc(ctx, vs)
	}
	var chips []h.HTMLComponent
	for _, v := range vs {
		if v.builder.hidden {
			continue
		}
		if v.builder.compFunc != nil {
			chips = append(chips, v.builder.compFunc(ctx, v))
			continue
		}
		label := v.Label
		if v.GroupLabel != "" {
			label = v.GroupLabel + " · " + label
		}
		chip := VChip(h.Span(label+": ").Class("text-medium-emphasis"), h.Strong(v.Text)).
			Size(SizeSmall).Variant(VariantTonal).
			Class("vx-aggregate", "vx-aggregate--"+strings.ReplaceAll(v.Name, " ", "-"))
		if v.Color != "" {
			chip.Color(v.Color)
		}
		chips = append(chips, chip)
	}
	if len(chips) == 0 {
		return nil
	}
	return h.Div(chips...).Class("d-flex flex-wrap ga-2 mb-2 rvq-listing-aggregates")
}

// Sum is the sum of col: SUM(col).
func Sum(col string) string { return "SUM(" + col + ")" }

// SumIf is the sum of col where cond: SUM(CASE WHEN cond THEN col END).
func SumIf(col, cond string) string { return "SUM(CASE WHEN " + cond + " THEN " + col + " END)" }

// Count is the count of the records: COUNT(*).
func Count() string { return "COUNT(*)" }

// CountIf is the count of the records where cond.
func CountIf(cond string) string { return "COUNT(CASE WHEN " + cond + " THEN 1 END)" }

// Avg is the average of col.
func Avg(col string) string { return "AVG(" + col + ")" }

// Min is the least of col.
func Min(col string) string { return "MIN(" + col + ")" }

// Max is the greatest of col.
func Max(col string) string { return "MAX(" + col + ")" }
