package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VFieldBuilder struct {
	VTagBuilder[*VFieldBuilder]
}

func VField(children ...h.HTMLComponent) *VFieldBuilder {
	return VTag(&VFieldBuilder{}, "v-field", children...)
}

func (b *VFieldBuilder) Label(v string) (r *VFieldBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VFieldBuilder) Id(v string) (r *VFieldBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VFieldBuilder) Focused(v bool) (r *VFieldBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Reverse(v bool) (r *VFieldBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Flat(v bool) (r *VFieldBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) AppendInnerIcon(v interface{}) (r *VFieldBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) BgColor(v string) (r *VFieldBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VFieldBuilder) Clearable(v bool) (r *VFieldBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) ClearIcon(v interface{}) (r *VFieldBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) Active(v bool) (r *VFieldBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) CenterAffix(v bool) (r *VFieldBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Color(v string) (r *VFieldBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFieldBuilder) BaseColor(v string) (r *VFieldBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFieldBuilder) Dirty(v bool) (r *VFieldBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Disabled(v bool) (r *VFieldBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Error(v bool) (r *VFieldBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) PersistentClear(v bool) (r *VFieldBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) PrependInnerIcon(v interface{}) (r *VFieldBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) SingleLine(v bool) (r *VFieldBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Variant(v interface{}) (r *VFieldBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) Loading(v interface{}) (r *VFieldBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) Rounded(v interface{}) (r *VFieldBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) Tile(v bool) (r *VFieldBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Theme(v string) (r *VFieldBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFieldBuilder) ModelValue(v interface{}) (r *VFieldBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) On(name string, value string) (r *VFieldBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFieldBuilder) Bind(name string, value string) (r *VFieldBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
