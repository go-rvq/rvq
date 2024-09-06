package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBottomNavigationBuilder struct {
	VTagBuilder[*VBottomNavigationBuilder]
}

func VBottomNavigation(children ...h.HTMLComponent) *VBottomNavigationBuilder {
	return VTag(&VBottomNavigationBuilder{}, "v-bottom-navigation", children...)
}

func (b *VBottomNavigationBuilder) BaseColor(v string) (r *VBottomNavigationBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBottomNavigationBuilder) BgColor(v string) (r *VBottomNavigationBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VBottomNavigationBuilder) Color(v string) (r *VBottomNavigationBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBottomNavigationBuilder) Grow(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":grow", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Mode(v string) (r *VBottomNavigationBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VBottomNavigationBuilder) Height(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Active(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Border(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Density(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Elevation(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Rounded(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Tile(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Name(v string) (r *VBottomNavigationBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VBottomNavigationBuilder) Order(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Absolute(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Tag(v string) (r *VBottomNavigationBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBottomNavigationBuilder) ModelValue(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Multiple(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Max(v int) (r *VBottomNavigationBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) SelectedClass(v string) (r *VBottomNavigationBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VBottomNavigationBuilder) Disabled(v bool) (r *VBottomNavigationBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBottomNavigationBuilder) Mandatory(v interface{}) (r *VBottomNavigationBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VBottomNavigationBuilder) Theme(v string) (r *VBottomNavigationBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBottomNavigationBuilder) On(name string, value string) (r *VBottomNavigationBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBottomNavigationBuilder) Bind(name string, value string) (r *VBottomNavigationBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
