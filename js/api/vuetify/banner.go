package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBannerBuilder struct {
	VTagBuilder[*VBannerBuilder]
}

func VBanner(children ...h.HTMLComponent) *VBannerBuilder {
	return VTag(&VBannerBuilder{}, "v-banner", children...)
}

func (b *VBannerBuilder) Text(v string) (r *VBannerBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VBannerBuilder) Border(v interface{}) (r *VBannerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Icon(v interface{}) (r *VBannerBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Density(v interface{}) (r *VBannerBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Height(v interface{}) (r *VBannerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MaxHeight(v interface{}) (r *VBannerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MaxWidth(v interface{}) (r *VBannerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MinHeight(v interface{}) (r *VBannerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MinWidth(v interface{}) (r *VBannerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Width(v interface{}) (r *VBannerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Elevation(v interface{}) (r *VBannerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Location(v interface{}) (r *VBannerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Position(v interface{}) (r *VBannerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Sticky(v bool) (r *VBannerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Rounded(v interface{}) (r *VBannerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Tile(v bool) (r *VBannerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Tag(v interface{}) (r *VBannerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Theme(v string) (r *VBannerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBannerBuilder) Color(v string) (r *VBannerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBannerBuilder) Stacked(v bool) (r *VBannerBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Avatar(v string) (r *VBannerBuilder) {
	b.Attr("avatar", v)
	return b
}

func (b *VBannerBuilder) BgColor(v string) (r *VBannerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VBannerBuilder) Mobile(v bool) (r *VBannerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) MobileBreakpoint(v interface{}) (r *VBannerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Lines(v interface{}) (r *VBannerBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) On(name string, value string) (r *VBannerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerBuilder) Bind(name string, value string) (r *VBannerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBannerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBannerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBannerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBannerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBannerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBannerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBannerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBannerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VBannerBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VBannerBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VBannerBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VBannerBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VBannerBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VBannerBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VBannerBuilder) SlotText(child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VBannerBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}

func (b *VBannerBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VBannerBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VBannerBuilder) SlotActions(child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VBannerBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VBannerBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}
