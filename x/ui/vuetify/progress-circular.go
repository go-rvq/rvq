package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VProgressCircularBuilder struct {
	VTagBuilder[*VProgressCircularBuilder]
}

func VProgressCircular(children ...h.HTMLComponent) *VProgressCircularBuilder {
	return VTag(&VProgressCircularBuilder{}, "v-progress-circular", children...)
}

func (b *VProgressCircularBuilder) BgColor(v string) (r *VProgressCircularBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VProgressCircularBuilder) Color(v string) (r *VProgressCircularBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VProgressCircularBuilder) ModelValue(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Rotate(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":rotate", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Width(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Size(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Tag(v string) (r *VProgressCircularBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VProgressCircularBuilder) Theme(v string) (r *VProgressCircularBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VProgressCircularBuilder) Indeterminate(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":indeterminate", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) On(name string, value string) (r *VProgressCircularBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VProgressCircularBuilder) Bind(name string, value string) (r *VProgressCircularBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
