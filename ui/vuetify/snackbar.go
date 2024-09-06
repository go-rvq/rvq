package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSnackbarBuilder struct {
	VTagBuilder[*VSnackbarBuilder]
}

func VSnackbar(children ...h.HTMLComponent) *VSnackbarBuilder {
	return VTag(&VSnackbarBuilder{}, "v-snackbar", children...)
}

func (b *VSnackbarBuilder) Activator(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Text(v string) (r *VSnackbarBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VSnackbarBuilder) MultiLine(v bool) (r *VSnackbarBuilder) {
	b.Attr(":multi-line", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Timer(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":timer", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Timeout(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":timeout", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Vertical(v bool) (r *VSnackbarBuilder) {
	b.Attr(":vertical", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Location(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Position(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Absolute(v bool) (r *VSnackbarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Rounded(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Tile(v bool) (r *VSnackbarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Color(v string) (r *VSnackbarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSnackbarBuilder) Variant(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Theme(v string) (r *VSnackbarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSnackbarBuilder) CloseOnBack(v bool) (r *VSnackbarBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Contained(v bool) (r *VSnackbarBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) ContentClass(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) ContentProps(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Disabled(v bool) (r *VSnackbarBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) Opacity(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) ModelValue(v bool) (r *VSnackbarBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) ZIndex(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Target(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) ActivatorProps(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) OpenOnClick(v bool) (r *VSnackbarBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) OpenOnHover(v bool) (r *VSnackbarBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) OpenOnFocus(v bool) (r *VSnackbarBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) CloseOnContentClick(v bool) (r *VSnackbarBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) CloseDelay(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) OpenDelay(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Height(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) MaxHeight(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) MaxWidth(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) MinHeight(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) MinWidth(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Width(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Eager(v bool) (r *VSnackbarBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VSnackbarBuilder) LocationStrategy(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Origin(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Offset(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Transition(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) Attach(v interface{}) (r *VSnackbarBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VSnackbarBuilder) On(name string, value string) (r *VSnackbarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSnackbarBuilder) Bind(name string, value string) (r *VSnackbarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
