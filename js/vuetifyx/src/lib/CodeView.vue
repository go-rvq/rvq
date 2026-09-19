<script setup lang="ts">
import { computed } from 'vue'
import Prism from 'prismjs'
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-markup' // html / xml
import 'prismjs/components/prism-yaml'
import { registerGad } from '@gad-lang/prism-gad'

// A readonly code viewer: Prism syntax highlighting, numbered lines, and per-line
// diff marks (added / removed) in theme colors. Used by the history JSON /
// plain-text diff handlers. The language is chosen by the caller; an unknown one
// falls back to escaped plain text.

registerGad(Prism)

const props = defineProps<{
  code?: string
  language?: string
  markLines?: number[] // 1-based line numbers to highlight
  markKind?: string // 'add' | 'del'
  // Intra-line change ranges (GoLand-style): 1-based line -> list of
  // [start, end, hunk] — start/end are rune offsets within the line and hunk is
  // the change region's index (-1 when there is no hunk grouping).
  changeRanges?: Record<number, [number, number, number][]>
}>()

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
  <pre class="vx-code"><code><span
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
