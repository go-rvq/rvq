package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBarNavIconBuilder struct {
	VTagBuilder[*VAppBarNavIconBuilder]
}

func VAppBarNavIcon(children ...h.HTMLComponent) *VAppBarNavIconBuilder {
	return VTag(&VAppBarNavIconBuilder{}, "v-app-bar-nav-icon", children...)
}

func (b *VAppBarNavIconBuilder) Symbol(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Text(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":text", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Flat(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Replace(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Border(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Icon(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Density(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Height(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) MaxHeight(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) MaxWidth(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) MinHeight(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) MinWidth(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Width(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Elevation(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Location(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Position(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Rounded(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Tile(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Tag(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Theme(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAppBarNavIconBuilder) Color(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAppBarNavIconBuilder) Variant(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Active(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) ActiveColor(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VAppBarNavIconBuilder) BaseColor(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VAppBarNavIconBuilder) PrependIcon(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) AppendIcon(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Block(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":block", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Readonly(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Slim(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Stacked(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) Ripple(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Value(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Disabled(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) SelectedClass(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VAppBarNavIconBuilder) Loading(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Href(v string) (r *VAppBarNavIconBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VAppBarNavIconBuilder) Exact(v bool) (r *VAppBarNavIconBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VAppBarNavIconBuilder) To(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) Size(v interface{}) (r *VAppBarNavIconBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VAppBarNavIconBuilder) On(name string, value string) (r *VAppBarNavIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBarNavIconBuilder) Bind(name string, value string) (r *VAppBarNavIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAppBarNavIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAppBarNavIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAppBarNavIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAppBarNavIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAppBarNavIconBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAppBarNavIconBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAppBarNavIconBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAppBarNavIconBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VAppBarNavIconBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VAppBarNavIconBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VAppBarNavIconBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VAppBarNavIconBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VAppBarNavIconBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VAppBarNavIconBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VAppBarNavIconBuilder) SlotAppend(child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VAppBarNavIconBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VAppBarNavIconBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VAppBarNavIconBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VAppBarNavIconBuilder) SlotLoader(child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VAppBarNavIconBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VAppBarNavIconBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}
