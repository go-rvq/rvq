package tiptap

// @snippet_begin(TipTapEditorHTMLComponent)
import (
	h "github.com/go-rvq/htmlgo"
)

type TipTapEditorBuilder struct {
	tag *h.HTMLTagBuilder
}

func TipTapEditor() (r *TipTapEditorBuilder) {
	r = &TipTapEditorBuilder{
		tag: h.Tag("tiptap-editor"),
	}

	return
}

func (b *TipTapEditorBuilder) Attr(vs ...interface{}) (r *TipTapEditorBuilder) {
	b.tag.Attr(vs...)
	return b
}

func (b *TipTapEditorBuilder) Write(ctx *h.Context) (err error) {
	return b.tag.Write(ctx)
}

// @snippet_end
