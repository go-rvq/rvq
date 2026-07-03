package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPullToRefreshBuilder struct {
	VTagBuilder[*VPullToRefreshBuilder]
}

func VPullToRefresh(children ...h.HTMLComponent) *VPullToRefreshBuilder {
	return VTag(&VPullToRefreshBuilder{}, "v-pull-to-refresh", children...)
}

func (b *VPullToRefreshBuilder) Disabled(v bool) (r *VPullToRefreshBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VPullToRefreshBuilder) PullDownThreshold(v interface{}) (r *VPullToRefreshBuilder) {
	b.Attr(":pull-down-threshold", h.JSONString(v))
	return b
}

func (b *VPullToRefreshBuilder) On(name string, value string) (r *VPullToRefreshBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPullToRefreshBuilder) Bind(name string, value string) (r *VPullToRefreshBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VPullToRefreshBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VPullToRefreshBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VPullToRefreshBuilder) Slot(name string, child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VPullToRefreshBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VPullToRefreshBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VPullToRefreshBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VPullToRefreshBuilder) SlotDefault(child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VPullToRefreshBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VPullToRefreshBuilder) SetSlotPullDownPanel(child ...h.HTMLComponent) {
	b.SetSlot("pullDownPanel", child...)
}

func (b *VPullToRefreshBuilder) SetScopedSlotPullDownPanel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("pullDownPanel", scope, child...)
}

func (b *VPullToRefreshBuilder) SlotPullDownPanel(child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetSlotPullDownPanel(child...)
	return b
}

func (b *VPullToRefreshBuilder) ScopedSlotPullDownPanel(scope string, child ...h.HTMLComponent) (r *VPullToRefreshBuilder) {
	b.SetScopedSlotPullDownPanel(scope, child...)
	return b
}
