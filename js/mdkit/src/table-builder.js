// mdkit-table-builder: the table of a Markdown text, built as a sheet — a
// Vue component, a dialog (open(table, done)) —: its cells written in place;
// a range chosen by dragging (or Shift and a click; a column by its letter,
// a row by its number), grown to the merges it touches; rows and columns put
// before and after it, or taken out; its cells merged into one, or one
// merge undone; the first row the header or not; the columns aligned, and
// sized by dragging the border of their letter — in pixels or in parts of
// the table (a double click: automatic) —; the
// text of a cell formatted (Markdown: bold, italic, strikethrough, code, a
// link, a picture — when `image` —, a footnote — when `footnote` —;
// Shift+Enter, a line break); its
// preview, as it is written — when `render` (the server's) —, as the
// editor's: a menu of modes, the sheet alone, side by side — a line between,
// their vertical scrolls together —, the preview alone. Done,
// the table's text (table-model toMarkdown: GFM, else HTML) given to done.
// Its words: `words` (the page's, the config's) and `barWords`
// (the rich text editor's: bold, italic, …).
import * as T from './table-model.js'
import { config, wordsOf } from './config.js'

const letters = (n) => {
  let s = ''
  for (n++; n > 0; n = Math.floor((n - 1) / 26)) s = String.fromCharCode(65 + ((n - 1) % 26)) + s
  return s
}

