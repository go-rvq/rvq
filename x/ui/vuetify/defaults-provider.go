package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDefaultsProviderBuilder struct {
	VTagBuilder[*VDefaultsProviderBuilder]
}

func VDefaultsProvider(children ...h.HTMLComponent) *VDefaultsProviderBuilder {
	return VTag(&VDefaultsProviderBuilder{}, "v-defaults-provider", children...)
}

func (b *VDefaultsProviderBuilder) Disabled(v bool) (r *VDefaultsProviderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDefaultsProviderBuilder) Reset(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":reset", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) Root(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":root", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) Scoped(v bool) (r *VDefaultsProviderBuilder) {
	b.Attr(":scoped", fmt.Sprint(v))
	return b
}

func (b *VDefaultsProviderBuilder) Defaults(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":defaults", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) On(name string, value string) (r *VDefaultsProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDefaultsProviderBuilder) Bind(name string, value string) (r *VDefaultsProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
