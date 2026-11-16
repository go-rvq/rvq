// The search of what a link or a card links (the link form, the card
// picker): what is typed and the filters of the page's search (its facets),
// asked a moment after the last change, its groups of options shown.
//   - a search (the page's meta's: metaSearch, or one given) is given q and
//     the values of the facets ({key: [values]}), and answers {groups,
//     facets} — or the groups alone (no facets) —;
//   - a facet: {key, label, kind — "chips" (all its items, chosen by a
//     click) or "autocomplete" (found as one types, chosen as chips) —,
//     items [{value, title, group}], default (the values at first)}; the
//     facets of each answer are those the values chosen call for (posts
//     chosen: their types, their tags), a value no longer an item's left out.

// DEBOUNCE is how long a search waits after the last key (ms).
export const DEBOUNCE = 300

// resultOf is the answer of a search as {groups, facets}: groups alone,
// facets null.
export function resultOf(out) {
  if (Array.isArray(out)) return { groups: out, facets: null }
  if (out && Array.isArray(out.groups)) return { groups: out.groups, facets: Array.isArray(out.facets) ? out.facets : null }
  return { groups: [], facets: null }
}

// facetValues are the values of the facets: the ones chosen, else their
// defaults.
export function facetValues(facets, values) {
  const out = {}
  for (const f of facets || []) out[f.key] = values[f.key] !== undefined ? values[f.key] : [...(f.default || [])]
  return out
}

// pruneValues are the values of the facets that are their items' (of a
// facet there); and whether any was left out.
export function pruneValues(facets, values) {
  const out = {}
  let pruned = false
  for (const f of facets || []) {
    const v = values[f.key]
    if (v === undefined) continue
    const items = new Set((f.items || []).map((it) => it.value))
    out[f.key] = v.filter((x) => items.has(x))
    if (out[f.key].length !== v.length) pruned = true
  }
  for (const k of Object.keys(values)) if (!(k in out) && (facets || []).every((f) => f.key !== k)) pruned = true
  return { values: out, pruned }
}

// itemsOf are the items of a facet, a header before each group of them
// ({header: true, title}).
export function itemsOf(f) {
  const out = []
  let group = null
  for (const it of f.items || []) {
    if (it.group && it.group !== group) out.push({ header: true, title: it.group, value: '\u0000' + it.group })
    group = it.group || null
    out.push(it)
  }
  return out
}

