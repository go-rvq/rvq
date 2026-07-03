package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperItemBuilder struct {
	VTagBuilder[*VStepperItemBuilder]
}

func VStepperItem(children ...h.HTMLComponent) *VStepperItemBuilder {
	return VTag(&VStepperItemBuilder{}, "v-stepper-item", children...)
}

func (b *VStepperItemBuilder) Error(v bool) (r *VStepperItemBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) Title(v string) (r *VStepperItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VStepperItemBuilder) Value(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Subtitle(v string) (r *VStepperItemBuilder) {
	b.Attr("subtitle", v)
	return b
}

func (b *VStepperItemBuilder) Disabled(v bool) (r *VStepperItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) Color(v string) (r *VStepperItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperItemBuilder) Rules(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Icon(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) SelectedClass(v string) (r *VStepperItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperItemBuilder) Ripple(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Complete(v bool) (r *VStepperItemBuilder) {
	b.Attr(":complete", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) CompleteIcon(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":complete-icon", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Editable(v bool) (r *VStepperItemBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) EditIcon(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":edit-icon", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) ErrorIcon(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":error-icon", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) On(name string, value string) (r *VStepperItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperItemBuilder) Bind(name string, value string) (r *VStepperItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VStepperItemBuilder) SetSlotIcon(child ...h.HTMLComponent) {
	b.SetSlot("icon", child...)
}

func (b *VStepperItemBuilder) SetScopedSlotIcon(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("icon", scope, child...)
}

func (b *VStepperItemBuilder) SlotIcon(child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetSlotIcon(child...)
	return b
}

func (b *VStepperItemBuilder) ScopedSlotIcon(scope string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetScopedSlotIcon(scope, child...)
	return b
}

func (b *VStepperItemBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VStepperItemBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VStepperItemBuilder) SlotTitle(child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VStepperItemBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VStepperItemBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VStepperItemBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VStepperItemBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VStepperItemBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VStepperItemBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}
