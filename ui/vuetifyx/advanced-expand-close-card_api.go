package vuetifyx

type (
	VXAdvancedExpandCloseCardTagGetter[T any] interface {
		VXAdvancedCloseCardTagGetter[T]
		GetVXAdvancedExpandCloseCardTagBuilder() *VXAdvancedExpandCloseCardTagBuilder[T]
	}

	VXAdvancedExpandCloseCardTagger interface {
		VXAdvancedCloseCardTagger
		SetExpandable(v bool)
		SetDialogOrDrawer(dialog bool)
	}
)

func (b *VXAdvancedCloseCardTagBuilder[T]) SetExpandable(v bool) {
	b.Attr("expandable", v)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetDialogOrDrawer(dialog bool) {
	if dialog {
		b.SetTag("vx-dialog")
	} else {
		b.SetTag("vx-navigation-drawer")
	}
}
