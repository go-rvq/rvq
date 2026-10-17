// <vx-gad-code>: the admin's editor of gad code (gad, gadx, gadt) — the
// editor (<vx-codemirror>) with what an IDE gives of the code in it:
//
//   - Format: the code as `gad fmt` writes it;
//   - Expand: the editor over the whole page, and back (Esc too);
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
import {StructureTree, type OutlineSym} from '@gad-lang/ide-vuetify/structure'
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
  minHeight: {type: String, default: '16rem'},
  maxHeight: {type: String, default: '40rem'},
  // the help of the editor (its slot), as <vx-codemirror>'s
  help: Boolean,
  helpMode: {type: String, default: 'dialog'},
  helpTitle: {type: String, default: 'Ajuda'},
  helpTooltip: {type: String, default: 'Ajuda'},
  // the words of its buttons and panels
  labels: {
    type: Object as PropType<Partial<Record<'format' | 'expand' | 'collapse' | 'structure' | 'doc' | 'refresh' | 'noStructure' | 'noDoc' | 'notFormatted', string>>>,
    default: () => ({}),
  },
} as const

const DEFAULT_LABELS = {
  format: 'Formatar',
  expand: 'Expandir',
  collapse: 'Recolher',
  structure: 'Estrutura',
  doc: 'Documentação',
  refresh: 'Atualizar',
  noStructure: 'Nenhuma declaração.',
  noDoc: 'Nenhuma documentação.',
  notFormatted: 'O código não pôde ser formatado: corrija o que está errado nele.',
}

const border = '1px solid rgba(var(--v-border-color), var(--v-border-opacity))'

export default defineComponent({
  name: 'VXGadCode',

  props: propsOptions,

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
      timer: 0 as number | ReturnType<typeof setTimeout>,
    }
  },

  computed: {
    words(): typeof DEFAULT_LABELS {
      return {...DEFAULT_LABELS, ...this.labels}
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
      ;(this.$refs.editor as {goto?: (line: number, column: number) => void} | undefined)?.goto?.(sym.line, sym.column)
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
    const height = this.expanded ? 'calc(100vh - 56px)' : undefined

    const toolbar = (
      <div style={{display: 'flex', alignItems: 'center', gap: '2px', padding: '2px 4px', borderBottom: border}}>
        <span style={{flex: '1 1 auto'}} class="text-caption text-medium-emphasis">
          {this.$slots.title?.()}
        </span>
        {!this.readonly && (
          <v-btn icon="mdi-auto-fix" size="small" variant="text" density="comfortable" title={w.format}
            loading={this.formatting} onClick={this.format} data-gad-code-format />
        )}
        <v-btn icon={this.expanded ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'} size="small" variant="text"
          density="comfortable" title={this.expanded ? w.collapse : w.expand}
          onClick={() => (this.expanded = !this.expanded)} data-gad-code-expand />
      </div>
    )

    const editor = h(
      CodeMirror,
      {
        ref: 'editor',
        modelValue: this.modelValue,
        'onUpdate:modelValue': (v: string) => this.$emit('update:modelValue', v),
        language: this.language,
        templateStart: this.templateStart,
        templateEnd: this.templateEnd,
        readonly: this.readonly,
        minHeight: height || this.minHeight,
        maxHeight: height || this.maxHeight,
        help: this.help,
        'onUpdate:help': (v: boolean) => this.$emit('update:help', v),
        helpMode: this.helpMode,
        helpTitle: this.helpTitle,
        helpTooltip: this.helpTooltip,
      },
      this.$slots.help ? {help: () => this.$slots.help?.()} : undefined,
    )

    const left = this.structureOpen ? (
      <div style={{flex: '0 0 260px', maxWidth: '260px', display: 'flex', flexDirection: 'column', borderRight: border, minHeight: 0}}>
        {this.panelHead(w.structure, this.refreshStructure, () => (this.structureOpen = false), 'mdi-chevron-left')}
        <div style={{flex: '1 1 auto', overflow: 'auto', maxHeight: height || this.maxHeight}} data-gad-code-structure>
          {h(StructureTree, {symbols: this.symbols, emptyText: w.noStructure, onNavigate: this.navigate})}
        </div>
      </div>
    ) : this.strip('mdi-file-tree-outline', w.structure, () => (this.structureOpen = true), 'left')

    const right = this.docOpen ? (
      <div style={{flex: '0 0 40%', maxWidth: '40%', display: 'flex', flexDirection: 'column', borderLeft: border, minHeight: 0}}>
        {this.panelHead(w.doc, this.refreshDoc, () => (this.docOpen = false), 'mdi-chevron-right')}
        <div style={{flex: '1 1 auto', overflow: 'auto', padding: '8px', maxHeight: height || this.maxHeight}} data-gad-code-doc>
          {this.docError ? (
            <pre class="text-error text-caption" style={{whiteSpace: 'pre-wrap'}}>{this.docError}</pre>
          ) : this.docHtml ? (
            <div class="vx-gad-code__doc user-doc" innerHTML={this.docHtml} />
          ) : (
            <div class="text-medium-emphasis text-caption">{w.noDoc}</div>
          )}
        </div>
      </div>
    ) : this.strip('mdi-book-open-variant', w.doc, () => (this.docOpen = true), 'right')

    return (
      <div
        class={['vx-gad-code', {'vx-gad-code--expanded': this.expanded}]}
        style={{
          border, borderRadius: '4px', display: 'flex', flexDirection: 'column', background: 'rgb(var(--v-theme-surface))',
          ...(this.expanded ? {position: 'fixed', inset: '0', zIndex: 2400, borderRadius: '0'} : {}),
        }}
      >
        {toolbar}
        <div style={{display: 'flex', alignItems: 'stretch', flex: '1 1 auto', minHeight: 0}}>
          {left}
          <div style={{flex: '1 1 auto', minWidth: 0}}>{editor}</div>
          {right}
        </div>
        {this.formatError && (
          <div class="text-error text-caption" style={{padding: '2px 8px', borderTop: border}} data-gad-code-format-error>
            {this.formatError}
          </div>
        )}
      </div>
    )
  },
})
