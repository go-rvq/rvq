package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBtnBuilder struct {
	VTagBuilder[*VBtnBuilder]
}

func VBtn(children ...h.HTMLComponent) *VBtnBuilder {
	return VTag(&VBtnBuilder{}, "v-btn", children...)
}

func (b *VBtnBuilder) Symbol(v interface{}) (r *VBtnBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Text(v interface{}) (r *VBtnBuilder) {
	b.Attr(":text", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Flat(v bool) (r *VBtnBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Replace(v bool) (r *VBtnBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Border(v interface{}) (r *VBtnBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Icon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Density(v interface{}) (r *VBtnBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Height(v interface{}) (r *VBtnBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MaxHeight(v interface{}) (r *VBtnBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MaxWidth(v interface{}) (r *VBtnBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MinHeight(v interface{}) (r *VBtnBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MinWidth(v interface{}) (r *VBtnBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Width(v interface{}) (r *VBtnBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Elevation(v interface{}) (r *VBtnBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Location(v interface{}) (r *VBtnBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Position(v interface{}) (r *VBtnBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Rounded(v interface{}) (r *VBtnBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Tile(v bool) (r *VBtnBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Tag(v interface{}) (r *VBtnBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Theme(v string) (r *VBtnBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnBuilder) Color(v string) (r *VBtnBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnBuilder) Variant(v interface{}) (r *VBtnBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Active(v bool) (r *VBtnBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) ActiveColor(v string) (r *VBtnBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VBtnBuilder) BaseColor(v string) (r *VBtnBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnBuilder) PrependIcon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) AppendIcon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Block(v bool) (r *VBtnBuilder) {
	b.Attr(":block", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Readonly(v bool) (r *VBtnBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Slim(v bool) (r *VBtnBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Stacked(v bool) (r *VBtnBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Ripple(v interface{}) (r *VBtnBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Value(v interface{}) (r *VBtnBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Disabled(v bool) (r *VBtnBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) SelectedClass(v string) (r *VBtnBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VBtnBuilder) Loading(v interface{}) (r *VBtnBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Href(v string) (r *VBtnBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VBtnBuilder) Exact(v bool) (r *VBtnBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) To(v interface{}) (r *VBtnBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Size(v interface{}) (r *VBtnBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) On(name string, value string) (r *VBtnBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnBuilder) Bind(name string, value string) (r *VBtnBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBtnBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBtnBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBtnBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBtnBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBtnBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBtnBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBtnBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBtnBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VBtnBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VBtnBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VBtnBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VBtnBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VBtnBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VBtnBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VBtnBuilder) SlotAppend(child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VBtnBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VBtnBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VBtnBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VBtnBuilder) SlotLoader(child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VBtnBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VBtnBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}
