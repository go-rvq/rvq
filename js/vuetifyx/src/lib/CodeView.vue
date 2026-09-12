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
}>()

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

const lines = computed(() => {
  const src = props.code ?? ''
  const lang = props.language || ''
  const grammar = lang ? (Prism.languages as any)[lang] : null
  return src.split('\n').map((ln) => {
    if (grammar) {
      try {
        return Prism.highlight(ln, grammar, lang)
      } catch (e) {
        /* fall through to escaped text */
      }
    }
    return escapeHtml(ln)
  })
})

const markSet = computed(() => new Set(props.markLines ?? []))
const markClass = computed(() => (props.markKind === 'del' ? 'vx-code-del' : 'vx-code-add'))
</script>

<template>
  <pre class="vx-code"><code><span
    v-for="(ln, i) in lines"
    :key="i"
    class="vx-code-line"
    :class="markSet.has(i + 1) ? markClass : ''"
  ><span class="vx-code-ln">{{ i + 1 }}</span><span class="vx-code-src" v-html="ln || ' '"></span></span></code></pre>
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
</style>
