package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXBtnBuilder struct {
	v.VTagBuilder[*VXBtnBuilder]
}

func VXBtn(text string, children ...h.HTMLComponent) *VXBtnBuilder {
	return v.VTag(&VXBtnBuilder{}, "vx-btn", children...).Text(text)
}

func (b *VXBtnBuilder) OnClick(eventFuncId string) (r *VXBtnBuilder) {
	b.Attr("@click", web.POST().EventFunc(eventFuncId).Go())
	return b
}

func (b *VXBtnBuilder) Symbol(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Flat(v bool) (r *VXBtnBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Active(v bool) (r *VXBtnBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) BaseColor(v string) (r *VXBtnBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VXBtnBuilder) PrependIcon(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) AppendIcon(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Block(v bool) (r *VXBtnBuilder) {
	b.Attr(":block", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Readonly(v bool) (r *VXBtnBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Slim(v bool) (r *VXBtnBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Stacked(v bool) (r *VXBtnBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Ripple(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Value(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Text(v string) (r *VXBtnBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VXBtnBuilder) Title(v string) (r *VXBtnBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VXBtnBuilder) Border(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Density(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Height(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) MaxHeight(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) MaxWidth(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) MinHeight(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) MinWidth(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Width(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Elevation(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Disabled(v bool) (r *VXBtnBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) SelectedClass(v string) (r *VXBtnBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VXBtnBuilder) Loading(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Location(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Position(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Rounded(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Tile(v bool) (r *VXBtnBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Href(v string) (r *VXBtnBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VXBtnBuilder) Replace(v bool) (r *VXBtnBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) Exact(v bool) (r *VXBtnBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VXBtnBuilder) To(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Size(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Tag(v string) (r *VXBtnBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VXBtnBuilder) Theme(v string) (r *VXBtnBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VXBtnBuilder) Color(v string) (r *VXBtnBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VXBtnBuilder) Variant(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) Icon(v interface{}) (r *VXBtnBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VXBtnBuilder) On(name string, value string) (r *VXBtnBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VXBtnBuilder) Bind(name string, value string) (r *VXBtnBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
