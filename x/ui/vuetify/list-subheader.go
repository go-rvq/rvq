package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListSubheaderBuilder struct {
	VTagBuilder[*VListSubheaderBuilder]
}

func VListSubheader(children ...h.HTMLComponent) *VListSubheaderBuilder {
	return VTag(&VListSubheaderBuilder{}, "v-list-subheader", children...)
}

func (b *VListSubheaderBuilder) Color(v string) (r *VListSubheaderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListSubheaderBuilder) Inset(v bool) (r *VListSubheaderBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VListSubheaderBuilder) Sticky(v bool) (r *VListSubheaderBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VListSubheaderBuilder) Title(v string) (r *VListSubheaderBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VListSubheaderBuilder) Tag(v string) (r *VListSubheaderBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListSubheaderBuilder) On(name string, value string) (r *VListSubheaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListSubheaderBuilder) Bind(name string, value string) (r *VListSubheaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
