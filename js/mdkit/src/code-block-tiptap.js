// The code blocks of the rich text editor (tiptap) edited by CodeMirror,
// with a choice of their language: withCodeMirror(CodeBlock) is the code
// block node — its HTML the same, pre > code.language-* — shown by a view of
// CodeMirror (loaded from where the page's meta mdkit-vendor says:
// "codemirror", window.RvqCodeMirror; meanwhile its text), its text and the
// document's kept the same both ways (ProseMirror's CodeMirror example),
// the arrows leaving it at its ends, Ctrl+Enter after it, undo and redo the
// editor's.
import { Selection, TextSelection } from '@tiptap/pm/state'
import { load } from './blocks.js'
import { wordsOf } from './config.js'

// CODE_LANGUAGES are the languages to choose, [name, label]: the name the
// one of the block's class (language-…), '' none (text)
export const CODE_LANGUAGES = [
  ['', 'Text'], ['bash', 'Bash'], ['css', 'CSS'], ['diff', 'Diff'], ['dockerfile', 'Dockerfile'],
  ['gad', 'Gad'], ['gadx', 'Gadx'], ['go', 'Go'], ['html', 'HTML'], ['ini', 'INI'],
  ['javascript', 'JavaScript'], ['json', 'JSON'], ['jsx', 'JSX'], ['markdown', 'Markdown'],
  ['python', 'Python'], ['rust', 'Rust'], ['sql', 'SQL'], ['toml', 'TOML'], ['typescript', 'TypeScript'],
  ['tsx', 'TSX'], ['xml', 'XML'], ['yaml', 'YAML'],
]

// style is the look of a code block being edited, once for the page.
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-code-edit-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-code-edit-style'
  st.textContent = '.mdkit-code-edit{margin:.75em 0;border:1px solid rgba(127,127,127,.3);border-radius:6px;overflow:hidden;background:rgba(127,127,127,.05)}' +
    '.mdkit-code-edit-bar{display:flex;align-items:center;gap:.5rem;padding:2px 6px;border-bottom:1px solid rgba(127,127,127,.2);font-size:.75rem}' +
    '.mdkit-code-edit-bar select{font:inherit;padding:1px 4px;border:1px solid rgba(127,127,127,.4);border-radius:4px;background:transparent;color:inherit}' +
    '.mdkit-code-edit-text{margin:0;padding:6px 10px;white-space:pre-wrap;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.8125rem}' +
    '.mdkit-code-edit .cm-editor{font-size:.8125rem}.mdkit-code-edit .cm-editor.cm-focused{outline:none}' +
    '.mdkit-code-edit .cm-gutters{background:transparent;border-right:1px solid rgba(127,127,127,.25)}' +
    '.mdkit-code-edit.ProseMirror-selectednode{outline:2px solid #3b82f6}'
  document.head.append(st)
}

// CodeBlockView is the view of a code block: CodeMirror and its language.
export class CodeBlockView {
  constructor({ node, view, getPos, editor }) {
    style()
    this.node = node
    this.view = view
    this.getPos = getPos
    this.editor = editor
    this.updating = false
    this.cm = null
    const w = { language: 'Language', ...wordsOf('code') }

    this.dom = document.createElement('div')
    this.dom.className = 'mdkit-code-edit'
    this.dom.setAttribute('data-code-edit', '')
    const bar = document.createElement('div')
    bar.className = 'mdkit-code-edit-bar'
    bar.contentEditable = 'false'
    this.select = document.createElement('select')
    this.select.title = w.language
    this.select.setAttribute('aria-label', w.language)
    this.select.setAttribute('data-code-language', '')
    const known = CODE_LANGUAGES.some(([n]) => n === (node.attrs.language || ''))
    for (const [name, label] of known ? CODE_LANGUAGES : [...CODE_LANGUAGES, [node.attrs.language, node.attrs.language]]) {
      const o = document.createElement('option')
      o.value = name
      o.textContent = name ? label : (w.text || label)
      this.select.append(o)
    }
    this.select.value = node.attrs.language || ''
    this.select.disabled = !editor.isEditable
    this.select.addEventListener('change', () => this.setLanguage(this.select.value))
    bar.append(this.select)
    this.host = document.createElement('div')
    // its text until CodeMirror is there
    const pre = document.createElement('pre')
    pre.className = 'mdkit-code-edit-text'
    pre.textContent = node.textContent
    this.host.append(pre)
    this.dom.append(bar, this.host)

    load('codemirror', 'RvqCodeMirror').then((CM) => this.mount(CM)).catch(() => {})
  }

  // mount puts CodeMirror in the place of the text
  mount(CM) {
    if (this.destroyed || this.cm) return
    this.CM = CM
    this.lang = new CM.Compartment()
    const editable = this.editor.isEditable
    this.cm = new CM.EditorView({
      doc: this.node.textContent,
      extensions: [
        CM.keymap.of([...this.keymap(), ...CM.defaultKeymap, CM.indentWithTab]),
        CM.lineNumbers(),
        CM.drawSelection(),
        CM.bracketMatching(),
        CM.syntaxHighlighting(CM.defaultHighlightStyle, { fallback: true }),
        CM.EditorView.editable.of(editable),
        CM.EditorState.readOnly.of(!editable),
        this.lang.of(this.languageOf(this.node.attrs.language)),
        CM.EditorView.updateListener.of((u) => this.forward(u)),
      ],
    })
    this.host.replaceChildren(this.cm.dom)
  }

