package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VChipGroupBuilder struct {
	VTagBuilder[*VChipGroupBuilder]
}

func VChipGroup(children ...h.HTMLComponent) *VChipGroupBuilder {
	return VTag(&VChipGroupBuilder{}, "v-chip-group", children...)
}

func (b *VChipGroupBuilder) Symbol(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) Column(v bool) (r *VChipGroupBuilder) {
	b.Attr(":column", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) Filter(v bool) (r *VChipGroupBuilder) {
	b.Attr(":filter", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) ValueComparator(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) CenterActive(v bool) (r *VChipGroupBuilder) {
	b.Attr(":center-active", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) Direction(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) NextIcon(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) PrevIcon(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) ShowArrows(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) Mobile(v bool) (r *VChipGroupBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) MobileBreakpoint(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) Tag(v string) (r *VChipGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VChipGroupBuilder) ModelValue(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) Multiple(v bool) (r *VChipGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) Max(v int) (r *VChipGroupBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) SelectedClass(v string) (r *VChipGroupBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VChipGroupBuilder) Disabled(v bool) (r *VChipGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VChipGroupBuilder) Mandatory(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) Theme(v string) (r *VChipGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VChipGroupBuilder) Color(v string) (r *VChipGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VChipGroupBuilder) Variant(v interface{}) (r *VChipGroupBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VChipGroupBuilder) On(name string, value string) (r *VChipGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VChipGroupBuilder) Bind(name string, value string) (r *VChipGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
