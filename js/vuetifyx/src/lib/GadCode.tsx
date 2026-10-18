// <vx-gad-code>: the admin's editor of gad code (gad, gadx, gadt) — the
// editor (<vx-codemirror>) with what an IDE gives of the code in it:
//
//   - in the bar of the editor: Format (the code as `gad fmt` writes it),
//     Undo and Redo (its history), Compare (the code with the value it
//     started from, side by side, scrolling as one — the merge view of the
//     IDE's Git panel —, the changes stepped through, a change reverted) and
//     Expand (the editor over the whole page, and back; Esc too);
//   - Structure, at the left: closed to a strip, it opens on the
//     declarations of the code as a tree (`gad symbols`, the IDE's
//     StructureTree); a click on one moves the cursor to it;
//   - Documentation, at the right: closed to a strip, it opens on the
//     documentation of the code (`gad doc`).
//
// What the code is asked of is served by the admin (rvq admin/gadcode):
// POST <api>{format,outline,doc} with the code and its language. The
// structure and the documentation follow the code, a moment after it changes,
// while they are open.
import {defineComponent, h, type PropType} from 'vue'
import {useVxMessages, type VxMessages} from '@/lib/locale'
import {StructureTree, type OutlineSym} from '@gad-lang/ide-vuetify/structure'
import {DiffMergeView} from '@gad-lang/git-vuetify/diff'
import {ResizeHandle, useStoredSize} from '@gad-lang/git-vuetify/resize'
import CodeMirror from '@/lib/CodeMirror'

const propsOptions = {
  modelValue: {type: String, default: ''},
  // gad, gadx or gadt
  language: {type: String, default: 'gad'},
  // the delimiters of a gadt's code
  templateStart: {type: String, default: ''},
  templateEnd: {type: String, default: ''},
  // where the code is asked of: <prefix>/__gadcode/
  api: {type: String, required: true as const},
  readonly: Boolean,
  // the value Compare compares the code with; the value it started from
  // when not said
  original: {type: String, default: undefined},
  minHeight: {type: String, default: '16rem'},
  maxHeight: {type: String, default: '40rem'},
  // the help of the editor (its slot), as <vx-codemirror>'s
  help: Boolean,
  helpMode: {type: String, default: 'dialog'},
  helpTitle: {type: String, default: undefined},
  helpTooltip: {type: String, default: undefined},
  // the words of its buttons and panels
  labels: {
    type: Object as PropType<Partial<Record<'format' | 'expand' | 'collapse' | 'structure' | 'doc' | 'refresh' | 'noStructure' | 'noDoc' | 'notFormatted' | 'undo' | 'redo' | 'compare' | 'edit' | 'prevChange' | 'nextChange' | 'revert', string>>>,
    default: () => ({}),
  },
} as const

// the words of its labels: the page's vue-i18n's (vx), the labels prop over
// them
type Labels = {
  format: string; expand: string; collapse: string; structure: string; doc: string; refresh: string
  noStructure: string; noDoc: string; notFormatted: string; undo: string; redo: string; compare: string
  edit: string; prevChange: string; nextChange: string; revert: string
}
const labelsOf = (w: VxMessages): Labels => ({
  format: w.format, expand: w.expand, collapse: w.collapse, structure: w.structure, doc: w.doc,
  refresh: w.refresh, noStructure: w.noStructure, noDoc: w.noDoc, notFormatted: w.notFormatted,
  undo: w.undo, redo: w.redo, compare: w.compare, edit: w.edit, prevChange: w.prevChange,
  nextChange: w.nextChange, revert: w.revert,
})

const border = '1px solid rgba(var(--v-border-color), var(--v-border-opacity))'

