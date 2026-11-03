package presets

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

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
