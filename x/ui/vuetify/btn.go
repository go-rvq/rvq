package vuetify

import (
	"github.com/qor5/web/v3"

	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBtnBuilder struct {
	VTagBuilder[*VBtnBuilder]
}

func VBtn(text string, children ...h.HTMLComponent) *VBtnBuilder {
	return VTag(&VBtnBuilder{}, "v-btn", children...).Text(text)
}

func (b *VBtnBuilder) OnClick(eventFuncId string) (r *VBtnBuilder) {
	b.Attr("@click", web.POST().EventFunc(eventFuncId).Go())
	return b
}

func (b *VBtnBuilder) Symbol(v interface{}) (r *VBtnBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Flat(v bool) (r *VBtnBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Active(v bool) (r *VBtnBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) BaseColor(v string) (r *VBtnBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnBuilder) PrependIcon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) AppendIcon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Block(v bool) (r *VBtnBuilder) {
	b.Attr(":block", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Readonly(v bool) (r *VBtnBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Slim(v bool) (r *VBtnBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Stacked(v bool) (r *VBtnBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Ripple(v interface{}) (r *VBtnBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Value(v interface{}) (r *VBtnBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Text(v string) (r *VBtnBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VBtnBuilder) Border(v interface{}) (r *VBtnBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Density(v interface{}) (r *VBtnBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Height(v interface{}) (r *VBtnBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MaxHeight(v interface{}) (r *VBtnBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MaxWidth(v interface{}) (r *VBtnBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MinHeight(v interface{}) (r *VBtnBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) MinWidth(v interface{}) (r *VBtnBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Width(v interface{}) (r *VBtnBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Elevation(v interface{}) (r *VBtnBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Disabled(v bool) (r *VBtnBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) SelectedClass(v string) (r *VBtnBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VBtnBuilder) Loading(v interface{}) (r *VBtnBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Location(v interface{}) (r *VBtnBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Position(v interface{}) (r *VBtnBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Rounded(v interface{}) (r *VBtnBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Tile(v bool) (r *VBtnBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Href(v string) (r *VBtnBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VBtnBuilder) Replace(v bool) (r *VBtnBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) Exact(v bool) (r *VBtnBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VBtnBuilder) To(v interface{}) (r *VBtnBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Size(v interface{}) (r *VBtnBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Tag(v string) (r *VBtnBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBtnBuilder) Theme(v string) (r *VBtnBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnBuilder) Color(v string) (r *VBtnBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnBuilder) Variant(v interface{}) (r *VBtnBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Icon(v interface{}) (r *VBtnBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) On(name string, value string) (r *VBtnBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnBuilder) Bind(name string, value string) (r *VBtnBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
