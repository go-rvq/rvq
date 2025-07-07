package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VConfirmEditBuilder struct {
	VTagBuilder[*VConfirmEditBuilder]
}

func VConfirmEdit(children ...h.HTMLComponent) *VConfirmEditBuilder {
	return VTag(&VConfirmEditBuilder{}, "v-confirm-edit", children...)
}

func (b *VConfirmEditBuilder) ModelValue(v interface{}) (r *VConfirmEditBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VConfirmEditBuilder) Color(v string) (r *VConfirmEditBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VConfirmEditBuilder) CancelText(v string) (r *VConfirmEditBuilder) {
	b.Attr("cancel-text", v)
	return b
}

func (b *VConfirmEditBuilder) OkText(v string) (r *VConfirmEditBuilder) {
	b.Attr("ok-text", v)
	return b
}

func (b *VConfirmEditBuilder) On(name string, value string) (r *VConfirmEditBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VConfirmEditBuilder) Bind(name string, value string) (r *VConfirmEditBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
