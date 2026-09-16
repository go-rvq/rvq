package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VXCodeMirrorBuilder renders the <vx-codemirror> component: a CodeMirror-based
// code editor (or readonly viewer) whose language is chosen by name from the
// component's language registry (see CodeMirror.tsx). An unknown language edits
// as plain text with line numbers.
type VXCodeMirrorBuilder struct {
	v.VTagBuilder[*VXCodeMirrorBuilder]
}

func VXCodeMirror(children ...h.HTMLComponent) *VXCodeMirrorBuilder {
	return v.VTag(&VXCodeMirrorBuilder{}, "vx-codemirror", children...)
}

// Language sets the grammar by name (e.g. "yaml", "json", "html").
func (b *VXCodeMirrorBuilder) Language(v string) *VXCodeMirrorBuilder {
	return b.Attr("language", v)
}

// Readonly renders a non-editable viewer (used for the live preview).
func (b *VXCodeMirrorBuilder) Readonly(v bool) *VXCodeMirrorBuilder {
	return b.Attr(":readonly", v)
}

func (b *VXCodeMirrorBuilder) MinHeight(v string) *VXCodeMirrorBuilder {
	return b.Attr("min-height", v)
}

func (b *VXCodeMirrorBuilder) MaxHeight(v string) *VXCodeMirrorBuilder {
	return b.Attr("max-height", v)
}

// HelpSlot fills the "help" slot: its presence adds a "?" button in the top-right
// corner (with a tooltip) that opens the help (a dialog or a side panel).
func (b *VXCodeMirrorBuilder) HelpSlot(children ...h.HTMLComponent) *VXCodeMirrorBuilder {
	return b.Children(web.Slot(children...).Name("help"))
}

// HelpMode picks how the help renders: "dialog" (default, a modal), "left" (a
// panel on the left of the editor) or "right" (a panel on the right).
func (b *VXCodeMirrorBuilder) HelpMode(v string) *VXCodeMirrorBuilder {
	return b.Attr("help-mode", v)
}

// HelpModel binds the help show/hide state (v-model:help) to expr.
func (b *VXCodeMirrorBuilder) HelpModel(expr string) *VXCodeMirrorBuilder {
	return b.Attr("v-model:help", expr)
}

func (b *VXCodeMirrorBuilder) HelpTitle(v string) *VXCodeMirrorBuilder {
	return b.Attr("help-title", v)
}

func (b *VXCodeMirrorBuilder) HelpTooltip(v string) *VXCodeMirrorBuilder {
	return b.Attr("help-tooltip", v)
}
