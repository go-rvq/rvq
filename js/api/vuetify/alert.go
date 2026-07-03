package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAlertBuilder struct {
	VTagBuilder[*VAlertBuilder]
}

func VAlert(children ...h.HTMLComponent) *VAlertBuilder {
	return VTag(&VAlertBuilder{}, "v-alert", children...)
}

func (b *VAlertBuilder) Title(v string) (r *VAlertBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VAlertBuilder) Text(v string) (r *VAlertBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAlertBuilder) Border(v interface{}) (r *VAlertBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) BorderColor(v string) (r *VAlertBuilder) {
	b.Attr("border-color", v)
	return b
}

func (b *VAlertBuilder) Closable(v bool) (r *VAlertBuilder) {
	b.Attr(":closable", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) CloseIcon(v interface{}) (r *VAlertBuilder) {
	b.Attr(":close-icon", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Type(v interface{}) (r *VAlertBuilder) {
	b.Attr(":type", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) CloseLabel(v string) (r *VAlertBuilder) {
	b.Attr("close-label", v)
	return b
}

func (b *VAlertBuilder) Icon(v interface{}) (r *VAlertBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) ModelValue(v bool) (r *VAlertBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Prominent(v bool) (r *VAlertBuilder) {
	b.Attr(":prominent", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Density(v interface{}) (r *VAlertBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Height(v interface{}) (r *VAlertBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MaxHeight(v interface{}) (r *VAlertBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MaxWidth(v interface{}) (r *VAlertBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MinHeight(v interface{}) (r *VAlertBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MinWidth(v interface{}) (r *VAlertBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Width(v interface{}) (r *VAlertBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Elevation(v interface{}) (r *VAlertBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Location(v interface{}) (r *VAlertBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Position(v interface{}) (r *VAlertBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Rounded(v interface{}) (r *VAlertBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Tile(v bool) (r *VAlertBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Tag(v interface{}) (r *VAlertBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Theme(v string) (r *VAlertBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAlertBuilder) Color(v string) (r *VAlertBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAlertBuilder) Variant(v interface{}) (r *VAlertBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) On(name string, value string) (r *VAlertBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAlertBuilder) Bind(name string, value string) (r *VAlertBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAlertBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAlertBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAlertBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAlertBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAlertBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAlertBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VAlertBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VAlertBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VAlertBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VAlertBuilder) SlotTitle(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VAlertBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VAlertBuilder) SlotText(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VAlertBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VAlertBuilder) SlotAppend(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VAlertBuilder) SetSlotClose(child ...h.HTMLComponent) {
	b.SetSlot("close", child...)
}

func (b *VAlertBuilder) SetScopedSlotClose(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("close", scope, child...)
}

func (b *VAlertBuilder) SlotClose(child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetSlotClose(child...)
	return b
}

func (b *VAlertBuilder) ScopedSlotClose(scope string, child ...h.HTMLComponent) (r *VAlertBuilder) {
	b.SetScopedSlotClose(scope, child...)
	return b
}
