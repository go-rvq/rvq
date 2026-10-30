// The CodeMirror of the code blocks of the rich text editor (built to rvq's
// js/vendor/codemirror: loaded by mdkit's code-block-tiptap.js when a code
// block is shown in an editor): window.RvqCodeMirror — its view, its state,
// its commands, the highlighting of its languages (languages[name]() its
// extensions; the names the ones of a block's language-* class).
import { EditorView, keymap, lineNumbers, highlightActiveLine, drawSelection } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { defaultKeymap, history, historyKeymap, indentWithTab, undo, redo } from '@codemirror/commands'
import { syntaxHighlighting, defaultHighlightStyle, StreamLanguage, indentUnit, bracketMatching } from '@codemirror/language'
import { json } from '@codemirror/lang-json'
import { javascript } from '@codemirror/lang-javascript'
import { html } from '@codemirror/lang-html'
import { css } from '@codemirror/lang-css'
import { go } from '@codemirror/lang-go'
import { python } from '@codemirror/lang-python'
import { sql } from '@codemirror/lang-sql'
import { markdown } from '@codemirror/lang-markdown'
import { yaml } from '@codemirror/lang-yaml'
import { xml } from '@codemirror/lang-xml'
import { shell } from '@codemirror/legacy-modes/mode/shell'
import { toml } from '@codemirror/legacy-modes/mode/toml'
import { dockerFile } from '@codemirror/legacy-modes/mode/dockerfile'
import { diff } from '@codemirror/legacy-modes/mode/diff'
import { properties } from '@codemirror/legacy-modes/mode/properties'
import { rust } from '@codemirror/legacy-modes/mode/rust'
import { gad, gadx } from '@gad-lang/codemirror-gad'

const legacy = (mode) => () => [StreamLanguage.define(mode)]

// languages are the extensions of each language, by the names a block's
// class may have (Prism's too)
const languages = {
  json: () => [json()],
  javascript: () => [javascript()], js: () => [javascript()],
  jsx: () => [javascript({ jsx: true })],
  typescript: () => [javascript({ typescript: true })], ts: () => [javascript({ typescript: true })],
  tsx: () => [javascript({ jsx: true, typescript: true })],
  html: () => [html()], markup: () => [html()], vue: () => [html()],
  xml: () => [xml()], svg: () => [xml()],
  css: () => [css()],
  go: () => [go()], golang: () => [go()],
  python: () => [python()], py: () => [python()],
  sql: () => [sql()],
  markdown: () => [markdown()], md: () => [markdown()],
  yaml: () => [yaml()], yml: () => [yaml()],
  bash: legacy(shell), sh: legacy(shell), shell: legacy(shell),
  toml: legacy(toml),
  docker: legacy(dockerFile), dockerfile: legacy(dockerFile),
  diff: legacy(diff),
  ini: legacy(properties), properties: legacy(properties),
  rust: legacy(rust), rs: legacy(rust),
  gad: () => [gad()], gadx: () => [gadx()],
}

window.RvqCodeMirror = {
  EditorView, EditorState, Compartment, keymap, lineNumbers, highlightActiveLine, drawSelection,
  defaultKeymap, history, historyKeymap, indentWithTab, undo, redo,
  syntaxHighlighting, defaultHighlightStyle, indentUnit, bracketMatching,
  languages,
}
