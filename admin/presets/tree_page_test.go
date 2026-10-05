package presets

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

type treePageDoc struct {
	ID   uint
	Body string
}

// treePageApp is a singleton whose page is a tree page: "a" with "a/b" under
// it, and "c".
func treePageApp(t *testing.T) (http.Handler, *TreePageBuilder) {
	t.Helper()
	b := New(i18n.New().SupportLanguages(language.English)).URIPrefix("/admin")
	mb := b.Model(&treePageDoc{}).URIName("docs").Singleton(true).Label("Docs")
	mb.Detailing().FetchFunc(func(any, ID, *web.EventContext) error { return nil })
	tp := mb.TreePage(func(*web.EventContext) []*TreePageNode {
		return []*TreePageNode{
			{ID: "a", Title: "Node A", Icon: "mdi-a", Children: []*TreePageNode{{ID: "a/b", Title: "Node B"}}},
			{ID: "c", Title: "Node C"},
		}
	})
	mb.Detailing("Body").Field("Body").ComponentFunc(func(_ *FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Text("shown:"+tp.Shown(ctx))).Attr("data-shown", true)
	})
	mux := http.NewServeMux()
	b.Build(mux)
	return mux, tp
}

func treePageGet(h http.Handler, target string) string {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", target, nil))
	return strings.NewReplacer(`<`, "<", `>`, ">", "&#39;", "'", `\"`, `"`, `&`, "&").Replace(w.Body.String())
}

// The page of a tree: each node at its address — the page, the node shown —;
// the breadcrumbs the model and the nodes above it, each linked, and it
// last, the title of the page; the menu with the tree, the node active and
// the way to it open; its events at the address of a node.
func TestTreePage(t *testing.T) {
	h, tp := treePageApp(t)
	if tp.Href("a/b") != "/admin/docs/a/b" || tp.Base() != "/admin/docs" {
		t.Errorf("Href %q Base %q", tp.Href("a/b"), tp.Base())
	}

	body := treePageGet(h, "/admin/docs/a/b")
	if !strings.Contains(body, "shown:a/b") {
		t.Fatalf("the node of the address not shown")
	}
	i := strings.Index(body, "rvq-page-breadcrumbs")
	if i < 0 {
		t.Fatal("no breadcrumbs")
	}
	bc := body[i:min(len(body), i+2000)]
	for _, want := range []string{"href='/admin/docs'", "Docs", "href='/admin/docs/a'", "Node A", "Node B"} {
		if !strings.Contains(bc, want) {
			t.Errorf("the breadcrumbs: no %q\n%s", want, bc)
		}
	}
	if strings.Contains(bc, "href='/admin/docs/a/b'") {
		t.Error("the node shown linked in its breadcrumbs")
	}
	if m := regexp.MustCompile(`<title>([^<]*)</title>`).FindStringSubmatch(body); m == nil || !strings.Contains(m[1], "Node B") {
		t.Errorf("the title: %v", m)
	}
	for _, want := range []string{`"value":"/admin/docs/a/b"`, `"prependIcon":"mdi-a"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the menu: no %q", want)
		}
	}
	if !regexp.MustCompile(`activated: \["/admin/docs/a/b"\]`).MatchString(body) {
		t.Error("the node not active in the menu")
	}
	if !regexp.MustCompile(`opened: \[[^\]]*"/admin/docs/a"`).MatchString(body) {
		t.Error("the way to the node not open")
	}

	// the page itself: its first node; a query is nothing
	if body := treePageGet(h, "/admin/docs?node=c"); !strings.Contains(body, "shown:a") || strings.Contains(body, "shown:c") {
		t.Error("the page itself: not its first node")
	}
	// the menu's event at the address of a node: the tree, as addresses
	r := httptest.NewRequest("POST", "/admin/docs/c?__execute_event__="+MenuChildrenEvent, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), `"/admin/docs/a/b"`) {
		t.Errorf("the menu's children: %.300s", w.Body.String())
	}
}

// The address of a node: its path, each part escaped.
func TestTreePageHref(t *testing.T) {
	for id, want := range map[string]string{"a": "/d/a", "a/b c": "/d/a/b%20c", "/x/?y/": "/d/x/%3Fy"} {
		if got := TreePageHref("/d/", id); got != want {
			t.Errorf("%q: %q, want %q", id, got, want)
		}
	}
	tree := []*TreePageNode{{ID: "a", Children: []*TreePageNode{{ID: "a/b"}}}}
	if n := TreePageFind(tree, "a/b"); n == nil || n.ID != "a/b" || TreePageFind(tree, "z") != nil {
		t.Error("TreePageFind")
	}
}
