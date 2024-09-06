package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTooltipBuilder struct {
	VTagBuilder[*VTooltipBuilder]
}

func VTooltip(children ...h.HTMLComponent) *VTooltipBuilder {
	return VTag(&VTooltipBuilder{}, "v-tooltip", children...)
}

func (b *VTooltipBuilder) Activator(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Id(v string) (r *VTooltipBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VTooltipBuilder) Text(v string) (r *VTooltipBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VTooltipBuilder) CloseOnBack(v bool) (r *VTooltipBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) Contained(v bool) (r *VTooltipBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) ContentClass(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) ContentProps(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Disabled(v bool) (r *VTooltipBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) Opacity(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) NoClickAnimation(v bool) (r *VTooltipBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) ModelValue(v bool) (r *VTooltipBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) Scrim(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) ZIndex(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Target(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) ActivatorProps(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) OpenOnClick(v bool) (r *VTooltipBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) OpenOnHover(v bool) (r *VTooltipBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) OpenOnFocus(v bool) (r *VTooltipBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) CloseOnContentClick(v bool) (r *VTooltipBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) CloseDelay(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) OpenDelay(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Height(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) MaxHeight(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) MaxWidth(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) MinHeight(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) MinWidth(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Width(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Eager(v bool) (r *VTooltipBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VTooltipBuilder) LocationStrategy(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Location(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Origin(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Offset(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) ScrollStrategy(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Theme(v string) (r *VTooltipBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTooltipBuilder) Transition(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) Attach(v interface{}) (r *VTooltipBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VTooltipBuilder) On(name string, value string) (r *VTooltipBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTooltipBuilder) Bind(name string, value string) (r *VTooltipBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
