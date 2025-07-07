package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTabsBuilder struct {
	VTagBuilder[*VTabsBuilder]
}

func VTabs(children ...h.HTMLComponent) *VTabsBuilder {
	return VTag(&VTabsBuilder{}, "v-tabs", children...)
}

func (b *VTabsBuilder) Symbol(v interface{}) (r *VTabsBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) AlignTabs(v interface{}) (r *VTabsBuilder) {
	b.Attr(":align-tabs", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Color(v string) (r *VTabsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTabsBuilder) FixedTabs(v bool) (r *VTabsBuilder) {
	b.Attr(":fixed-tabs", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Items(v interface{}) (r *VTabsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Stacked(v bool) (r *VTabsBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) BgColor(v string) (r *VTabsBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTabsBuilder) Grow(v bool) (r *VTabsBuilder) {
	b.Attr(":grow", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Height(v interface{}) (r *VTabsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) HideSlider(v bool) (r *VTabsBuilder) {
	b.Attr(":hide-slider", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) SliderColor(v string) (r *VTabsBuilder) {
	b.Attr("slider-color", v)
	return b
}

func (b *VTabsBuilder) CenterActive(v bool) (r *VTabsBuilder) {
	b.Attr(":center-active", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Direction(v interface{}) (r *VTabsBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) NextIcon(v interface{}) (r *VTabsBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) PrevIcon(v interface{}) (r *VTabsBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) ShowArrows(v interface{}) (r *VTabsBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Mobile(v bool) (r *VTabsBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) MobileBreakpoint(v interface{}) (r *VTabsBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Tag(v string) (r *VTabsBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTabsBuilder) ModelValue(v interface{}) (r *VTabsBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Multiple(v bool) (r *VTabsBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Max(v int) (r *VTabsBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) SelectedClass(v string) (r *VTabsBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsBuilder) Disabled(v bool) (r *VTabsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsBuilder) Mandatory(v interface{}) (r *VTabsBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) Density(v interface{}) (r *VTabsBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTabsBuilder) On(name string, value string) (r *VTabsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsBuilder) Bind(name string, value string) (r *VTabsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
