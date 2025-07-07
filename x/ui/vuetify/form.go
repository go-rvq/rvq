package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFormBuilder struct {
	VTagBuilder[*VFormBuilder]
}

func VForm(children ...h.HTMLComponent) *VFormBuilder {
	return VTag(&VFormBuilder{}, "v-form", children...)
}

func (b *VFormBuilder) ModelValue(v bool) (r *VFormBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) Disabled(v bool) (r *VFormBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) FastFail(v bool) (r *VFormBuilder) {
	b.Attr(":fast-fail", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) Readonly(v bool) (r *VFormBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) ValidateOn(v interface{}) (r *VFormBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VFormBuilder) On(name string, value string) (r *VFormBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFormBuilder) Bind(name string, value string) (r *VFormBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
