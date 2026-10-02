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
	ID    string `json:"value"`
	Title string `json:"title"`
	Icon  string `json:"prependIcon,omitempty"`
	// Href is the URL of its document: the item is a link to it.
	Href     string  `json:"href,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

// Tree is the tree of the documentation: the side menu's, in its order — its
// groups, the models and the pages in them —, each model with its actions and
// the models nested in it.
func (b *Builder) Tree(ctx *web.EventContext) []*Node {
	return b.groupNodes(b.p.MenuTree(), ctx)
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
			nodes = append(nodes, &Node{ID: "pages/" + strings.Trim(page.Path(), "/"),
				Title: page.TTitle(ctx.Context()), Icon: "mdi-file-document-outline"})
		}
	}
	return
}

// canList says the request may see the listing of mb.
func (b *Builder) canList(mb *presets.ModelBuilder, ctx *web.EventContext) bool {
	return !mb.Permissioner().ReqLister(ctx.R).Denied()
}

// modelTree is the node of mb, its document at id: its actions, and the
// models nested in it under it.
func (b *Builder) modelTree(mb *presets.ModelBuilder, id string, ctx *web.EventContext) *Node {
	n := &Node{ID: id, Title: mb.TTitlePlural(ctx.Context()), Icon: mb.GetMenuIcon()}
	for _, a := range detailingActions(mb) {
		n.Children = append(n.Children, &Node{ID: id + "/actions/" + a.Name(),
			Title: a.RequestTitle(mb, ctx.Context()), Icon: "mdi-gesture-tap"})
	}
	for _, a := range mb.Listing().GetBulkActions() {
		n.Children = append(n.Children, &Node{ID: id + "/actions/" + a.Name(),
			Title: a.RequestTitle(ctx.Context()), Icon: "mdi-checkbox-multiple-marked-outline"})
	}
	for _, c := range mb.Children() {
		n.Children = append(n.Children, b.modelTree(c, id+"/children/"+childKey(c), ctx))
	}
	return n
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

// ancestors are the ids of the nodes above id, from the root: the nodes the
// tree opens to show it.
func ancestors(nodes []*Node, id string) []string {
	for _, n := range nodes {
		if n.ID == id {
			return []string{}
		}
		if a := ancestors(n.Children, id); a != nil {
			return append([]string{n.ID}, a...)
		}
	}
	return nil
}
