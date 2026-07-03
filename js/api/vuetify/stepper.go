package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperBuilder struct {
	VTagBuilder[*VStepperBuilder]
}

func VStepper(children ...h.HTMLComponent) *VStepperBuilder {
	return VTag(&VStepperBuilder{}, "v-stepper", children...)
}

func (b *VStepperBuilder) Flat(v bool) (r *VStepperBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Border(v interface{}) (r *VStepperBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) ModelValue(v interface{}) (r *VStepperBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Height(v interface{}) (r *VStepperBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MaxHeight(v interface{}) (r *VStepperBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MaxWidth(v interface{}) (r *VStepperBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MinHeight(v interface{}) (r *VStepperBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MinWidth(v interface{}) (r *VStepperBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Width(v interface{}) (r *VStepperBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Elevation(v interface{}) (r *VStepperBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Location(v interface{}) (r *VStepperBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Position(v interface{}) (r *VStepperBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Rounded(v interface{}) (r *VStepperBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Tile(v bool) (r *VStepperBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Tag(v interface{}) (r *VStepperBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Theme(v string) (r *VStepperBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperBuilder) Color(v string) (r *VStepperBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperBuilder) Disabled(v bool) (r *VStepperBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) SelectedClass(v string) (r *VStepperBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperBuilder) Max(v interface{}) (r *VStepperBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) BgColor(v string) (r *VStepperBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperBuilder) Mobile(v bool) (r *VStepperBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) MobileBreakpoint(v interface{}) (r *VStepperBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Multiple(v bool) (r *VStepperBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Mandatory(v interface{}) (r *VStepperBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Items(v interface{}) (r *VStepperBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) ItemTitle(v string) (r *VStepperBuilder) {
	b.Attr("item-title", v)
	return b
}

func (b *VStepperBuilder) ItemValue(v string) (r *VStepperBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VStepperBuilder) HideActions(v bool) (r *VStepperBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) AltLabels(v bool) (r *VStepperBuilder) {
	b.Attr(":alt-labels", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) CompleteIcon(v interface{}) (r *VStepperBuilder) {
	b.Attr(":complete-icon", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) EditIcon(v interface{}) (r *VStepperBuilder) {
	b.Attr(":edit-icon", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Editable(v bool) (r *VStepperBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) ErrorIcon(v interface{}) (r *VStepperBuilder) {
	b.Attr(":error-icon", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) NonLinear(v bool) (r *VStepperBuilder) {
	b.Attr(":non-linear", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) PrevText(v string) (r *VStepperBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperBuilder) NextText(v string) (r *VStepperBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperBuilder) On(name string, value string) (r *VStepperBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperBuilder) Bind(name string, value string) (r *VStepperBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VStepperBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VStepperBuilder) SlotActions(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VStepperBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VStepperBuilder) SlotHeader(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotHeaderItem(child ...h.HTMLComponent) {
	b.SetSlot("header-item", child...)
}

func (b *VStepperBuilder) SetScopedSlotHeaderItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header-item", scope, child...)
}

func (b *VStepperBuilder) SlotHeaderItem(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotHeaderItem(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotHeaderItem(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotHeaderItem(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotIcon(child ...h.HTMLComponent) {
	b.SetSlot("icon", child...)
}

func (b *VStepperBuilder) SetScopedSlotIcon(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("icon", scope, child...)
}

func (b *VStepperBuilder) SlotIcon(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotIcon(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotIcon(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotIcon(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VStepperBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VStepperBuilder) SlotTitle(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VStepperBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VStepperBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VStepperBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VStepperBuilder) SlotItem(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VStepperBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VStepperBuilder) SlotPrev(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VStepperBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VStepperBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VStepperBuilder) SlotNext(child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VStepperBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VStepperBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
