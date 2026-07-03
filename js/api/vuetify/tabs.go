package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTabsBuilder struct {
	VTagBuilder[*VTabsBuilder]
}

func VTabs(children ...h.HTMLComponent) *VTabsBuilder {
	return VTag(&VTabsBuilder{}, "v-tabs", children...)
}

func (b *VTabsBuilder) Symbol(v interface{}) (r *VTabsBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Tag(v interface{}) (r *VTabsBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Items(v interface{}) (r *VTabsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) BgColor(v string) (r *VTabsBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTabsBuilder) Disabled(v bool) (r *VTabsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Multiple(v bool) (r *VTabsBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Mandatory(v interface{}) (r *VTabsBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Density(v interface{}) (r *VTabsBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Height(v interface{}) (r *VTabsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Color(v string) (r *VTabsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTabsBuilder) ModelValue(v interface{}) (r *VTabsBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Direction(v interface{}) (r *VTabsBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Max(v interface{}) (r *VTabsBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Mobile(v bool) (r *VTabsBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) MobileBreakpoint(v interface{}) (r *VTabsBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) PrevIcon(v interface{}) (r *VTabsBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) NextIcon(v interface{}) (r *VTabsBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) SelectedClass(v string) (r *VTabsBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsBuilder) Stacked(v bool) (r *VTabsBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) CenterActive(v bool) (r *VTabsBuilder) {
	b.Attr(":center-active", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) ShowArrows(v interface{}) (r *VTabsBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) AlignTabs(v interface{}) (r *VTabsBuilder) {
	b.Attr(":align-tabs", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) FixedTabs(v bool) (r *VTabsBuilder) {
	b.Attr(":fixed-tabs", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Grow(v bool) (r *VTabsBuilder) {
	b.Attr(":grow", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) HideSlider(v bool) (r *VTabsBuilder) {
	b.Attr(":hide-slider", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) SliderColor(v string) (r *VTabsBuilder) {
	b.Attr("slider-color", v)
	return b
}

func (b *VTabsBuilder) On(name string, value string) (r *VTabsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsBuilder) Bind(name string, value string) (r *VTabsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTabsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTabsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTabsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTabsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTabsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTabsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTabsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTabsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VTabsBuilder) SetSlotTab(child ...h.HTMLComponent) {
	b.SetSlot("tab", child...)
}

func (b *VTabsBuilder) SetScopedSlotTab(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tab", scope, child...)
}

func (b *VTabsBuilder) SlotTab(child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetSlotTab(child...)
	return b
}

func (b *VTabsBuilder) ScopedSlotTab(scope string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetScopedSlotTab(scope, child...)
	return b
}

func (b *VTabsBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VTabsBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VTabsBuilder) SlotItem(child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VTabsBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VTabsBuilder) SetSlotWindow(child ...h.HTMLComponent) {
	b.SetSlot("window", child...)
}

func (b *VTabsBuilder) SetScopedSlotWindow(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("window", scope, child...)
}

func (b *VTabsBuilder) SlotWindow(child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetSlotWindow(child...)
	return b
}

func (b *VTabsBuilder) ScopedSlotWindow(scope string, child ...h.HTMLComponent) (r *VTabsBuilder) {
	b.SetScopedSlotWindow(scope, child...)
	return b
}
