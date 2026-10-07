<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Prism from 'prismjs'
// (the core has markup — html, xml, svg —, css, clike and javascript)
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-ini'
import 'prismjs/components/prism-toml'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-diff'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-docker'
import 'prismjs/components/prism-makefile'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-jsx'
import 'prismjs/components/prism-tsx'
// the plugins: they work on the element Prism highlights (highlightElement)
import 'prismjs/plugins/line-numbers/prism-line-numbers'
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'
import 'prismjs/plugins/match-braces/prism-match-braces'
import 'prismjs/plugins/match-braces/prism-match-braces.css'
import 'prismjs/plugins/toolbar/prism-toolbar'
import 'prismjs/plugins/toolbar/prism-toolbar.css'
import 'prismjs/plugins/show-language/prism-show-language'
import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard'
import 'prismjs/plugins/diff-highlight/prism-diff-highlight'
import 'prismjs/plugins/diff-highlight/prism-diff-highlight.css'
import { registerGad, registerGadx, registerGadTemplate } from '@gad-lang/prism-gad'

// A readonly code viewer, the one of the admin: a document's code (a Markdown
// fence), a file's diff, a record's JSON.
//
// By default the code is highlighted as Prism highlights an element, with its
// plugins: line numbers, matching braces, a toolbar with the language and a
// copy button, and — `diff` — a diff whose lines are highlighted in the
// language of the file (diff-highlight). The element is the component's own,
// built here and not by Vue, since the plugins rearrange it (the toolbar wraps
// it).
//
// With diff marks (markKind, markLines, changeRanges — the history's
// comparisons), it is a line per row instead, its marked lines tinted and the
// changed spans of a line given to the `changed` slot.

// registerGad first: gadx and gadt embed it
registerGad(Prism)
registerGadx(Prism)
registerGadTemplate(Prism)
const L = Prism.languages as any
L.golang = L.go
L.gadtemplate = L.gadt

const props = withDefaults(defineProps<{
  code?: string
  language?: string
  markLines?: number[] // 1-based line numbers to highlight
  markKind?: string // 'add' | 'del'
  // Intra-line change ranges (GoLand-style): 1-based line -> list of
  // [start, end, hunk] — start/end are rune offsets within the line and hunk is
  // the change region's index (-1 when there is no hunk grouping).
  changeRanges?: Record<number, [number, number, number][]>
  // code is a diff (git's), its lines highlighted in language
  diff?: boolean
  lineNumbers?: boolean
  matchBraces?: boolean
  // what the toolbar says: the language (default its name), the copy button
  label?: string
  copyText?: string
  copiedText?: string
  copyErrorText?: string
}>(), {
  lineNumbers: true,
  matchBraces: true,
  copyText: 'Copy',
  copiedText: 'Copied!',
  copyErrorText: 'Press Ctrl+C to copy',
})

// marked is the history's mode: a row per line, with its marks
const marked = computed(() => props.markKind != null || props.markLines != null || props.changeRanges != null)

// LABELS are the names a language is shown by, where Prism's own would not do
const LABELS: Record<string, string> = {
  gad: 'Gad', gadx: 'Gadx', gadt: 'Gad template', gadtemplate: 'Gad template',
  sh: 'Shell', shell: 'Shell', bash: 'Bash', yml: 'YAML', yaml: 'YAML', json: 'JSON',
  md: 'Markdown', markdown: 'Markdown', html: 'HTML', markup: 'HTML', xml: 'XML', svg: 'SVG',
  js: 'JavaScript', javascript: 'JavaScript', ts: 'TypeScript', typescript: 'TypeScript',
  go: 'Go', golang: 'Go', sql: 'SQL', css: 'CSS', toml: 'TOML', ini: 'INI', py: 'Python',
  python: 'Python', docker: 'Dockerfile', dockerfile: 'Dockerfile', makefile: 'Makefile', diff: 'Diff',
}

const host = ref<HTMLElement>()

