package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VValidationBuilder struct {
	VTagBuilder[*VValidationBuilder]
}

func VValidation(children ...h.HTMLComponent) *VValidationBuilder {
	return VTag(&VValidationBuilder{}, "v-validation", children...)
}

func (b *VValidationBuilder) ModelValue(v interface{}) (r *VValidationBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) Error(v bool) (r *VValidationBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VValidationBuilder) Name(v string) (r *VValidationBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VValidationBuilder) Readonly(v bool) (r *VValidationBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VValidationBuilder) Disabled(v bool) (r *VValidationBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VValidationBuilder) Label(v string) (r *VValidationBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VValidationBuilder) ErrorMessages(v interface{}) (r *VValidationBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) MaxErrors(v interface{}) (r *VValidationBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) Rules(v interface{}) (r *VValidationBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) ValidateOn(v interface{}) (r *VValidationBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) ValidationValue(v interface{}) (r *VValidationBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VValidationBuilder) Focused(v bool) (r *VValidationBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VValidationBuilder) On(name string, value string) (r *VValidationBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VValidationBuilder) Bind(name string, value string) (r *VValidationBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VValidationBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VValidationBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VValidationBuilder) Slot(name string, child ...h.HTMLComponent) (r *VValidationBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VValidationBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VValidationBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VValidationBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VValidationBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VValidationBuilder) SlotDefault(child ...h.HTMLComponent) (r *VValidationBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VValidationBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VValidationBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
