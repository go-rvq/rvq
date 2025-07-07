package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSpeedDialBuilder struct {
	VTagBuilder[*VSpeedDialBuilder]
}

func VSpeedDial(children ...h.HTMLComponent) *VSpeedDialBuilder {
	return VTag(&VSpeedDialBuilder{}, "v-speed-dial", children...)
}

func (b *VSpeedDialBuilder) Activator(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Id(v string) (r *VSpeedDialBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSpeedDialBuilder) CloseOnBack(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) Contained(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) ContentClass(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) ContentProps(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Disabled(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) Opacity(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) NoClickAnimation(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) ModelValue(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) Persistent(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) Scrim(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) ZIndex(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Target(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) ActivatorProps(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) OpenOnClick(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) OpenOnHover(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) OpenOnFocus(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) CloseOnContentClick(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) CloseDelay(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) OpenDelay(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Height(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) MaxHeight(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) MaxWidth(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) MinHeight(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) MinWidth(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Width(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Eager(v bool) (r *VSpeedDialBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VSpeedDialBuilder) LocationStrategy(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Location(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Origin(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Offset(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) ScrollStrategy(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Theme(v string) (r *VSpeedDialBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSpeedDialBuilder) Transition(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) Attach(v interface{}) (r *VSpeedDialBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VSpeedDialBuilder) On(name string, value string) (r *VSpeedDialBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSpeedDialBuilder) Bind(name string, value string) (r *VSpeedDialBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
