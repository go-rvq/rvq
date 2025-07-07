package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/tag"
)

type (
	VXAdvancedExpandCloseCardTagBuilder[T any] struct {
		VXAdvancedCloseCardTagBuilder[T]
	}

	VXAdvancedExpandCloseCardBuilder struct {
		VXAdvancedExpandCloseCardTagBuilder[*VXAdvancedExpandCloseCardBuilder]
	}
)

func VXAdvancedExpandCloseCardTag[T VXAdvancedExpandCloseCardTagGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	vtb := dot.GetVXAdvancedExpandCloseCardTagBuilder()
	vtb.TagBuilder = *tag.NewTag(dot, name, children...).GetTagBuilder()
	return dot
}

func (b *VXAdvancedExpandCloseCardTagBuilder[T]) GetVXAdvancedExpandCloseCardTagBuilder() *VXAdvancedExpandCloseCardTagBuilder[T] {
	return b
}

func (b *VXAdvancedCloseCardTagBuilder[T]) Expandable(v bool) T {
	b.SetExpandable(v)
	return b.Dot()
}

func (b *VXAdvancedCloseCardTagBuilder[T]) DialogOrDrawer(dialog bool) T {
	b.SetDialogOrDrawer(dialog)
	return b.Dot()
}

func VXAdvancedExpandCloseCard(tag string, children ...h.HTMLComponent) *VXAdvancedExpandCloseCardBuilder {
	return VXAdvancedExpandCloseCardTag(&VXAdvancedExpandCloseCardBuilder{}, tag, children...)
}
