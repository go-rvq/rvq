package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBottomSheetBuilder struct {
	VTagBuilder[*VBottomSheetBuilder]
}

func VBottomSheet(children ...h.HTMLComponent) *VBottomSheetBuilder {
	return VTag(&VBottomSheetBuilder{}, "v-bottom-sheet", children...)
}

func (b *VBottomSheetBuilder) Activator(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Inset(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Fullscreen(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":fullscreen", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) RetainFocus(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":retain-focus", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Scrollable(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Absolute(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) CloseOnBack(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Contained(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) ContentClass(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) ContentProps(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Disabled(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Opacity(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) NoClickAnimation(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) ModelValue(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Persistent(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) Scrim(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) ZIndex(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Target(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) ActivatorProps(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) OpenOnClick(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) OpenOnHover(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) OpenOnFocus(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) CloseOnContentClick(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) CloseDelay(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) OpenDelay(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Height(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) MaxHeight(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) MaxWidth(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) MinHeight(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) MinWidth(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Width(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Eager(v bool) (r *VBottomSheetBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VBottomSheetBuilder) LocationStrategy(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Location(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Origin(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Offset(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) ScrollStrategy(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Theme(v string) (r *VBottomSheetBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBottomSheetBuilder) Transition(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) Attach(v interface{}) (r *VBottomSheetBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VBottomSheetBuilder) On(name string, value string) (r *VBottomSheetBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBottomSheetBuilder) Bind(name string, value string) (r *VBottomSheetBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
