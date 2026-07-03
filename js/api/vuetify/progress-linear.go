package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VProgressLinearBuilder struct {
	VTagBuilder[*VProgressLinearBuilder]
}

func VProgressLinear(children ...h.HTMLComponent) *VProgressLinearBuilder {
	return VTag(&VProgressLinearBuilder{}, "v-progress-linear", children...)
}

func (b *VProgressLinearBuilder) ModelValue(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Reverse(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Height(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Location(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Absolute(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Rounded(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Tile(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Tag(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Theme(v string) (r *VProgressLinearBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VProgressLinearBuilder) Color(v string) (r *VProgressLinearBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VProgressLinearBuilder) Active(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Max(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) BgColor(v string) (r *VProgressLinearBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VProgressLinearBuilder) Opacity(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Indeterminate(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":indeterminate", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Stream(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":stream", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) BgOpacity(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":bg-opacity", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) BufferValue(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":buffer-value", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) BufferColor(v string) (r *VProgressLinearBuilder) {
	b.Attr("buffer-color", v)
	return b
}

func (b *VProgressLinearBuilder) BufferOpacity(v interface{}) (r *VProgressLinearBuilder) {
	b.Attr(":buffer-opacity", h.JSONString(v))
	return b
}

func (b *VProgressLinearBuilder) Clickable(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":clickable", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) Striped(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":striped", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) RoundedBar(v bool) (r *VProgressLinearBuilder) {
	b.Attr(":rounded-bar", fmt.Sprint(v))
	return b
}

func (b *VProgressLinearBuilder) On(name string, value string) (r *VProgressLinearBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VProgressLinearBuilder) Bind(name string, value string) (r *VProgressLinearBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VProgressLinearBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VProgressLinearBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VProgressLinearBuilder) Slot(name string, child ...h.HTMLComponent) (r *VProgressLinearBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VProgressLinearBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VProgressLinearBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VProgressLinearBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VProgressLinearBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VProgressLinearBuilder) SlotDefault(child ...h.HTMLComponent) (r *VProgressLinearBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VProgressLinearBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VProgressLinearBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
