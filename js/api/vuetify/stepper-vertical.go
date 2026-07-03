package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperVerticalBuilder struct {
	VTagBuilder[*VStepperVerticalBuilder]
}

func VStepperVertical(children ...h.HTMLComponent) *VStepperVerticalBuilder {
	return VTag(&VStepperVerticalBuilder{}, "v-stepper-vertical", children...)
}

func (b *VStepperVerticalBuilder) Flat(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) ModelValue(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Elevation(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Rounded(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Tile(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Tag(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Theme(v string) (r *VStepperVerticalBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperVerticalBuilder) Color(v string) (r *VStepperVerticalBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalBuilder) Variant(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Readonly(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Ripple(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Disabled(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) SelectedClass(v string) (r *VStepperVerticalBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperVerticalBuilder) Max(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) BgColor(v string) (r *VStepperVerticalBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperVerticalBuilder) Mobile(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) MobileBreakpoint(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Multiple(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Mandatory(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Eager(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Items(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) ExpandIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) CollapseIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) ItemTitle(v string) (r *VStepperVerticalBuilder) {
	b.Attr("item-title", v)
	return b
}

func (b *VStepperVerticalBuilder) ItemValue(v string) (r *VStepperVerticalBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VStepperVerticalBuilder) HideActions(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) AltLabels(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":alt-labels", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) CompleteIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":complete-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) EditIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":edit-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Editable(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) ErrorIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":error-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) NonLinear(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":non-linear", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) PrevText(v string) (r *VStepperVerticalBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperVerticalBuilder) NextText(v string) (r *VStepperVerticalBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperVerticalBuilder) Focusable(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) On(name string, value string) (r *VStepperVerticalBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalBuilder) Bind(name string, value string) (r *VStepperVerticalBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperVerticalBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperVerticalBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperVerticalBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotActions(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotIcon(child ...h.HTMLComponent) {
	b.SetSlot("icon", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotIcon(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("icon", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotIcon(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotIcon(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotIcon(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotIcon(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotTitle(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotPrev(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VStepperVerticalBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VStepperVerticalBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VStepperVerticalBuilder) SlotNext(child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VStepperVerticalBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VStepperVerticalBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