export default defineComponent({
  name: 'VXGadCode',

  props: propsOptions,

  setup() {
    // the widths of the structure and of the documentation, the height of
    // the editor: dragged by their edges, kept across visits
    return {
      vx: useVxMessages(),
      structureW: useStoredSize('vx-gad-code.structure', 260, 140, 1600),
      docW: useStoredSize('vx-gad-code.doc', 420, 160, 2400),
      editorH: useStoredSize('vx-gad-code.height', 360, 120, 6000),
    }
  },

  emits: {
    'update:modelValue': (_value: string) => true,
    'update:help': (_value: boolean) => true,
  },

  data() {
    return {
      expanded: false,
      structureOpen: false,
      docOpen: false,
      symbols: [] as OutlineSym[],
      docHtml: '',
      docError: '',
      formatError: '',
      formatting: false,
      comparing: false,
      canUndo: false,
      canRedo: false,
      // the value the code started from (Compare)
      initial: this.modelValue,
      timer: 0 as number | ReturnType<typeof setTimeout>,
    }
  },

  computed: {
    words(): Labels {
      return {...labelsOf(this.vx), ...this.labels} as Labels
    },
    code(): {source: string; language: string; templateStart: string; templateEnd: string} {
      return {
        source: this.modelValue,
        language: this.language,
        templateStart: this.templateStart,
        templateEnd: this.templateEnd,
      }
    },
  },

  watch: {
    modelValue() {
      this.formatError = ''
      if (!this.structureOpen && !this.docOpen) return
      clearTimeout(this.timer)
      this.timer = setTimeout(() => this.refresh(), 500)
    },
    structureOpen(v: boolean) {
      if (v) this.refreshStructure()
      this.remeasure()
    },
    docOpen(v: boolean) {
      if (v) this.refreshDoc()
      this.remeasure()
    },
    expanded(v: boolean) {
      document.documentElement.style.overflow = v ? 'hidden' : ''
      this.remeasure()
    },
    comparing() {
      this.canUndo = this.canRedo = false
      this.remeasure()
    },
  },

  mounted() {
    window.addEventListener('keydown', this.onKey)
  },

  beforeUnmount() {
    clearTimeout(this.timer)
    window.removeEventListener('keydown', this.onKey)
    if (this.expanded) document.documentElement.style.overflow = ''
  },

  methods: {
    async ask<T>(op: string): Promise<T> {
      const r = await fetch(this.api + op, {
        method: 'POST',
        credentials: 'same-origin',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(this.code),
      })
      if (!r.ok) throw new Error(await r.text())
      return (await r.json()) as T
    },
    refresh() {
      if (this.structureOpen) this.refreshStructure()
      if (this.docOpen) this.refreshDoc()
    },
    async refreshStructure() {
      try {
        this.symbols = await this.ask<OutlineSym[]>('outline')
      } catch {
        this.symbols = []
      }
    },
    async refreshDoc() {
      try {
        const r = await this.ask<{html: string; error?: string}>('doc')
        this.docHtml = r.html || ''
        this.docError = r.error || ''
      } catch (e) {
        this.docHtml = ''
        this.docError = e instanceof Error ? e.message : String(e)
      }
    },
    async format() {
      if (this.readonly || this.formatting) return
      this.formatting = true
      try {
        const r = await this.ask<{ok: boolean; source: string}>('format')
        if (r.ok) {
          this.formatError = ''
          if (r.source !== this.modelValue) this.$emit('update:modelValue', r.source)
        } else {
          this.formatError = this.words.notFormatted
        }
      } catch (e) {
        this.formatError = e instanceof Error ? e.message : String(e)
      } finally {
        this.formatting = false
      }
    },
    // remeasure has the editor measure its lines once its width changed
    remeasure() {
      this.$nextTick(() => (this.$refs.editor as {measure?: () => void} | undefined)?.measure?.())
    },
    navigate(sym: OutlineSym) {
      if (this.comparing) this.comparing = false
      this.$nextTick(() =>
        (this.$refs.editor as {goto?: (line: number, column: number) => void} | undefined)?.goto?.(sym.line, sym.column),
      )
    },
    // the editor in use: the code's, or the current side of the comparison
    active(): {undo?: () => void; redo?: () => void; next?: () => void; prev?: () => void} | undefined {
      return (this.comparing ? this.$refs.merge : this.$refs.editor) as never
    },
    onHistory(st: {canUndo: boolean; canRedo: boolean}) {
      this.canUndo = st.canUndo
      this.canRedo = st.canRedo
    },
    dark(): boolean {
      return !!(this as unknown as {$vuetify?: {theme?: {current?: {value?: {dark?: boolean}}}}}).$vuetify?.theme?.current?.value?.dark
    },
    onKey(e: KeyboardEvent) {
      if (e.key === 'Escape' && this.expanded) this.expanded = false
    },

    // the strip a closed panel is: its icon and its name, written down
    strip(icon: string, label: string, open: () => void, side: 'left' | 'right') {
      return (
        <div
          class="vx-gad-code__strip"
          title={label}
          onClick={open}
          style={{
            flex: '0 0 26px', cursor: 'pointer', display: 'flex', flexDirection: 'column', alignItems: 'center',
            gap: '6px', paddingTop: '6px', [side === 'left' ? 'borderRight' : 'borderLeft']: border,
            userSelect: 'none',
          }}
        >
          <v-icon size="16">{icon}</v-icon>
          <span class="text-caption" style={{writingMode: 'vertical-rl', transform: side === 'left' ? 'rotate(180deg)' : ''}}>
            {label}
          </span>
        </div>
      )
    },

    // the head of an open panel: its name, refresh and close
    panelHead(label: string, refresh: () => void, close: () => void, closeIcon: string) {
      return (
        <div style={{display: 'flex', alignItems: 'center', padding: '2px 4px 2px 8px', borderBottom: border}}>
          <span class="text-caption font-weight-medium" style={{flex: '1 1 auto'}}>{label}</span>
          <v-btn icon={'mdi-refresh'} size="x-small" variant="text" density="comfortable" title={this.words.refresh} onClick={refresh} />
          <v-btn icon={closeIcon} size="x-small" variant="text" density="comfortable" title={label} onClick={close} />
        </div>
      )
    },
  },

  render() {
    const w = this.words
    const height = this.expanded ? 'calc(100vh - 40px)' : `${this.editorH.size.value}px`
    const btn = (icon: string, title: string, onClick: () => void, extra: Record<string, unknown> = {}) => (
      <v-btn icon={icon} size="small" variant="text" density="comfortable" title={title} onClick={onClick} {...extra} />
    )

    // the bar of the editor: its title, and what is done to the code
    const bar = (
      <div style={{display: 'flex', alignItems: 'center', gap: '2px', padding: '2px 4px', borderBottom: border}} data-gad-code-bar>
        <span style={{flex: '1 1 auto', minWidth: 0}} class="text-caption text-medium-emphasis">
          {this.$slots.title?.()}
        </span>
        {!this.readonly && !this.comparing &&
          btn('mdi-auto-fix', w.format, this.format, {loading: this.formatting, 'data-gad-code-format': true})}
        {!this.readonly && btn('mdi-undo', w.undo, () => this.active()?.undo?.(), {disabled: !this.canUndo, 'data-gad-code-undo': true})}
        {!this.readonly && btn('mdi-redo', w.redo, () => this.active()?.redo?.(), {disabled: !this.canRedo, 'data-gad-code-redo': true})}
        {this.comparing && btn('mdi-arrow-up', w.prevChange, () => this.active()?.prev?.(), {'data-gad-code-prev': true})}
        {this.comparing && btn('mdi-arrow-down', w.nextChange, () => this.active()?.next?.(), {'data-gad-code-next': true})}
        {btn(this.comparing ? 'mdi-code-braces' : 'mdi-compare-horizontal', this.comparing ? w.edit : w.compare,
          () => (this.comparing = !this.comparing), {'data-gad-code-compare': true, color: this.comparing ? 'primary' : undefined})}
        {btn(this.expanded ? 'mdi-arrow-collapse' : 'mdi-arrow-expand', this.expanded ? w.collapse : w.expand,
          () => (this.expanded = !this.expanded), {'data-gad-code-expand': true})}
      </div>
    )

    const editor = this.comparing
      ? h(
          'div',
          {style: {height, overflow: 'auto'}, 'data-gad-code-compare-view': true},
          [
            h(DiffMergeView as never, {
              ref: 'merge',
              old: this.original ?? this.initial,
              current: this.modelValue,
              language: this.language,
              editable: !this.readonly,
              dark: this.dark(),
              revertTitle: w.revert,
              onChange: (v: string) => this.$emit('update:modelValue', v),
              onHistory: this.onHistory,
            }),
          ],
        )
      : h(
          CodeMirror,
          {
            ref: 'editor',
            modelValue: this.modelValue,
            'onUpdate:modelValue': (v: string) => this.$emit('update:modelValue', v),
            onHistory: this.onHistory,
            language: this.language,
            templateStart: this.templateStart,
            templateEnd: this.templateEnd,
            readonly: this.readonly,
            minHeight: height,
            maxHeight: height,
            help: this.help,
            'onUpdate:help': (v: boolean) => this.$emit('update:help', v),
            helpMode: this.helpMode,
            helpTitle: this.helpTitle ?? this.vx.help,
            helpTooltip: this.helpTooltip ?? this.vx.help,
          },
          this.$slots.help ? {help: () => this.$slots.help?.()} : undefined,
        )

    const sW = this.structureW, dW = this.docW, eH = this.editorH
    const left = this.structureOpen ? [
      <div style={{flex: `0 1 ${sW.size.value}px`, display: 'flex', flexDirection: 'column', borderRight: border, minHeight: 0, minWidth: '120px'}}>
        {this.panelHead(w.structure, this.refreshStructure, () => (this.structureOpen = false), 'mdi-chevron-left')}
        <div style={{flex: '1 1 auto', overflow: 'auto', maxHeight: height}} data-gad-code-structure>
          {h(StructureTree, {symbols: this.symbols, emptyText: w.noStructure, onNavigate: this.navigate})}
        </div>
      </div>,
      h(ResizeHandle, {onStart: sW.begin, onDrag: (dx: number) => sW.drag(dx), onReset: sW.reset, 'data-gad-code-resize-structure': true}),
    ] : this.strip('mdi-file-tree-outline', w.structure, () => (this.structureOpen = true), 'left')

    const right = this.docOpen ? [
      h(ResizeHandle, {onStart: dW.begin, onDrag: (dx: number) => dW.drag(dx, -1), onReset: dW.reset, 'data-gad-code-resize-doc': true}),
      <div style={{flex: `0 1 ${dW.size.value}px`, display: 'flex', flexDirection: 'column', borderLeft: border, minHeight: 0, minWidth: '120px'}}>
        {this.panelHead(w.doc, this.refreshDoc, () => (this.docOpen = false), 'mdi-chevron-right')}
        <div style={{flex: '1 1 auto', overflow: 'auto', padding: '8px', maxHeight: height}} data-gad-code-doc>
          {this.docError ? (
            <pre class="text-error text-caption" style={{whiteSpace: 'pre-wrap'}}>{this.docError}</pre>
          ) : this.docHtml ? (
            <div class="vx-gad-code__doc user-doc" innerHTML={this.docHtml} />
          ) : (
            <div class="text-medium-emphasis text-caption">{w.noDoc}</div>
          )}
        </div>
      </div>,
    ] : this.strip('mdi-book-open-variant', w.doc, () => (this.docOpen = true), 'right')

    const body = (
      <div
        class={['vx-gad-code', {'vx-gad-code--expanded': this.expanded}]}
        style={{
          border, borderRadius: '4px', display: 'flex', alignItems: 'stretch', background: 'rgb(var(--v-theme-surface))',
          ...(this.expanded ? {position: 'fixed', inset: '0', zIndex: 2400, borderRadius: '0'} : {}),
        }}
      >
        {left}
        <div style={{flex: '1 1 0', minWidth: '200px', display: 'flex', flexDirection: 'column'}}>
          {bar}
          {editor}
          {this.formatError && (
            <div class="text-error text-caption" style={{padding: '2px 8px', borderTop: border}} data-gad-code-format-error>
              {this.formatError}
            </div>
          )}
        </div>
        {right}
      </div>
    )
    // the height dragged by the bottom edge (over the page, the page's)
    if (this.expanded) return body
    return (
      <div>
        {body}
        {h(ResizeHandle, {direction: 'row', onStart: eH.begin, onDrag: (dy: number) => eH.drag(dy), onReset: eH.reset,
          onPointerup: () => this.remeasure(), 'data-gad-code-resize-height': true})}
      </div>
    )
  },
})