  languageOf(name) {
    const f = this.CM && this.CM.languages[(name || '').toLowerCase()]
    return f ? f() : []
  }

  // setLanguage makes the block's language name
  setLanguage(name) {
    const pos = this.getPos()
    if (typeof pos !== 'number') return
    this.view.dispatch(this.view.state.tr.setNodeMarkup(pos, undefined, { ...this.node.attrs, language: name || null }))
  }

  // forward: what changed in CodeMirror — its text, its selection — in the
  // document
  forward(update) {
    if (this.updating || !this.cm.hasFocus) return
    let offset = this.getPos() + 1
    const { main } = update.state.selection
    const selFrom = offset + main.from
    const selTo = offset + main.to
    const pmSel = this.view.state.selection
    if (update.docChanged || pmSel.from !== selFrom || pmSel.to !== selTo) {
      const tr = this.view.state.tr
      update.changes.iterChanges((fromA, toA, fromB, toB, text) => {
        if (text.length) tr.replaceWith(offset + fromA, offset + toA, this.view.state.schema.text(text.toString()))
        else tr.delete(offset + fromA, offset + toA)
        offset += (toB - fromB) - (toA - fromA)
      })
      tr.setSelection(TextSelection.create(tr.doc, selFrom, selTo))
      this.view.dispatch(tr)
    }
  }

  // setSelection: the document's selection in the block, CodeMirror's
  setSelection(anchor, head) {
    if (!this.cm) return
    this.cm.focus()
    this.updating = true
    this.cm.dispatch({ selection: { anchor, head } })
    this.updating = false
  }

  keymap() {
    const run = (f) => () => f()
    return [
      { key: 'ArrowUp', run: this.escape('line', -1) },
      { key: 'ArrowLeft', run: this.escape('char', -1) },
      { key: 'ArrowDown', run: this.escape('line', 1) },
      { key: 'ArrowRight', run: this.escape('char', 1) },
      { key: 'Ctrl-Enter', run: run(() => this.exit()) },
      { key: 'Mod-Enter', run: run(() => this.exit()) },
      { key: 'Mod-z', run: run(() => this.editor.commands.undo()) },
      { key: 'Mod-y', run: run(() => this.editor.commands.redo()) },
      { key: 'Shift-Mod-z', run: run(() => this.editor.commands.redo()) },
    ]
  }

  // exit leaves the block for a paragraph after it
  exit() {
    if (!this.editor.commands.exitCode()) return false
    this.view.focus()
    return true
  }

  // escape: an arrow at an end of the block leaves it
  escape(unit, dir) {
    return (cm) => {
      const { state } = cm
      let { main } = state.selection
      if (!main.empty) return false
      if (unit === 'line') main = state.doc.lineAt(main.head)
      if (dir < 0 ? main.from > 0 : main.to < state.doc.length) return false
      const targetPos = this.getPos() + (dir < 0 ? 0 : this.node.nodeSize)
      const selection = Selection.near(this.view.state.doc.resolve(targetPos), dir)
      this.view.dispatch(this.view.state.tr.setSelection(selection).scrollIntoView())
      this.view.focus()
      return true
    }
  }

  update(node) {
    if (node.type !== this.node.type) return false
    const lang = node.attrs.language !== this.node.attrs.language
    this.node = node
    if (lang) {
      this.select.value = node.attrs.language || ''
      if (this.cm) this.cm.dispatch({ effects: this.lang.reconfigure(this.languageOf(node.attrs.language)) })
    }
    if (this.updating) return true
    if (!this.cm) {
      this.host.firstChild && (this.host.firstChild.textContent = node.textContent)
      return true
    }
    const newText = node.textContent
    const curText = this.cm.state.doc.toString()
    if (newText !== curText) {
      let start = 0
      let curEnd = curText.length
      let newEnd = newText.length
      while (start < curEnd && curText.charCodeAt(start) === newText.charCodeAt(start)) ++start
      while (curEnd > start && newEnd > start && curText.charCodeAt(curEnd - 1) === newText.charCodeAt(newEnd - 1)) {
        curEnd--
        newEnd--
      }
      this.updating = true
      this.cm.dispatch({ changes: { from: start, to: curEnd, insert: newText.slice(start, newEnd) } })
      this.updating = false
    }
    return true
  }

  selectNode() {
    if (this.cm) this.cm.focus()
  }

  stopEvent() {
    return true
  }

  ignoreMutation() {
    return true
  }

  destroy() {
    this.destroyed = true
    if (this.cm) this.cm.destroy()
  }
}

// withCodeMirror is the code block node shown by CodeMirror (its content
// text only).
export function withCodeMirror(CodeBlock) {
  return CodeBlock.extend({
    addNodeView() {
      return (props) => new CodeBlockView(props)
    },
  })
}
