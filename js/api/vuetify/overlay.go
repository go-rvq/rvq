package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VOverlayBuilder struct {
	VTagBuilder[*VOverlayBuilder]
}

func VOverlay(children ...h.HTMLComponent) *VOverlayBuilder {
	return VTag(&VOverlayBuilder{}, "v-overlay", children...)
}

func (b *VOverlayBuilder) ModelValue(v bool) (r *VOverlayBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Height(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) MaxHeight(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) MaxWidth(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) MinHeight(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) MinWidth(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Width(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Location(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Absolute(v bool) (r *VOverlayBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Theme(v string) (r *VOverlayBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VOverlayBuilder) Disabled(v bool) (r *VOverlayBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Transition(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Activator(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) CloseOnBack(v bool) (r *VOverlayBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Contained(v bool) (r *VOverlayBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) ContentClass(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) ContentProps(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Opacity(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) NoClickAnimation(v bool) (r *VOverlayBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Persistent(v bool) (r *VOverlayBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) Scrim(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) ZIndex(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Target(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) ActivatorProps(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) OpenOnClick(v bool) (r *VOverlayBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) OpenOnHover(v bool) (r *VOverlayBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) OpenOnFocus(v bool) (r *VOverlayBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) CloseOnContentClick(v bool) (r *VOverlayBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) CloseDelay(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) OpenDelay(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Eager(v bool) (r *VOverlayBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VOverlayBuilder) LocationStrategy(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Origin(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Offset(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) ScrollStrategy(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) Attach(v interface{}) (r *VOverlayBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VOverlayBuilder) On(name string, value string) (r *VOverlayBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VOverlayBuilder) Bind(name string, value string) (r *VOverlayBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VOverlayBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VOverlayBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VOverlayBuilder) Slot(name string, child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VOverlayBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VOverlayBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VOverlayBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VOverlayBuilder) SlotDefault(child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VOverlayBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VOverlayBuilder) SetSlotActivator(child ...h.HTMLComponent) {
	b.SetSlot("activator", child...)
}

func (b *VOverlayBuilder) SetScopedSlotActivator(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("activator", scope, child...)
}

func (b *VOverlayBuilder) SlotActivator(child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetSlotActivator(child...)
	return b
}

func (b *VOverlayBuilder) ScopedSlotActivator(scope string, child ...h.HTMLComponent) (r *VOverlayBuilder) {
	b.SetScopedSlotActivator(scope, child...)
	return b
}
