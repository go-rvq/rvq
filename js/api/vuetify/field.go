package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFieldBuilder struct {
	VTagBuilder[*VFieldBuilder]
}

func VField(children ...h.HTMLComponent) *VFieldBuilder {
	return VTag(&VFieldBuilder{}, "v-field", children...)
}

func (b *VFieldBuilder) Flat(v bool) (r *VFieldBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Reverse(v bool) (r *VFieldBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Error(v bool) (r *VFieldBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Label(v string) (r *VFieldBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VFieldBuilder) Theme(v string) (r *VFieldBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFieldBuilder) ID(v string) (r *VFieldBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VFieldBuilder) BaseColor(v string) (r *VFieldBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFieldBuilder) BgColor(v string) (r *VFieldBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VFieldBuilder) Disabled(v bool) (r *VFieldBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
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

func (b *VFieldBuilder) Color(v string) (r *VFieldBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFieldBuilder) Variant(v interface{}) (r *VFieldBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) ModelValue(v interface{}) (r *VFieldBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) CenterAffix(v bool) (r *VFieldBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) Glow(v bool) (r *VFieldBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) IconColor(v interface{}) (r *VFieldBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VFieldBuilder) Focused(v bool) (r *VFieldBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VFieldBuilder) AppendInnerIcon(v interface{}) (r *VFieldBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
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

func (b *VFieldBuilder) Dirty(v bool) (r *VFieldBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
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

func (b *VFieldBuilder) Loading(v interface{}) (r *VFieldBuilder) {
	b.Attr(":loading", h.JSONString(v))
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

func (b *VFieldBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFieldBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFieldBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFieldBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VFieldBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VFieldBuilder) SlotClear(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VFieldBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VFieldBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VFieldBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VFieldBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VFieldBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VFieldBuilder) SlotLabel(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VFieldBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VFieldBuilder) SlotLoader(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VFieldBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFieldBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFieldBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFieldBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFieldBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
