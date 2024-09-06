package vuetify

import (
	h "github.com/theplant/htmlgo"

	"fmt"
)

type VToolbarTitleBuilder struct {
	VTagBuilder[*VToolbarTitleBuilder]
}

func VToolbarTitle(text string, children ...h.HTMLComponent) *VToolbarTitleBuilder {
	return VTag(&VToolbarTitleBuilder{}, "v-toolbar-title", children...).Text(text)
}

func (b *VToolbarTitleBuilder) Text(v string) (r *VToolbarTitleBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VToolbarTitleBuilder) Tag(v string) (r *VToolbarTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VToolbarTitleBuilder) On(name string, value string) (r *VToolbarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarTitleBuilder) Bind(name string, value string) (r *VToolbarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
