package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardItemBuilder struct {
	VTagBuilder[*VCardItemBuilder]
}

func VCardItem(children ...h.HTMLComponent) *VCardItemBuilder {
	return VTag(&VCardItemBuilder{}, "v-card-item", children...)
}

func (b *VCardItemBuilder) Title(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) Subtitle(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) AppendAvatar(v string) (r *VCardItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VCardItemBuilder) AppendIcon(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) PrependAvatar(v string) (r *VCardItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VCardItemBuilder) PrependIcon(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) Density(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) On(name string, value string) (r *VCardItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardItemBuilder) Bind(name string, value string) (r *VCardItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
