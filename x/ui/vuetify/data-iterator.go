package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
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

func (b *VDataIteratorBuilder) Loading(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) Items(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ItemValue(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ItemSelectable(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ReturnObject(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) ShowSelect(v bool) (r *VDataIteratorBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataIteratorBuilder) SelectStrategy(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) Page(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ModelValue(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataIteratorBuilder) ValueComparator(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
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

func (b *VDataIteratorBuilder) Tag(v string) (r *VDataIteratorBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDataIteratorBuilder) Transition(v interface{}) (r *VDataIteratorBuilder) {
	b.Attr(":transition", h.JSONString(v))
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
