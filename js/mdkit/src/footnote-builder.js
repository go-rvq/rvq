// mdkit-footnote-builder: the footnotes of a text, a Vue component — a
// dialog, open(adapter, options) —: its notes listed (their numbers, labels,
// texts, references; a note with no reference, a reference with no note);
// a reference to one put at the cursor, a new one made (its label filled: a
// random one, "n4k2x" — the adapter's next() —, or a name), its text written — formatted, previewed as the
// editor's (the text alone, side by side, the preview alone) —, renamed,
// deleted (asked in the dialog). It works on the adapter, the editor's —
// the Markdown's (footnote-model) or the rich text's —:
//
//	{ kind: 'markdown' | 'rich', notes(), next(), create(text, label), insert(label),
//	  edit(label, text), rename(from, to), remove(label), done(),
//	  pointed(), point(label) }
//
// Opened from a reference (options.ref), its note chosen and marked; another
// (or a new one) may be made the reference's (point).
//
// done, the editor changed once (organized: the notes at the end, in
// order); cancelled, as it was. With options.pick (a cell of a table), the
// label of the note chosen given to options.picked, the dialog closed.
// Its words: `words` (the config's) and `barWords` (bold, …).
import { cleanLabel } from './footnote-model.js'
import { config, wordsOf } from './config.js'

