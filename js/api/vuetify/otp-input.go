package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VOtpInputBuilder struct {
	VTagBuilder[*VOtpInputBuilder]
}

func VOtpInput(children ...h.HTMLComponent) *VOtpInputBuilder {
	return VTag(&VOtpInputBuilder{}, "v-otp-input", children...)
}

func (b *VOtpInputBuilder) Length(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Type(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":type", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) ModelValue(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Error(v bool) (r *VOtpInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Height(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) MaxHeight(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) MaxWidth(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) MinHeight(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) MinWidth(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Width(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Rounded(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Theme(v string) (r *VOtpInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VOtpInputBuilder) Color(v string) (r *VOtpInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VOtpInputBuilder) Variant(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) BaseColor(v string) (r *VOtpInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VOtpInputBuilder) Disabled(v bool) (r *VOtpInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Loading(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Label(v string) (r *VOtpInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VOtpInputBuilder) BgColor(v string) (r *VOtpInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VOtpInputBuilder) Divider(v string) (r *VOtpInputBuilder) {
	b.Attr("divider", v)
	return b
}

func (b *VOtpInputBuilder) Placeholder(v string) (r *VOtpInputBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VOtpInputBuilder) Focused(v bool) (r *VOtpInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Autofocus(v bool) (r *VOtpInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) FocusAll(v bool) (r *VOtpInputBuilder) {
	b.Attr(":focus-all", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) On(name string, value string) (r *VOtpInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VOtpInputBuilder) Bind(name string, value string) (r *VOtpInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VOtpInputBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VOtpInputBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VOtpInputBuilder) Slot(name string, child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VOtpInputBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VOtpInputBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VOtpInputBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VOtpInputBuilder) SlotDefault(child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VOtpInputBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VOtpInputBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VOtpInputBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VOtpInputBuilder) SlotLoader(child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VOtpInputBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VOtpInputBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}
