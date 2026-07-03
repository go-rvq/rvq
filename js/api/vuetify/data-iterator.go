package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataIteratorBuilder struct {
	VTagBuilder[*VDataIteratorBuilder]
}

func VDataIterator(children ...h.HTMLComponent) *VDataIteratorBuilder {
	return VTag(&VDataIteratorBuilder{}, "v-data-iterator", children...)
}

func (b *VDataIteratorBuilder) Search(v string) (r *VDataIteratorBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataIteratorBuilder) Tag(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) Items(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) FilterMode(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) NoFilter(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) CustomFilter(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) CustomKeyFilter(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) FilterKeys(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) SelectStrategy(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ItemValue(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ReturnObject(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) ValueComparator(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ModelValue(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) Transition(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) Loading(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) Page(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ItemSelectable(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ShowSelect(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) SortBy(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) MultiSort(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) MustSort(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) CustomKeySort(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ItemsPerPage(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":items-per-page", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ExpandOnClick(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) ShowExpand(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) Expanded(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) GroupBy(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) On(name string, value string) (r *VDataIteratorBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataIteratorBuilder) Bind(name string, value string) (r *VDataIteratorBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataIteratorBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataIteratorBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataIteratorBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataIteratorBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDataIteratorBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDataIteratorBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VDataIteratorBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VDataIteratorBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VDataIteratorBuilder) SlotHeader(child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}

func (b *VDataIteratorBuilder) SetSlotFooter(child ...h.HTMLComponent) {
	b.SetSlot("footer", child...)
}

func (b *VDataIteratorBuilder) SetScopedSlotFooter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("footer", scope, child...)
}

func (b *VDataIteratorBuilder) SlotFooter(child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlotFooter(child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlotFooter(scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlotFooter(scope, child...)
	return b
}

func (b *VDataIteratorBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDataIteratorBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDataIteratorBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDataIteratorBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VDataIteratorBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VDataIteratorBuilder) SlotNoData(child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VDataIteratorBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VDataIteratorBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}