// render builds the element and highlights it: Prism's plugins do the rest.
function render() {
  const el = host.value
  if (!el) return
  const lang = (props.language || '').toLowerCase()
  const known = lang && L[lang] ? lang : ''
  // a diff is diff-<language> (diff-highlight), or a plain diff
  const prismLang = props.diff ? (known ? 'diff-' + known : 'diff') : known || 'none'

  const pre = document.createElement('pre')
  pre.className = 'language-' + prismLang
  if (props.lineNumbers) pre.classList.add('line-numbers')
  if (props.matchBraces) pre.classList.add('match-braces')
  if (props.diff) pre.classList.add('diff-highlight')
  const label = props.label || (props.diff ? 'Diff' + (known ? ' · ' + (LABELS[known] || known) : '') : LABELS[lang] || lang)
  if (label) pre.setAttribute('data-language', label)
  pre.setAttribute('data-prismjs-copy', props.copyText)
  pre.setAttribute('data-prismjs-copy-success', props.copiedText)
  pre.setAttribute('data-prismjs-copy-error', props.copyErrorText)

  const code = document.createElement('code')
  code.className = 'language-' + prismLang
  code.textContent = (props.code ?? '').replace(/\n$/, '')
  pre.appendChild(code)

  el.replaceChildren(pre)
  Prism.highlightElement(code)
}

onMounted(() => { if (!marked.value) render() })
watch(() => [props.code, props.language, props.diff, props.lineNumbers, props.matchBraces, props.label,
  props.copyText, props.copiedText, props.copyErrorText], () => { if (!marked.value) render() })

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

// wrapRanges emphasizes the given [start,end) rune ranges inside a line of
// intra-line changed content is rendered as separate segments (see lineSegments),
// so a caller can make each changed span interactive through the `changed` slot.

// A rendered line is a list of segments: 'h' = unchanged, syntax-highlighted HTML;
// 'c' = a changed span (plain text + its rune offsets), rendered by the `changed`
// slot so a caller can make it interactive (e.g. a partial-revert toggle).
type Seg = { t: 'h'; html: string } | { t: 'c'; text: string; start: number; end: number; hunk: number }

function decodeEntities(s: string): string {
  return s
    .replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'").replace(/&amp;/g, '&')
}

// lineSegments splits a line of Prism-highlighted HTML into segments at the change
// range boundaries: unchanged runs keep their HTML (syntax colors), changed runs
// become plain text (decoded) for the slot. Tag markup inside a changed run is
// dropped (the changed text is what the slot needs).
function lineSegments(htmlStr: string, ranges: [number, number, number][] | undefined): Seg[] {
  if (!ranges || !ranges.length) return [{ t: 'h', html: htmlStr }]
  // The hunk index of the range covering pos, or null when none covers it.
  const hunkAt = (pos: number): number | null => {
    for (const r of ranges) if (pos >= r[0] && pos < r[1]) return r[2]
    return null
  }
  const segs: Seg[] = []
  let pos = 0
  let curH = '' // accumulating unchanged HTML
  let curC = '' // accumulating changed plain text
  let curCStart = 0
  let curHunk = -1
  const flushH = () => { if (curH) { segs.push({ t: 'h', html: curH }); curH = '' } }
  const flushC = () => { if (curC) { segs.push({ t: 'c', text: curC, start: curCStart, end: pos, hunk: curHunk }); curC = '' } }
  for (let i = 0; i < htmlStr.length; ) {
    const ch = htmlStr[i]
    if (ch === '<') {
      const end = htmlStr.indexOf('>', i)
      const tag = end === -1 ? htmlStr.slice(i) : htmlStr.slice(i, end + 1)
      // Tags belong to the HTML (unchanged) stream; inside a changed run they are
      // skipped (only the plain text is collected for the slot).
      if (!curC) curH += tag
      i += tag.length
      continue
    }
    let tok = ch
    if (ch === '&') {
      const e = htmlStr.indexOf(';', i)
      if (e !== -1 && e - i <= 10) tok = htmlStr.slice(i, e + 1)
    }
    const hk = hunkAt(pos)
    if (hk !== null) {
      // A change region: flush any pending unchanged run, and split when the hunk
      // index changes so each segment maps to one hunk.
      if (!curC || hk !== curHunk) { flushC(); flushH(); curCStart = pos; curHunk = hk }
      curC += decodeEntities(tok)
    } else {
      if (curC) flushC()
      curH += tok
    }
    pos += 1
    i += tok.length
  }
  flushH()
  flushC()
  return segs
}

