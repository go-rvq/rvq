package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VLocaleProviderBuilder struct {
	VTagBuilder[*VLocaleProviderBuilder]
}

func VLocaleProvider(children ...h.HTMLComponent) *VLocaleProviderBuilder {
	return VTag(&VLocaleProviderBuilder{}, "v-locale-provider", children...)
}

func (b *VLocaleProviderBuilder) Locale(v string) (r *VLocaleProviderBuilder) {
	b.Attr("locale", v)
	return b
}

func (b *VLocaleProviderBuilder) FallbackLocale(v string) (r *VLocaleProviderBuilder) {
	b.Attr("fallback-locale", v)
	return b
}

func (b *VLocaleProviderBuilder) Messages(v interface{}) (r *VLocaleProviderBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VLocaleProviderBuilder) Rtl(v bool) (r *VLocaleProviderBuilder) {
	b.Attr(":rtl", fmt.Sprint(v))
	return b
}

func (b *VLocaleProviderBuilder) On(name string, value string) (r *VLocaleProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLocaleProviderBuilder) Bind(name string, value string) (r *VLocaleProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
