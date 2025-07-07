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

func (b *VListBuilder) ExpandIcon(v string) (r *VListBuilder) {
	b.Attr("expand-icon", v)
	return b
}

func (b *VListBuilder) CollapseIcon(v string) (r *VListBuilder) {
	b.Attr("collapse-icon", v)
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

func (b *VListBuilder) Nav(v bool) (r *VListBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
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

func (b *VListBuilder) Activated(v interface{}) (r *VListBuilder) {
	b.Attr(":activated", h.JSONString(v))
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

func (b *VListBuilder) Items(v interface{}) (r *VListBuilder) {
	b.Attr(":items", h.JSONString(v))
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

func (b *VListBuilder) Tag(v string) (r *VListBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListBuilder) Theme(v string) (r *VListBuilder) {
	b.Attr("theme", v)
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
