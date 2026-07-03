package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VInfiniteScrollBuilder struct {
	VTagBuilder[*VInfiniteScrollBuilder]
}

func VInfiniteScroll(children ...h.HTMLComponent) *VInfiniteScrollBuilder {
	return VTag(&VInfiniteScrollBuilder{}, "v-infinite-scroll", children...)
}

func (b *VInfiniteScrollBuilder) Tag(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Mode(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Height(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MaxHeight(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MaxWidth(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MinHeight(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MinWidth(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Width(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Color(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VInfiniteScrollBuilder) Direction(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Margin(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":margin", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Side(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":side", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) LoadMoreText(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("load-more-text", v)
	return b
}

func (b *VInfiniteScrollBuilder) EmptyText(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("empty-text", v)
	return b
}

func (b *VInfiniteScrollBuilder) On(name string, value string) (r *VInfiniteScrollBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VInfiniteScrollBuilder) Bind(name string, value string) (r *VInfiniteScrollBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VInfiniteScrollBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VInfiniteScrollBuilder) Slot(name string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VInfiniteScrollBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VInfiniteScrollBuilder) SlotDefault(child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlotLoading(child ...h.HTMLComponent) {
	b.SetSlot("loading", child...)
}

func (b *VInfiniteScrollBuilder) SetScopedSlotLoading(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loading", scope, child...)
}

func (b *VInfiniteScrollBuilder) SlotLoading(child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlotLoading(child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlotLoading(scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlotLoading(scope, child...)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlotError(child ...h.HTMLComponent) {
	b.SetSlot("error", child...)
}

func (b *VInfiniteScrollBuilder) SetScopedSlotError(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("error", scope, child...)
}

func (b *VInfiniteScrollBuilder) SlotError(child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlotError(child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlotError(scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlotError(scope, child...)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlotEmpty(child ...h.HTMLComponent) {
	b.SetSlot("empty", child...)
}

func (b *VInfiniteScrollBuilder) SetScopedSlotEmpty(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("empty", scope, child...)
}

func (b *VInfiniteScrollBuilder) SlotEmpty(child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlotEmpty(child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlotEmpty(scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlotEmpty(scope, child...)
	return b
}

func (b *VInfiniteScrollBuilder) SetSlotLoadMore(child ...h.HTMLComponent) {
	b.SetSlot("load-more", child...)
}

func (b *VInfiniteScrollBuilder) SetScopedSlotLoadMore(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("load-more", scope, child...)
}

func (b *VInfiniteScrollBuilder) SlotLoadMore(child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetSlotLoadMore(child...)
	return b
}

func (b *VInfiniteScrollBuilder) ScopedSlotLoadMore(scope string, child ...h.HTMLComponent) (r *VInfiniteScrollBuilder) {
	b.SetScopedSlotLoadMore(scope, child...)
	return b
}
