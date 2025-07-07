package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPickerTitleBuilder struct {
	VTagBuilder[*VPickerTitleBuilder]
}

func VPickerTitle(children ...h.HTMLComponent) *VPickerTitleBuilder {
	return VTag(&VPickerTitleBuilder{}, "v-picker-title", children...)
}

func (b *VPickerTitleBuilder) Tag(v string) (r *VPickerTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VPickerTitleBuilder) On(name string, value string) (r *VPickerTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPickerTitleBuilder) Bind(name string, value string) (r *VPickerTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
