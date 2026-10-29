// mdkit-markdown-editor: a text in Markdown, a Vue component —
// `<mdkit-markdown-editor v-model="body" preview="/api/issues/preview">` —,
// EasyMDE: its box the Markdown itself (a template's hints, <!-- … -->, seen
// there, as GitHub's), a bar of buttons — bold, italic, strikethrough,
// headings of every level, lists, quote, code, link, table (built as a
// sheet: mdkit-table-builder), footnote (mdkit-footnote-builder; a
// reference in the box marked, Ctrl+click opens its note), files of
// Markdown included (their footnotes made their own first), rule, undo, redo
// (its history of changes), full screen —, its preview — a menu of modes:
// the text alone, side by side (split), the preview alone (total) — rendered
// by the server (`preview`: the address that renders it —
// the same rendering and sanitizing as where it is shown —; none, no
// preview), and the help (?): its shortcuts, tips and the guide of Markdown
// in the page's language (words.helpUrl).
// Its pictures — chosen (its button, or a cell of a table), pasted, dropped
// —, `![name](address)` in the text, as the rich text editor's: with `upload`
// (an issue there is) each sent at once, its address; else, with `images`,
// each given to the page by the event `file` (a File named `inline-…`), the
// text naming it `attachment:<name>` (the server puts its address there;
// the preview shows it from the browser); without either, none.
// Its words: `words` (the page's markdown words) and `barWords` (the rich
// text editor's: bold, italic, …, notImage, uploadError).
import EasyMDE from 'easymde'
import css from './vendor/easymde-css.js'
import TableBuilder from './table-builder.js'
import FootnoteBuilder from './footnote-builder.js'
import { tableAt } from './table-model.js'
import * as F from './footnote-model.js'
import { config, wordsOf } from './config.js'

// cursorMark is the place of the cursor in the copy of the text the
// footnotes builder works on (a character of no use, taken out at the end).
const cursorMark = '\uE000'
// refMark is the place of the reference clicked (Ctrl+click) in that copy:
// the builder may point it at another note.
const refMark = '\uE001'
// mac: the shortcuts are of Cmd (⌘), not Ctrl
const mac = typeof navigator !== 'undefined' && /Mac/.test(navigator.platform)

// shortcuts are the keys of the editor (EasyMDE's, CodeMirror's), each with
// the key of its words and its word by default: what the help shows.
const shortcuts = [
  ['Mod-B', 'bold', 'Bold'], ['Mod-I', 'italic', 'Italic'], ['Mod-K', 'link', 'Link'],
  ['Mod-H', 'headingSmaller', 'Smaller heading'], ['Shift-Mod-H', 'headingBigger', 'Bigger heading'],
  ['Ctrl-Alt-1 … 6', 'headingLevels', 'Heading 1 … 6'], ['Mod-E', 'cleanBlock', 'Clear the block'],
  ["Mod-'", 'quote', 'Quote'], ['Mod-L', 'bullets', 'List'], ['Mod-Alt-L', 'numbers', 'Numbered list'],
  ['Shift-Mod-L', 'taskList', 'Task list'], ['Mod-Alt-C', 'code', 'Code'], ['Mod-Alt-I', 'image', 'Picture', 'image'],
  ['Mod-Z', 'undo', 'Undo'], ['Mod-Y · Shift-Mod-Z', 'redo', 'Redo'],
  ['Mod-P', 'preview', 'Preview'], ['F9', 'sideBySide', 'Side by side'], ['F11', 'fullscreen', 'Full screen'],
  ['Esc', 'exitFullscreen', 'Leave full screen'],
]
// keyText is a shortcut as the user's keyboard says it
const keyText = (k) => k.replace(/Mod/g, mac ? '⌘' : 'Ctrl').replace(/Alt/g, mac ? '⌥' : 'Alt').replace(/-/g, '+')

const unmark = (s) => s.split(cursorMark).join('').split(refMark).join('')

const imageTypes = /^image\/(png|jpe?g|gif|webp)$/

// inlineURL is the address the preview's server is told for a picture not
// sent yet (attachment:<name>): one it keeps, the browser's put back after.
const inlineURL = (name) => 'https://inline.invalid/' + name

