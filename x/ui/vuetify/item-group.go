package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VItemGroupBuilder struct {
	VTagBuilder[*VItemGroupBuilder]
}

func VItemGroup(children ...h.HTMLComponent) *VItemGroupBuilder {
	return VTag(&VItemGroupBuilder{}, "v-item-group", children...)
}

func (b *VItemGroupBuilder) ModelValue(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) Multiple(v bool) (r *VItemGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VItemGroupBuilder) Max(v int) (r *VItemGroupBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VItemGroupBuilder) SelectedClass(v string) (r *VItemGroupBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VItemGroupBuilder) Disabled(v bool) (r *VItemGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VItemGroupBuilder) Mandatory(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) Tag(v string) (r *VItemGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VItemGroupBuilder) Theme(v string) (r *VItemGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VItemGroupBuilder) On(name string, value string) (r *VItemGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VItemGroupBuilder) Bind(name string, value string) (r *VItemGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
