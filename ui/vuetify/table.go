package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTableBuilder struct {
	VTagBuilder[*VTableBuilder]
}

func (b *VTableBuilder) FixedHeader(v bool) (r *VTableBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) FixedFooter(v bool) (r *VTableBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) Height(v interface{}) (r *VTableBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTableBuilder) Hover(v bool) (r *VTableBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) Density(v interface{}) (r *VTableBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTableBuilder) Tag(v string) (r *VTableBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTableBuilder) Theme(v string) (r *VTableBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTableBuilder) On(name string, value string) (r *VTableBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTableBuilder) Bind(name string, value string) (r *VTableBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