// LinkFacets are the filters of a search: a facet of chips, its chips; one of
// a list (autocomplete), its values chosen as chips and a button that opens
// its dialog — as the selector of the record of the menu: its items found as
// one types (a moment after the last key), each checked, in their groups.
export const LinkFacets = {
  name: 'MdkitLinkFacets',
  props: {
    facets: { type: Array, default: () => [] },
    // modelValue: the values of the facets, {key: [values]}
    modelValue: { type: Object, default: () => ({}) },
    words: { type: Object, default: () => ({}) },
  },
  emits: ['update:modelValue'],
  data: () => ({ dialog: false, current: null, chosen: [], dq: '', shownQ: '', timer: null }),
  computed: {
    // shown are the items of the facet of the dialog matching its search,
    // their groups' headers before them
    shown() {
      if (!this.current) return []
      const q = this.shownQ.trim().toLowerCase()
      const items = (this.current.items || []).filter((it) => !q || String(it.title).toLowerCase().includes(q))
      return itemsOf({ items })
    },
  },
  watch: {
    // a moment after the last key
    dq(v) {
      clearTimeout(this.timer)
      this.timer = setTimeout(() => (this.shownQ = v || ''), DEBOUNCE)
    },
  },
  template: `
<div v-if="facets.length" class="d-flex flex-column ga-2 mt-3" data-link-facets>
  <div v-for="f in facets" :key="f.key" class="d-flex align-center ga-2 flex-wrap" :data-link-facet="f.key">
    <span class="text-body-2 text-medium-emphasis">{{ f.label }}</span>
    <v-chip-group v-if="f.kind === 'chips'" :model-value="valueOf(f)" @update:model-value="set(f, $event)" multiple column selected-class="text-primary">
      <v-chip v-for="it in f.items" :key="it.value" :value="it.value" size="small" variant="outlined" filter
        :data-link-facet-item="it.value">{{ it.title }}</v-chip>
    </v-chip-group>
    <template v-else>
      <v-chip v-for="v in valueOf(f)" :key="v" size="small" color="primary" variant="tonal" closable
        @click:close="set(f, valueOf(f).filter((x) => x !== v))" :data-link-facet-chosen="v">{{ titleOf(f, v) }}</v-chip>
      <v-btn size="small" variant="tonal" prepend-icon="mdi-filter-variant" @click="open(f)" :data-link-facet-open="f.key">
        {{ words.facetChoose || 'Choose' }}</v-btn>
    </template>
  </div>
  <v-dialog v-model="dialog" max-width="520" scrollable data-link-facet-dialog>
    <v-card v-if="current">
      <v-card-title>{{ current.label }}</v-card-title>
      <v-card-text>
        <v-text-field v-model="dq" autofocus clearable density="compact" variant="outlined" hide-details class="mb-2"
          prepend-inner-icon="mdi-magnify" :placeholder="words.facetSearch || 'Search'" data-link-facet-q></v-text-field>
        <v-list density="compact" data-link-facet-items>
          <template v-for="it in shown" :key="it.value">
            <v-list-subheader v-if="it.header" :data-link-facet-group="it.title">{{ it.title }}</v-list-subheader>
            <v-list-item v-else @click="toggle(it.value)" :data-link-facet-item="it.value">
              <template #prepend><v-checkbox-btn :model-value="chosen.includes(it.value)" density="compact"></v-checkbox-btn></template>
              <v-list-item-title>{{ it.title }}</v-list-item-title>
            </v-list-item>
          </template>
          <div v-if="!shown.length" class="text-body-2 text-medium-emphasis pa-2">{{ words.none || 'Nothing found.' }}</div>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-btn variant="text" @click="chosen = []" data-link-facet-clear>{{ words.facetClear || 'Clear' }}</v-btn>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="dialog = false">{{ words.cancel || 'Cancel' }}</v-btn>
        <v-btn color="primary" variant="flat" @click="apply" data-link-facet-apply>{{ words.apply || 'Apply' }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</div>`,
  methods: {
    items: itemsOf,
    valueOf(f) {
      const v = this.modelValue[f.key]
      return v !== undefined ? v : [...(f.default || [])]
    },
    titleOf(f, v) {
      const it = (f.items || []).find((x) => x.value === v)
      return it ? it.title : v
    },
    set(f, v) {
      this.$emit('update:modelValue', { ...this.modelValue, [f.key]: (v || []).filter((x) => !String(x).startsWith('\u0000')) })
    },
    open(f) {
      this.current = f
      this.chosen = [...this.valueOf(f)]
      this.dq = ''
      this.shownQ = ''
      this.dialog = true
    },
    toggle(v) {
      this.chosen = this.chosen.includes(v) ? this.chosen.filter((x) => x !== v) : [...this.chosen, v]
    },
    apply() {
      this.set(this.current, this.chosen)
      this.dialog = false
    },
  },
}

// flatten are the options of a tree, by their targets.
export function flatten(items, out = {}) {
  for (const it of items || []) {
    out[it.target] = it
    flatten(it.children, out)
  }
  return out
}

// isTree says a group is a tree: an option with children (or its kind).
export const isTree = (g) => g.kind === 'tree' || (g.items || []).some((it) => it.children && it.children.length)

