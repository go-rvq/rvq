package presets

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// The values of a facet: none posted, nil (no filter); posted empty (f_k=""),
// none chosen; the values posted.
func TestLinkCardFacetValues(t *testing.T) {
	ctx := func(form url.Values) *web.EventContext {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return &web.EventContext{R: r}
	}
	if v := LinkCardFacetValues(ctx(url.Values{"q": {"x"}}), "kind"); v != nil {
		t.Errorf("not posted: %#v", v)
	}
	if v := LinkCardFacetValues(ctx(url.Values{"f_kind": {""}}), "kind"); v == nil || len(v) != 0 {
		t.Errorf("posted empty: %#v", v)
	}
	if v := LinkCardFacetValues(ctx(url.Values{"f_kind": {"pages", "posts"}}), "kind"); !reflect.DeepEqual(v, []string{"pages", "posts"}) {
		t.Errorf("posted: %#v", v)
	}
}

// The search asked for its facets answers {groups, facets}: the
// application's, and its searchers' groups even without a text; not asked,
// the groups alone, none without a text.
func TestLinkCardSearchFacets(t *testing.T) {
	b := &Builder{}
	b.LinkCardSearcher(func(ctx *web.EventContext, q string) []*LinkCardOptionGroup {
		if len(LinkCardFacetValues(ctx, "kind")) == 0 {
			return nil
		}
		return []*LinkCardOptionGroup{{Kind: "posts", Items: []LinkCardOption{{Title: "A", Target: "/a"}}}}
	})
	b.LinkCardFacets(func(ctx *web.EventContext) []*LinkCardFacet {
		return []*LinkCardFacet{{Key: "kind", Kind: LinkCardFacetChips, Items: []LinkCardFacetItem{{Value: "posts", Title: "Posts"}}}}
	})
	call := func(form url.Values) any {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := b.linkCardSearch(&web.EventContext{R: r})
		if err != nil {
			t.Fatal(err)
		}
		return res.Data
	}
	got, ok := call(url.Values{"facets": {"1"}, "f_kind": {"posts"}}).(*linkCardSearchResult)
	if !ok || len(got.Facets) != 1 || got.Facets[0].Key != "kind" || len(got.Groups) != 1 {
		t.Fatalf("with facets: %#v", got)
	}
	if g, ok := call(url.Values{"q": {"x"}, "f_kind": {"posts"}}).([]*LinkCardOptionGroup); !ok || len(g) != 0 {
		t.Errorf("without facets, one letter: %#v", g)
	}
}

// Without the records (LinkCardRecords(false)), only the application's
// options: the search of the records is not asked.
func TestLinkCardSearchNoRecords(t *testing.T) {
	b := &Builder{}
	b.LinkCardRecords(false)
	b.LinkCardSearcher(func(ctx *web.EventContext, q string) []*LinkCardOptionGroup {
		return []*LinkCardOptionGroup{{Kind: "pages", Items: []LinkCardOption{{Title: "A", Target: "/a"}}}}
	})
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url.Values{"q": {"abc"}}.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := b.linkCardSearch(&web.EventContext{R: r})
	if err != nil {
		t.Fatal(err)
	}
	if g, ok := res.Data.([]*LinkCardOptionGroup); !ok || len(g) != 1 || g[0].Kind != "pages" {
		t.Errorf("got %#v", res.Data)
	}
}
