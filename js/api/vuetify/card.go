package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardBuilder struct {
	VTagBuilder[*VCardBuilder]
}

func VCard(children ...h.HTMLComponent) *VCardBuilder {
	return VTag(&VCardBuilder{}, "v-card", children...)
}

func (b *VCardBuilder) Title(v interface{}) (r *VCardBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Text(v interface{}) (r *VCardBuilder) {
	b.Attr(":text", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Flat(v bool) (r *VCardBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Replace(v bool) (r *VCardBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Link(v bool) (r *VCardBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Border(v interface{}) (r *VCardBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Density(v interface{}) (r *VCardBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Height(v interface{}) (r *VCardBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MaxHeight(v interface{}) (r *VCardBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MaxWidth(v interface{}) (r *VCardBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MinHeight(v interface{}) (r *VCardBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MinWidth(v interface{}) (r *VCardBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Width(v interface{}) (r *VCardBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Elevation(v interface{}) (r *VCardBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Location(v interface{}) (r *VCardBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Position(v interface{}) (r *VCardBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Rounded(v interface{}) (r *VCardBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Tile(v bool) (r *VCardBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Tag(v interface{}) (r *VCardBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Theme(v string) (r *VCardBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCardBuilder) Color(v string) (r *VCardBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCardBuilder) Variant(v interface{}) (r *VCardBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Image(v string) (r *VCardBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VCardBuilder) PrependIcon(v interface{}) (r *VCardBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VCardBuilder) AppendIcon(v interface{}) (r *VCardBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Ripple(v interface{}) (r *VCardBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Disabled(v bool) (r *VCardBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Loading(v interface{}) (r *VCardBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Href(v string) (r *VCardBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VCardBuilder) Exact(v bool) (r *VCardBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) To(v interface{}) (r *VCardBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Subtitle(v interface{}) (r *VCardBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VCardBuilder) AppendAvatar(v string) (r *VCardBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VCardBuilder) Hover(v bool) (r *VCardBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) PrependAvatar(v string) (r *VCardBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VCardBuilder) On(name string, value string) (r *VCardBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardBuilder) Bind(name string, value string) (r *VCardBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VCardBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VCardBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VCardBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VCardBuilder) SlotAppend(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VCardBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VCardBuilder) SlotTitle(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VCardBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VCardBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VCardBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VCardBuilder) SlotActions(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VCardBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VCardBuilder) SlotText(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VCardBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VCardBuilder) SlotLoader(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotImage(child ...h.HTMLComponent) {
	b.SetSlot("image", child...)
}

func (b *VCardBuilder) SetScopedSlotImage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("image", scope, child...)
}

func (b *VCardBuilder) SlotImage(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotImage(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotImage(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotImage(scope, child...)
	return b
}

func (b *VCardBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VCardBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VCardBuilder) SlotItem(child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VCardBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VCardBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}