// LinkOptions are the options found, in their groups — a group a tree (a
// page under its parents, a post under its parents and the page of its
// type: v-treeview, open; one there only to hold others dimmed) or a list
// —; pick, the option chosen.
export const LinkOptions = {
  name: 'MdkitLinkOptions',
  props: {
    groups: { type: Array, default: () => [] },
    words: { type: Object, default: () => ({}) },
    // headers: each group's header (none: in a tab of its own)
    headers: { type: Boolean, default: true },
  },
  emits: ['pick'],
  computed: {
    byTarget() {
      const out = {}
      for (const g of this.groups) flatten(g.items, out)
      return out
    },
  },
  template: `
<div data-link-options>
  <template v-for="(g, i) in groups" :key="i">
    <div v-if="headers && (g.title || words[g.kind])" class="d-flex align-center ga-2 px-4 pt-3 pb-1 text-subtitle-2 text-medium-emphasis"
      :class="{ 'border-t': i > 0 }" :data-link-group="g.kind || g.title">
      <v-icon v-if="g.icon" :icon="g.icon" size="small"></v-icon>{{ g.title || words[g.kind] }}</div>
    <v-treeview v-if="isTree(g)" :items="g.items" item-title="title" item-value="target" item-children="children"
      open-all activatable active-strategy="single-independent" density="compact" :active="[]"
      @update:activated="(v) => v && v.length && pick(v[0])" data-link-tree>
      <template #prepend="{ item }"><v-icon :icon="item.icon || 'mdi-file-document-outline'" size="small"
        :class="{ 'text-medium-emphasis': item.context }"></v-icon></template>
      <template #title="{ item }">
        <div :class="{ 'text-medium-emphasis': item.context }" :data-link-option="item.target">
          <div class="text-body-2">{{ item.title }}</div>
          <div class="text-caption text-medium-emphasis text-truncate">{{ item.target }}</div>
          <div v-if="item.info && item.info.length" class="d-flex flex-wrap ga-3 text-caption text-medium-emphasis"><span v-for="(d, j) in (item.info || [])" :key="j" class="mdkit-link-info d-inline-flex align-center ga-1" :title="d.hint"
            data-link-info><v-icon :icon="d.icon || 'mdi-calendar'" size="x-small"></v-icon>{{ d.text }}</span></div>
        </div>
      </template>
    </v-treeview>
    <template v-else>
      <v-list-item v-for="it in g.items" :key="it.target" :title="it.title"
        :prepend-icon="it.icon || g.icon || 'mdi-file-document-outline'" @click="pick(it.target)" :data-link-option="it.target">
        <template #subtitle><div class="text-truncate">{{ it.sub && it.sub !== g.title ? it.sub : it.target }}</div>
          <div v-if="it.info && it.info.length" class="d-flex flex-wrap ga-3"><span v-for="(d, j) in (it.info || [])" :key="j" class="mdkit-link-info d-inline-flex align-center ga-1" :title="d.hint"
            data-link-info><v-icon :icon="d.icon || 'mdi-calendar'" size="x-small"></v-icon>{{ d.text }}</span></div></template>
      </v-list-item>
    </template>
  </template>
</div>`,
  methods: {
    isTree,
    pick(target) {
      const it = this.byTarget[target]
      if (it) this.$emit('pick', it)
    },
  },
}

// foundOf is how many options of a group were found (not those only holding
// others): its count, else its options'.
export function foundOf(g) {
  if (g && g.count) return g.count
  let n = 0
  const walk = (items) => {
    for (const it of items || []) {
      if (!it.context) n++
      walk(it.children)
    }
  }
  walk(g && g.items)
  return n
}

// countText is what the chip of a group's tab says: how many, "+" when
// more were found than shown.
export const countText = (g) => String(foundOf(g)) + (g && g.more ? '+' : '')

// tabOf is the key of a group's tab.
export const tabOf = (g) => g.kind || g.title

