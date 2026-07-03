package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTreeviewBuilder struct {
	VTagBuilder[*VTreeviewBuilder]
}

func VTreeview(children ...h.HTMLComponent) *VTreeviewBuilder {
	return VTag(&VTreeviewBuilder{}, "v-treeview", children...)
}

func (b *VTreeviewBuilder) Search(v string) (r *VTreeviewBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VTreeviewBuilder) Tag(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Activated(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":activated", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Theme(v string) (r *VTreeviewBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTreeviewBuilder) Items(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) FilterMode(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) NoFilter(v bool) (r *VTreeviewBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) CustomFilter(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) CustomKeyFilter(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) FilterKeys(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) BaseColor(v string) (r *VTreeviewBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTreeviewBuilder) ActiveColor(v string) (r *VTreeviewBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VTreeviewBuilder) ActiveClass(v string) (r *VTreeviewBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VTreeviewBuilder) BgColor(v string) (r *VTreeviewBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTreeviewBuilder) Disabled(v bool) (r *VTreeviewBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) ExpandIcon(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) CollapseIcon(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Lines(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Slim(v bool) (r *VTreeviewBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Activatable(v bool) (r *VTreeviewBuilder) {
	b.Attr(":activatable", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Selectable(v bool) (r *VTreeviewBuilder) {
	b.Attr(":selectable", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Opened(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":opened", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Selected(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":selected", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Mandatory(v bool) (r *VTreeviewBuilder) {
	b.Attr(":mandatory", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) ActiveStrategy(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":active-strategy", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) SelectStrategy(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Border(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Density(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Height(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) MaxHeight(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) MaxWidth(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) MinHeight(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) MinWidth(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Width(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Elevation(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ItemTitle(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":item-title", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ItemValue(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ItemChildren(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":item-children", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ItemProps(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":item-props", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ReturnObject(v bool) (r *VTreeviewBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) ValueComparator(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Rounded(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Tile(v bool) (r *VTreeviewBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Color(v string) (r *VTreeviewBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTreeviewBuilder) Variant(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) ModelValue(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) OpenOnClick(v bool) (r *VTreeviewBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Fluid(v bool) (r *VTreeviewBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) FalseIcon(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) TrueIcon(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) OpenAll(v bool) (r *VTreeviewBuilder) {
	b.Attr(":open-all", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) LoadingIcon(v string) (r *VTreeviewBuilder) {
	b.Attr("loading-icon", v)
	return b
}

func (b *VTreeviewBuilder) IndeterminateIcon(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":indeterminate-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) SelectedColor(v string) (r *VTreeviewBuilder) {
	b.Attr("selected-color", v)
	return b
}

func (b *VTreeviewBuilder) LoadChildren(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":load-children", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) On(name string, value string) (r *VTreeviewBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTreeviewBuilder) Bind(name string, value string) (r *VTreeviewBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTreeviewBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTreeviewBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTreeviewBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VTreeviewBuilder) SlotTitle(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VTreeviewBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VTreeviewBuilder) SlotAppend(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VTreeviewBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTreeviewBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VTreeviewBuilder) SlotItem(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotDivider(child ...h.HTMLComponent) {
	b.SetSlot("divider", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotDivider(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("divider", scope, child...)
}

func (b *VTreeviewBuilder) SlotDivider(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotDivider(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotDivider(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotDivider(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotSubheader(child ...h.HTMLComponent) {
	b.SetSlot("subheader", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotSubheader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subheader", scope, child...)
}

func (b *VTreeviewBuilder) SlotSubheader(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotSubheader(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotSubheader(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotSubheader(scope, child...)
	return b
}

func (b *VTreeviewBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VTreeviewBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VTreeviewBuilder) SlotHeader(child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VTreeviewBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VTreeviewBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}
