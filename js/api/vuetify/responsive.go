package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VResponsiveBuilder struct {
	VTagBuilder[*VResponsiveBuilder]
}

func VResponsive(children ...h.HTMLComponent) *VResponsiveBuilder {
	return VTag(&VResponsiveBuilder{}, "v-responsive", children...)
}

func (b *VResponsiveBuilder) Height(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MaxHeight(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MaxWidth(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MinHeight(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MinWidth(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) Width(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) Inline(v bool) (r *VResponsiveBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VResponsiveBuilder) ContentClass(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) AspectRatio(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":aspect-ratio", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) On(name string, value string) (r *VResponsiveBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VResponsiveBuilder) Bind(name string, value string) (r *VResponsiveBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VResponsiveBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VResponsiveBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VResponsiveBuilder) Slot(name string, child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VResponsiveBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VResponsiveBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VResponsiveBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VResponsiveBuilder) SlotDefault(child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VResponsiveBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VResponsiveBuilder) SetSlotAdditional(child ...h.HTMLComponent) {
	b.SetSlot("additional", child...)
}

func (b *VResponsiveBuilder) SetScopedSlotAdditional(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("additional", scope, child...)
}

func (b *VResponsiveBuilder) SlotAdditional(child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetSlotAdditional(child...)
	return b
}

func (b *VResponsiveBuilder) ScopedSlotAdditional(scope string, child ...h.HTMLComponent) (r *VResponsiveBuilder) {
	b.SetScopedSlotAdditional(scope, child...)
	return b
}
