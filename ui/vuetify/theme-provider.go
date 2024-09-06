package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VThemeProviderBuilder struct {
	VTagBuilder[*VThemeProviderBuilder]
}

func VThemeProvider(children ...h.HTMLComponent) *VThemeProviderBuilder {
	return VTag(&VThemeProviderBuilder{}, "v-theme-provider", children...)
}

func (b *VThemeProviderBuilder) WithBackground(v bool) (r *VThemeProviderBuilder) {
	b.Attr(":with-background", fmt.Sprint(v))
	return b
}

func (b *VThemeProviderBuilder) Theme(v string) (r *VThemeProviderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VThemeProviderBuilder) Tag(v string) (r *VThemeProviderBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VThemeProviderBuilder) On(name string, value string) (r *VThemeProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VThemeProviderBuilder) Bind(name string, value string) (r *VThemeProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
