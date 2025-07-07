package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VAppBarTitleBuilder struct {
	VTagBuilder[*VAppBarTitleBuilder]
}

func VAppBarTitle(children ...h.HTMLComponent) *VAppBarTitleBuilder {
	return VTag(&VAppBarTitleBuilder{}, "v-app-bar-title", children...)
}

func (b *VAppBarTitleBuilder) Text(v string) (r *VAppBarTitleBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAppBarTitleBuilder) Tag(v string) (r *VAppBarTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAppBarTitleBuilder) On(name string, value string) (r *VAppBarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBarTitleBuilder) Bind(name string, value string) (r *VAppBarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
