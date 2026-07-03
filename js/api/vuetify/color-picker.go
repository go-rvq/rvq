package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VColorPickerBuilder struct {
	VTagBuilder[*VColorPickerBuilder]
}

func VColorPicker(children ...h.HTMLComponent) *VColorPickerBuilder {
	return VTag(&VColorPickerBuilder{}, "v-color-picker", children...)
}

func (b *VColorPickerBuilder) Title(v string) (r *VColorPickerBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VColorPickerBuilder) Border(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) ModelValue(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Height(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) MaxHeight(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) MaxWidth(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) MinHeight(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) MinWidth(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Width(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Elevation(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Location(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Position(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Rounded(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Tile(v bool) (r *VColorPickerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) Tag(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Theme(v string) (r *VColorPickerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VColorPickerBuilder) Color(v string) (r *VColorPickerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VColorPickerBuilder) Disabled(v bool) (r *VColorPickerBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) BgColor(v string) (r *VColorPickerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VColorPickerBuilder) Mode(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Divided(v bool) (r *VColorPickerBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) HideHeader(v bool) (r *VColorPickerBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) CanvasHeight(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":canvas-height", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) DotSize(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":dot-size", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) HideCanvas(v bool) (r *VColorPickerBuilder) {
	b.Attr(":hide-canvas", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) HideSliders(v bool) (r *VColorPickerBuilder) {
	b.Attr(":hide-sliders", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) HideInputs(v bool) (r *VColorPickerBuilder) {
	b.Attr(":hide-inputs", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) Modes(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":modes", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) ShowSwatches(v bool) (r *VColorPickerBuilder) {
	b.Attr(":show-swatches", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) SwatchesMaxHeight(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":swatches-max-height", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) Landscape(v bool) (r *VColorPickerBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VColorPickerBuilder) Swatches(v interface{}) (r *VColorPickerBuilder) {
	b.Attr(":swatches", h.JSONString(v))
	return b
}

func (b *VColorPickerBuilder) On(name string, value string) (r *VColorPickerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VColorPickerBuilder) Bind(name string, value string) (r *VColorPickerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VColorPickerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VColorPickerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VColorPickerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VColorPickerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VColorPickerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VColorPickerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}
