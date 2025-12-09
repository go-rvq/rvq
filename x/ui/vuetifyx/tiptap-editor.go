package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXTipTapEditorBuilder struct {
	v.VTagBuilder[*VXTipTapEditorBuilder]
}

func VXTipTapEditor(children ...h.HTMLComponent) *VXTipTapEditorBuilder {
	return v.VTag(&VXTipTapEditorBuilder{}, "vx-tiptap-editor", children...)
}

func (b *VXTipTapEditorBuilder) Label(v string) *VXTipTapEditorBuilder {
	return b.Attr("label", v)
}

func (b *VXTipTapEditorBuilder) Hint(v string) *VXTipTapEditorBuilder {
	return b.Attr("hint", v)
}

func (b *VXTipTapEditorBuilder) Error(v bool) *VXTipTapEditorBuilder {
	return b.Attr("error", v)
}

func (b *VXTipTapEditorBuilder) Lang(v string) *VXTipTapEditorBuilder {
	return b.Attr("lang", v)
}

func (b *VXTipTapEditorBuilder) Output(v string) *VXTipTapEditorBuilder {
	return b.Attr("output", v)
}

func (b *VXTipTapEditorBuilder) Dense(v bool) *VXTipTapEditorBuilder {
	return b.Attr("dense", v)
}

func (b *VXTipTapEditorBuilder) HideToolbar(v bool) *VXTipTapEditorBuilder {
	return b.Attr("hide-toolbar", v)
}

func (b *VXTipTapEditorBuilder) DisableToolbar(v bool) *VXTipTapEditorBuilder {
	return b.Attr("disable-toolbar", v)
}

func (b *VXTipTapEditorBuilder) Readonly(v bool) *VXTipTapEditorBuilder {
	return b.Attr("readonly", v)
}

func (b *VXTipTapEditorBuilder) ReadonlyClass(v string) *VXTipTapEditorBuilder {
	return b.Attr("readonly-class", v)
}

func (b *VXTipTapEditorBuilder) ErrorMessages(v []string) *VXTipTapEditorBuilder {
	return b.Attr(":error-messages", v)
}

func (b *VXTipTapEditorBuilder) Headings(v []int) *VXTipTapEditorBuilder {
	return b.Attr("headings", v)
}

func (b *VXTipTapEditorBuilder) MaxWidth(v int) *VXTipTapEditorBuilder {
	return b.Attr("max-width", v)
}

func (b *VXTipTapEditorBuilder) MaxHeight(v int) *VXTipTapEditorBuilder {
	return b.Attr("max-height", v)
}

func (b *VXTipTapEditorBuilder) Translator(v any) *VXTipTapEditorBuilder {
	return b.Attr("translator", v)
}

func (b *VXTipTapEditorBuilder) Template(v bool) *VXTipTapEditorBuilder {
	return b.Attr("template", v)
}

func (b *VXTipTapEditorBuilder) PreviewSlot(comp h.HTMLComponent) *VXTipTapEditorBuilder {
	return b.AppendChild(web.Slot(comp).Name("preview").Scope("{editor}"))
}

func (b *VXTipTapEditorBuilder) HelpSlot(comp h.HTMLComponent) *VXTipTapEditorBuilder {
	return b.AppendChild(web.Slot(comp).Name("help").Scope("{editor}"))
}
