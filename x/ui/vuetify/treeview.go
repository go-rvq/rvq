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

func (b *VTreeviewBuilder) OpenAll(v bool) (r *VTreeviewBuilder) {
	b.Attr(":open-all", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) Search(v string) (r *VTreeviewBuilder) {
	b.Attr("search", v)
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

func (b *VTreeviewBuilder) LoadingIcon(v string) (r *VTreeviewBuilder) {
	b.Attr("loading-icon", v)
	return b
}

func (b *VTreeviewBuilder) Selectable(v bool) (r *VTreeviewBuilder) {
	b.Attr(":selectable", fmt.Sprint(v))
	return b
}

func (b *VTreeviewBuilder) LoadChildren(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":load-children", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Items(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) SelectStrategy(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
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

func (b *VTreeviewBuilder) ExpandIcon(v string) (r *VTreeviewBuilder) {
	b.Attr("expand-icon", v)
	return b
}

func (b *VTreeviewBuilder) CollapseIcon(v string) (r *VTreeviewBuilder) {
	b.Attr("collapse-icon", v)
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

func (b *VTreeviewBuilder) Opened(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":opened", h.JSONString(v))
	return b
}

func (b *VTreeviewBuilder) Activated(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":activated", h.JSONString(v))
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

func (b *VTreeviewBuilder) OpenStrategy(v interface{}) (r *VTreeviewBuilder) {
	b.Attr(":open-strategy", h.JSONString(v))
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

func (b *VTreeviewBuilder) ItemType(v string) (r *VTreeviewBuilder) {
	b.Attr("item-type", v)
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

func (b *VTreeviewBuilder) Tag(v string) (r *VTreeviewBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTreeviewBuilder) Theme(v string) (r *VTreeviewBuilder) {
	b.Attr("theme", v)
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

func (b *VTreeviewBuilder) On(name string, value string) (r *VTreeviewBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTreeviewBuilder) Bind(name string, value string) (r *VTreeviewBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
