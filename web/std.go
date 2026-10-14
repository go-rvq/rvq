package web

import "strings"

// StdContentClass default class for standard content like Post Body, Page Content etc.
const StdContentClass = "std-content"

// StdContentEditorSelector is the content of the rich text editor
// (vx-tiptap-editor) as it is edited: it looks as it is read. The editor has
// no theme of its own (vuetify-pro-tiptap's markdown theme "std", which it
// has not: no colors, no marks of its).
const StdContentEditorSelector = ".vx-tiptap-editor .vuetify-pro-tiptap-editor__content.markdown-theme-std .tiptap.ProseMirror"

// StdContentCSS is the look of a text written — the user docs', a body of a
// page or a post, read or edited — in the elements of selectors: its
// paragraphs and lists spaced, the lists numbered or bulleted (the reset of
// the admin's styles takes their marks), its headings, pictures, quotes,
// tables and inline code.
func StdContentCSS(selectors ...string) string {
	s := ":is(" + strings.Join(selectors, ", ") + ")"
	return strings.NewReplacer("S ", s+" ").Replace(`
S img { max-width: 100%; height: auto; border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity)); border-radius: 6px; margin: 8px 0; }
S h1 { font-size: 1.6rem; margin: 0 0 8px; }
S h2 { font-size: 1.25rem; margin: 20px 0 8px; }
S h3 { font-size: 1.1rem; margin: 16px 0 6px; }
S p, S ul, S ol, S blockquote, S table { margin-top: 0; margin-bottom: 10px; }
S ul, S ol { padding-left: 24px; }
S ul { list-style: disc; }
S ol { list-style: decimal; }
S ul ul { list-style: circle; }
S li > p { margin-bottom: 4px; }
S blockquote { border-left: 3px solid rgba(var(--v-border-color), var(--v-border-opacity)); padding-left: 12px; color: rgba(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity)); }
/* inline code; a block of code has a look of its own */
S :not(pre) > code { background: rgba(var(--v-theme-on-surface), 0.06); padding: 1px 4px; border-radius: 4px; }
S table { border-collapse: collapse; }
S th, S td { border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity)); padding: 4px 8px; }
`)
}

// StdContentStyles are the styles of the standard content: the class
// StdContentClass (a body read) and the editor's content
// (StdContentEditorSelector), as the user docs (StdContentCSS).
func StdContentStyles() ComponentsPack {
	return ComponentsPack(StdContentCSS("."+StdContentClass, StdContentEditorSelector))
}
