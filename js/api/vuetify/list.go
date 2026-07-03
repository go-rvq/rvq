package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListBuilder struct {
	VTagBuilder[*VListBuilder]
}

func VList(children ...h.HTMLComponent) *VListBuilder {
	return VTag(&VListBuilder{}, "v-list", children...)
}

func (b *VListBuilder) Tag(v interface{}) (r *VListBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListBuilder) Nav(v bool) (r *VListBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Activated(v interface{}) (r *VListBuilder) {
	b.Attr(":activated", h.JSONString(v))
	return b
}

func (b *VListBuilder) Theme(v string) (r *VListBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VListBuilder) Items(v interface{}) (r *VListBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VListBuilder) BaseColor(v string) (r *VListBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VListBuilder) ActiveColor(v string) (r *VListBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VListBuilder) ActiveClass(v string) (r *VListBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VListBuilder) BgColor(v string) (r *VListBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VListBuilder) Disabled(v bool) (r *VListBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) ExpandIcon(v interface{}) (r *VListBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VListBuilder) CollapseIcon(v interface{}) (r *VListBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VListBuilder) Lines(v interface{}) (r *VListBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VListBuilder) Slim(v bool) (r *VListBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Activatable(v bool) (r *VListBuilder) {
	b.Attr(":activatable", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Selectable(v bool) (r *VListBuilder) {
	b.Attr(":selectable", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Opened(v interface{}) (r *VListBuilder) {
	b.Attr(":opened", h.JSONString(v))
	return b
}

func (b *VListBuilder) Selected(v interface{}) (r *VListBuilder) {
	b.Attr(":selected", h.JSONString(v))
	return b
}

func (b *VListBuilder) Mandatory(v bool) (r *VListBuilder) {
	b.Attr(":mandatory", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) ActiveStrategy(v interface{}) (r *VListBuilder) {
	b.Attr(":active-strategy", h.JSONString(v))
	return b
}

func (b *VListBuilder) SelectStrategy(v interface{}) (r *VListBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VListBuilder) OpenStrategy(v interface{}) (r *VListBuilder) {
	b.Attr(":open-strategy", h.JSONString(v))
	return b
}

func (b *VListBuilder) Border(v interface{}) (r *VListBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VListBuilder) Density(v interface{}) (r *VListBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VListBuilder) Height(v interface{}) (r *VListBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VListBuilder) MaxHeight(v interface{}) (r *VListBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VListBuilder) MaxWidth(v interface{}) (r *VListBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VListBuilder) MinHeight(v interface{}) (r *VListBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VListBuilder) MinWidth(v interface{}) (r *VListBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VListBuilder) Width(v interface{}) (r *VListBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VListBuilder) Elevation(v interface{}) (r *VListBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VListBuilder) ItemType(v string) (r *VListBuilder) {
	b.Attr("item-type", v)
	return b
}

func (b *VListBuilder) ItemTitle(v interface{}) (r *VListBuilder) {
	b.Attr(":item-title", h.JSONString(v))
	return b
}

func (b *VListBuilder) ItemValue(v interface{}) (r *VListBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VListBuilder) ItemChildren(v interface{}) (r *VListBuilder) {
	b.Attr(":item-children", h.JSONString(v))
	return b
}

func (b *VListBuilder) ItemProps(v interface{}) (r *VListBuilder) {
	b.Attr(":item-props", h.JSONString(v))
	return b
}

func (b *VListBuilder) ReturnObject(v bool) (r *VListBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) ValueComparator(v interface{}) (r *VListBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VListBuilder) Rounded(v interface{}) (r *VListBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VListBuilder) Tile(v bool) (r *VListBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Color(v string) (r *VListBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListBuilder) Variant(v interface{}) (r *VListBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VListBuilder) On(name string, value string) (r *VListBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListBuilder) Bind(name string, value string) (r *VListBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VListBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VListBuilder) SlotTitle(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VListBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VListBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VListBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VListBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VListBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VListBuilder) SlotAppend(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VListBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VListBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VListBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VListBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VListBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VListBuilder) SlotItem(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VListBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotDivider(child ...h.HTMLComponent) {
	b.SetSlot("divider", child...)
}

func (b *VListBuilder) SetScopedSlotDivider(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("divider", scope, child...)
}

func (b *VListBuilder) SlotDivider(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotDivider(child...)
	return b
}

func (b *VListBuilder) ScopedSlotDivider(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotDivider(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotSubheader(child ...h.HTMLComponent) {
	b.SetSlot("subheader", child...)
}

func (b *VListBuilder) SetScopedSlotSubheader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subheader", scope, child...)
}

func (b *VListBuilder) SlotSubheader(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotSubheader(child...)
	return b
}

func (b *VListBuilder) ScopedSlotSubheader(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotSubheader(scope, child...)
	return b
}

func (b *VListBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VListBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VListBuilder) SlotHeader(child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VListBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VListBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}
