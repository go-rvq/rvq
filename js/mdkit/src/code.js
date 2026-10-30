// The blocks of code of a text the page shows (x/markup's fences, the rich
// text editor's code blocks: pre > code.language-*) highlighted as the
// admin's code viewer (vuetifyx vx-code) highlights them: Prism, in their
// language, with line numbers, and a toolbar with the language and a copy
// button. Prism is the page's (the admin's) or, none, the one where the
// page's meta mdkit-vendor says ("prism"), loaded only when there is a block.
// What is in an editor (contenteditable), a diagram (mermaid) or already the
// code viewer's is not touched.
import { load, own } from './blocks.js'
import { wordsOf } from './config.js'

export const codeSel = 'pre > code[class*="language-"]'

// LABELS are the names a language is shown by (vx-code's), where its own
// would not do.
const LABELS = {
  gad: 'Gad', gadx: 'Gadx', gadt: 'Gad template', gadtemplate: 'Gad template',
  sh: 'Shell', shell: 'Shell', bash: 'Bash', yml: 'YAML', yaml: 'YAML', json: 'JSON',
  md: 'Markdown', markdown: 'Markdown', html: 'HTML', markup: 'HTML', xml: 'XML', svg: 'SVG',
  js: 'JavaScript', javascript: 'JavaScript', ts: 'TypeScript', typescript: 'TypeScript',
  go: 'Go', golang: 'Go', sql: 'SQL', css: 'CSS', toml: 'TOML', ini: 'INI', py: 'Python',
  python: 'Python', docker: 'Dockerfile', dockerfile: 'Dockerfile', makefile: 'Makefile', diff: 'Diff',
  jsx: 'JSX', tsx: 'TSX', plaintext: 'Text', text: 'Text', txt: 'Text',
}

// languageOf is the language of a code element (its language-* class),
// lower case; '' none.
export function languageOf(code) {
  const m = /(?:^|\s)language-([\w+-]+)/.exec(code.className || '')
  return m ? m[1].toLowerCase() : ''
}

// labelOf is the name a language is shown by.
export function labelOf(lang) {
  return LABELS[lang] || lang
}

// COPY are the words of the copy button by the page's language, when the
// page gives none (wordsOf('code')).
const COPY = {
  en: { copy: 'Copy', copied: 'Copied!', copyError: 'Press Ctrl+C to copy' },
  pt: { copy: 'Copiar', copied: 'Copiado!', copyError: 'Tecle Ctrl+C para copiar' },
  es: { copy: 'Copiar', copied: '¡Copiado!', copyError: 'Pulse Ctrl+C para copiar' },
}

function copyWords() {
  const lang = ((typeof document !== 'undefined' && document.documentElement.lang) || 'en').slice(0, 2).toLowerCase()
  return { ...(COPY[lang] || COPY.en), ...wordsOf('code') }
}

