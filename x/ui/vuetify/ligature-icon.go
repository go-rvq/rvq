package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLigatureIconBuilder struct {
	VTagBuilder[*VLigatureIconBuilder]
}

func VLigatureIcon(children ...h.HTMLComponent) *VLigatureIconBuilder {
	return VTag(&VLigatureIconBuilder{}, "v-ligature-icon", children...)
}

func (b *VLigatureIconBuilder) Icon(v interface{}) (r *VLigatureIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VLigatureIconBuilder) Tag(v string) (r *VLigatureIconBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VLigatureIconBuilder) On(name string, value string) (r *VLigatureIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLigatureIconBuilder) Bind(name string, value string) (r *VLigatureIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
