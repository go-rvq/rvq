package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VXTreeDraggableBuilder is a tree whose nodes are dragged and dropped —
// among their siblings, or into another node, as a child — as the folders of a
// file system (vx-tree-draggable). Its items are changed in place; "change"
// tells them, after each drop. A node's content is the slot "item" ({ item,
// depth }), its title (TitleKey) without it.
type VXTreeDraggableBuilder struct {
	v.VTagBuilder[*VXTreeDraggableBuilder]
}

// VXTreeDraggable is a tree of nodes dragged and dropped.
func VXTreeDraggable(children ...h.HTMLComponent) *VXTreeDraggableBuilder {
	return v.VTag(&VXTreeDraggableBuilder{}, "vx-tree-draggable", children...)
}

// Items are the nodes, each with its children (ChildrenKey).
func (b *VXTreeDraggableBuilder) Items(v interface{}) *VXTreeDraggableBuilder {
	return b.Attr(":items", h.JSONString(v))
}

// ItemsVar is the expression of the nodes (a reactive array: changed in
// place).
func (b *VXTreeDraggableBuilder) ItemsVar(v string) *VXTreeDraggableBuilder {
	return b.Attr(":items", v)
}

// ItemKey is the key of a node's id ("id").
func (b *VXTreeDraggableBuilder) ItemKey(v string) *VXTreeDraggableBuilder {
	return b.Attr("item-key", v)
}

// TitleKey is the key of a node's title ("title").
func (b *VXTreeDraggableBuilder) TitleKey(v string) *VXTreeDraggableBuilder {
	return b.Attr("title-key", v)
}

// ChildrenKey is the key of a node's children ("children").
func (b *VXTreeDraggableBuilder) ChildrenKey(v string) *VXTreeDraggableBuilder {
	return b.Attr("children-key", v)
}

// Group is the group of the drag: trees of the same group exchange nodes.
func (b *VXTreeDraggableBuilder) Group(v string) *VXTreeDraggableBuilder {
	return b.Attr("group", v)
}

// Disabled turns the drag off.
func (b *VXTreeDraggableBuilder) Disabled(v bool) *VXTreeDraggableBuilder {
	return b.Attr(":disabled", h.JSONString(v))
}
