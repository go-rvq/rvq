package userdocs

import (
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/iancoleman/strcase"
)

// Node is a node of the tree of the documentation: a group of the menu, a
// model and the models nested in it, an action, a page. Its ID is its
// document's path without the file (DocFile).
type Node struct {
	ID       string  `json:"value"`
	Title    string  `json:"title"`
	Icon     string  `json:"prependIcon,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

// Tree is the tree of the documentation: the side menu's, in its order — its
// groups, the models and the pages in them —, each model with its actions and
// the models nested in it; the custom documents where they were put (Custom).
func (b *Builder) Tree(ctx *web.EventContext) []*Node {
	return b.withCustom(b.groupNodes(b.p.MenuTree(), ctx), ctx)
}

func (b *Builder) groupNodes(g *presets.MenuGroupBuilder, ctx *web.EventContext) (nodes []*Node) {
	for _, it := range g.Items() {
		switch it.Type {
		case presets.MenuItemGroup:
			sub := it.Group()
			if sub == nil {
				continue
			}
			children := b.groupNodes(sub, ctx)
			if len(children) == 0 {
				continue
			}
			// its permissions, after what it holds
			children = append(children, permissionsNode("groups/"+sub.Name(), GetMessages(ctx.Context())))
			nodes = append(nodes, &Node{ID: "groups/" + sub.Name(), Title: sub.TTitle(ctx.Context()),
				Icon: "mdi-folder-outline", Children: children})
		case presets.MenuItemModel:
			mb, _ := it.Value.(*presets.ModelBuilder)
			if mb == nil || mb == b.mb || !mb.IsInMenu() || !b.canList(mb, ctx) {
				continue
			}
			nodes = append(nodes, b.modelTree(mb, mb.MenuID(), ctx))
		case presets.MenuItemPage:
			page, _ := it.Value.(*presets.HttpPageBuilder)
			if page == nil {
				continue
			}
			id := "pages/" + strings.Trim(page.Path(), "/")
			nodes = append(nodes, &Node{ID: id, Title: page.TTitle(ctx.Context()), Icon: "mdi-file-document-outline",
				Children: []*Node{permissionsNode(id, GetMessages(ctx.Context()))}})
		}
	}
	return
}

// canList says the request may see the listing of mb.
func (b *Builder) canList(mb *presets.ModelBuilder, ctx *web.EventContext) bool {
	return !mb.Permissioner().ReqLister(ctx.R).Denied()
}

// modelTree is the node of mb, its document at id: its forms (formNodes);
// its actions, under Actions; its pages, under Pages; the models nested in it.
func (b *Builder) modelTree(mb *presets.ModelBuilder, id string, ctx *web.EventContext) *Node {
	msgs := GetMessages(ctx.Context())
	n := &Node{ID: id, Title: mb.TTitleAuto(ctx.Context()), Icon: mb.GetMenuIcon()}

	// its forms: of a new record, of an edit, the detail
	n.Children = append(n.Children, formNodes(mb, id, msgs)...)
	// its permissions: of its records, fields, sections, actions, pages
	n.Children = append(n.Children, permissionsNode(id, msgs))

	actions := &Node{ID: id + "/actions", Title: msgs.Actions, Icon: "mdi-gesture-tap"}
	for _, a := range detailingActions(mb) {
		actions.Children = append(actions.Children, &Node{ID: id + "/actions/" + a.Name(),
			Title: a.RequestTitle(mb, ctx.Context())})
	}
	for _, a := range mb.Listing().GetBulkActions() {
		actions.Children = append(actions.Children, &Node{ID: id + "/actions/" + a.Name(),
			Title: a.RequestTitle(ctx.Context())})
	}
	for _, a := range rowMenuItems(mb) {
		if findNode(actions.Children, id+"/actions/"+a.Name()) == nil {
			actions.Children = append(actions.Children, &Node{ID: id + "/actions/" + a.Name(),
				Title: a.TTitle(ctx.Context())})
		}
	}
	if len(actions.Children) > 0 {
		n.Children = append(n.Children, actions)
	}

	pages := &Node{ID: id + "/pages", Title: msgs.Pages, Icon: "mdi-file-document-multiple-outline"}
	for _, p := range modelPages(mb) {
		pages.Children = append(pages.Children, &Node{ID: id + "/pages/" + pageKey(p.Path()),
			Title: p.TTitle(ctx.Context())})
	}
	if len(pages.Children) > 0 {
		n.Children = append(n.Children, pages)
	}

	for _, c := range mb.Children() {
		n.Children = append(n.Children, b.modelTree(c, id+"/children/"+childKey(c), ctx))
	}
	return n
}

// modelPages are the pages of mb: of its listing, and of its detail — none
// when it has no detail: asking for it would make one.
func modelPages(mb *presets.ModelBuilder) (pages []*presets.HttpPageBuilder) {
	pages = append(pages, mb.Listing().PagesRegistrator().HttpPages()...)
	if mb.HasDetailing() {
		pages = append(pages, mb.Detailing().PagesRegistrator().HttpPages()...)
	}
	return
}

// pageKey is the name of a page in the path of its document: its path, its
// "/" as "_" (report, counters_detail).
func pageKey(path string) string {
	return strings.ReplaceAll(strings.Trim(path, "/{}"), "/", "_")
}

// rowMenuItems are the items of the menu of a record of mb's listing — the
// nested models it opens left out: they are nodes of their own; the deletion
// left out where there is none. A singleton has no listing: no items.
func rowMenuItems(mb *presets.ModelBuilder) (items []*presets.RowMenuItemBuilder) {
	if mb.GetSingleton() {
		return
	}
	for _, it := range mb.Listing().RowMenu().Items() {
		if it.Child() == nil && (it.Name() != "Delete" || !mb.DeletingDisabled()) {
			items = append(items, it)
		}
	}
	return
}

// detailingActions are the actions of the detail of mb — none when it has no
// detail: asking for it would make one.
func detailingActions(mb *presets.ModelBuilder) []*presets.ActionBuilder {
	if !mb.HasDetailing() {
		return nil
	}
	return mb.Detailing().GetActions()
}

// childKey is the name of a model nested in another, in the path of its
// documents: its id (the field it is, Modules), else its URI's name, in
// snake case (modules).
func childKey(mb *presets.ModelBuilder) string {
	if id := mb.MenuID(); id != "" && !strings.Contains(id, ".") {
		return strcase.ToSnake(id)
	}
	return strcase.ToSnake(mb.UriName())
}

// modelNode is the node of mb: its id, under its parents'.
func (b *Builder) modelNode(mb *presets.ModelBuilder) string {
	if p := mb.Parent(); p != nil {
		return b.modelNode(p) + "/children/" + childKey(mb)
	}
	return mb.MenuID()
}

// findNode is the node id of nodes, or nil.
func findNode(nodes []*Node, id string) *Node {
	for _, n := range nodes {
		if n.ID == id {
			return n
		}
		if f := findNode(n.Children, id); f != nil {
			return f
		}
	}
	return nil
}
