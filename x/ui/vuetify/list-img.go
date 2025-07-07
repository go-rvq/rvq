package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListImgBuilder struct {
	VTagBuilder[*VListImgBuilder]
}

func VListImg(children ...h.HTMLComponent) *VListImgBuilder {
	return VTag(&VListImgBuilder{}, "v-list-img", children...)
}

func (b *VListImgBuilder) Tag(v string) (r *VListImgBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListImgBuilder) On(name string, value string) (r *VListImgBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListImgBuilder) Bind(name string, value string) (r *VListImgBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
