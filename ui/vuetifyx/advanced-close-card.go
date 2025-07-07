package vuetifyx

import (
	"github.com/qor5/web/v3/tag"
	h "github.com/theplant/htmlgo"
)

type (
	VXAdvancedCloseCardTagBuilder[T any] struct {
		VXAdvancedCardTagBuilder[T]
	}
)

func VXAdvancedCloseCardTag[T VXAdvancedCloseCardTagGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	vtb := dot.GetVXAdvancedCloseCardTagBuilder()
	vtb.TagBuilder = *tag.NewTag(dot, name, children...).GetTagBuilder()
	return dot
}

func (b *VXAdvancedCloseCardTagBuilder[T]) GetVXAdvancedCloseCardTagBuilder() *VXAdvancedCloseCardTagBuilder[T] {
	return b
}

func (b *VXAdvancedCloseCardTagBuilder[T]) Closable(v bool) T {
	b.SetClosable(v)
	return b.Dot()
}

func (b *VXAdvancedCloseCardTagBuilder[T]) CloseIcon(s string) T {
	b.SetCloseIcon(s)
	return b.Dot()
}

func (b *VXAdvancedCloseCardTagBuilder[T]) Width(v string) T {
	b.SetWidth(v)
	return b.Dot()
}

func (b *VXAdvancedCloseCardTagBuilder[T]) VModel(s string) T {
	b.SetVModel(s)
	return b.Dot()
}

func (b *VXAdvancedCloseCardTagBuilder[T]) ModelValue(s string) T {
	b.SetModelValue(s)
	return b.Dot()
}
