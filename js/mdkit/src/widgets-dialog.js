// The dialogs of the buttons and of the tabs of a text — Vue components
// (open(value, done)): done given the value made ({align, items}), or null
// to take it out (one there) — (widgets-model.js):
//   - mdkit-buttons-dialog: a list of buttons, each its link (an address, a
//     page of the site or a record of the admin: the link picker's; the
//     title it has, the one the page shows when it is given none), title,
//     description (its tooltip), color of the theme, icons left and right;
//     the line's alignment; a preview;
//   - mdkit-tabs-dialog: tabs, each its title, color, icons and content —
//     written in the editor given (the text's own: `editor`, `editorProps`,
//     its v-model the content) —; the bar's alignment.
import { wordsOf } from './config.js'
import { COLORS, buttonsHTML } from './widgets-model.js'
import { cardTarget } from './link-cards.js'
import { drawButtons, widgetStyle } from './widgets.js'
import { vendorOf } from './blocks.js'
import LinkPicker from './link-picker.js'

const WORDS = {
  buttons: 'Buttons', tabs: 'Tabs', align: 'Alignment', start: 'Left', center: 'Center', end: 'Right', stretch: 'Fill the line',
  add: 'Add', addTab: 'Add a tab', link: 'Link', choose: 'Choose…', title: 'Title', titleHint: 'Empty: the title of the page',
  description: 'Description (tooltip)', color: 'Color', iconLeft: 'Icon on the left', iconRight: 'Icon on the right',
  iconHint: 'A Material Design icon: mdi-…', icons: 'The icons', up: 'Up', down: 'Down', left: 'To the left', right: 'To the right',
  remove: 'Remove', removeAll: 'Remove', cancel: 'Cancel', insert: 'Insert', update: 'Update', preview: 'Preview',
  newTab: 'In a new tab', tab: 'Tab', content: 'Content', noLink: 'No link chosen',
  default: 'Default', primary: 'Primary', secondary: 'Secondary', success: 'Success', info: 'Information', warning: 'Warning', error: 'Error',
}

const ALIGN_ICONS = { start: 'mdi-format-align-left', center: 'mdi-format-align-center', end: 'mdi-format-align-right', stretch: 'mdi-format-align-justify' }

const wordsMixin = {
  props: { words: { type: Object, default: () => wordsOf('widgets') } },
  computed: {
    w() { return { ...WORDS, ...this.words } },
    colors() { return COLORS.map((c) => ({ value: c, title: this.w[c || 'default'] || c || 'Default' })) },
    aligns() { return Object.keys(ALIGN_ICONS).map((a) => ({ value: a, icon: ALIGN_ICONS[a], title: this.w[a] })) },
  },
}

// lookFields are the fields of an item's look: its color, its icons.
const lookFields = `
<div class="d-flex flex-wrap ga-3 mt-3">
  <v-select v-model="it.color" :items="colors" item-title="title" item-value="value" :label="w.color" density="compact"
    variant="outlined" hide-details style="flex:1 1 10rem" data-widget-color>
    <template #selection="{ item }"><span class="mdkit-swatch" :data-color="item.value || undefined"></span>{{ item.title }}</template>
    <template #item="{ props, item }"><v-list-item v-bind="props"><template #prepend><span class="mdkit-swatch" :data-color="item.value || undefined"></span></template></v-list-item></template>
  </v-select>
  <v-text-field v-model.trim="it.iconLeft" :label="w.iconLeft" :prepend-inner-icon="validIcon(it.iconLeft) || 'mdi-shape-outline'"
    placeholder="mdi-…" density="compact" variant="outlined" hide-details clearable style="flex:1 1 10rem" data-widget-icon-left></v-text-field>
  <v-text-field v-model.trim="it.iconRight" :label="w.iconRight" :append-inner-icon="validIcon(it.iconRight) || 'mdi-shape-outline'"
    placeholder="mdi-…" density="compact" variant="outlined" hide-details clearable style="flex:1 1 10rem" data-widget-icon-right></v-text-field>
</div>`