// style is the look of a block, once for the page: the code viewer's box,
// Prism coy's colors (and lighter ones on a dark page).
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-code-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-code-style'
  const p = 'pre.mdkit-code'
  st.textContent =
    '.mdkit-code-host{margin:.75em 0 1em}' +
    '.mdkit-code-host div.code-toolbar{position:relative}' +
    p + '{position:relative;float:none;margin:0;padding:10px 12px;max-height:none;overflow:auto;' +
    'background:rgba(127,127,127,.07);border:1px solid rgba(127,127,127,.28);border-radius:6px;box-shadow:none;' +
    'font-size:.8125rem;line-height:1.5;tab-size:4;white-space:pre;word-wrap:normal;color:inherit}' +
    p + '::before,' + p + '::after{display:none}' +
    p + '>code{display:block;background:none;padding:0;border:0;border-radius:0;font-size:inherit;' +
    'font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;white-space:inherit;color:inherit;box-shadow:none;max-height:none;overflow:visible}' +
    p + '.line-numbers.line-numbers{padding-left:3.8em;counter-reset:linenumber}' +
    p + '.line-numbers>code{position:relative;white-space:inherit}' +
    p + ' .line-numbers-rows{position:absolute;pointer-events:none;top:0;font-size:100%;left:-3.8em;width:3em;' +
    'letter-spacing:-1px;border-right:1px solid rgba(127,127,127,.35);user-select:none}' +
    p + ' .line-numbers-rows>span{display:block;counter-increment:linenumber}' +
    p + ' .line-numbers-rows>span::before{content:counter(linenumber);color:rgba(127,127,127,.8);display:block;padding-right:.8em;text-align:right}' +
    '.mdkit-code-host div.code-toolbar>.toolbar{position:absolute;z-index:1;top:.3em;right:.3em;opacity:0;transition:opacity .3s ease-in-out;display:flex;gap:.25em}' +
    '.mdkit-code-host div.code-toolbar:hover>.toolbar,.mdkit-code-host div.code-toolbar:focus-within>.toolbar{opacity:1}' +
    '.mdkit-code-host div.code-toolbar>.toolbar>.toolbar-item>a,.mdkit-code-host div.code-toolbar>.toolbar>.toolbar-item>button,' +
    '.mdkit-code-host div.code-toolbar>.toolbar>.toolbar-item>span{color:inherit;font-size:.75em;padding:0 .5em;' +
    'background:rgba(127,127,127,.18);box-shadow:none;border:0;border-radius:.5em;cursor:pointer;font-family:inherit;line-height:1.8}' +
    '.mdkit-code-host div.code-toolbar>.toolbar>.toolbar-item>span{cursor:default}' +
    // coy's colors
    p + ' .token.comment,' + p + ' .token.block-comment,' + p + ' .token.prolog,' + p + ' .token.doctype,' + p + ' .token.cdata{color:#7d8b99}' +
    p + ' .token.punctuation{color:#5f6364}' +
    p + ' .token.property,' + p + ' .token.tag,' + p + ' .token.boolean,' + p + ' .token.number,' + p + ' .token.function-name,' +
    p + ' .token.constant,' + p + ' .token.symbol,' + p + ' .token.deleted{color:#c92c2c}' +
    p + ' .token.selector,' + p + ' .token.attr-name,' + p + ' .token.string,' + p + ' .token.char,' + p + ' .token.function,' +
    p + ' .token.builtin,' + p + ' .token.inserted{color:#2f9c0a}' +
    p + ' .token.operator,' + p + ' .token.entity,' + p + ' .token.url,' + p + ' .token.variable{color:#a67f59;background:none}' +
    p + ' .token.atrule,' + p + ' .token.attr-value,' + p + ' .token.keyword,' + p + ' .token.class-name{color:#1990b8}' +
    p + ' .token.regex,' + p + ' .token.important{color:#e90}' +
    p + ' .token.important,' + p + ' .token.bold{font-weight:bold}' +
    p + ' .token.italic{font-style:italic}' +
    // a dark page (vuetify's dark theme, or data-theme=dark)
    ['.v-theme--dark', '[data-theme="dark"]'].map((d) =>
      d + ' ' + p + ' .token.comment,' + d + ' ' + p + ' .token.prolog{color:#8b9bab}' +
      d + ' ' + p + ' .token.punctuation{color:#a9b1b6}' +
      d + ' ' + p + ' .token.property,' + d + ' ' + p + ' .token.tag,' + d + ' ' + p + ' .token.boolean,' + d + ' ' + p + ' .token.number,' +
      d + ' ' + p + ' .token.constant,' + d + ' ' + p + ' .token.symbol,' + d + ' ' + p + ' .token.deleted{color:#f07178}' +
      d + ' ' + p + ' .token.selector,' + d + ' ' + p + ' .token.attr-name,' + d + ' ' + p + ' .token.string,' + d + ' ' + p + ' .token.char,' +
      d + ' ' + p + ' .token.function,' + d + ' ' + p + ' .token.builtin,' + d + ' ' + p + ' .token.inserted{color:#9ccc65}' +
      d + ' ' + p + ' .token.operator,' + d + ' ' + p + ' .token.entity,' + d + ' ' + p + ' .token.url,' + d + ' ' + p + ' .token.variable{color:#d7a86e}' +
      d + ' ' + p + ' .token.atrule,' + d + ' ' + p + ' .token.attr-value,' + d + ' ' + p + ' .token.keyword,' + d + ' ' + p + ' .token.class-name{color:#5ec4ff}',
    ).join('')
  document.head.append(st)
}

// ours says code is a block of a text to highlight: not a diagram (but its
// source shown: pre[data-fence]), not highlighted yet, not in an editor nor
// the code viewer's.
const ours = (code) => {
  const lang = languageOf(code)
  const source = code.parentElement && code.parentElement.hasAttribute('data-fence')
  return (source || (lang !== 'mermaid' && lang !== 'math')) && !code.hasAttribute('data-drawn') && own(code) &&
    !code.closest('.vx-code-host,.code-toolbar,.mdkit-code-host')
}

// prepare makes the pre of code the code viewer's (its classes, its label,
// its copy words) in a host of its own; the language Prism knows, else none.
export function prepare(code, Prism, words = copyWords()) {
  const pre = code.parentElement
  const lang = languageOf(code)
  const known = Prism && Prism.languages[lang] ? lang : 'none'
  code.className = code.className.replace(/(?:^|\s)language-[\w+-]+/g, '').trim()
  code.classList.add('language-' + known)
  pre.classList.add('mdkit-code', 'line-numbers', 'language-' + known)
  code.setAttribute('data-language', lang)
  const label = labelOf(lang)
  if (label && lang !== 'none') pre.setAttribute('data-language', label)
  pre.setAttribute('data-prismjs-copy', words.copy)
  pre.setAttribute('data-prismjs-copy-success', words.copied)
  pre.setAttribute('data-prismjs-copy-error', words.copyError)
  // its text without the last line end (no empty numbered line)
  code.textContent = (code.textContent || '').replace(/\n$/, '')
  if (!pre.parentElement.classList.contains('mdkit-code-host')) {
    const host = document.createElement('div')
    host.className = 'mdkit-code-host'
    pre.replaceWith(host)
    host.append(pre)
  }
  return pre
}

// drawCode highlights the blocks of code under root not highlighted yet.
export async function drawCode(root, assets) {
  const codes = [...root.querySelectorAll(codeSel)].filter(ours)
  if (!codes.length) return
  codes.forEach((c) => c.setAttribute('data-drawn', 'asked'))
  style()
  let Prism = null
  try {
    Prism = await load('prism', 'Prism', assets)
    if (!Prism || !Prism.highlightElement) Prism = null
  } catch (e) {
    Prism = null
  }
  const words = copyWords()
  for (const code of codes) {
    prepare(code, Prism, words)
    if (!Prism) {
      code.setAttribute('data-drawn', 'none')
      continue
    }
    try {
      Prism.highlightElement(code)
      code.setAttribute('data-drawn', 'drawn')
    } catch (e) {
      code.setAttribute('data-drawn', 'error')
    }
  }
}
