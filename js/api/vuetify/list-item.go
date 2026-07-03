package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemBuilder struct {
	VTagBuilder[*VListItemBuilder]
}

func VListItem(children ...h.HTMLComponent) *VListItemBuilder {
	return VTag(&VListItemBuilder{}, "v-list-item", children...)
}

func (b *VListItemBuilder) Replace(v bool) (r *VListItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Link(v bool) (r *VListItemBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Tag(v interface{}) (r *VListItemBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Nav(v bool) (r *VListItemBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Title(v interface{}) (r *VListItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Theme(v string) (r *VListItemBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VListItemBuilder) Value(v interface{}) (r *VListItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Exact(v bool) (r *VListItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Subtitle(v interface{}) (r *VListItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) BaseColor(v string) (r *VListItemBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VListItemBuilder) ActiveColor(v string) (r *VListItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VListItemBuilder) ActiveClass(v string) (r *VListItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VListItemBuilder) Disabled(v bool) (r *VListItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Lines(v interface{}) (r *VListItemBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Slim(v bool) (r *VListItemBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Border(v interface{}) (r *VListItemBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Density(v interface{}) (r *VListItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Height(v interface{}) (r *VListItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MaxHeight(v interface{}) (r *VListItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MaxWidth(v interface{}) (r *VListItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MinHeight(v interface{}) (r *VListItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MinWidth(v interface{}) (r *VListItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Width(v interface{}) (r *VListItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Elevation(v interface{}) (r *VListItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Rounded(v interface{}) (r *VListItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Tile(v bool) (r *VListItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Color(v string) (r *VListItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListItemBuilder) Variant(v interface{}) (r *VListItemBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) AppendIcon(v interface{}) (r *VListItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) PrependIcon(v interface{}) (r *VListItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Active(v bool) (r *VListItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Href(v string) (r *VListItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VListItemBuilder) To(v interface{}) (r *VListItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Ripple(v interface{}) (r *VListItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) AppendAvatar(v string) (r *VListItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VListItemBuilder) PrependAvatar(v string) (r *VListItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VListItemBuilder) On(name string, value string) (r *VListItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemBuilder) Bind(name string, value string) (r *VListItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListItemBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VListItemBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VListItemBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VListItemBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VListItemBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VListItemBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VListItemBuilder) SlotAppend(child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VListItemBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VListItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VListItemBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VListItemBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VListItemBuilder) SlotTitle(child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VListItemBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VListItemBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VListItemBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VListItemBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VListItemBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VListItemBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}
