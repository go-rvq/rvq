package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VXMarkdownEditorBuilder is vx-markdown-editor: a text in Markdown — rvq
// mdkit's editor: tables, footnotes, files included, pictures, preview, help
// —, its words vx's (Mdkit…).
type VXMarkdownEditorBuilder struct {
	v.VTagBuilder[*VXMarkdownEditorBuilder]
}

func VXMarkdownEditor(children ...h.HTMLComponent) *VXMarkdownEditorBuilder {
	return v.VTag(&VXMarkdownEditorBuilder{}, "vx-markdown-editor", children...)
}

func (b *VXMarkdownEditorBuilder) Label(v string) *VXMarkdownEditorBuilder {
	return b.Attr("label", v)
}

// HelpHref is the address of the guide of the editor: a ? beside its label
// (none: no ?).
func (b *VXMarkdownEditorBuilder) HelpHref(v string) *VXMarkdownEditorBuilder {
	return b.Attr("help-href", v)
}

func (b *VXMarkdownEditorBuilder) Hint(v string) *VXMarkdownEditorBuilder {
	return b.Attr("hint", v)
}

func (b *VXMarkdownEditorBuilder) Placeholder(v string) *VXMarkdownEditorBuilder {
	return b.Attr("placeholder", v)
}

func (b *VXMarkdownEditorBuilder) Rows(v int) *VXMarkdownEditorBuilder {
	return b.Attr(":rows", v)
}

func (b *VXMarkdownEditorBuilder) Disabled(v bool) *VXMarkdownEditorBuilder {
	return b.Attr(":disabled", v)
}

func (b *VXMarkdownEditorBuilder) ErrorMessages(v []string) *VXMarkdownEditorBuilder {
	return b.Attr(":error-messages", v)
}

// Preview is what renders the preview: a script of a function of the text
// giving its HTML (a promise), or an address (POST {format, body} → {html}).
func (b *VXMarkdownEditorBuilder) Preview(script string) *VXMarkdownEditorBuilder {
	return b.Attr(":preview", script)
}

// Upload is what sends a picture: a script of a function of the FormData
// (its "file") giving { url } or { message } (a promise); none, no picture.
func (b *VXMarkdownEditorBuilder) Upload(script string) *VXMarkdownEditorBuilder {
	return b.Attr(":upload", script)
}