// LinkResults are the options found in a card — titled Recent when nothing
// was searched (the groups all recent) —, each group in a tab (its icon,
// its name, a chip of how many it found); its slot "before", above the tabs
// (an address typed).
export const LinkResults = {
  name: 'MdkitLinkResults',
  components: { MdkitLinkOptions: LinkOptions },
  props: {
    groups: { type: Array, default: () => [] },
    words: { type: Object, default: () => ({}) },
  },
  emits: ['pick'],
  data: () => ({ tab: null }),
  computed: {
    recent() { return this.groups.length > 0 && this.groups.every((g) => g.recent) },
    current() { return this.groups.find((g) => tabOf(g) === this.tab) || this.groups[0] || null },
  },
  watch: {
    // the tab chosen kept while its group is there; else the first
    groups: {
      immediate: true,
      handler(gs) {
        if (!gs.some((g) => tabOf(g) === this.tab)) this.tab = gs.length ? tabOf(gs[0]) : null
      },
    },
  },
  template: `
<v-card variant="outlined" class="mt-3" data-link-results :data-link-recent="recent ? '' : null">
  <v-card-title class="d-flex align-center ga-2 text-subtitle-1 py-2">
    <v-icon :icon="recent ? 'mdi-history' : 'mdi-magnify'" size="small"></v-icon>
    {{ recent ? (words.recent || 'Recent') : (words.results || 'Search results') }}</v-card-title>
  <v-divider></v-divider>
  <slot name="before"></slot>
  <template v-if="groups.length">
    <v-tabs v-model="tab" density="compact" color="primary" show-arrows data-link-tabs>
      <v-tab v-for="g in groups" :key="tabOf(g)" :value="tabOf(g)" class="text-none" :data-link-tab="tabOf(g)">
        <v-icon v-if="g.icon" :icon="g.icon" size="small" class="me-1"></v-icon>{{ g.title || words[g.kind] || g.kind }}
        <v-chip size="x-small" class="ms-2" :color="tabOf(g) === tab ? 'primary' : undefined" variant="tonal"
          :data-link-count="tabOf(g)">{{ countText(g) }}</v-chip>
      </v-tab>
    </v-tabs>
    <v-divider></v-divider>
    <div style="max-height:20rem;overflow:auto">
      <mdkit-link-options v-if="current" :key="tabOf(current)" :groups="[current]" :headers="false" :words="words"
        @pick="(it) => $emit('pick', it)"></mdkit-link-options>
    </div>
  </template>
</v-card>`,
  methods: { tabOf, countText },
}

// searchMixin is the search of a picker: its finder (a computed of the
// component), q, its facets and their values, its groups (without a text,
// those the server answers: the latest); runSearch() asks
// again, a moment after (now: at once) — not search(): the pickers' prop of
// their search is called so —; skip(q) says a text is not searched (an
// address typed).
export const searchMixin = {
  components: { MdkitLinkFacets: LinkFacets, MdkitLinkOptions: LinkOptions, MdkitLinkResults: LinkResults },
  data: () => ({ q: '', groups: [], loading: false, seq: 0, timer: null, facets: [], values: {}, faceted: false }),
  computed: {
    // filtered: a facet with a value other than its default (a search
    // without a text finds by them)
    filtered() {
      return this.facets.some((f) => {
        const v = this.values[f.key]
        if (v === undefined) return false
        const d = f.default || []
        return v.length !== d.length || v.some((x) => !d.includes(x))
      })
    },
    searched() { return (this.q || '').trim().length >= 2 || this.filtered },
  },
  watch: {
    values: { deep: true, handler() { this.runSearch() } },
  },
  methods: {
    // resetSearch empties the search, its facets asked again
    resetSearch() {
      this.q = ''
      this.groups = []
      this.values = {}
      this.facets = []
      this.faceted = false
      this.runSearch(true)
    },
    skip() { return false },
    runSearch(now = false) {
      clearTimeout(this.timer)
      const q = (this.q || '').trim()
      if (!this.finder || this.skip(q) || (!this.faceted && q.length < 2 && !now)) {
        this.groups = []
        this.loading = false
        return
      }
      const n = ++this.seq
      this.loading = true
      const run = async () => {
        let res = { groups: [], facets: null }
        try {
          res = resultOf(await this.finder(q, facetValues(this.facets, this.values)))
        } catch (e) {
          res = { groups: [], facets: null }
        }
        if (n !== this.seq) return
        if (res.facets) {
          this.faceted = true
          this.facets = res.facets
          const { values, pruned } = pruneValues(res.facets, this.values)
          if (pruned) {
            // a value no longer there: asked again without it (the watch)
            this.values = values
            return
          }
        }
        // what the server answers: without a text, the latest (recent)
        this.groups = res.groups.filter((g) => g && g.items && g.items.length)
        this.loading = false
      }
      if (now) run()
      else this.timer = setTimeout(run, DEBOUNCE)
    },
  },
}
