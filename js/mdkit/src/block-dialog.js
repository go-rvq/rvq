// mdkit-block-dialog: a formula or a diagram, written — a Vue component, a
// dialog (open(kind, value, done)) —: a formula's LaTeX, in the line or on a
// line of its own, a diagram's Mermaid source; its preview as it is written
// (KaTeX, Mermaid: loaded from where the page's meta mdkit-vendor says);
// done given {text, block}, or null to take it out (one there).
import { wordsOf } from './config.js'
import { drawMath } from './math-tiptap.js'
import { drawInto } from './diagram-tiptap.js'

const sample = 'graph TD\n  A[Start] --> B{Ok?}\n  B -->|Yes| C[Done]\n  B -->|No| A'

export default {
  name: 'MdkitBlockDialog',
  props: {
    words: { type: Object, default: () => wordsOf('blocks') },
  },
  data: () => ({ show: false, kind: 'math', text: '', block: false, there: false, done: null, invalid: false, timer: null }),
  computed: {
    title() {
      const w = this.words
      return this.kind === 'math' ? (w.formulaTitle || 'Formula (LaTeX)') : (w.diagramTitle || 'Diagram (Mermaid)')
    },
    guide() {
      return this.kind === 'math' ? 'https://katex.org/docs/supported.html' : 'https://mermaid.js.org/intro/syntax-reference.html'
    },
  },
  watch: {
    text() { this.schedule() },
    block() { this.schedule() },
  },
  template: `
<v-dialog v-model="show" :max-width="kind === 'math' ? 640 : 900" data-block-dialog>
  <v-card>
    <v-card-title class="d-flex align-center ga-2">
      <v-icon :icon="kind === 'math' ? 'mdi-sigma' : 'mdi-sitemap-outline'"></v-icon>{{ title }}
      <v-spacer></v-spacer>
      <v-btn size="small" variant="text" :href="guide" target="_blank" rel="noopener" prepend-icon="mdi-book-open-variant">{{ words.guide || 'Guide' }}</v-btn>
    </v-card-title>
    <v-card-text>
      <div class="d-flex ga-4 flex-wrap">
        <div style="flex: 1 1 18rem; min-width: 0">
          <v-textarea v-model="text" :label="kind === 'math' ? (words.latex || 'LaTeX') : (words.source || 'Source (Mermaid)')"
            :rows="kind === 'math' ? 3 : 10" auto-grow autofocus variant="outlined" hide-details style="font-family: monospace"
            data-block-text></v-textarea>
          <v-btn-toggle v-if="kind === 'math'" v-model="block" mandatory density="compact" variant="outlined" class="mt-3" data-block-kind>
            <v-btn :value="false" size="small">{{ words.inline || 'In the line' }}</v-btn>
            <v-btn :value="true" size="small">{{ words.block || 'On a line of its own' }}</v-btn>
          </v-btn-toggle>
        </div>
        <div style="flex: 1 1 18rem; min-width: 0">
          <div class="text-caption text-medium-emphasis">{{ words.preview || 'Preview' }}</div>
          <div ref="preview" class="pa-3 border rounded" :class="kind === 'diagram' || block ? 'text-center' : ''"
            style="min-height: 4rem; overflow: auto" data-block-preview></div>
          <div v-if="invalid" class="text-caption text-error mt-1" data-block-invalid>{{ words.invalid || 'It does not draw: check the source.' }}</div>
        </div>
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn v-if="there" color="error" variant="text" @click="finish(null)" data-block-remove>{{ words.remove || 'Remove' }}</v-btn>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ words.cancel || 'Cancel' }}</v-btn>
      <v-btn color="primary" variant="flat" :disabled="!text.trim()" @click="finish({ text: text.trim(), block })" data-block-apply>
        {{ there ? (words.update || 'Update') : (words.insert || 'Insert') }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>`,
  methods: {
    // open shows the dialog of kind (math, diagram) with value — {text,
    // block, there: one in the text} —, done given what is written
    open(kind, value, done) {
      this.kind = kind
      this.text = (value && value.text) || (kind === 'diagram' ? sample : '')
      this.block = !!(value && value.block)
      this.there = !!(value && value.there)
      this.done = done
      this.invalid = false
      this.show = true
      this.schedule(0)
    },
    // schedule draws the preview, a moment after the last key
    schedule(wait = 300) {
      clearTimeout(this.timer)
      this.timer = setTimeout(async () => {
        const el = this.$refs.preview
        if (!this.show || !el) return
        this.invalid = this.kind === 'math' ? !(await drawMath(el, this.text, this.block)) : !(await drawInto(el, this.text))
      }, wait)
    },
    finish(value) {
      this.show = false
      if (this.done) this.done(value)
    },
  },
}