// style is the sheet's look, once for the page.
function style() {
  if (document.getElementById('mdkit-table-builder-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-table-builder-style'
  st.textContent =
    '.htb-grid{overflow:auto;max-height:60vh;border:1px solid #d1d5db;border-radius:4px;user-select:none}' +
    '.htb-grid table{border-collapse:collapse;min-width:100%}' +
    '.htb-grid td,.htb-grid th{border:1px solid #d1d5db;padding:0;vertical-align:top;min-width:3rem}' +
    '.htb-grid td{min-width:7rem}' +
    '.htb-grid table.htb-sized td{min-width:3rem}' +
    // the handle that sizes a column: the right border of its letter
    '.htb-grid th.htb-col{position:relative}' +
    '.htb-resize{position:absolute;top:0;right:-4px;width:8px;height:100%;cursor:col-resize;z-index:1}' +
    '.htb-resize:hover,.htb-resizing .htb-resize{background:#3b82f6;opacity:.5}' +
    '.htb-width{color:#2563eb;margin-left:.35rem}' +
    '.htb-grid th.htb-ref{background:#f3f4f6;font-weight:500;font-size:.75rem;color:#6b7280;min-width:2rem;padding:.15rem .4rem;cursor:pointer;text-align:center}' +
    '.htb-grid td.htb-head{background:#f9fafb}' +
    '.htb-grid td.htb-head textarea{font-weight:600}' +
    '.htb-grid td.htb-sel{background:#dbeafe}' +
    // side by side: the preview by the sheet, a vertical line between
    '.htb-split{border-left:1px solid #d1d5db;margin-left:.75rem;padding-top:0!important}' +
    '.htb-grid textarea{display:block;width:100%;border:0;outline:0;resize:none;background:transparent;padding:.35rem .5rem;font:inherit;line-height:1.4;min-height:2.1rem;user-select:text}'
  document.head.append(st)
}

export default {
  name: 'MdkitTableBuilder',
  props: {
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('table') },
    // the words of the buttons of the text (bold, italic, …)
    barWords: { type: Object, default: () => wordsOf('rich') },
    // image, when pictures may be put in: the Markdown of pictures the user
    // chooses (a promise; '' none) — the Markdown editor's
    image: { type: Function, default: null },
    // footnote, when notes may be put in: the label of a note chosen or made
    // (a promise) — the Markdown editor's footnotes builder
    footnote: { type: Function, default: null },
    // render, when there is a preview: the HTML of a text (a promise)
    render: { type: Function, default: null },
  },
  data: () => ({
    show: false, editing: false, t: T.create(), sel: { r1: 0, c1: 0, r2: 0, c2: 0 },
    anchor: { r: 0, c: 0 }, dragging: false, done: null, cell: null,
    mode: 'sheet', previewHTML: '', following: null, unit: 'px', resizing: null,
  }),
  computed: {
    // the class of a text shown (config: the page's)
    textClass: () => config.textClass,
    // the buttons of the text of a cell
    formats() {
      return [
        { key: 'bold', icon: 'mdi-format-bold' }, { key: 'italic', icon: 'mdi-format-italic' },
        { key: 'strike', icon: 'mdi-format-strikethrough' }, { key: 'code', icon: 'mdi-code-tags' },
        { key: 'link', icon: 'mdi-link-variant' },
        ...(this.image ? [{ key: 'image', icon: 'mdi-image-plus' }] : []),
        ...(this.footnote ? [{ key: 'footnote', icon: 'mdi-format-superscript' }] : []),
      ]
    },
    // the modes of the preview
    modes() {
      const w = this.words
      return [
        { key: 'sheet', icon: 'mdi-table-edit', title: w.modeSheet || 'Sheet' },
        { key: 'split', icon: 'mdi-view-split-vertical', title: w.modeSplit || 'Side by side' },
        { key: 'total', icon: 'mdi-eye-outline', title: w.modePreview || 'Preview only' },
      ]
    },
    html() { return !T.gfm(this.t) },
    // sized: a column has its width
    sized() { return this.t.widths.some((w) => w) },
    // the sheet as wide as the parts of its columns ask (100%), else as its
    // columns
    sheetStyle() {
      if (T.percent(this.t)) return 'width:100%'
      return this.sized ? 'width:auto' : 'min-width:100%'
    },
    multi() { return this.sel.r1 !== this.sel.r2 || this.sel.c1 !== this.sel.c2 },
    merged() { return this.t.merges.some((m) => m.r >= this.sel.r1 && m.r <= this.sel.r2 && m.c >= this.sel.c1 && m.c <= this.sel.c2) },
    align() { return this.t.aligns[this.sel.c1] || 'left' },
  },
  template: `
<v-dialog v-model="show" :max-width="mode === 'split' ? 1400 : 960" scrollable data-table-builder>
  <v-card>
    <v-card-title class="d-flex align-center">
      <v-icon icon="mdi-table" class="me-2"></v-icon>{{ words.title || 'Table' }}
      <v-spacer></v-spacer>
      <span class="text-caption text-medium-emphasis">{{ t.rows }} × {{ t.cols }}</span>
    </v-card-title>
    <v-toolbar density="compact" color="grey-lighten-4" class="px-2 flex-wrap" style="height:auto">
      <v-btn icon="mdi-table-row-plus-before" size="small" variant="text" :title="words.rowAbove || 'Row above'" data-table-op="row-above" @click="op('rowAbove')"></v-btn>
      <v-btn icon="mdi-table-row-plus-after" size="small" variant="text" :title="words.rowBelow || 'Row below'" data-table-op="row-below" @click="op('rowBelow')"></v-btn>
      <v-btn icon="mdi-table-row-remove" size="small" variant="text" :title="words.rowRemove || 'Remove rows'" :disabled="t.rows <= 1" data-table-op="row-remove" @click="op('rowRemove')"></v-btn>
      <v-divider vertical class="mx-1"></v-divider>
      <v-btn icon="mdi-table-column-plus-before" size="small" variant="text" :title="words.colLeft || 'Column left'" data-table-op="col-left" @click="op('colLeft')"></v-btn>
      <v-btn icon="mdi-table-column-plus-after" size="small" variant="text" :title="words.colRight || 'Column right'" data-table-op="col-right" @click="op('colRight')"></v-btn>
      <v-btn icon="mdi-table-column-remove" size="small" variant="text" :title="words.colRemove || 'Remove columns'" :disabled="t.cols <= 1" data-table-op="col-remove" @click="op('colRemove')"></v-btn>
      <v-divider vertical class="mx-1"></v-divider>
      <v-btn icon="mdi-table-merge-cells" size="small" variant="text" :title="words.merge || 'Merge cells'" :disabled="!multi" data-table-op="merge" @click="op('merge')"></v-btn>
      <v-btn icon="mdi-table-split-cell" size="small" variant="text" :title="words.unmerge || 'Unmerge cells'" :disabled="!merged" data-table-op="unmerge" @click="op('unmerge')"></v-btn>
      <v-divider vertical class="mx-1"></v-divider>
      <v-btn v-for="f in formats" :key="f.key" :icon="f.icon" size="small" variant="text" :title="(f.key === 'footnote' ? words.footnote : barWords[f.key]) || f.key"
        :disabled="!cell" :data-table-format="f.key" @mousedown.prevent @click="format(f.key)"></v-btn>
      <v-divider vertical class="mx-1"></v-divider>
      <v-btn-toggle :model-value="align" density="compact" variant="text" mandatory @update:model-value="setAlign">
        <v-btn value="left" icon="mdi-format-align-left" size="small" :title="words.alignLeft || 'Align left'"></v-btn>
        <v-btn value="center" icon="mdi-format-align-center" size="small" :title="words.alignCenter || 'Center'"></v-btn>
        <v-btn value="right" icon="mdi-format-align-right" size="small" :title="words.alignRight || 'Align right'"></v-btn>
      </v-btn-toggle>
      <v-divider vertical class="mx-1"></v-divider>
      <v-menu v-if="render" location="bottom">
        <template #activator="{ props }">
          <v-btn v-bind="props" :icon="modes.find((m) => m.key === mode).icon" size="small" variant="text"
            :active="mode !== 'sheet'" :title="words.preview || 'Preview'" data-table-op="preview"></v-btn>
        </template>
        <v-list density="compact">
          <v-list-item v-for="m in modes" :key="m.key" :prepend-icon="m.icon" :title="m.title" :active="mode === m.key"
            :data-table-mode="m.key" @click="setMode(m.key)"></v-list-item>
        </v-list>
      </v-menu>
      <v-divider vertical class="mx-1"></v-divider>
      <v-btn-toggle :model-value="unit" density="compact" variant="text" mandatory @update:model-value="setUnit" data-table-unit>
        <v-btn value="px" size="small" :title="words.widthFixed || 'Widths in pixels'">px</v-btn>
        <v-btn value="%" size="small" :title="words.widthPercent || 'Widths in parts of the table'">%</v-btn>
      </v-btn-toggle>
      <v-btn icon="mdi-arrow-expand-horizontal" size="small" variant="text" :disabled="!sized"
        :title="words.widthAuto || 'Automatic widths'" data-table-op="width-auto" @click="t.widths = t.widths.map(() => '')"></v-btn>
      <v-divider vertical class="mx-1"></v-divider>
      <v-switch v-model="t.header" color="primary" density="compact" hide-details class="flex-grow-0 ms-2"
        :label="words.header || 'Header row'" data-table-header></v-switch>
    </v-toolbar>
    <v-card-text class="pt-3">
      <div class="d-flex" :class="mode === 'split' ? '' : 'flex-column ga-3'">
      <div v-show="mode !== 'total'" ref="grid" class="htb-grid flex-grow-1" :style="mode === 'split' ? 'flex-basis:0;min-width:0' : ''"
        @scroll="follow('grid')" @mouseup="dragging = false" @mouseleave="dragging = false"
        @dragstart.prevent @drop.prevent data-table-grid>
        <table ref="sheet" :class="{ 'htb-sized': sized }" :style="sheetStyle">
          <colgroup>
            <col style="width:2.5rem">
            <col v-for="c in t.cols" :key="c" :style="colStyle(c - 1)">
          </colgroup>
          <tr>
            <th class="htb-ref"></th>
            <th v-for="c in t.cols" :key="c" class="htb-ref htb-col" :class="{ 'htb-resizing': resizing === c - 1 }"
              :data-col="c - 1" @mousedown.prevent="pick(0, c - 1, t.rows - 1, c - 1)">
              {{ letter(c - 1) }}<span v-if="t.widths[c - 1]" class="htb-width">{{ widthLabel(t.widths[c - 1]) }}</span>
              <div class="htb-resize" :title="words.resize || 'Drag to size the column; a double click: automatic'"
                data-table-resize @mousedown.prevent.stop="startResize($event, c - 1)" @dblclick.stop="t.widths[c - 1] = ''"></div>
            </th>
          </tr>
          <tr v-for="(row, r) in t.cells" :key="r">
            <th class="htb-ref" @mousedown.prevent="pick(r, 0, r, t.cols - 1)">{{ r + 1 }}</th>
            <template v-for="(cell, c) in row" :key="c">
              <td v-if="!covered(r, c)" :colspan="spanOf(r, c).cs" :rowspan="spanOf(r, c).rs"
                :class="{ 'htb-sel': inSel(r, c), 'htb-head': t.header && r === 0 }"
                @mousedown="down($event, r, c)" @mouseenter="enter(r, c)" :data-cell="r + ',' + c">
                <textarea rows="1" v-model="t.cells[r][c]" :style="{ textAlign: t.aligns[c] || 'left' }"
                  @focus="focus(r, c, $event.target)" @input="grow($event.target)"
                  @keydown.enter.exact.prevent="move(r, c, 1)"
                  @keydown.ctrl.b.prevent="format('bold')" @keydown.meta.b.prevent="format('bold')"
                  @keydown.ctrl.i.prevent="format('italic')" @keydown.meta.i.prevent="format('italic')"></textarea>
              </td>
            </template>
          </tr>
        </table>
      </div>
      <div v-if="mode !== 'sheet'" class="flex-grow-1" :style="mode === 'split' ? 'flex-basis:0;min-width:0' : ''" data-table-preview>
        <div v-if="mode === 'total'" class="text-caption text-medium-emphasis mb-1">{{ words.preview || 'Preview' }}</div>
        <div ref="preview" class="pa-3" :class="[textClass, mode === 'split' ? 'htb-split' : 'border rounded']"
          style="overflow:auto;max-height:60vh" @scroll="follow('preview')" v-html="previewHTML"></div>
      </div>
      </div>
      <div v-if="html" class="text-caption text-medium-emphasis mt-2" data-table-html>
        <v-icon icon="mdi-information-outline" size="small"></v-icon>
        {{ words.htmlNote || 'With merged cells, or no header row, the table is written as HTML.' }}
      </div>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ words.cancel || 'Cancel' }}</v-btn>
      <v-btn color="primary" variant="flat" data-table-done @click="finish">{{ editing ? (words.update || 'Update') : (words.insert || 'Insert') }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>`,
  watch: {
    // the preview, a moment after the last change
    t: {
      deep: true,
      handler() { if (this.mode !== 'sheet') this.schedule() },
    },
  },
  mounted() {
    style()
  },
  methods: {
    // open shows the builder with table (a copy; none, a new one of 3 × 3),
    // done given the text of the table when it is finished
    open(table, done) {
      this.editing = !!table
      this.t = T.normalize(table ? JSON.parse(JSON.stringify(table)) : T.create(3, 3))
      this.unit = T.percent(this.t) ? '%' : 'px'
      this.done = done
      this.cell = null
      this.pick(0, 0, 0, 0)
      this.show = true
      // the boxes as tall as their texts, once drawn
      setTimeout(() => document.querySelectorAll('[data-table-grid] textarea').forEach((el) => this.grow(el)), 50)
    },
    letter: letters,
    covered(r, c) { return T.covered(this.t, r, c) },
    spanOf(r, c) { return T.mergeAt(this.t, r, c) || { rs: 1, cs: 1 } },
    inSel(r, c) { return r >= this.sel.r1 && r <= this.sel.r2 && c >= this.sel.c1 && c <= this.sel.c2 },
    // pick chooses the range, grown to the merges it touches
    pick(r1, c1, r2, c2) {
      this.anchor = { r: r1, c: c1 }
      this.sel = T.span(this.t, r1, c1, r2, c2)
    },
    down(e, r, c) {
      if (e.shiftKey) {
        e.preventDefault()
        this.sel = T.span(this.t, this.anchor.r, this.anchor.c, r, c)
        return
      }
      this.pick(r, c, r, c)
      this.dragging = true
    },
    // enter grows the range being dragged; out of its first cell, a range of
    // cells (not of text: the box left)
    enter(r, c) {
      if (!this.dragging) return
      this.sel = T.span(this.t, this.anchor.r, this.anchor.c, r, c)
      if ((r !== this.anchor.r || c !== this.anchor.c) && document.activeElement && document.activeElement.closest('[data-table-grid]')) {
        document.activeElement.blur()
        const ws = window.getSelection()
        if (ws) ws.removeAllRanges()
      }
    },
    focus(r, c, el) {
      this.cell = { r, c, el }
      if (!this.dragging && !this.inSel(r, c)) this.pick(r, c, r, c)
    },
    // grow makes the box of a cell as tall as its text
    grow(el) {
      el.style.height = 'auto'
      el.style.height = el.scrollHeight + 'px'
    },
    // format puts the marks of key (Markdown) around the text chosen in the
    // cell being written — or takes them off —, the text chosen after; a
    // link, [text](https://), its address chosen
    async format(key) {
      if (!this.cell) return
      const { r, c, el } = this.cell
      let a = el.selectionStart
      let b = el.selectionEnd
      if (key === 'footnote') {
        // a reference to a note, at the caret of the cell
        const label = this.footnote ? await this.footnote() : ''
        if (!label) return
        const v = this.t.cells[r][c] || ''
        const ref = '[^' + label + ']'
        this.t.cells[r][c] = v.slice(0, a) + ref + v.slice(b)
        this.$nextTick(() => {
          el.focus()
          el.setSelectionRange(a + ref.length, a + ref.length)
        })
        return
      }
      if (key === 'image') {
        // the pictures, at the caret of the cell
        const md = this.image ? await this.image() : ''
        if (!md) return
        const v = this.t.cells[r][c] || ''
        this.t.cells[r][c] = v.slice(0, a) + md + v.slice(b)
        this.$nextTick(() => {
          el.focus()
          el.setSelectionRange(a + md.length, a + md.length)
          this.grow(el)
        })
        return
      }
      const v = this.t.cells[r][c] || ''
      const marks = { bold: '**', italic: '_', strike: '~~', code: '`' }
      let next
      let from
      let to
      if (key === 'link') {
        const label = v.slice(a, b) || (this.barWords.link || 'link')
        next = v.slice(0, a) + '[' + label + '](https://)' + v.slice(b)
        from = a + label.length + 3
        to = from + 8
      } else {
        const m = marks[key]
        if (v.slice(a - m.length, a) === m && v.slice(b, b + m.length) === m) {
          next = v.slice(0, a - m.length) + v.slice(a, b) + v.slice(b + m.length)
          from = a - m.length
          to = b - m.length
        } else {
          next = v.slice(0, a) + m + v.slice(a, b) + m + v.slice(b)
          from = a + m.length
          to = b + m.length
        }
      }
      this.t.cells[r][c] = next
      this.$nextTick(() => {
        el.focus()
        el.setSelectionRange(from, to)
        this.grow(el)
      })
    },
    // move goes to the cell dr rows below (Enter)
    move(r, c, dr) {
      const to = Math.min(this.t.rows - 1, r + this.spanOf(r, c).rs * dr)
      const el = this.$el && document.querySelector(`[data-table-builder] [data-cell="${to},${c}"] textarea`)
      if (el) el.focus()
    },
    op(name) {
      const t = this.t
      const s = this.sel
      switch (name) {
        case 'rowAbove': T.insertRow(t, s.r1); this.pick(s.r1, s.c1, s.r1, s.c2); break
        case 'rowBelow': T.insertRow(t, s.r2 + 1); this.pick(s.r2 + 1, s.c1, s.r2 + 1, s.c2); break
        case 'rowRemove':
          for (let r = s.r2; r >= s.r1 && t.rows > 1; r--) T.removeRow(t, r)
          this.pick(Math.min(s.r1, t.rows - 1), 0, Math.min(s.r1, t.rows - 1), 0)
          break
        case 'colLeft': T.insertCol(t, s.c1); this.pick(s.r1, s.c1, s.r2, s.c1); break
        case 'colRight': T.insertCol(t, s.c2 + 1); this.pick(s.r1, s.c2 + 1, s.r2, s.c2 + 1); break
        case 'colRemove':
          for (let c = s.c2; c >= s.c1 && t.cols > 1; c--) T.removeCol(t, c)
          this.pick(0, Math.min(s.c1, t.cols - 1), 0, Math.min(s.c1, t.cols - 1))
          break
        case 'merge': if (T.merge(t, s.r1, s.c1, s.r2, s.c2)) this.pick(s.r1, s.c1, s.r1, s.c1); break
        case 'unmerge':
          for (const m of t.merges.filter((m) => m.r >= s.r1 && m.r <= s.r2 && m.c >= s.c1 && m.c <= s.c2)) T.unmerge(t, m.r, m.c)
          break
      }
    },
    colStyle(c) {
      const w = this.t.widths[c]
      return w ? 'width:' + (/%$/.test(w) ? w : w + 'px') : ''
    },
    widthLabel: (w) => (/%$/.test(w) ? w : w + 'px'),
    // dataWidth is the width of the columns of the sheet (not of its numbers)
    dataWidth() {
      const sheet = this.$refs.sheet
      const nums = sheet && sheet.querySelector('th.htb-ref')
      return sheet ? sheet.offsetWidth - (nums ? nums.offsetWidth : 0) : 1
    },
    // startResize sizes the column c as the pointer goes: its width the
    // pixels, or the part of the sheet's (unit)
    startResize(e, c) {
      const th = e.target.closest('th')
      const x0 = e.clientX
      const w0 = th.offsetWidth
      const total = this.dataWidth()
      this.resizing = c
      const move = (ev) => {
        const w = Math.max(32, w0 + ev.clientX - x0)
        this.t.widths[c] = this.unit === '%' ? Math.max(1, Math.min(100, Math.round((w / total) * 100))) + '%' : String(Math.round(w))
      }
      const up = () => {
        this.resizing = null
        window.removeEventListener('mousemove', move)
        window.removeEventListener('mouseup', up)
      }
      window.addEventListener('mousemove', move)
      window.addEventListener('mouseup', up)
    },
    // setUnit makes the widths there are of unit: as they are drawn now
    setUnit(unit) {
      if (unit === this.unit) return
      const total = this.dataWidth()
      const ths = this.$refs.sheet ? this.$refs.sheet.querySelectorAll('th.htb-col') : []
      this.t.widths = this.t.widths.map((w, c) => {
        if (!w || !ths[c]) return w
        const px = ths[c].offsetWidth
        return unit === '%' ? Math.max(1, Math.round((px / total) * 100)) + '%' : String(px)
      })
      this.unit = unit
    },
    // follow scrolls the other side as from scrolled, side by side: the same
    // part of its height (one moved, not both by each other)
    follow(from) {
      if (this.mode !== 'split' || this.following === from) return
      const a = from === 'grid' ? this.$refs.grid : this.$refs.preview
      const b = from === 'grid' ? this.$refs.preview : this.$refs.grid
      if (!a || !b) return
      const max = a.scrollHeight - a.clientHeight
      this.following = from === 'grid' ? 'preview' : 'grid'
      b.scrollTop = max > 0 ? (a.scrollTop / max) * (b.scrollHeight - b.clientHeight) : 0
      requestAnimationFrame(() => { this.following = null })
    },
    // setMode shows the sheet, the preview, or both side by side
    setMode(m) {
      this.mode = m
      if (m !== 'sheet') this.schedule(0)
    },
    // schedule renders the table's text (the server's), after wait ms
    schedule(wait = 250) {
      clearTimeout(this.timer)
      this.timer = setTimeout(async () => {
        if (this.render) this.previewHTML = await this.render(T.toMarkdown(this.t))
      }, wait)
    },
    setAlign(a) {
      for (let c = this.sel.c1; c <= this.sel.c2; c++) this.t.aligns[c] = a === 'left' ? '' : a
    },
    finish() {
      const text = T.toMarkdown(this.t)
      this.show = false
      if (this.done) this.done(text)
    },
  },
}
