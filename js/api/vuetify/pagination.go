package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPaginationBuilder struct {
	VTagBuilder[*VPaginationBuilder]
}

func VPagination(children ...h.HTMLComponent) *VPaginationBuilder {
	return VTag(&VPaginationBuilder{}, "v-pagination", children...)
}

func (b *VPaginationBuilder) Length(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Border(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Start(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":start", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) ModelValue(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Density(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Elevation(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Rounded(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Tile(v bool) (r *VPaginationBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VPaginationBuilder) Tag(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) Theme(v string) (r *VPaginationBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VPaginationBuilder) Color(v string) (r *VPaginationBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VPaginationBuilder) Variant(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) ActiveColor(v string) (r *VPaginationBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VPaginationBuilder) Disabled(v bool) (r *VPaginationBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VPaginationBuilder) Size(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) AriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("aria-label", v)
	return b
}

func (b *VPaginationBuilder) NextIcon(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) PrevIcon(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) TotalVisible(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":total-visible", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) FirstIcon(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":first-icon", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) LastIcon(v interface{}) (r *VPaginationBuilder) {
	b.Attr(":last-icon", h.JSONString(v))
	return b
}

func (b *VPaginationBuilder) PageAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("page-aria-label", v)
	return b
}

func (b *VPaginationBuilder) CurrentPageAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("current-page-aria-label", v)
	return b
}

func (b *VPaginationBuilder) FirstAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("first-aria-label", v)
	return b
}

func (b *VPaginationBuilder) PreviousAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("previous-aria-label", v)
	return b
}

func (b *VPaginationBuilder) NextAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("next-aria-label", v)
	return b
}

func (b *VPaginationBuilder) LastAriaLabel(v string) (r *VPaginationBuilder) {
	b.Attr("last-aria-label", v)
	return b
}

func (b *VPaginationBuilder) Ellipsis(v string) (r *VPaginationBuilder) {
	b.Attr("ellipsis", v)
	return b
}

func (b *VPaginationBuilder) ShowFirstLastPage(v bool) (r *VPaginationBuilder) {
	b.Attr(":show-first-last-page", fmt.Sprint(v))
	return b
}

func (b *VPaginationBuilder) On(name string, value string) (r *VPaginationBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPaginationBuilder) Bind(name string, value string) (r *VPaginationBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VPaginationBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VPaginationBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VPaginationBuilder) Slot(name string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VPaginationBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VPaginationBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VPaginationBuilder) SlotItem(child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VPaginationBuilder) SetSlotFirst(child ...h.HTMLComponent) {
	b.SetSlot("first", child...)
}

func (b *VPaginationBuilder) SetScopedSlotFirst(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("first", scope, child...)
}

func (b *VPaginationBuilder) SlotFirst(child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlotFirst(child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlotFirst(scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlotFirst(scope, child...)
	return b
}

func (b *VPaginationBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VPaginationBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VPaginationBuilder) SlotPrev(child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VPaginationBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VPaginationBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VPaginationBuilder) SlotNext(child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}

func (b *VPaginationBuilder) SetSlotLast(child ...h.HTMLComponent) {
	b.SetSlot("last", child...)
}

func (b *VPaginationBuilder) SetScopedSlotLast(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("last", scope, child...)
}

func (b *VPaginationBuilder) SlotLast(child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetSlotLast(child...)
	return b
}

func (b *VPaginationBuilder) ScopedSlotLast(scope string, child ...h.HTMLComponent) (r *VPaginationBuilder) {
	b.SetScopedSlotLast(scope, child...)
	return b
}