const lines = computed(() => {
  const src = props.code ?? ''
  const lang = props.language || ''
  const grammar = lang ? (Prism.languages as any)[lang] : null
  const ranges = props.changeRanges || {}
  return src.split('\n').map((ln, i) => {
    let html: string
    if (grammar) {
      try {
        html = Prism.highlight(ln, grammar, lang)
      } catch (e) {
        html = escapeHtml(ln)
      }
    } else {
      html = escapeHtml(ln)
    }
    return { no: i + 1, segments: lineSegments(html, ranges[i + 1]) }
  })
})

const chgClass = computed(() => (props.markKind === 'del' ? 'vx-code-chg-del' : 'vx-code-chg-add'))

const markSet = computed(() => new Set(props.markLines ?? []))
const markClass = computed(() => (props.markKind === 'del' ? 'vx-code-del' : 'vx-code-add'))
</script>

<template>
  <div v-if="!marked" ref="host" class="vx-code-host"></div>
  <pre v-else class="vx-code"><code><span
    v-for="ln in lines"
    :key="ln.no"
    class="vx-code-line"
    :class="markSet.has(ln.no) ? markClass : ''"
  ><span class="vx-code-ln">{{ ln.no }}</span><span class="vx-code-src"><template
      v-for="(seg, si) in ln.segments" :key="si"
    ><span v-if="seg.t === 'h'" v-html="seg.html"></span><slot
        v-else name="changed"
        :text="seg.text" :line="ln.no" :start="seg.start" :end="seg.end" :hunk="seg.hunk"
        :kind="markKind" :added="markKind !== 'del'"
      ><mark :class="chgClass">{{ seg.text }}</mark></slot></template></span></span></code></pre>
</template>

<style scoped>
.vx-code {
  margin: 0;
  padding: 6px 4px;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.5;
  background: rgba(var(--v-theme-on-surface), 0.04);
  border-radius: 4px;
}
.vx-code code {
  display: block;
  font-family: 'Roboto Mono', ui-monospace, SFMono-Regular, Menlo, monospace;
}
.vx-code-line {
  display: flex;
}
.vx-code-ln {
  flex: 0 0 auto;
  width: 3em;
  padding-right: 1em;
  text-align: right;
  color: rgba(var(--v-theme-on-surface), 0.4);
  user-select: none;
}
.vx-code-src {
  flex: 1 1 auto;
  white-space: pre;
}
/* alternated, a marked line keeping its tint */
.vx-code-line:nth-child(even):not(.vx-code-add):not(.vx-code-del) {
  background: rgba(var(--v-theme-on-surface), 0.035);
}
.vx-code-add {
  background: rgba(var(--v-theme-success), 0.16);
}
.vx-code-del {
  background: rgba(var(--v-theme-error), 0.16);
}
/* Intra-line changed content (stronger than the whole-line tint). */
.vx-code-chg-add {
  background: rgba(var(--v-theme-success), 0.42);
  border-radius: 2px;
}
.vx-code-chg-del {
  background: rgba(var(--v-theme-error), 0.42);
  border-radius: 2px;
}
</style>

<!-- Not scoped: the interactive `changed` slot content is rendered in the
     parent's scope (partial-revert hunk toggles), so its classes must be global.
     Nested under .vx-code for enough specificity to beat the Prism theme rules. -->