// style is the editor's look, once for the page: EasyMDE's, its buttons
// the page's icons (mdi), its preview as where the text is shown.
function style() {
  if (document.getElementById('mdkit-markdown-editor-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-markdown-editor-style'
  st.textContent = css +
    '.mdkit-markdown-editor .editor-toolbar{background:#f5f5f5;border-color:#d1d5db;opacity:1}' +
    '.mdkit-markdown-editor .editor-toolbar button{font-size:1.15rem;width:32px;height:32px}' +
    '.mdkit-markdown-editor .editor-toolbar .easymde-dropdown-content button{width:auto;font-size:.9rem;white-space:nowrap;padding:0 .5rem;display:flex;align-items:center;gap:.5rem}' +
    // a menu over the text (EasyMDE's z-index is the lines' own: they hid it)
    '.mdkit-markdown-editor .easymde-dropdown-content{z-index:20;padding:4px;background:#fff;border-radius:4px}' +
    // the items of a menu (headings, the preview's modes): their names too
    '.mdkit-markdown-editor .editor-toolbar .easymde-dropdown-content button::after{content:attr(title);font-family:inherit;font-size:.85rem}' +
    '.mdkit-markdown-editor .CodeMirror{border-color:#d1d5db}' +
    // a reference of a footnote in the box
    '.mdkit-markdown-editor .cm-footnote{color:#6d28d9;background:#f5f3ff;border-radius:3px}' +
    // the headings in the box: a little larger, as the text shown has them
    // (EasyMDE's are a page's titles)
    '.mdkit-markdown-editor .cm-s-easymde .cm-header{line-height:1.4;margin:0}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-1{font-size:1.4em}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-2{font-size:1.25em}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-3{font-size:1.15em}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-4{font-size:1.05em}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-5{font-size:1em}' +
    '.mdkit-markdown-editor .cm-s-easymde .cm-header-6{font-size:.95em}' +
    '.mdkit-markdown-editor.invalid .CodeMirror,.mdkit-markdown-editor.invalid .editor-toolbar{border-color:#b00020}' +
    '.mdkit-markdown-editor.disabled{opacity:.6;pointer-events:none}' +
    '.mdkit-markdown-editor .editor-preview,.mdkit-markdown-editor .editor-preview-side{background:#fff}'
  document.head.append(st)
}

// normal makes the lines of the selection text again: no heading
function normal(editor) {
  const cm = editor.codemirror
  const from = cm.getCursor('from').line
  const to = cm.getCursor('to').line
  for (let i = from; i <= to; i++) {
    const m = /^#{1,6} +/.exec(cm.getLine(i))
    if (m) cm.replaceRange('', { line: i, ch: 0 }, { line: i, ch: m[0].length })
  }
  cm.focus()
}

// button is a button of the bar: its action, its icon (mdi), its title
const button = (name, action, icon, title) => ({ name, action, className: 'mdi ' + icon, title })

// item is an item of a menu of the bar: done, the focus back in the box —
// the menu (open while it has the focus) closed
const item = (name, action, icon, title) => button(name, (editor) => {
  action(editor)
  editor.codemirror.focus()
}, icon, title)

export default {
  name: 'MdkitMarkdownEditor',
  props: {
    modelValue: { type: String, default: '' },
    // preview: the address that renders a text (POST {format, body} →
    // {html}), or a function of the text giving its HTML (a promise)
    preview: { type: [String, Function], default: '' },
    rows: { type: [Number, String], default: 8 },
    placeholder: { type: String, default: '' },
    maxlength: { type: [Number, String], default: null },
    disabled: { type: Boolean, default: false },
    invalid: { type: Boolean, default: false },
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('markdown') },
    // the words of the buttons it shares with the rich text editor
    barWords: { type: Object, default: () => wordsOf('rich') },
    inputId: { type: String, default: '' },
    path: { type: String, default: '' },
    // its pictures: sent at once (an address, or a function of the FormData
    // that sends it: its response, or what it says — { url } or { message }
    // —), `token` the secret link's; or given to the page
    upload: { type: [String, Function], default: '' },
    token: { type: String, default: '' },
    images: { type: Boolean, default: false },
  },
  emits: ['update:modelValue', 'input', 'file'],
  data: () => ({ notice: '', helping: false }),
  computed: {
    canImage() { return !!this.upload || this.images },
    // the shortcuts of the help, in the user's words and keyboard (the
    // picture's only where pictures go in)
    shortcutList() {
      const w = this.words
      const b = this.barWords
      return shortcuts.filter((s) => s[3] !== 'image' || this.canImage).map(([keys, key, word]) => ({
        key, keys: keyText(keys), label: w[key] || b[key] || word,
      }))
    },
    // the tips of the help
    tips() {
      const w = this.words
      return [
        ...(this.canImage ? [w.tipPictures || 'Paste or drop pictures to put them in the text.'] : []),
        w.tipFootnote || 'Ctrl+click (⌘+click on a Mac) on a reference [^…] opens its note.',
        w.tipLists || 'In a list, Enter goes on with it; Tab and Shift+Tab indent it.',
        w.tipTables || 'A table under the cursor (or selected) opens in the table builder.',
      ]
    },
  },
  components: { MdkitTableBuilder: TableBuilder, MdkitFootnoteBuilder: FootnoteBuilder },
  template: `
<div class="mdkit-markdown-editor" :class="{ invalid, disabled }" data-markdown-editor>
  <textarea ref="area" :id="inputId || undefined" :data-path="path || undefined"></textarea>
  <mdkit-table-builder ref="tables" :image="canImage ? imageText : null" :render="preview ? renderHTML : null"
    :footnote="footnotePick"></mdkit-table-builder>
  <mdkit-footnote-builder ref="footnotes" :render="preview ? renderHTML : null"></mdkit-footnote-builder>
  <v-dialog v-model="helping" max-width="620" scrollable data-markdown-help>
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon icon="mdi-keyboard-outline" class="me-2"></v-icon>{{ words.shortcuts || 'Keyboard shortcuts' }}
      </v-card-title>
      <v-card-text>
        <v-table density="compact">
          <tbody>
            <tr v-for="s in shortcutList" :key="s.key" :data-markdown-shortcut="s.key">
              <td>{{ s.label }}</td>
              <td class="text-end"><kbd class="text-caption">{{ s.keys }}</kbd></td>
            </tr>
          </tbody>
        </v-table>
        <ul class="text-body-2 mt-4 ps-4" style="list-style:disc" data-markdown-tips>
          <li v-for="tip in tips" :key="tip">{{ tip }}</li>
        </ul>
      </v-card-text>
      <v-card-actions>
        <v-btn v-if="words.helpUrl" :href="words.helpUrl" target="_blank" rel="noopener" color="primary" variant="tonal"
          prepend-icon="mdi-book-open-variant" data-markdown-guide>{{ words.guide || 'Markdown guide' }}</v-btn>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="helping = false">{{ words.close || 'Close' }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <div v-if="notice" class="text-caption text-error mt-1" data-markdown-notice>{{ notice }}</div>
  <div v-if="words.hint" class="text-caption text-medium-emphasis mt-1">{{ words.hint }}</div>
</div>`,
  mounted() {
    style()
    // the pictures not sent yet: their names, the browser's addresses
    this.pending = new Map()
    this.seq = 0
    const w = this.words
    const b = this.barWords
    const t = (key, fallback) => w[key] || b[key] || fallback
    const level = (n) => (w.headingLevel || 'Heading {n}').replace('{n}', n)
    const toolbar = [
      button('bold', EasyMDE.toggleBold, 'mdi-format-bold', t('bold', 'Bold')),
      button('italic', EasyMDE.toggleItalic, 'mdi-format-italic', t('italic', 'Italic')),
      button('strikethrough', EasyMDE.toggleStrikethrough, 'mdi-format-strikethrough', t('strike', 'Strikethrough')),
      {
        name: 'heading', className: 'mdi mdi-format-header-pound', title: t('heading', 'Heading'),
        children: [
          ...[1, 2, 3, 4, 5, 6].map((n) => item('h' + n, EasyMDE['toggleHeading' + n], 'mdi-format-header-' + n, level(n))),
          item('h0', normal, 'mdi-format-text', t('paragraph', 'Normal text')),
        ],
      },
      '|',
      button('unordered-list', EasyMDE.toggleUnorderedList, 'mdi-format-list-bulleted', t('bullets', 'List')),
      button('ordered-list', EasyMDE.toggleOrderedList, 'mdi-format-list-numbered', t('numbers', 'Numbered list')),
      button('quote', EasyMDE.toggleBlockquote, 'mdi-format-quote-close', t('quote', 'Quote')),
      button('code', EasyMDE.toggleCodeBlock, 'mdi-code-braces', t('code', 'Code')),
      '|',
      button('link', EasyMDE.drawLink, 'mdi-link-variant', t('link', 'Link')),
      ...(this.canImage ? [button('image', () => this.chooseImages(), 'mdi-image-plus', t('image', 'Picture'))] : []),
      // the table under the cursor, or a new one: built as a sheet
      button('table', (editor) => this.openTable(editor), 'mdi-table', t('table', 'Table')),
      // the footnotes: built in a dialog
      button('footnote', (editor) => this.openFootnotes(editor), 'mdi-format-superscript', t('footnote', 'Footnote')),
      // files of Markdown put at the cursor, their footnotes their own
      button('include', (editor) => this.includeFiles(editor), 'mdi-file-document-plus-outline', t('include', 'Include Markdown files')),
      button('horizontal-rule', EasyMDE.drawHorizontalRule, 'mdi-minus', t('rule', 'Horizontal line')),
      '|',
      // the history of changes
      button('undo', EasyMDE.undo, 'mdi-undo', t('undo', 'Undo')),
      button('redo', EasyMDE.redo, 'mdi-redo', t('redo', 'Redo')),
    ]
    if (this.preview) {
      // the preview's modes: the text alone, side by side (split), the
      // preview alone (total)
      const mode = (m) => (editor) => {
        if (editor.isSideBySideActive() !== (m === 'split')) EasyMDE.toggleSideBySide(editor)
        if (editor.isPreviewActive() !== (m === 'total')) EasyMDE.togglePreview(editor)
      }
      toolbar.push('|', {
        name: 'preview-modes', className: 'mdi mdi-eye-outline', title: t('preview', 'Preview'), noDisable: true,
        children: [
          { ...item('mode-write', mode('write'), 'mdi-pencil-outline', t('write', 'Write')), noDisable: true },
          { ...item('mode-split', mode('split'), 'mdi-view-split-vertical', t('sideBySide', 'Side by side')), noDisable: true },
          { ...item('mode-total', mode('total'), 'mdi-eye-outline', t('previewOnly', 'Preview only')), noDisable: true },
        ],
      })
    }
    toolbar.push({ ...button('fullscreen', EasyMDE.toggleFullScreen, 'mdi-fullscreen', t('fullscreen', 'Full screen')), noDisable: true, noMobile: true })
    // the help: the shortcuts, the tips, the guide of Markdown
    toolbar.push('|', {
      ...button('guide', () => { this.helping = true }, 'mdi-help-circle-outline', t('help', 'Markdown help')),
      noDisable: true,
    })
    const max = Number(this.maxlength) || 0
    this.mde = new EasyMDE({
      element: this.$refs.area,
      initialValue: this.modelValue || '',
      placeholder: this.placeholder,
      toolbar,
      // the picture of EasyMDE (an address to write): ours, below
      shortcuts: { drawImage: null },
      autoDownloadFontAwesome: false,
      spellChecker: false,
      nativeSpellcheck: true,
      forceSync: true,
      minHeight: `${Math.max(Number(this.rows) || 8, 3) * 1.5}em`,
      unorderedListStyle: '-',
      // side by side in its place (full screen: its own button)
      sideBySideFullscreen: false,
      previewClass: ['editor-preview', config.textClass],
      // the server's rendering, when it comes (null: what is there stays)
      previewRender: (text, el) => {
        this.render(text, el)
        return null
      },
      // its size, against the most it may have
      status: max ? [{
        className: 'chars',
        defaultValue: (el) => { el.textContent = `${(this.modelValue || '').length} / ${max}` },
        onUpdate: (el) => { el.textContent = `${this.mde ? this.mde.value().length : 0} / ${max}` },
      }] : false,
    })
    // what was set at first is not a change to undo
    this.mde.codemirror.clearHistory()
    this.mde.codemirror.on('change', () => {
      const v = this.mde.value()
      if (v === this.modelValue) return
      this.$emit('update:modelValue', v)
      this.$emit('input', v)
    })
    if (this.disabled) this.mde.codemirror.setOption('readOnly', true)
    // a picture chosen (its shortcut), where pictures go in
    if (this.canImage) this.mde.codemirror.addKeyMap({ [mac ? 'Cmd-Alt-I' : 'Ctrl-Alt-I']: () => this.chooseImages() })
    // pictures pasted or dropped
    const take = (cm, e, files) => {
      const list = Array.from(files || []).filter((f) => imageTypes.test(f.type))
      if (!list.length || !this.canImage) return
      e.preventDefault()
      if (e.type === 'drop') {
        const at = cm.coordsChar({ left: e.clientX, top: e.clientY }, 'window')
        cm.setCursor(at)
      }
      this.insertImages(list)
    }
    // the references of footnotes marked; Ctrl+click opens the note
    this.mde.codemirror.addOverlay({
      token(stream) {
        if (stream.match(/\[\^[^\]\s]+\]/)) return 'footnote'
        while (stream.next() != null && !stream.match(/\[\^/, false));
        return null
      },
    })
    this.mde.codemirror.on('mousedown', (cm, e) => {
      if (!(e.ctrlKey || e.metaKey)) return
      const at = cm.coordsChar({ left: e.clientX, top: e.clientY }, 'window')
      const line = cm.getLine(at.line)
      for (const m of line.matchAll(/\[\^([^\]\s]+)\]/g)) {
        if (at.ch >= m.index && at.ch <= m.index + m[0].length) {
          e.preventDefault()
          this.openFootnotes(this.mde, m[1], cm.indexFromPos({ line: at.line, ch: m.index }))
          return
        }
      }
    })
    this.mde.codemirror.on('paste', (cm, e) => take(cm, e, e.clipboardData && e.clipboardData.files))
    this.mde.codemirror.on('drop', (cm, e) => take(cm, e, e.dataTransfer && e.dataTransfer.files))
  },
  beforeUnmount() {
    clearTimeout(this.timer)
    if (this.mde) {
      this.mde.toTextArea()
      this.mde = null
    }
  },
  watch: {
    // changed from outside (the form sent, made empty): the box again — a
    // change of its history —, the preview closed when empty
    modelValue(v) {
      if (!this.mde || (v || '') === this.mde.value()) return
      this.mde.value(v || '')
      if (!v && this.mde.isPreviewActive()) EasyMDE.togglePreview(this.mde)
    },
    disabled(v) {
      if (this.mde) this.mde.codemirror.setOption('readOnly', v)
    },
  },
  methods: {
    // openTable opens the table builder: with the table selected, or under
    // the cursor — of GFM or HTML, made the sheet; then put back in its
    // place —, or a new one put at the cursor, a line
    // of its own (a change of the history: undone as one)
    openTable(editor) {
      const cm = editor.codemirror
      const lines = cm.getValue().split('\n')
      if (this.tableMark) this.tableMark.clear()
      // the table of the selection (any of its lines), else of the cursor's
      let found = null
      if (cm.somethingSelected()) {
        const from = cm.getCursor('from').line
        const to = cm.getCursor('to').line
        for (let l = from; l <= to && !found; l++) found = tableAt(lines, l)
      } else {
        found = tableAt(lines, cm.getCursor().line)
      }
      // the table followed as the text changes (its notes made meanwhile)
      const mark = found && cm.markText({ line: found.from, ch: 0 }, { line: found.to, ch: lines[found.to].length },
        { inclusiveLeft: true, inclusiveRight: true })
      this.tableMark = mark
      this.$refs.tables.open(found && found.table, (text) => {
        const at = mark && mark.find()
        if (mark) mark.clear()
        this.tableMark = null
        if (at) {
          cm.replaceRange(text, at.from, at.to)
        } else {
          const cur = cm.getCursor()
          const before = cm.getLine(cur.line).slice(0, cur.ch).trim() ? '\n\n' : (cur.line > 0 && cm.getLine(cur.line - 1).trim() ? '\n' : '')
          cm.replaceSelection(before + text + '\n\n')
        }
        cm.focus()
      })
    },
    // apply makes the text of cm out, changing only its part that changed (a
    // change of its history; the marks elsewhere kept: a table being
    // built); the cursor at cursor (an index), when one
    apply(cm, out, cursor) {
      const cur = cm.getValue()
      if (out !== cur) {
        let a = 0
        while (a < cur.length && a < out.length && cur[a] === out[a]) a++
        let b = 0
        while (b < cur.length - a && b < out.length - a && cur[cur.length - 1 - b] === out[out.length - 1 - b]) b++
        cm.replaceRange(out.slice(a, out.length - b), cm.posFromIndex(a), cm.posFromIndex(cur.length - b))
      }
      if (cursor >= 0) cm.setCursor(cm.posFromIndex(cursor))
      cm.focus()
    },
    // footnoteAdapter is the footnotes builder's way to the text of editor:
    // a copy (the cursor marked in it), changed by footnote-model; done, the
    // editor changed once, organized (the notes at the end, in order,
    // renumbered) — unless pick (a cell of a table: the text kept where it
    // is) —
    footnoteAdapter(editor, pick = false, refAt = null) {
      const cm = editor.codemirror
      const at = cm.indexFromPos(cm.getCursor())
      const before = cm.getValue()
      // the marks of the cursor and of the reference clicked, from the last
      const marks = [[at, cursorMark], ...(refAt != null ? [[refAt, refMark]] : [])].sort((a, b) => b[0] - a[0])
      let text = before
      for (const [i, m] of marks) text = text.slice(0, i) + m + text.slice(i)
      const refRe = new RegExp(refMark + '\\[\\^([^\\]\\s]+)\\]')
      const fresh = (label) => {
        const l = F.cleanLabel(label)
        return l && !F.parse(text).notes.some((n) => n.label === l) ? l : F.newLabel(text)
      }
      return {
        kind: 'markdown',
        notes: () => F.parse(text).notes.map((n) => ({ ...n, text: unmark(n.text) })),
        // the note of the reference clicked; point makes it another's
        pointed: () => { const m = refRe.exec(text); return m ? m[1] : null },
        point: (label) => { text = text.replace(refRe, refMark + '[^' + label + ']') },
        next: () => F.newLabel(text),
        create: (body, label) => {
          const l = fresh(label)
          text = F.edit(text, l, body)
          return l
        },
        insert: (label) => {
          const i = text.indexOf(cursorMark)
          text = text.slice(0, i) + '[^' + label + ']' + text.slice(i)
        },
        edit: (label, body) => { text = F.edit(text, label, body) },
        rename: (from, to) => {
          const r = F.rename(text, from, to)
          if (r === false) return false
          text = r
          return true
        },
        remove: (label) => { text = F.remove(text, label) },
        done: () => {
          const out = (pick ? text : F.organize(text, true)).split(refMark).join('')
          const cursor = pick ? -1 : out.indexOf(cursorMark)
          this.apply(cm, unmark(out), cursor)
        },
      }
    },
    // openFootnotes opens the footnotes builder on the text (label: the note
    // chosen)
    openFootnotes(editor, label, refAt = null) {
      this.$refs.footnotes.open(this.footnoteAdapter(editor, false, refAt), { label, ref: refAt != null })
    },
    // footnotePick is the label of a note chosen or made (a cell of a table:
    // the table builder's)
    footnotePick() {
      return new Promise((resolve) => {
        this.$refs.footnotes.open(this.footnoteAdapter(this.mde, true), { pick: true, picked: resolve })
      })
    },
    // pickFiles is the files of Markdown (or text) the user chooses
    pickFiles() {
      return new Promise((resolve) => {
        const input = document.createElement('input')
        input.type = 'file'
        input.accept = '.md,.markdown,.mdown,.txt,text/markdown,text/plain'
        input.multiple = true
        input.onchange = () => resolve(Array.from(input.files || []))
        input.click()
      })
    },
    // includeFiles puts the files chosen at the cursor, one after the other
    // — the labels of each file's footnotes made its own first (its name
    // their prefix: none the text's, nor another file's) —, the notes then
    // organized; a change of the history
    async includeFiles(editor) {
      const files = await this.pickFiles()
      if (!files.length) return
      const cm = editor.codemirror
      const at = cm.indexFromPos(cm.getCursor('from'))
      const end = cm.indexFromPos(cm.getCursor('to'))
      const cur = cm.getValue()
      let all = cur
      const parts = []
      for (const f of files) {
        const body = (await f.text()).replace(/\r\n?/g, '\n').trim()
        if (!body) continue
        const map = F.uniqueLabels(body, all + '\n' + parts.join('\n'), f.name.replace(/\.[^.]*$/, ''))
        parts.push(F.relabel(body, map))
      }
      if (!parts.length) return
      const head = cur.slice(0, at)
      const sep = head.trim() ? (head.endsWith('\n\n') ? '' : head.endsWith('\n') ? '\n' : '\n\n') : ''
      const text = head + sep + parts.join('\n\n') + cursorMark + '\n\n' + cur.slice(end).replace(/^\n+/, '')
      const out = F.organize(text, true)
      this.apply(cm, out.split(cursorMark).join(''), out.indexOf(cursorMark))
    },
    // render puts the text rendered by the server in el, a moment after the
    // last change (side by side, as it is written); the pictures not sent
    // yet, the browser's
    render(text, el) {
      clearTimeout(this.timer)
      this.timer = setTimeout(async () => { el.innerHTML = await this.renderHTML(text) }, 200)
    },
    // renderHTML is text rendered by the server (preview), the pictures not
    // sent yet the browser's (the table builder's preview too)
    async renderHTML(text) {
      let body = text || ''
      this.pending.forEach((url, name) => { body = body.split('attachment:' + name).join(inlineURL(name)) })
      let html = ''
      if (typeof this.preview === 'function') {
        html = (await Promise.resolve(this.preview(body)).catch(() => '')) || ''
      } else {
        const res = await fetch(this.preview, {
          method: 'POST', credentials: 'same-origin', headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ format: 'markdown', body }),
        }).catch(() => null)
        const out = res ? await res.json().catch(() => ({})) : {}
        html = out.html || ''
      }
      this.pending.forEach((url, name) => { html = html.split(inlineURL(name)).join(url) })
      return html
    },
    // pick is the pictures the user chooses
    pick() {
      return new Promise((resolve) => {
        const input = document.createElement('input')
        input.type = 'file'
        input.accept = 'image/png,image/jpeg,image/gif,image/webp'
        input.multiple = true
        input.onchange = () => resolve(Array.from(input.files || []))
        input.click()
      })
    },
    async chooseImages() {
      this.insertImages(await this.pick())
    },
    // insertImages puts the pictures at the cursor, each on a line of its own
    async insertImages(files) {
      const text = await this.imagesText(files)
      if (!text) return
      const cm = this.mde.codemirror
      const cur = cm.getCursor()
      const before = cm.getLine(cur.line).slice(0, cur.ch).trim() ? '\n' : ''
      cm.replaceSelection(before + text + '\n')
      cm.focus()
    },
    // imageText is the Markdown of pictures the user chooses (a cell of a
    // table: the table builder's)
    async imageText() {
      return this.imagesText(await this.pick(), ' ')
    },
    // imagesText is the Markdown of the pictures — each sent, or named as an
    // attachment of what is being sent —, joined by sep
    async imagesText(files, sep = '\n') {
      const out = []
      for (const f of files) {
        const md = await this.imageOf(f)
        if (md) out.push(md)
      }
      return out.join(sep)
    },
    async imageOf(file) {
      const b = this.barWords
      if (!imageTypes.test(file.type)) {
        this.notice = b.notImage || 'Not a picture.'
        return ''
      }
      this.notice = ''
      const alt = file.name.replace(/\.[a-z0-9]+$/i, '').replace(/[[\]]/g, '')
      if (this.upload) {
        const body = new FormData()
        body.append('file', file, file.name)
        if (this.token) body.append('token', this.token)
        try {
          const res = typeof this.upload === 'function'
            ? await this.upload(body)
            : await fetch(this.upload, { method: 'POST', credentials: 'same-origin', body })
          // a response (fetch's), or what it says
          const sent = res && typeof res.json === 'function'
          const out = sent ? await res.json().catch(() => ({})) : (res || {})
          if ((sent && !res.ok) || !out.url) {
            this.notice = out.message || b.uploadError || 'The picture could not be sent.'
            return ''
          }
          return `![${alt}](${out.url})`
        } catch (e) {
          this.notice = b.uploadError || 'The picture could not be sent.'
          return ''
        }
      }
      const ext = (file.name.match(/\.[a-z0-9]+$/i) || [''])[0].toLowerCase()
      const name = 'inline-' + (++this.seq) + '-' + Date.now().toString(36) + ext
      const named = new File([file], name, { type: file.type })
      this.pending.set(name, URL.createObjectURL(named))
      this.$emit('file', named)
      return `![${alt}](attachment:${name})`
    },
  },
}
