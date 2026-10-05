package presets

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-rvq/rvq/web"
)

// A TreePage is the page of a model — its detail; a singleton's — made of a
// tree of nodes, each at its own address under the page: base/<id>, the id
// its path ("/admin/docs/guides/policies"). One page, the node of its address
// shown, and all of it following the node:
//
//   - the menu: the tree under the model's item — with the page, the node
//     shown active and the way to it open; when the item is opened,
//     otherwise —, each node opening its address without a reload;
//   - the address: base/<id> serves the page, and its events;
//   - the breadcrumbs: the model, then the nodes above the one shown, each
//     its link; the node shown, last, the title of the page.
//
// The documentation of the admin is one (userdocs); any menu of its own
// nodes — dynamic — is another.

// TreePageNode is a node of a TreePage: its ID — its path under the page,
// "guides/policies" —, its title, its icon, the nodes under it.
type TreePageNode struct {
	ID       string
	Title    string
	Icon     string
	Children []*TreePageNode
}

// TreePageBuilder is the tree page of a model (ModelBuilder.TreePage).
type TreePageBuilder struct {
	mb   *ModelBuilder
	tree func(ctx *web.EventContext) []*TreePageNode
}

// treePageParam is the wildcard of the node in the address of a node.
const treePageParam = "node"

type treePageNodeKey struct{}

// TreePage makes the page of the model the tree of nodes of tree (TreePage),
// the request's: the nodes it may see, in its language.
func (mb *ModelBuilder) TreePage(tree func(ctx *web.EventContext) []*TreePageNode) *TreePageBuilder {
	b := &TreePageBuilder{mb: mb, tree: tree}
	mb.MenuChildren(b.menuChildren)
	d := mb.Detailing()
	d.Breadcrumbs(func(_ any, ctx *web.EventContext, bc *BreadcrumbsBuilder) {
		path := treePagePath(b.tree(ctx), b.Shown(ctx))
		if len(path) == 0 {
			return
		}
		bc.Append(&Breadcrumb{Label: mb.TTitle(ctx.Context()), URI: b.Base()})
		for _, n := range path[:len(path)-1] {
			bc.Append(&Breadcrumb{Label: n.Title, URI: b.Href(n.ID)})
		}
	})
	d.PageTitleFunc(func(_ any, ctx *web.EventContext) string {
		if n := TreePageFind(b.tree(ctx), b.Shown(ctx)); n != nil {
			return n.Title
		}
		return ""
	})
	mb.p.MuxSetup(func(_ string, mux *http.ServeMux) {
		base := b.Base()
		mux.Handle(base+"/{"+treePageParam+"...}", treePageHandler(base, mux))
	})
	return b
}

// Base is the address of the page: base, of the first node.
func (b *TreePageBuilder) Base() string { return b.mb.Info().ListingHref() }

// Href is the address of the node id: base/<id>.
func (b *TreePageBuilder) Href(id string) string { return TreePageHref(b.Base(), id) }

// Shown is the node the request shows: of its address, else the first.
func (b *TreePageBuilder) Shown(ctx *web.EventContext) string {
	if id := TreePageNodeOf(ctx.R); id != "" {
		return id
	}
	if tree := b.tree(ctx); len(tree) > 0 {
		return tree[0].ID
	}
	return ""
}

// Tree is the tree of the request.
func (b *TreePageBuilder) Tree(ctx *web.EventContext) []*TreePageNode { return b.tree(ctx) }

// menuChildren are the nodes as the children of the model's item, each its
// address; the one shown active.
func (b *TreePageBuilder) menuChildren(ctx *web.EventContext) ([]*MenuNode, string) {
	var nodes func(tree []*TreePageNode) []*MenuNode
	nodes = func(tree []*TreePageNode) (r []*MenuNode) {
		for _, n := range tree {
			href := b.Href(n.ID)
			props := map[string]any{"href": href}
			if n.Icon != "" {
				props["prependIcon"] = n.Icon
			}
			r = append(r, &MenuNode{Title: n.Title, Value: href, Props: props, Children: nodes(n.Children)})
		}
		return
	}
	return nodes(b.tree(ctx)), b.Href(b.Shown(ctx))
}

// TreePageHref is the address of the node id of the tree page at base: its
// path under it, each part escaped.
func TreePageHref(base, id string) string {
	parts := strings.Split(strings.Trim(id, "/"), "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.TrimSuffix(base, "/") + "/" + strings.Join(parts, "/")
}

// TreePageNodeOf is the node of the address of r, of a tree page; "" at the
// page itself.
func TreePageNodeOf(r *http.Request) string {
	n, _ := r.Context().Value(treePageNodeKey{}).(string)
	return n
}

// treePageHandler serves base/<node> — the page and its events — as the page
// at base, the node of the address told (TreePageNodeOf).
func treePageHandler(base string, mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(context.WithValue(r.Context(), treePageNodeKey{}, r.PathValue(treePageParam)))
		r2.URL.Path, r2.URL.RawPath = base, ""
		mux.ServeHTTP(w, r2)
	})
}

// TreePageFind is the node id of tree; nil when none.
func TreePageFind(tree []*TreePageNode, id string) *TreePageNode {
	if p := treePagePath(tree, id); len(p) > 0 {
		return p[len(p)-1]
	}
	return nil
}

// treePagePath are the nodes from the root of tree to the node id, it the
// last; none when it is not in tree.
func treePagePath(tree []*TreePageNode, id string) []*TreePageNode {
	for _, n := range tree {
		if n.ID == id {
			return []*TreePageNode{n}
		}
		if p := treePagePath(n.Children, id); p != nil {
			return append([]*TreePageNode{n}, p...)
		}
	}
	return nil
}
