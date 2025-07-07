package vuetifyx

type (
	VXAdvancedCloseCardTagGetter[T any] interface {
		VXAdvancedCardTagGetter[T]
		GetVXAdvancedCloseCardTagBuilder() *VXAdvancedCloseCardTagBuilder[T]
	}

	VXAdvancedCloseCardTagger interface {
		VXAdvancedCardTagger
		SetSetup(s string)
		SetClosable(v bool)
		SetCloseIcon(s string)
		SetWidth(v string)
		SetVModel(s string)
		SetModelValue(s string)
	}
)

func (b *VXAdvancedCloseCardTagBuilder[T]) SetSetup(s string) {
	b.Attr(":setup", `() => {`+s+`}`)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetClosable(v bool) {
	b.Attr("closable", v)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetCloseIcon(s string) {
	b.Attr("close-icon", s)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetWidth(v string) {
	b.Attr("width", v)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetVModel(s string) {
	b.Attr("v-model", s)
}

func (b *VXAdvancedCloseCardTagBuilder[T]) SetModelValue(s string) {
	b.Attr("model-value", s)
}
