package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VListGroupBuilder struct {
	VTagBuilder[*VListGroupBuilder]
}

func VListGroup(children ...h.HTMLComponent) *VListGroupBuilder {
	return VTag(&VListGroupBuilder{}, "v-list-group", children...)
}

func (b *VListGroupBuilder) ActiveColor(v string) (r *VListGroupBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VListGroupBuilder) BaseColor(v string) (r *VListGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VListGroupBuilder) Color(v string) (r *VListGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListGroupBuilder) CollapseIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) ExpandIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) PrependIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) AppendIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) Fluid(v bool) (r *VListGroupBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VListGroupBuilder) Subgroup(v bool) (r *VListGroupBuilder) {
	b.Attr(":subgroup", fmt.Sprint(v))
	return b
}

func (b *VListGroupBuilder) Title(v string) (r *VListGroupBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VListGroupBuilder) Value(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) Tag(v string) (r *VListGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListGroupBuilder) On(name string, value string) (r *VListGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListGroupBuilder) Bind(name string, value string) (r *VListGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
