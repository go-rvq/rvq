package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperVerticalItemBuilder struct {
	VTagBuilder[*VStepperVerticalItemBuilder]
}

func VStepperVerticalItem(children ...h.HTMLComponent) *VStepperVerticalItemBuilder {
	return VTag(&VStepperVerticalItemBuilder{}, "v-stepper-vertical-item", children...)
}

func (b *VStepperVerticalItemBuilder) Tag(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Error(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Title(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Value(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Text(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Eager(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Subtitle(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("subtitle", v)
	return b
}

func (b *VStepperVerticalItemBuilder) BgColor(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Disabled(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) ExpandIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) CollapseIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Height(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) MaxHeight(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) MaxWidth(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) MinHeight(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) MinWidth(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Width(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Elevation(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Rounded(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Tile(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Color(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Static(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Readonly(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Rules(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Icon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) HideActions(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) SelectedClass(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Focusable(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Ripple(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Complete(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":complete", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) CompleteIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":complete-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Editable(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) EditIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":edit-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) ErrorIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":error-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) On(name string, value string) (r *VStepperVerticalItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalItemBuilder) Bind(name string, value string) (r *VStepperVerticalItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperVerticalItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperVerticalItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotIcon(child ...h.HTMLComponent) {
	b.SetSlot("icon", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotIcon(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("icon", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotIcon(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotIcon(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotIcon(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotIcon(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotTitle(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotText(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotPrev(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotNext(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}

func (b *VStepperVerticalItemBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VStepperVerticalItemBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VStepperVerticalItemBuilder) SlotActions(child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VStepperVerticalItemBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VStepperVerticalItemBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}
