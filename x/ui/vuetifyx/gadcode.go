package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VXGadCodeBuilder renders <vx-gad-code>: the admin's editor of gad code (gad,
// gadx, gadt) — <vx-codemirror> with a button that formats the code (gad fmt),
// one that expands it over the whole page, the structure of the code as a tree
// at the left (a click moves the cursor to a declaration) and its
// documentation (gad doc) at the right, both closed to a strip that opens.
// What the code is asked of is served by the admin (rvq admin/gadcode, its
// API).
type VXGadCodeBuilder struct {
	v.VTagBuilder[*VXGadCodeBuilder]
}

func VXGadCode(children ...h.HTMLComponent) *VXGadCodeBuilder {
	return v.VTag(&VXGadCodeBuilder{}, "vx-gad-code", children...)
}

// Language is gad, gadx or gadt.
func (b *VXGadCodeBuilder) Language(v string) *VXGadCodeBuilder {
	return b.Attr("language", v)
}

// TemplateDelimiters are the delimiters of a gadt's code: `{%` / `%}` unless
// set.
func (b *VXGadCodeBuilder) TemplateDelimiters(start, end string) *VXGadCodeBuilder {
	return b.Attr("template-start", start).Attr("template-end", end)
}

// API is where the code is asked of: gadcode.API(the admin's prefix).
func (b *VXGadCodeBuilder) API(v string) *VXGadCodeBuilder {
	return b.Attr("api", v)
}

func (b *VXGadCodeBuilder) Readonly(v bool) *VXGadCodeBuilder {
	return b.Attr(":readonly", v)
}

func (b *VXGadCodeBuilder) MinHeight(v string) *VXGadCodeBuilder {
	return b.Attr("min-height", v)
}

func (b *VXGadCodeBuilder) MaxHeight(v string) *VXGadCodeBuilder {
	return b.Attr("max-height", v)
}

// Labels are the words of its buttons and panels (format, expand, collapse,
// structure, doc, refresh, noStructure, noDoc, notFormatted).
func (b *VXGadCodeBuilder) Labels(v map[string]string) *VXGadCodeBuilder {
	return b.Attr(":labels", h.JSONString(v))
}

// TitleSlot is what its toolbar says at the left (the field's label).
func (b *VXGadCodeBuilder) TitleSlot(children ...h.HTMLComponent) *VXGadCodeBuilder {
	return b.Children(web.Slot(children...).Name("title"))
}

// HelpSlot is the help of the editor, as VXCodeMirror's.
func (b *VXGadCodeBuilder) HelpSlot(children ...h.HTMLComponent) *VXGadCodeBuilder {
	return b.Children(web.Slot(children...).Name("help"))
}

// HelpMode is how the help renders: "dialog", "left" or "right".
func (b *VXGadCodeBuilder) HelpMode(v string) *VXGadCodeBuilder {
	return b.Attr("help-mode", v)
}
