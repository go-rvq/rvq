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

func (b *VOtpInputBuilder) Autofocus(v bool) (r *VOtpInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Divider(v string) (r *VOtpInputBuilder) {
	b.Attr("divider", v)
	return b
}

func (b *VOtpInputBuilder) FocusAll(v bool) (r *VOtpInputBuilder) {
	b.Attr(":focus-all", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Label(v string) (r *VOtpInputBuilder) {
	b.Attr("label", v)
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

func (b *VOtpInputBuilder) Placeholder(v string) (r *VOtpInputBuilder) {
	b.Attr("placeholder", v)
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

func (b *VOtpInputBuilder) Focused(v bool) (r *VOtpInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) BgColor(v string) (r *VOtpInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VOtpInputBuilder) Color(v string) (r *VOtpInputBuilder) {
	b.Attr("color", v)
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

func (b *VOtpInputBuilder) Error(v bool) (r *VOtpInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VOtpInputBuilder) Variant(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VOtpInputBuilder) Loading(v interface{}) (r *VOtpInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
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

func (b *VOtpInputBuilder) On(name string, value string) (r *VOtpInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VOtpInputBuilder) Bind(name string, value string) (r *VOtpInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