<style>
.vx-code .vx-hunk {
  cursor: pointer;
  border-radius: 3px;
  padding: 0 2px;
  position: relative;
  transition: box-shadow 0.1s ease;
}
.vx-code .vx-hunk-add {
  background: rgba(var(--v-theme-success), 0.3);
}
/* A deletion (excluded content): red tint with a red border. */
.vx-code .vx-hunk-del {
  background: rgba(var(--v-theme-error), 0.3);
  outline: 1px solid rgb(var(--v-theme-error));
}
/* Selected (marked for revert): super-highlight. */
.vx-code .vx-hunk-sel {
  background: rgba(var(--v-theme-primary), 0.28);
  outline: 2px solid rgb(var(--v-theme-primary));
  box-shadow: 0 0 0 3px rgba(var(--v-theme-primary), 0.25);
}
/* A deletion marked for revert: red super-highlight (matches its red border). */
.vx-code .vx-hunk-del.vx-hunk-sel {
  background: rgba(var(--v-theme-error), 0.28);
  outline: 2px solid rgb(var(--v-theme-error));
  box-shadow: 0 0 0 3px rgba(var(--v-theme-error), 0.25);
}
/* The ✓ takes its span's border color: primary by default, red on a deletion. */
.vx-code .vx-hunk-check {
  color: rgb(var(--v-theme-primary));
}
.vx-code .vx-hunk-del .vx-hunk-check {
  color: rgb(var(--v-theme-error));
}
</style>

<!-- The highlighted element is built by render(), not by Vue: its rules are
     global, under .vx-code-host — specific enough to beat the Prism theme the
     editor brings (coy: a pre with its own padding, ::before/::after shadows, a
     striped code) — and in the theme's colors, light and dark. -->
<style>
.vx-code-host {
  margin: 8px 0 12px;
}
.vx-code-host div.code-toolbar {
  position: relative;
}
.vx-code-host pre[class*='language-'] {
  position: relative;
  float: none;
  margin: 0;
  padding: 10px 12px;
  max-height: none;
  overflow: auto;
  background: rgba(var(--v-theme-on-surface), 0.04);
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
  box-shadow: none;
  font-size: 0.8125rem;
  line-height: 1.5;
  tab-size: 4;
  white-space: pre;
  word-wrap: normal;
}
/* the numbers in the pre's padding, not over the code: the coy theme puts
   them at the code's left (its .line-numbers.line-numbers rules), with a
   padding of its own on the code */
.vx-code-host pre[class*='language-'].line-numbers.line-numbers {
  padding-left: 3.8em;
}
.vx-code-host pre[class*='language-'].line-numbers.line-numbers > code {
  padding-left: 0;
}
.vx-code-host pre[class*='language-'].line-numbers.line-numbers .line-numbers-rows {
  left: -3.8em;
}
.vx-code-host pre[class*='language-']::before,
.vx-code-host pre[class*='language-']::after {
  content: none;
  display: none;
  box-shadow: none;
}
.vx-code-host pre[class*='language-'] > code[class*='language-'] {
  position: relative;
  display: block;
  height: auto;
  max-height: none;
  overflow: visible;
  padding: 0;
  margin: 0;
  border: 0;
  box-shadow: none;
  background: none;
  font-family: 'Roboto Mono', ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: inherit;
  color: rgb(var(--v-theme-on-surface));
  white-space: inherit;
  text-shadow: none;
}
/* the lines, lightly alternated: a band every other line (the line height is
   1.5em, so two lines are 3em), as wide as the longest line */
