package vuetify

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/tag"
)

type (
	VTagBuilderGetter[T any] interface {
		tag.TagBuilderGetter[T]
		GetVTagBuilder() *VTagBuilder[T]
	}

	VTagBuilder[T any] struct {
		tag.TagBuilder[T]
	}
)

func VTag[T VTagBuilderGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	vtb := dot.GetVTagBuilder()
	vtb.TagBuilder = *tag.NewTag(dot, name, children...).GetTagBuilder()
	return dot
}

func (b *VTagBuilder[T]) GetVTagBuilder() *VTagBuilder[T] {
	return b
}

func (t *VTagBuilder[T]) RawWidth(v interface{}) T {
	return t.Attr(":width", v)
}

func (t *VTagBuilder[T]) RawHeight(v interface{}) T {
	return t.Attr(":height", v)
}

func (t *VTagBuilder[T]) RawClass(v interface{}) T {
	return t.Attr(":class", v)
}

func (t *VTagBuilder[T]) FormField(formKey string, v interface{}) T {
	return t.Attr(web.VField(formKey, v)...)
}
