package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSvgIconBuilder struct {
	VTagBuilder[*VSvgIconBuilder]
}

func VSvgIcon(children ...h.HTMLComponent) *VSvgIconBuilder {
	return VTag(&VSvgIconBuilder{}, "v-svg-icon", children...)
}

func (b *VSvgIconBuilder) Icon(v interface{}) (r *VSvgIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VSvgIconBuilder) Tag(v string) (r *VSvgIconBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSvgIconBuilder) On(name string, value string) (r *VSvgIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSvgIconBuilder) Bind(name string, value string) (r *VSvgIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