.vx-code-host pre[class*='language-'] > code[class*='language-'] {
  min-width: 100%;
  width: max-content;
  background-image: linear-gradient(transparent 50%, rgba(var(--v-theme-on-surface), 0.035) 50%);
  background-size: 100% 3em;
  background-origin: content-box;
}
.vx-code-host .line-numbers .line-numbers-rows {
  border-right: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-code-host .line-numbers-rows > span::before {
  color: rgba(var(--v-theme-on-surface), 0.4);
}
/* the toolbar: the language and the copy button, at the top right */
.vx-code-host div.code-toolbar > .toolbar {
  top: 6px;
  right: 8px;
}
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item + .toolbar-item {
  margin-left: 4px;
}
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item > button,
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item > span,
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item > a {
  color: rgba(var(--v-theme-on-surface), 0.7);
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  box-shadow: none;
  border-radius: 4px;
  padding: 3px 9px;
  font-size: 0.8125rem;
  line-height: 1.4;
  cursor: default;
}
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item > button {
  cursor: pointer;
}
.vx-code-host div.code-toolbar > .toolbar > .toolbar-item > button:hover {
  color: rgb(var(--v-theme-primary));
}
/* the diff: its lines tinted in the theme's colors */
.vx-code-host pre.diff-highlight > code .token.deleted:not(.prefix),
.vx-code-host pre > code.diff-highlight .token.deleted:not(.prefix) {
  background-color: rgba(var(--v-theme-error), 0.14);
}
.vx-code-host pre.diff-highlight > code .token.inserted:not(.prefix),
.vx-code-host pre > code.diff-highlight .token.inserted:not(.prefix) {
  background-color: rgba(var(--v-theme-success), 0.14);
}
.vx-code-host .token.coord {
  color: rgb(var(--v-theme-info));
}
/* matching braces */
.vx-code-host .token.punctuation.brace-hover,
.vx-code-host .token.punctuation.brace-selected {
  outline: 1px solid rgba(var(--v-theme-primary), 0.7);
  border-radius: 2px;
}
/* the tokens: none with the background the coy theme gives some (operators,
   entities, urls), in either mode */
.vx-code-host .token,
.vx-code .token {
  background: none;
}
.vx-code-host .token.comment, .vx-code-host .token.prolog, .vx-code-host .token.doctype, .vx-code-host .token.cdata { color: #6e7781; font-style: italic; }
.vx-code-host .token.punctuation { color: #656d76; }
.vx-code-host .token.keyword, .vx-code-host .token.boolean, .vx-code-host .token.atrule, .vx-code-host .token.important { color: #cf222e; }
.vx-code-host .token.string, .vx-code-host .token.char, .vx-code-host .token.attr-value, .vx-code-host .token.regex { color: #0a7d33; }
.vx-code-host .token.number, .vx-code-host .token.constant, .vx-code-host .token.symbol { color: #0550ae; }
.vx-code-host .token.function, .vx-code-host .token.class-name { color: #8250df; }
.vx-code-host .token.builtin, .vx-code-host .token.type, .vx-code-host .token.tag, .vx-code-host .token.attr-name, .vx-code-host .token.selector, .vx-code-host .token.property { color: #953800; }
.vx-code-host .token.operator, .vx-code-host .token.entity, .vx-code-host .token.url, .vx-code-host .token.variable { color: inherit; }
.v-theme--dark .vx-code-host .token.comment, .v-theme--dark .vx-code-host .token.prolog, .v-theme--dark .vx-code-host .token.doctype, .v-theme--dark .vx-code-host .token.cdata { color: #8b949e; }
.v-theme--dark .vx-code-host .token.punctuation { color: #9aa0a6; }
.v-theme--dark .vx-code-host .token.keyword, .v-theme--dark .vx-code-host .token.boolean, .v-theme--dark .vx-code-host .token.atrule, .v-theme--dark .vx-code-host .token.important { color: #ff7b72; }
.v-theme--dark .vx-code-host .token.string, .v-theme--dark .vx-code-host .token.char, .v-theme--dark .vx-code-host .token.attr-value, .v-theme--dark .vx-code-host .token.regex { color: #7ee787; }
.v-theme--dark .vx-code-host .token.number, .v-theme--dark .vx-code-host .token.constant, .v-theme--dark .vx-code-host .token.symbol { color: #79c0ff; }
.v-theme--dark .vx-code-host .token.function, .v-theme--dark .vx-code-host .token.class-name { color: #d2a8ff; }
.v-theme--dark .vx-code-host .token.builtin, .v-theme--dark .vx-code-host .token.type, .v-theme--dark .vx-code-host .token.tag, .v-theme--dark .vx-code-host .token.attr-name, .v-theme--dark .vx-code-host .token.selector, .v-theme--dark .vx-code-host .token.property { color: #ffa657; }
</style>
