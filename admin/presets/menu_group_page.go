package presets

import (
	"context"
	"fmt"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// A group of the menu has a path — the segment of its models' URIs
// (/admin/orgs/{org}/finance/carteiras) — and a page there:
//
//   - its own, when it is made of one (PageMenuGroup): the group is the page,
//     its name, title and icon the page's — a dashboard of an area, say;
//   - else an index (groupIndex): the items of the group the user may see —
//     the menu's —, a tree with their icons, titles (links) and
//     descriptions.
//
// The title of the group in the menu, in the rail and in the breadcrumbs is
// a link to that page.

// PageMenuGroup is the group of the menu that is page: its name is the
// segment of the page's path ("/finance": "finance") — its key in the menu,
// its segment in its models' URIs, the resource of its permissions —, its
// title and icon the page's; none of them is set otherwise (the group's
// Title, TitleFunc, Icon and URIName panic: two names of a thing are a
// confusion). The page is a page of the builder, or of the detail of a
// model's record (Detailing().PagesRegistrator()) for a group of the
// resources of that record (the organization's). MenuGroup(name) is the same
// group after it.
func (b *Builder) PageMenuGroup(page *HttpPageBuilder) *MenuGroupBuilder {
	name := strings.TrimPrefix(page.path, "/")
	if name == "" || strings.ContainsAny(name, "/{}") {
		panic(fmt.Sprintf("presets: the page of a menu group has a path of one segment, not %q", page.path))
	}
	b.MenuTree()
	if it := b.menuItems[menuKey(MenuItemGroup, name)]; it != nil {
		if g := it.Group(); g != nil {
			switch {
			case g.page == page:
				return g
			case g.page != nil:
				panic(fmt.Sprintf("presets: the menu group %q is of another page", name))
			case g.title != nil || g.icon != "" || g.uriName != nil:
				panic(fmt.Sprintf("presets: the menu group %q has a title, an icon or an URI name of its own: of a page, they are the page's", name))
			}
			g.page = page
			page.notInMenu = true
			return g
		}
	}
	g := b.MenuGroup(name)
	g.page = page
	// the group is the page: no entry of its own
	page.notInMenu = true
	return g
}

// GetPage is the page the group is (PageMenuGroup), nil none.
func (b *MenuGroupBuilder) GetPage() *HttpPageBuilder { return b.page }

// notOfPage panics when the group is a page's: what of it is the page's is
// not set (PageMenuGroup).
func (b *MenuGroupBuilder) notOfPage(what string) {
	if b.page != nil {
		panic(fmt.Sprintf("presets: %s of the menu group %q: of a page (PageMenuGroup), it is the page's", what, b.name))
	}
}

// rootHref is the address of the page of the group of the root (not of a
// record's): its page's, or its index's; "" when it has none.
func (b *MenuGroupBuilder) rootHref() string {
	if b == nil || b.b == nil {
		return ""
	}
	if b.page != nil {
		if b.page.owner != nil {
			return ""
		}
		return strings.TrimSuffix(b.b.prefix, "/") + b.page.path
	}
	if !b.indexRoot {
		return ""
	}
	return strings.TrimSuffix(b.b.prefix, "/") + "/" + b.URIPath()
}

// recordHref is the address of the page of the group of the resources of
// the record id of mb (its parents): its page's, or its index's; "" when it
// has none.
func (b *MenuGroupBuilder) recordHref(mb *ModelBuilder, id string, parents IDSlice) string {
	if b == nil {
		return ""
	}
	if b.page != nil {
		if b.page.owner != mb {
			return ""
		}
		return mb.Info().DetailingHref(id, parents...) + b.page.path
	}
	if b.indexRecords == nil || !b.indexRecords[mb] {
		return ""
	}
	return mb.Info().DetailingHref(id, parents...) + "/" + b.URIPath()
}

// setupGroupPages registers the index of each group with a path and no page
// of its own (groupIndex): at the root, a group of the root's items; under
// a record of a model, a group of its children's. Before the routes (Build).
func (b *Builder) setupGroupPages() {
	if b.groupPagesSet {
		return
	}
	b.groupPagesSet = true
	b.MenuTree()
	for _, it := range b.menuItems {
		g := it.Group()
		if g == nil || g.page != nil || g.URIPath() == "" {
			continue
		}
		// at the root: a group of a model of no parent, of a page of the
		// builder or of another group
		if g.hasRootItems() {
			path := "/" + g.URIPath()
			if b.pagesRegistrator.GetHttpPage(path) == nil {
				g.indexRoot = true
				g := g
				b.pagesRegistrator.New(HttpPage(path).InMenu(false).
					TitleFunc(g.TTitle).DescriptionFunc(g.TDescription)).
					Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
						return g.index(ctx, func(n *MenuNode) bool { return n.Value == menuKey(MenuItemGroup, g.name) })
					}).Build()
			}
		}
	}
	// under a record: the groups of its children
	for _, mb := range b.models {
		seen := map[*MenuGroupBuilder]bool{}
		for _, child := range mb.children {
			g := child.MenuGroup()
			if g == nil || seen[g] || g.page != nil || g.URIPath() == "" {
				continue
			}
			seen[g] = true
			path := "/" + g.URIPath()
			if mb.Detailing().PagesRegistrator().GetHttpPage(path) != nil {
				continue
			}
			if g.indexRecords == nil {
				g.indexRecords = map[*ModelBuilder]bool{}
			}
			g.indexRecords[mb] = true
			g, mb := g, mb
			mb.Detailing().PagesRegistrator().New(HttpPage(path).InMenu(false).
				TitleFunc(g.TTitle).DescriptionFunc(g.TDescription)).
				RecordLayout(func(ctx *web.EventContext, mid model.ID, _ []model.ID, _ func() (any, error)) (r web.PageResponse, err error) {
					suffix := ":" + g.name
					key := mb.menuRecordKey(ParentsModelID(ctx.R))
					return g.index(ctx, func(n *MenuNode) bool {
						return strings.HasPrefix(n.Value, key) && strings.HasSuffix(n.Value, suffix)
					})
				}).Build()
		}
	}
}

