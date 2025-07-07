package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTreeviewGroupBuilder struct {
	VTagBuilder[*VTreeviewGroupBuilder]
}

func VTreeviewGroup(children ...h.HTMLComponent) *VTreeviewGroupBuilder {
	return VTag(&VTreeviewGroupBuilder{}, "v-treeview-group", children...)
}

func (b *VTreeviewGroupBuilder) ActiveColor(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VTreeviewGroupBuilder) BaseColor(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTreeviewGroupBuilder) Color(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTreeviewGroupBuilder) CollapseIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) ExpandIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) PrependIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) AppendIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) Fluid(v bool) (r *VTreeviewGroupBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VTreeviewGroupBuilder) Title(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VTreeviewGroupBuilder) Value(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) Tag(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTreeviewGroupBuilder) On(name string, value string) (r *VTreeviewGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTreeviewGroupBuilder) Bind(name string, value string) (r *VTreeviewGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
