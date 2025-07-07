package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VNavigationDrawerBuilder struct {
	VTagBuilder[*VNavigationDrawerBuilder]
}

func VNavigationDrawer(children ...h.HTMLComponent) *VNavigationDrawerBuilder {
	return VTag(&VNavigationDrawerBuilder{}, "v-navigation-drawer", children...)
}

func (b *VNavigationDrawerBuilder) Image(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VNavigationDrawerBuilder) Color(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("color", v)
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

func (b *VNavigationDrawerBuilder) Floating(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) ModelValue(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Permanent(v bool) (r *VNavigationDrawerBuilder) {
	return b.Attr("permanent", v)
}

func (b *VNavigationDrawerBuilder) Rail(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":rail", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) RailWidth(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":rail-width", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Scrim(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Temporary(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":temporary", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Persistent(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Touchless(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":touchless", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Width(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Location(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Sticky(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) Border(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":border", h.JSONString(v))
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

func (b *VNavigationDrawerBuilder) Mobile(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VNavigationDrawerBuilder) MobileBreakpoint(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Elevation(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Name(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VNavigationDrawerBuilder) Order(v interface{}) (r *VNavigationDrawerBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VNavigationDrawerBuilder) Absolute(v bool) (r *VNavigationDrawerBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
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

func (b *VNavigationDrawerBuilder) Tag(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VNavigationDrawerBuilder) Theme(v string) (r *VNavigationDrawerBuilder) {
	b.Attr("theme", v)
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
