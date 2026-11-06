package presets

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/data"
	"github.com/go-rvq/rvq/web"
)

func TestListingAggregates(t *testing.T) {
	lb := &ListingBuilder{}
	lb.AggregateBase("valor", "t.v")
	lb.Aggregate("total").Label("Total").Expr(Sum("valor"))
	lb.Aggregate("neg").Expr(SumIf("valor", "valor < 0")).
		When(func(ctx *web.EventContext) bool { return ctx.R.FormValue("tipo") != "pos" })
	lb.Aggregate("pos").Expr(SumIf("valor", "valor > 0")).Serialize(false)

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/?tipo=pos", nil)}
	q := lb.aggregateQuery(ctx, url.Values{})
	if q == nil || len(q.Items) != 2 || q.Items[0].Expr != "SUM(valor)" || q.Items[1].Name != "pos" {
		t.Fatalf("the query: %+v (neg out by When)", q)
	}
	if q.Items[1].Expr != "SUM(CASE WHEN valor > 0 THEN valor END)" {
		t.Fatalf("SumIf %q", q.Items[1].Expr)
	}
	if lb.aggregateQuery(ctx, url.Values{ParamAggregates: {"0"}}) != nil {
		t.Fatal("?aggregates=0 asks for none")
	}

	vs := lb.aggregateValues(ctx, []*data.AggregateRow{{Name: "total", Value: "12.5"}, {Name: "pos", Null: true}})
	if len(vs) != 2 || vs[0].Label != "Total" || vs[0].Text != "12.5" || vs[1].Text != "0" {
		t.Fatalf("values %+v %+v", vs[0], vs[1])
	}
	b, err := json.Marshal(vs)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"total":{"value":"12.5"}}` {
		t.Fatalf("serialized %s (pos is Serialize(false))", b)
	}
	if lb.aggregatesComponent(ctx, vs) == nil {
		t.Fatal("no chips")
	}
}

// The chip of an aggregate: its label, a margin, its values — of its groups
// (by currency) in one chip when combined (CombineGroups), " · " between
// them, each in its color, its group its title; else a chip of each group,
// the group before the label. The chip of the color of all, when one.
func TestListingAggregateChip(t *testing.T) {
	lb := &ListingBuilder{}
	a := lb.Aggregate("total").Label("Total").CombineGroups(true)
	n := lb.Aggregate("neg").Label("Pagamentos")
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	got := h.MustString(lb.aggregatesComponent(ctx, AggregateValues{
		{Name: "total", Label: "Total", Group: "USD", GroupLabel: "USD", Text: "$ 100,00", Color: "info", builder: a},
		{Name: "neg", Label: "Pagamentos", Group: "USD", GroupLabel: "USD", Text: "$ -5,00", Color: "error", builder: n},
		{Name: "total", Label: "Total", Group: "BRL", GroupLabel: "BRL", Text: "R$ -500,00", Color: "error", builder: a},
	}), context.TODO())
	if c := strings.Count(got, "<v-chip"); c != 2 {
		t.Fatalf("one chip per aggregate: %d\n%s", c, got)
	}
	total := got[:strings.LastIndex(got, "<v-chip")]
	for _, want := range []string{"Total:</span>", "text-medium-emphasis me-1", "$ 100,00", "R$ -500,00",
		"·</span>", "text-info", "text-error", "title='USD'", "title='BRL'"} {
		if !strings.Contains(total, want) {
			t.Errorf("the chip of the total: no %q\n%s", want, total)
		}
	}
	if strings.Index(total, "$ 100,00") > strings.Index(total, "R$ -500,00") {
		t.Error("the values out of their order")
	}
	if strings.Contains(total, " color=") {
		t.Errorf("a chip of values of two colors has one:\n%s", total)
	}
	// combined, one value: no group before the label, its group its title
	one := h.MustString(lb.aggregatesComponent(ctx, AggregateValues{
		{Name: "total", Label: "Total", Group: "USD", GroupLabel: "USD", Text: "$ 100,00", Color: "info", builder: a},
	}), context.TODO())
	if !strings.Contains(one, ">Total:</span>") || !strings.Contains(one, "title='USD'") || !strings.Contains(one, "color='info'") ||
		strings.Contains(one, "<strong") {
		t.Errorf("combined, one value:\n%s", one)
	}
	// not combined: a chip of its group, the group before the label; one
	// color, the chip's too
	if neg := got[strings.LastIndex(got, "<v-chip"):]; !strings.Contains(neg, "USD · Pagamentos:</span>") ||
		!strings.Contains(neg, "$ -5,00") || !strings.Contains(neg, "color='error'") {
		t.Errorf("the chip of the payments:\n%s", neg)
	}
}