// hasRootItems says the group has an item of the root: a model of no
// parent, a page of the builder, a group with one.
func (b *MenuGroupBuilder) hasRootItems() bool {
	for _, it := range b.items {
		switch v := it.Value.(type) {
		case *ModelBuilder:
			if v.parent == nil {
				return true
			}
		case *HttpPageBuilder:
			if v.owner == nil {
				return true
			}
		case *MenuGroupBuilder:
			if v.hasRootItems() {
				return true
			}
		}
	}
	return false
}

// menuDescribedKey says, in a context, that the nodes of the menu carry the
// descriptions of their items (subtitle): the index of a group.
type menuDescribedKey struct{}

func menuDescribed(ctx context.Context) bool {
	v, _ := ctx.Value(menuDescribedKey{}).(bool)
	return v
}

// describe puts in props the description desc of an item, when the nodes are
// of an index (menuDescribed).
func describe(ctx context.Context, props map[string]any, desc string) {
	if desc != "" && menuDescribed(ctx) {
		props["subtitle"] = desc
	}
}

// findMenuNode is the first node of nodes (or under them) that is.
func findMenuNode(nodes []*MenuNode, is func(n *MenuNode) bool) *MenuNode {
	for _, n := range nodes {
		if is(n) {
			return n
		}
		if f := findMenuNode(n.Children, is); f != nil {
			return f
		}
	}
	return nil
}

// index is the page of the group with no page of its own: its node of the
// menu of the request (is), its items — the menu's, the ones the user may
// see — as a tree, their descriptions under their titles; none, not found.
func (b *MenuGroupBuilder) index(ctx *web.EventContext, is func(n *MenuNode) bool) (r web.PageResponse, err error) {
	dctx := *ctx
	dctx.R = ctx.R.WithContext(context.WithValue(ctx.R.Context(), menuDescribedKey{}, true))
	nodes, _, _ := b.b.menuNodes(&dctx)
	n := findMenuNode(nodes, is)
	if n == nil || len(n.Children) == 0 {
		return b.b.DefaultNotFoundPageFunc(ctx)
	}
	r.PageTitle = b.TTitle(ctx.Context())
	body := h.Div().Class("pa-4").Attr("data-menu-group-index", b.name)
	if d := b.TDescription(ctx.Context()); d != "" {
		body.AppendChildren(h.Div(h.Text(d)).Class("text-body-2 text-medium-emphasis mb-4"))
	}
	body.AppendChildren(MenuIndexTree(n.Children))
	r.Body = body
	return
}

// MenuIndexTree is nodes of the menu as a tree of an index: open, each its
// icon, its title — a link to its page —, its description (subtitle).
func MenuIndexTree(nodes []*MenuNode) h.HTMLComponent {
	return VTreeview(
		h.Template(
			h.Img("").Attr(":src", "item.props.prependImage").
				Attr("v-if", "item.props && item.props.prependImage").
				Attr("alt", "").Style("width:24px;height:24px;object-fit:contain"),
			VIcon("").Attr(":icon", "item.props.prependIcon").
				Attr("v-else-if", "item.props && item.props.prependIcon"),
		).Attr("v-slot:prepend", "{ item }"),
		h.Template(
			h.A(h.Text("{{ item.title }}")).Attr(":href", "item.props.href || item.props.groupHref").
				Class("text-decoration-none text-high-emphasis font-weight-medium").
				Attr("v-if", "item.props && (item.props.href || item.props.groupHref)").
				Attr("@click.stop.prevent", "plaid().vars(vars).pushStateURL(item.props.href || item.props.groupHref).go()"),
			h.Span("{{ item.title }}").Class("font-weight-medium").Attr("v-else", true),
		).Attr("v-slot:title", "{ item }"),
	).
		Attr(":items", h.JSONString(indexNodes(nodes))).
		ItemTitle("title").
		ItemValue("value").
		ItemChildren("children").
		ItemProps(true).
		OpenAll(true).
		Lines("two").
		Class("rvq-menu-group-index")
}

// indexNodes are nodes of the menu for an index: the chips, the lazy
// loading and the behavior of the side menu left out.
func indexNodes(nodes []*MenuNode) []*MenuNode {
	out := make([]*MenuNode, 0, len(nodes))
	for _, n := range nodes {
		props := map[string]any{}
		for _, k := range []string{"href", "groupHref", "prependIcon", "prependImage", "subtitle"} {
			if v, ok := n.Props[k]; ok {
				props[k] = v
			}
		}
		out = append(out, &MenuNode{Title: n.Title, Value: n.Value, Props: props, Children: indexNodes(n.Children)})
	}
	return out
}

// groupLink is the title of a group in the side menu as a link to its page
// (props groupHref): its click opens the page, not the group.
func groupLink() *h.HTMLTagBuilder {
	return h.A(h.Text("{{ item.title }}")).Attr(":href", "item.props.groupHref").
		Class("rvq-menu-group-link text-decoration-none").Style("color: inherit").
		Attr("@click.stop.prevent", "plaid().vars(vars).pushStateURL(item.props.groupHref).go()")
}
