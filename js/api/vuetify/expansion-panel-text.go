package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpansionPanelTextBuilder struct {
	VTagBuilder[*VExpansionPanelTextBuilder]
}

func VExpansionPanelText(children ...h.HTMLComponent) *VExpansionPanelTextBuilder {
	return VTag(&VExpansionPanelTextBuilder{}, "v-expansion-panel-text", children...)
}

func (b *VExpansionPanelTextBuilder) Eager(v bool) (r *VExpansionPanelTextBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTextBuilder) On(name string, value string) (r *VExpansionPanelTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpansionPanelTextBuilder) Bind(name string, value string) (r *VExpansionPanelTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VExpansionPanelTextBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VExpansionPanelTextBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VExpansionPanelTextBuilder) Slot(name string, child ...h.HTMLComponent) (r *VExpansionPanelTextBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VExpansionPanelTextBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VExpansionPanelTextBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VExpansionPanelTextBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VExpansionPanelTextBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VExpansionPanelTextBuilder) SlotDefault(child ...h.HTMLComponent) (r *VExpansionPanelTextBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VExpansionPanelTextBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VExpansionPanelTextBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
