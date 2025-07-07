package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VInputBuilder struct {
	VTagBuilder[*VInputBuilder]
}

func VInput(children ...h.HTMLComponent) *VInputBuilder {
	return VTag(&VInputBuilder{}, "v-input", children...)
}

func (b *VInputBuilder) Id(v string) (r *VInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VInputBuilder) AppendIcon(v interface{}) (r *VInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VInputBuilder) CenterAffix(v bool) (r *VInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) PrependIcon(v interface{}) (r *VInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VInputBuilder) HideSpinButtons(v bool) (r *VInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Hint(v string) (r *VInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VInputBuilder) PersistentHint(v bool) (r *VInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Messages(v interface{}) (r *VInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Direction(v interface{}) (r *VInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Density(v interface{}) (r *VInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MaxWidth(v interface{}) (r *VInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MinWidth(v interface{}) (r *VInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Width(v interface{}) (r *VInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Theme(v string) (r *VInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VInputBuilder) Disabled(v bool) (r *VInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Error(v bool) (r *VInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) ErrorMessages(v interface{}) (r *VInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MaxErrors(v interface{}) (r *VInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Name(v string) (r *VInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VInputBuilder) Label(v string) (r *VInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VInputBuilder) Readonly(v bool) (r *VInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Rules(v interface{}) (r *VInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ModelValue(v interface{}) (r *VInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ValidateOn(v interface{}) (r *VInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ValidationValue(v interface{}) (r *VInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Focused(v bool) (r *VInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) HideDetails(v interface{}) (r *VInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VInputBuilder) On(name string, value string) (r *VInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VInputBuilder) Bind(name string, value string) (r *VInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
