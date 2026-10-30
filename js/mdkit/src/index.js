// mdkit: the editors of a text — Markdown (EasyMDE: its bar, preview modes,
// pictures, files included, help), rich (tiptap) — and their builders —
// tables as a sheet, footnotes (their map, relabeled for texts put
// together) —, Vue components made with Vuetify (the page's), their words
// the page's (configure, or window.mdkitWords). For a site (its bundle:
// window.RvqMdKit, the pack of rvq x/ui/mdkit) and for the admin (vuetifyx:
// the plugins of vx-tiptap-editor, the MarkdownBuilder); the link cards of
// a text drawn (watchLinkCards).
import RichEditor from './rich-editor.js'
import MarkdownEditor from './markdown-editor.js'
import TableBuilder from './table-builder.js'
import FootnoteBuilder from './footnote-builder.js'
import LinkCardPicker from './link-card-picker.js'
import LinkPicker, { LinkForm } from './link-picker.js'
import { config, configure, wordsOf } from './config.js'
import { textStyle } from './text-style.js'
import { footnoteBalloons } from './balloons.js'

export { RichEditor, MarkdownEditor, TableBuilder, FootnoteBuilder, config, configure, wordsOf, textStyle, footnoteBalloons }
export * as tables from './table-model.js'
export * as footnotes from './footnote-model.js'
export { FootnoteRef, Footnotes, numbers as footnoteNumbers } from './footnote-tiptap.js'
export { finalize as finalizeFootnotes, noteHTML } from './footnote-rich.js'
export { tiptapFootnoteAdapter } from './footnote-adapter.js'
// the link cards of a text (a page, a record linked on a line of its own)
export { linkCards, watchLinkCards, metaResolver, drawCard, safeHref, cardClass, cardTarget } from './link-cards.js'
export { LinkCard } from './link-card-tiptap.js'
// the formulas and diagrams of a text (KaTeX, Mermaid)
export { watchBlocks, drawBlocks, vendorOf, load as loadVendor } from './blocks.js'
export { Diagram, drawInto as drawDiagramInto } from './diagram-tiptap.js'
export { InlineMath, BlockMath, MathCommands, drawMath } from './math-tiptap.js'
export { default as BlockDialog } from './block-dialog.js'
export { default as LinkCardPicker, metaSearch } from './link-card-picker.js'
export { default as LinkPicker, LinkForm, hrefOf } from './link-picker.js'

// components are the components, by their names after the prefix
export const components = {
  'rich-editor': RichEditor,
  'markdown-editor': MarkdownEditor,
  'table-builder': TableBuilder,
  'footnote-builder': FootnoteBuilder,
  'link-card-picker': LinkCardPicker,
  'link-picker': LinkPicker,
  'link-form': LinkForm,
}

// install gives app the components, named prefix-… (mdkit-rich-editor, …)
export function install(app, { prefix = 'mdkit' } = {}) {
  Object.keys(components).forEach((name) => app.component(prefix + '-' + name, components[name]))
}
export { inlineHTML } from './rich-cells.js'
