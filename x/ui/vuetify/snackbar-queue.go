package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSnackbarQueueBuilder struct {
	VTagBuilder[*VSnackbarQueueBuilder]
}

func VSnackbarQueue(children ...h.HTMLComponent) *VSnackbarQueueBuilder {
	return VTag(&VSnackbarQueueBuilder{}, "v-snackbar-queue", children...)
}

func (b *VSnackbarQueueBuilder) Activator(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Text(v string) (r *VSnackbarQueueBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VSnackbarQueueBuilder) MultiLine(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":multi-line", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Timer(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":timer", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Timeout(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":timeout", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Vertical(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":vertical", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Location(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Position(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Absolute(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Rounded(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Tile(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Color(v string) (r *VSnackbarQueueBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSnackbarQueueBuilder) Variant(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Theme(v string) (r *VSnackbarQueueBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSnackbarQueueBuilder) CloseOnBack(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Contained(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) ContentClass(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) ContentProps(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Disabled(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) Opacity(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) ModelValue(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) ZIndex(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Target(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) ActivatorProps(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) OpenOnClick(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) OpenOnHover(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) OpenOnFocus(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) CloseOnContentClick(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) CloseDelay(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) OpenDelay(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Height(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) MaxHeight(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) MaxWidth(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) MinHeight(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) MinWidth(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Width(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Eager(v bool) (r *VSnackbarQueueBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VSnackbarQueueBuilder) LocationStrategy(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Origin(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Offset(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Transition(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Attach(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) Closable(v interface{}) (r *VSnackbarQueueBuilder) {
	b.Attr(":closable", h.JSONString(v))
	return b
}

func (b *VSnackbarQueueBuilder) CloseText(v string) (r *VSnackbarQueueBuilder) {
	b.Attr("close-text", v)
	return b
}

func (b *VSnackbarQueueBuilder) On(name string, value string) (r *VSnackbarQueueBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSnackbarQueueBuilder) Bind(name string, value string) (r *VSnackbarQueueBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