// style is the builder's look, once for the page.
function style() {
  if (document.getElementById('mdkit-footnote-builder-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-footnote-builder-style'
  st.textContent =
    '.hfb-list{border-right:1px solid #e5e7eb;min-width:16rem;max-width:18rem;overflow:auto;max-height:60vh}' +
    '.hfb-text{width:100%;min-height:9rem;border:1px solid #d1d5db;border-radius:4px;padding:.5rem;font:inherit;resize:vertical;outline:0}' +
    '.hfb-text:focus{border-color:#3b82f6}' +
    '.hfb-split{border-left:1px solid #d1d5db;margin-left:.75rem;padding-left:.75rem}' +
    '.hfb-preview{overflow:auto;max-height:40vh}'
  document.head.append(st)
}

export default {
  name: 'MdkitFootnoteBuilder',
  props: {
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('footnote') },
    barWords: { type: Object, default: () => wordsOf('rich') },
    // render, when there is a preview: the HTML of a text of Markdown (a promise)
    render: { type: Function, default: null },
    // noteEditor, the editor of the text of a note of a rich text (the rich
    // text editor: mdkit-rich-editor)
    noteEditor: { type: Object, default: null },
  },
  data: () => ({
    show: false, a: null, notes: [], sel: null, label: '', text: '', isNew: false,
    asking: false, notice: '', mode: 'write', previewHTML: '', pick: false, picked: null,
    ref: false, pointed: null,
  }),
  computed: {
    // the class of a text shown (config: the page's)
    textClass: () => config.textClass,
    // the note chosen (none: a new one)
    note() { return this.notes.find((n) => n.label === this.sel) || null },
    rich() { return this.a && this.a.kind === 'rich' },
    formats: () => [
      { key: 'bold', icon: 'mdi-format-bold', m: '**' }, { key: 'italic', icon: 'mdi-format-italic', m: '_' },
      { key: 'strike', icon: 'mdi-format-strikethrough', m: '~~' }, { key: 'code', icon: 'mdi-code-tags', m: '`' },
      { key: 'link', icon: 'mdi-link-variant' },
    ],
    modes() {
      const w = this.words
      return [
        { key: 'write', icon: 'mdi-pencil-outline', title: w.modeWrite || 'Write' },
        { key: 'split', icon: 'mdi-view-split-vertical', title: w.modeSplit || 'Side by side' },
        { key: 'total', icon: 'mdi-eye-outline', title: w.modePreview || 'Preview only' },
      ]
    },
  },
  template: `
<v-dialog v-model="show" max-width="1100" scrollable data-footnote-builder>
  <v-card>
    <v-card-title class="d-flex align-center">
      <v-icon icon="mdi-format-superscript" class="me-2"></v-icon>{{ words.title || 'Footnotes' }}
      <v-spacer></v-spacer>
      <v-btn size="small" variant="tonal" color="primary" prepend-icon="mdi-plus" data-footnote-new @click="startNew">
        {{ words.add || 'New note' }}
      </v-btn>
    </v-card-title>
    <v-divider></v-divider>
    <v-card-text class="pa-0">
      <div class="d-flex">
        <div class="hfb-list" data-footnote-list>
          <div v-if="!notes.length" class="text-caption text-medium-emphasis pa-3">{{ words.empty || 'No footnotes yet.' }}</div>
          <v-list density="compact" class="py-0">
            <v-list-item v-for="n in notes" :key="n.label" :active="sel === n.label && !isNew" @click="choose(n.label)"
              :data-footnote-item="n.label">
              <template #prepend>
                <v-chip size="x-small" class="me-2" :color="n.number ? 'primary' : undefined" label>{{ n.number || '–' }}</v-chip>
              </template>
              <v-list-item-title class="text-body-2">[^{{ n.label }}]</v-list-item-title>
              <v-list-item-subtitle class="text-caption">{{ snippet(n.text) || (words.noNote || 'No note') }}</v-list-item-subtitle>
              <template #append>
                <v-icon v-if="ref && n.label === pointed" icon="mdi-target" color="primary" size="small" class="me-1"
                  :title="words.here || 'This reference'" data-footnote-pointed></v-icon>
                <v-icon v-if="!n.defined" icon="mdi-alert-outline" color="error" size="small" :title="words.noNote || 'No note'"></v-icon>
                <v-icon v-else-if="!n.refs" icon="mdi-link-variant-off" color="warning" size="small" :title="words.noRef || 'No reference'"></v-icon>
                <span v-else class="text-caption text-medium-emphasis" :title="words.refs || 'References'">{{ n.refs }}×</span>
              </template>
            </v-list-item>
          </v-list>
        </div>
        <div class="flex-grow-1 pa-3" style="min-width:0">
          <template v-if="isNew || note">
            <div class="d-flex ga-2 align-center mb-2">
              <v-text-field v-model="label" density="compact" variant="outlined" hide-details style="max-width:14rem"
                :label="words.label || 'Label'" prefix="[^" suffix="]" data-footnote-label
                @update:model-value="label = clean(label)"></v-text-field>
              <v-btn v-if="!isNew && label !== sel" size="small" variant="tonal" data-footnote-rename @click="rename">
                {{ words.rename || 'Rename' }}
              </v-btn>
              <v-spacer></v-spacer>
              <template v-if="!isNew">
                <template v-if="asking">
                  <span class="text-caption">{{ words.removeAsk || 'Delete this note and its references?' }}</span>
                  <v-btn size="small" color="error" variant="flat" data-footnote-remove-yes @click="remove">{{ words.removeYes || 'Delete' }}</v-btn>
                  <v-btn size="small" variant="text" @click="asking = false">{{ words.cancel || 'Cancel' }}</v-btn>
                </template>
                <v-btn v-else size="small" variant="text" color="error" prepend-icon="mdi-delete-outline" data-footnote-remove
                  @click="asking = true">{{ words.remove || 'Delete' }}</v-btn>
              </template>
            </div>
            <div v-if="!rich" class="d-flex align-center mb-1">
              <v-btn v-for="f in formats" :key="f.key" :icon="f.icon" size="small" variant="text" density="comfortable"
                :title="barWords[f.key] || f.key" :disabled="mode === 'total'" @mousedown.prevent @click="format(f)"></v-btn>
              <v-spacer></v-spacer>
              <v-btn-toggle v-if="render" v-model="mode" density="compact" variant="text" mandatory @update:model-value="schedule(0)">
                <v-btn v-for="m in modes" :key="m.key" :value="m.key" :icon="m.icon" size="small" :title="m.title" :data-footnote-mode="m.key"></v-btn>
              </v-btn-toggle>
            </div>
            <div v-if="rich">
              <component :is="noteEditor" v-model="text" :footnotes="false" data-footnote-text></component>
            </div>
            <div v-else class="d-flex">
              <textarea v-show="mode !== 'total'" ref="area" v-model="text" class="hfb-text flex-grow-1" :style="mode === 'split' ? 'flex-basis:0;min-width:0' : ''"
                :placeholder="words.text || 'The note'" data-footnote-text @input="schedule()"></textarea>
              <div v-if="mode !== 'write'" class="hfb-preview flex-grow-1" :class="[textClass, mode === 'split' ? 'hfb-split' : 'border rounded pa-2']"
                :style="mode === 'split' ? 'flex-basis:0;min-width:0' : ''" v-html="previewHTML" data-footnote-preview></div>
            </div>
            <div class="d-flex ga-2 mt-3">
              <v-btn v-if="ref && (isNew || sel !== pointed)" color="primary" variant="flat" prepend-icon="mdi-target"
                data-footnote-point @click="usePointed">{{ words.useHere || 'Use it in this reference' }}</v-btn>
              <v-btn color="primary" variant="tonal" prepend-icon="mdi-format-superscript" data-footnote-insert @click="insert">
                {{ words.insert || 'Insert the reference' }}
              </v-btn>
              <v-btn v-if="!isNew && !pick" variant="text" data-footnote-save @click="save">{{ words.save || 'Keep the note' }}</v-btn>
            </div>
            <div v-if="notice" class="text-caption text-error mt-2" data-footnote-notice>{{ notice }}</div>
          </template>
          <div v-else class="text-caption text-medium-emphasis">{{ words.choose || 'Choose a note, or make a new one.' }}</div>
        </div>
      </div>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ words.cancel || 'Cancel' }}</v-btn>
      <v-btn v-if="!pick" color="primary" variant="flat" data-footnote-done @click="finish">{{ words.done || 'Done' }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>`,
  mounted() {
    style()
  },
  methods: {
    // open shows the notes of the adapter; options: label (the note chosen),
    // pick and picked (a cell of a table: the note's label given)
    open(adapter, options = {}) {
      this.a = adapter
      this.pick = !!options.pick
      this.picked = options.picked || null
      this.ref = !!options.ref
      this.asking = false
      this.notice = ''
      this.mode = 'write'
      this.reload()
      if (options.label && this.notes.some((n) => n.label === options.label)) this.choose(options.label)
      else if (!this.notes.length) this.startNew()
      else this.sel = null
      this.show = true
    },
    reload() {
      this.notes = this.a.notes()
      this.pointed = this.ref && this.a.pointed ? this.a.pointed() : null
    },
    clean: cleanLabel,
    snippet: (t) => String(t || '').replace(/<[^>]*>/g, '').replace(/\s+/g, ' ').trim().slice(0, 60),
    choose(label) {
      const n = this.notes.find((x) => x.label === label)
      this.sel = label
      this.isNew = false
      this.label = label
      this.text = n ? n.text : ''
      this.asking = false
      this.notice = ''
      this.schedule(0)
    },
    startNew() {
      this.sel = null
      this.isNew = true
      this.label = this.a.next()
      this.text = ''
      this.asking = false
      this.notice = ''
      this.previewHTML = ''
    },
    // insert puts a reference to the note at the cursor — a new one made
    // first, the text of a note kept —
    insert() {
      let label = this.label
      if (this.isNew) {
        label = this.a.create(this.text, this.label)
      } else {
        this.a.edit(this.sel, this.text)
      }
      if (this.pick) {
        this.show = false
        this.a.done()
        if (this.picked) this.picked(label)
        return
      }
      this.a.insert(label)
      this.reload()
      this.choose(label)
    },
    // usePointed makes the note chosen (a new one made first) the one of the
    // reference clicked
    usePointed() {
      let label = this.sel
      if (this.isNew) label = this.a.create(this.text, this.label)
      else this.a.edit(this.sel, this.text)
      this.a.point(label)
      this.reload()
      this.choose(label)
    },
    save() {
      this.a.edit(this.sel, this.text)
      this.reload()
      this.choose(this.sel)
    },
    rename() {
      if (!this.a.rename(this.sel, this.label)) {
        this.notice = this.words.taken || 'This label is taken.'
        return
      }
      this.a.edit(this.label, this.text)
      this.reload()
      this.choose(this.label)
    },
    remove() {
      this.a.remove(this.sel)
      this.reload()
      this.sel = null
      this.isNew = false
      this.asking = false
    },
    finish() {
      if (!this.isNew && this.sel) this.a.edit(this.sel, this.text)
      this.show = false
      this.a.done()
    },
    // format puts the marks of f around the text chosen in the note (a
    // link, [text](https://))
    format(f) {
      const el = this.$refs.area
      if (!el) return
      const v = this.text || ''
      const a = el.selectionStart
      const b = el.selectionEnd
      let next
      let from
      let to
      if (f.key === 'link') {
        const label = v.slice(a, b) || (this.barWords.link || 'link')
        next = v.slice(0, a) + '[' + label + '](https://)' + v.slice(b)
        from = a + label.length + 3
        to = from + 8
      } else if (v.slice(a - f.m.length, a) === f.m && v.slice(b, b + f.m.length) === f.m) {
        next = v.slice(0, a - f.m.length) + v.slice(a, b) + v.slice(b + f.m.length)
        from = a - f.m.length
        to = b - f.m.length
      } else {
        next = v.slice(0, a) + f.m + v.slice(a, b) + f.m + v.slice(b)
        from = a + f.m.length
        to = b + f.m.length
      }
      this.text = next
      this.$nextTick(() => {
        el.focus()
        el.setSelectionRange(from, to)
      })
      this.schedule()
    },
    // schedule renders the note (the server's), after wait ms
    schedule(wait = 250) {
      if (!this.render || this.mode === 'write' || this.rich) return
      clearTimeout(this.timer)
      this.timer = setTimeout(async () => { this.previewHTML = await this.render(this.text || '') }, wait)
    },
  },
}
