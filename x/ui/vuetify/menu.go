package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VMenuBuilder struct {
	VTagBuilder[*VMenuBuilder]
}

func VMenu(children ...h.HTMLComponent) *VMenuBuilder {
	return VTag(&VMenuBuilder{}, "v-menu", children...)
}

func (b *VMenuBuilder) Activator(v interface{}) (r *VMenuBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Id(v string) (r *VMenuBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VMenuBuilder) CloseOnBack(v bool) (r *VMenuBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) Contained(v bool) (r *VMenuBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) ContentClass(v interface{}) (r *VMenuBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) ContentProps(v interface{}) (r *VMenuBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Disabled(v bool) (r *VMenuBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) Opacity(v interface{}) (r *VMenuBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) NoClickAnimation(v bool) (r *VMenuBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) ModelValue(v bool) (r *VMenuBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) Persistent(v bool) (r *VMenuBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) Scrim(v interface{}) (r *VMenuBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) ZIndex(v interface{}) (r *VMenuBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Target(v interface{}) (r *VMenuBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) ActivatorProps(v interface{}) (r *VMenuBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) OpenOnClick(v bool) (r *VMenuBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) OpenOnHover(v bool) (r *VMenuBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) OpenOnFocus(v bool) (r *VMenuBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) CloseOnContentClick(v bool) (r *VMenuBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) CloseDelay(v interface{}) (r *VMenuBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) OpenDelay(v interface{}) (r *VMenuBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Height(v interface{}) (r *VMenuBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) MaxHeight(v interface{}) (r *VMenuBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) MaxWidth(v interface{}) (r *VMenuBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) MinHeight(v interface{}) (r *VMenuBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) MinWidth(v interface{}) (r *VMenuBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Width(v interface{}) (r *VMenuBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Eager(v bool) (r *VMenuBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VMenuBuilder) LocationStrategy(v interface{}) (r *VMenuBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Location(v interface{}) (r *VMenuBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Origin(v interface{}) (r *VMenuBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Offset(v interface{}) (r *VMenuBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) ScrollStrategy(v interface{}) (r *VMenuBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Theme(v string) (r *VMenuBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VMenuBuilder) Transition(v interface{}) (r *VMenuBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) Attach(v interface{}) (r *VMenuBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VMenuBuilder) On(name string, value string) (r *VMenuBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VMenuBuilder) Bind(name string, value string) (r *VMenuBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
