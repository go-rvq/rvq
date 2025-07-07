package vuetify

import (
	"fmt"

	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type VDialogBuilder struct {
	VTagBuilder[*VDialogBuilder]
}

func VDialog(children ...h.HTMLComponent) *VDialogBuilder {
	return VTag(&VDialogBuilder{}, "v-dialog", children...)
}

func (b *VDialogBuilder) Activator(v interface{}) (r *VDialogBuilder) {
	b.Attr(":activator", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Fullscreen(v any) (r *VDialogBuilder) {
	b.Attr(":fullscreen", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) RetainFocus(v bool) (r *VDialogBuilder) {
	b.Attr(":retain-focus", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Scrollable(v any) (r *VDialogBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Absolute(v bool) (r *VDialogBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) CloseOnBack(v bool) (r *VDialogBuilder) {
	b.Attr(":close-on-back", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Contained(v bool) (r *VDialogBuilder) {
	b.Attr(":contained", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) ContentClass(v interface{}) (r *VDialogBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) ContentProps(v interface{}) (r *VDialogBuilder) {
	b.Attr(":content-props", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Disabled(v bool) (r *VDialogBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Opacity(v interface{}) (r *VDialogBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) NoClickAnimation(v bool) (r *VDialogBuilder) {
	b.Attr(":no-click-animation", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) ModelValue(v bool) (r *VDialogBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Persistent(v bool) (r *VDialogBuilder) {
	b.Attr(":persistent", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) Scrim(v interface{}) (r *VDialogBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) ZIndex(v interface{}) (r *VDialogBuilder) {
	b.Attr(":z-index", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Target(v interface{}) (r *VDialogBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) ActivatorProps(v interface{}) (r *VDialogBuilder) {
	b.Attr(":activator-props", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) OpenOnClick(v bool) (r *VDialogBuilder) {
	b.Attr(":open-on-click", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) OpenOnHover(v bool) (r *VDialogBuilder) {
	b.Attr(":open-on-hover", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) OpenOnFocus(v bool) (r *VDialogBuilder) {
	b.Attr(":open-on-focus", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) CloseOnContentClick(v bool) (r *VDialogBuilder) {
	b.Attr(":close-on-content-click", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) CloseDelay(v interface{}) (r *VDialogBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) OpenDelay(v interface{}) (r *VDialogBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Height(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":height", v)
	return b
}

func (b *VDialogBuilder) MaxHeight(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":max-height", v)
	return b
}

func (b *VDialogBuilder) MaxWidth(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":max-width", v)
	return b
}

func (b *VDialogBuilder) MinHeight(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":min-height", v)
	return b
}

func (b *VDialogBuilder) MinWidth(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":min-width", v)
	return b
}

func (b *VDialogBuilder) Width(v interface{}) (r *VDialogBuilder) {
	if va, ok := v.(web.Var); ok {
		v = string(va)
	} else {
		v = h.JSONString(v)
	}
	b.Attr(":width", v)
	return b
}

func (b *VDialogBuilder) Eager(v bool) (r *VDialogBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VDialogBuilder) LocationStrategy(v interface{}) (r *VDialogBuilder) {
	b.Attr(":location-strategy", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Location(v interface{}) (r *VDialogBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Origin(v interface{}) (r *VDialogBuilder) {
	b.Attr(":origin", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Offset(v interface{}) (r *VDialogBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) ScrollStrategy(v interface{}) (r *VDialogBuilder) {
	b.Attr(":scroll-strategy", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Theme(v string) (r *VDialogBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDialogBuilder) Transition(v interface{}) (r *VDialogBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) Attach(v interface{}) (r *VDialogBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VDialogBuilder) On(name string, value string) (r *VDialogBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogBuilder) Bind(name string, value string) (r *VDialogBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
