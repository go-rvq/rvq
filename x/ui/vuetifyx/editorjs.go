package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/ui/vuetify"
)

type EditorJsBuilder struct {
	vuetify.VTagBuilder[*EditorJsBuilder]
}

func EditorJS() *EditorJsBuilder {
	return vuetify.VTag(&EditorJsBuilder{}, "vx-editorjs")
}

func (b *EditorJsBuilder) Label(text string) *EditorJsBuilder {
	return b.Attr("label", text)
}

func (b *EditorJsBuilder) ErrorMessages(msgs []string) *EditorJsBuilder {
	if len(msgs) > 0 {
		return b.Attr(":error-messages", h.JSONString(msgs))
	}
	return b
}
