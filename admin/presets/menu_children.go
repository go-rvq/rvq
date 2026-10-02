package presets

import (
	"encoding/json"

	"github.com/go-rvq/rvq/web"
)

// MenuNode is one node of the data-driven side menu (a VTreeview item): a group
// (has Children) or a leaf (Value is the target path). Props carries the item's
// component props (the prepend icon, and href on leaves).
type MenuNode struct {
	Title    string         `json:"title"`
	Value    string         `json:"value"`
	Props    map[string]any `json:"props,omitempty"`
	Children []*MenuNode    `json:"children,omitempty"`
	// lazy: its children come when it is opened — an empty array, which is
	// what tells VTreeview to load them
	lazy bool
}

func (n *MenuNode) MarshalJSON() ([]byte, error) {
	type node MenuNode
	if n.lazy && len(n.Children) == 0 {
		return json.Marshal(struct {
			*node
			Children []*MenuNode `json:"children"`
		}{(*node)(n), []*MenuNode{}})
	}
	return json.Marshal((*node)(n))
}

// has says the node or one under it has the value v.
func (n *MenuNode) has(v string) bool {
	return n.Value == v || menuPath(n.Children, v) != nil
}

// menuPath are the values of the nodes above the one of value v, from the
// top — nil when no node has it.
func menuPath(nodes []*MenuNode, v string) []string {
	for _, n := range nodes {
		if n.Value == v {
			return []string{}
		}
		if p := menuPath(n.Children, v); p != nil {
			return append([]string{n.Value}, p...)
		}
	}
	return nil
}

// MenuChildrenEvent is the event of a model that answers the children of its
// node of the menu (ModelBuilder.MenuChildren), in the response's data.
const MenuChildrenEvent = "presets_MenuChildren"

// MenuChildrenFunc are the children of a model's node of the menu, and the
// value of the one the request shows ("" for none).
type MenuChildrenFunc func(ctx *web.EventContext) (children []*MenuNode, active string)

// MenuChildren gives the node of the model in the menu children: those of f,
// loaded when the node is opened (MenuChildrenEvent), or with the page when
// the page is the model's — the one it shows active, the way to it open. A
// child whose value is a path (href in its props) opens that page, as an item
// of the menu does.
func (mb *ModelBuilder) MenuChildren(f MenuChildrenFunc) *ModelBuilder {
	mb.menuChildren = f
	mb.RegisterEventFunc(MenuChildrenEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		children, _ := f(ctx)
		if children == nil {
			children = []*MenuNode{}
		}
		r.Data = children
		return
	})
	return mb
}

// menuLoadChildrenScript loads the children of a node of the menu, from the
// event its props name (loadChildren).
const menuLoadChildrenScript = `(item) => {
	const n = item && item.raw ? item.raw : item;
	const l = n && n.props && n.props.loadChildren;
	if (!l) return Promise.resolve();
	return plaid().vars(vars).eventFunc(l.event).url(l.url).go().then((r) => {
		n.children.splice(0, n.children.length, ...((r && r.data) || []));
	});
}`
