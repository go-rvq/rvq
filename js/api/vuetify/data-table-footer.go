package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableFooterBuilder struct {
	VTagBuilder[*VDataTableFooterBuilder]
}

func VDataTableFooter(children ...h.HTMLComponent) *VDataTableFooterBuilder {
	return VTag(&VDataTableFooterBuilder{}, "v-data-table-footer", children...)
}

func (b *VDataTableFooterBuilder) PrevIcon(v interface{}) (r *VDataTableFooterBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VDataTableFooterBuilder) NextIcon(v interface{}) (r *VDataTableFooterBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VDataTableFooterBuilder) FirstIcon(v interface{}) (r *VDataTableFooterBuilder) {
	b.Attr(":first-icon", h.JSONString(v))
	return b
}

func (b *VDataTableFooterBuilder) LastIcon(v interface{}) (r *VDataTableFooterBuilder) {
	b.Attr(":last-icon", h.JSONString(v))
	return b
}

func (b *VDataTableFooterBuilder) ItemsPerPageText(v string) (r *VDataTableFooterBuilder) {
	b.Attr("items-per-page-text", v)
	return b
}

func (b *VDataTableFooterBuilder) PageText(v string) (r *VDataTableFooterBuilder) {
	b.Attr("page-text", v)
	return b
}

func (b *VDataTableFooterBuilder) FirstPageLabel(v string) (r *VDataTableFooterBuilder) {
	b.Attr("first-page-label", v)
	return b
}

func (b *VDataTableFooterBuilder) PrevPageLabel(v string) (r *VDataTableFooterBuilder) {
	b.Attr("prev-page-label", v)
	return b
}

func (b *VDataTableFooterBuilder) NextPageLabel(v string) (r *VDataTableFooterBuilder) {
	b.Attr("next-page-label", v)
	return b
}

func (b *VDataTableFooterBuilder) LastPageLabel(v string) (r *VDataTableFooterBuilder) {
	b.Attr("last-page-label", v)
	return b
}

func (b *VDataTableFooterBuilder) ItemsPerPageOptions(v interface{}) (r *VDataTableFooterBuilder) {
	b.Attr(":items-per-page-options", h.JSONString(v))
	return b
}

func (b *VDataTableFooterBuilder) ShowCurrentPage(v bool) (r *VDataTableFooterBuilder) {
	b.Attr(":show-current-page", fmt.Sprint(v))
	return b
}

func (b *VDataTableFooterBuilder) On(name string, value string) (r *VDataTableFooterBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableFooterBuilder) Bind(name string, value string) (r *VDataTableFooterBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableFooterBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableFooterBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableFooterBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableFooterBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableFooterBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableFooterBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableFooterBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VDataTableFooterBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VDataTableFooterBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VDataTableFooterBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VDataTableFooterBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableFooterBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}