const alignField = `
<div class="d-flex align-center ga-3 flex-wrap">
  <span class="text-body-2">{{ w.align }}</span>
  <v-btn-toggle v-model="align" mandatory density="compact" variant="outlined" divided data-widget-align>
    <v-btn v-for="a in aligns" :key="a.value" :value="a.value" size="small" :title="a.title" :data-align="a.value"><v-icon :icon="a.icon"></v-icon></v-btn>
  </v-btn-toggle>
  <v-spacer></v-spacer>
  <a href="https://pictogrammers.com/library/mdi/" target="_blank" rel="noopener" class="text-caption">{{ w.icons }} <v-icon icon="mdi-open-in-new" size="x-small"></v-icon></a>
</div>`

const methods = {
  validIcon: (i) => (/^mdi-[a-z0-9-]+$/.test(i || '') ? i : ''),
  move(list, i, d) {
    const j = i + d
    if (j < 0 || j >= list.length) return
    const [x] = list.splice(i, 1)
    list.splice(j, 0, x)
  },
  // clean is an item's look, only what is valid
  clean(it) {
    return { ...it, color: COLORS.includes(it.color) ? it.color : '', iconLeft: this.validIcon(it.iconLeft), iconRight: this.validIcon(it.iconRight) }
  },
}

// swatchStyle is the look of the color swatches of the selects.
function swatchStyle() {
  widgetStyle()
  if (typeof document === 'undefined' || document.getElementById('mdkit-swatch-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-swatch-style'
  st.textContent = '.mdkit-swatch{display:inline-block;width:.9rem;height:.9rem;border-radius:50%;margin-right:.5rem;vertical-align:middle;' +
    'background:rgb(var(--mdkit-c,200,200,200));border:1px solid rgba(127,127,127,.5)}'
  document.head.append(st)
}

const newButton = () => ({ href: '', target: '', title: '', description: '', color: '', iconLeft: '', iconRight: '', newTab: false, hint: '' })

export const ButtonsDialog = {
  name: 'MdkitButtonsDialog',
  mixins: [wordsMixin],
  components: { MdkitLinkPicker: LinkPicker },
  props: { linkWords: { type: Object, default: () => wordsOf('link') } },
  data: () => ({ show: false, align: 'start', items: [], there: false, done: null }),
  computed: {
    value() { return { align: this.align, items: this.items.map((it) => this.clean(it)) } },
    previewHTML() { return buttonsHTML(this.value) },
  },
  watch: {
    previewHTML() { this.$nextTick(() => this.drawPreview()) },
  },
  template: `
<v-dialog v-model="show" max-width="860" scrollable data-buttons-dialog>
  <v-card>
    <v-card-title class="d-flex align-center ga-2"><v-icon icon="mdi-gesture-tap-button"></v-icon>{{ w.buttons }}</v-card-title>
    <v-divider></v-divider>
    <v-card-text style="max-height:70vh">
      ${alignField}
      <v-card v-for="(it, i) in items" :key="i" variant="outlined" class="mt-4 pa-3" data-widget-item>
        <div class="d-flex align-center ga-2">
          <v-icon icon="mdi-link-variant" size="small"></v-icon>
          <span class="text-body-2 text-truncate" style="flex:1 1 auto;min-width:0" data-widget-href>{{ it.href || w.noLink }}</span>
          <v-btn size="small" variant="tonal" color="primary" prepend-icon="mdi-magnify" @click="pick(i)" data-widget-pick>{{ w.choose }}</v-btn>
          <v-btn icon="mdi-arrow-up" size="small" variant="text" :disabled="i === 0" :title="w.up" @click="move(items, i, -1)"></v-btn>
          <v-btn icon="mdi-arrow-down" size="small" variant="text" :disabled="i === items.length - 1" :title="w.down" @click="move(items, i, 1)"></v-btn>
          <v-btn icon="mdi-delete-outline" size="small" variant="text" color="error" :title="w.remove" @click="items.splice(i, 1)" data-widget-remove></v-btn>
        </div>
        <div class="d-flex flex-wrap ga-3 mt-3">
          <v-text-field v-model="it.title" :label="w.title" :placeholder="it.hint || w.titleHint" persistent-placeholder density="compact"
            variant="outlined" hide-details style="flex:1 1 14rem" data-widget-title></v-text-field>
          <v-text-field v-model="it.description" :label="w.description" density="compact" variant="outlined" hide-details
            style="flex:2 1 18rem" data-widget-description></v-text-field>
        </div>
        ${lookFields}
        <v-switch v-model="it.newTab" :label="w.newTab" color="primary" density="compact" hide-details class="mt-1" data-widget-newtab></v-switch>
      </v-card>
      <v-btn class="mt-4" variant="tonal" prepend-icon="mdi-plus" @click="add" data-widget-add>{{ w.add }}</v-btn>
      <div class="text-caption text-medium-emphasis mt-4">{{ w.preview }}</div>
      <div ref="preview" class="pa-3 border rounded" data-widget-preview></div>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-actions>
      <v-btn v-if="there" color="error" variant="text" @click="finish(null)" data-widget-delete>{{ w.removeAll }}</v-btn>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ w.cancel }}</v-btn>
      <v-btn color="primary" variant="flat" :disabled="!value.items.some((it) => it.href)" @click="finish(value)" data-widget-apply>
        {{ there ? w.update : w.insert }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
<mdkit-link-picker ref="links" :words="linkWords"></mdkit-link-picker>`,
  methods: {
    ...methods,
    // open shows the buttons of value (none: one to make), done given the
    // value made or null
    open(value, done) {
      swatchStyle()
      this.align = (value && value.align) || 'start'
      this.items = value && value.items && value.items.length ? value.items.map((it) => ({ ...newButton(), ...it })) : [newButton()]
      this.there = !!(value && value.items && value.items.length)
      this.done = done
      this.show = true
      if (!value || !value.items || !value.items.length) this.$nextTick(() => this.pick(0))
      else this.$nextTick(() => this.drawPreview())
    },
    add() {
      this.items.push(newButton())
      this.pick(this.items.length - 1)
    },
    // pick chooses the link of item i (the link picker): its address, its
    // target (a card's: its title shown when it has none), its new tab
    pick(i) {
      const it = this.items[i]
      this.$refs.links.open({ href: it.href, target: it.newTab ? '_blank' : '' }, (link) => {
        if (!link) return
        it.href = link.href
        it.target = cardTarget(link.href, typeof location !== 'undefined' ? location.origin : '') || ''
        it.newTab = link.target === '_blank'
        it.hint = link.title || ''
      })
    },
    drawPreview() {
      const el = this.$refs.preview
      if (!el) return
      el.innerHTML = this.previewHTML
      drawButtons(el, vendorOf(), undefined, true)
    },
    finish(value) {
      this.show = false
      if (this.done) this.done(value && { align: value.align, items: value.items.filter((it) => it.href).map(({ hint, ...it }) => it) })
    },
  },
}

const newTab = (n) => ({ title: n ? 'Tab ' + n : '', color: '', iconLeft: '', iconRight: '', body: '' })

export const TabsDialog = {
  name: 'MdkitTabsDialog',
  mixins: [wordsMixin],
  props: {
    // editor is the editor of a tab's content (the text's own), its
    // v-model the content; editorProps, its props
    editor: { type: [Object, String], default: null },
    editorProps: { type: Object, default: () => ({}) },
  },
  data: () => ({ show: false, align: 'start', items: [], active: 0, there: false, done: null, seq: 0 }),
  computed: {
    it() { return this.items[this.active] || null },
  },
  template: `
<v-dialog v-model="show" max-width="1100" scrollable data-tabs-dialog>
  <v-card>
    <v-card-title class="d-flex align-center ga-2"><v-icon icon="mdi-tab"></v-icon>{{ w.tabs }}</v-card-title>
    <v-divider></v-divider>
    <v-card-text style="max-height:78vh">
      ${alignField}
      <div class="d-flex align-center mt-3">
        <v-tabs v-model="active" show-arrows density="compact" style="flex:1 1 auto;min-width:0" data-tabs-bar>
          <v-tab v-for="(t, i) in items" :key="i" :value="i" class="text-none" :data-tab="i">
            <v-icon v-if="validIcon(t.iconLeft)" :icon="t.iconLeft" class="me-1" size="small"></v-icon>
            <span class="mdkit-swatch" v-if="t.color" :data-color="t.color"></span>{{ t.title || (w.tab + ' ' + (i + 1)) }}
            <v-icon v-if="validIcon(t.iconRight)" :icon="t.iconRight" class="ms-1" size="small"></v-icon>
          </v-tab>
        </v-tabs>
        <v-btn variant="tonal" size="small" prepend-icon="mdi-plus" @click="add" data-tabs-add>{{ w.addTab }}</v-btn>
      </div>
      <div v-if="it" :key="seq + ':' + active" class="mt-3" data-tabs-item>
        <div class="d-flex align-center ga-2">
          <v-text-field v-model="it.title" :label="w.title" density="compact" variant="outlined" hide-details style="flex:1 1 auto" data-tabs-title></v-text-field>
          <v-btn icon="mdi-arrow-left" size="small" variant="text" :disabled="active === 0" :title="w.left" @click="shift(-1)"></v-btn>
          <v-btn icon="mdi-arrow-right" size="small" variant="text" :disabled="active === items.length - 1" :title="w.right" @click="shift(1)"></v-btn>
          <v-btn icon="mdi-delete-outline" size="small" variant="text" color="error" :disabled="items.length < 2" :title="w.remove" @click="removeTab" data-tabs-remove></v-btn>
        </div>
        ${lookFields}
        <div class="text-caption text-medium-emphasis mt-4 mb-1">{{ w.content }}</div>
        <component v-if="editor" :is="editor" v-bind="editorProps" v-model="it.body" data-tabs-editor></component>
        <v-textarea v-else v-model="it.body" auto-grow rows="6" variant="outlined" hide-details></v-textarea>
      </div>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-actions>
      <v-btn v-if="there" color="error" variant="text" @click="finish(null)" data-widget-delete>{{ w.removeAll }}</v-btn>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ w.cancel }}</v-btn>
      <v-btn color="primary" variant="flat" :disabled="!items.length" @click="finish(value())" data-widget-apply>{{ there ? w.update : w.insert }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>`,
  methods: {
    ...methods,
    // open shows the tabs of value (none: two to make), done given the value
    // made or null
    open(value, done) {
      swatchStyle()
      this.align = (value && value.align) || 'start'
      const items = value && value.items && value.items.length ? value.items : null
      this.items = items ? items.map((it) => ({ ...newTab(0), ...it })) : [newTab(1), newTab(2)]
      if (!items) this.items.forEach((t, i) => { t.title = (this.w.tab || 'Tab') + ' ' + (i + 1) })
      this.there = !!items
      this.active = 0
      this.seq++
      this.done = done
      this.show = true
    },
    value() { return { align: this.align, items: this.items.map((it) => this.clean(it)) } },
    add() {
      const t = newTab(0)
      t.title = (this.w.tab || 'Tab') + ' ' + (this.items.length + 1)
      this.items.push(t)
      this.active = this.items.length - 1
    },
    shift(d) {
      this.move(this.items, this.active, d)
      this.active += d
      this.seq++
    },
    removeTab() {
      if (this.items.length < 2) return
      this.items.splice(this.active, 1)
      this.active = Math.max(0, this.active - 1)
      this.seq++
    },
    finish(value) {
      this.show = false
      if (this.done) this.done(value)
    },
  },
}
