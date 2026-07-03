package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VNavigationDrawerBuilder struct {
	VTagBuilder[*VNavigationDrawerBuilder]
}

func VNavigationDrawer(children ...h.HTMLComponent) *VNavigationDrawerBuilder {
	return VTag(&VNavigationDrawerBuilder{}, "v-navigation-drawer", children...)
}

func (b *VNavigationDrawerBuilder) Tag(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Name(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VNavigationDrawerBuilder) Theme(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VNavigationDrawerBuilder) Border(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Width(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Elevation(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Rounded(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Tile(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Color(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VNavigationDrawerBuilder) Absolute(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) ModelValue(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Persistent(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Scrim(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) CloseDelay(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) OpenDelay(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Location(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Mobile(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Image(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VNavigationDrawerBuilder) Order(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) MobileBreakpoint(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Sticky(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Floating(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) DisableResizeWatcher(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":disable-resize-watcher", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) DisableRouteWatcher(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":disable-route-watcher", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) ExpandOnHover(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":expand-on-hover", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Permanent(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":permanent", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Rail(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":rail", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) RailWidth(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":rail-width", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Temporary(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":temporary", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Touchless(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":touchless", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) On(name string, value string) (r *VNavigationDrawerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VNavigationDrawerBuilder) Bind(name string, value string) (r *VNavigationDrawerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VNavigationDrawerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VNavigationDrawerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VNavigationDrawerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VNavigationDrawerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VNavigationDrawerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VNavigationDrawerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VNavigationDrawerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VNavigationDrawerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VNavigationDrawerBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VNavigationDrawerBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VNavigationDrawerBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VNavigationDrawerBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VNavigationDrawerBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VNavigationDrawerBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VNavigationDrawerBuilder) SlotAppend(child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VNavigationDrawerBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VNavigationDrawerBuilder) SetSlotImage(child ...h.HTMLComponent) {
	b.SetSlot("image", child...)
}

func (b *VNavigationDrawerBuilder) SetScopedSlotImage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("image", scope, child...)
}

func (b *VNavigationDrawerBuilder) SlotImage(child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetSlotImage(child...)
	return b
}

func (b *VNavigationDrawerBuilder) ScopedSlotImage(scope string, child ...h.HTMLComponent) (r *VNavigationDrawerBuilder) {
	b.SetScopedSlotImage(scope, child...)
	return b
}
