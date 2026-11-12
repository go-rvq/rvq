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

export const LinkFacets = {
  name: 'MdkitLinkFacets',
  props: {
    facets: { type: Array, default: () => [] },
    // modelValue: the values of the facets, {key: [values]}
    modelValue: { type: Object, default: () => ({}) },
  },
  emits: ['update:modelValue'],
  template: `
<div v-if="facets.length" class="d-flex flex-column ga-2 mt-3" data-link-facets>
  <template v-for="f in facets" :key="f.key">
    <div v-if="f.kind === 'chips'" class="d-flex align-center ga-2 flex-wrap" :data-link-facet="f.key">
      <span class="text-body-2 text-medium-emphasis">{{ f.label }}</span>
      <v-chip-group :model-value="valueOf(f)" @update:model-value="set(f, $event)" multiple column selected-class="text-primary">
        <v-chip v-for="it in f.items" :key="it.value" :value="it.value" size="small" variant="outlined" filter
          :data-link-facet-item="it.value">{{ it.title }}</v-chip>
      </v-chip-group>
    </div>
    <v-autocomplete v-else :model-value="valueOf(f)" @update:model-value="set(f, $event)" :items="items(f)" item-title="title"
      item-value="value" :item-props="(it) => (it.header ? { disabled: true } : {})" :label="f.label" multiple chips
      closable-chips clearable density="compact" variant="outlined" hide-details :data-link-facet="f.key">
      <template #item="{ props, item }">
        <v-list-subheader v-if="item.raw.header" :data-link-facet-group="item.raw.title">{{ item.raw.title }}</v-list-subheader>
        <v-list-item v-else v-bind="props" :data-link-facet-item="item.raw.value"></v-list-item>
      </template>
    </v-autocomplete>
  </template>
</div>`,
  methods: {
    items: itemsOf,
    valueOf(f) {
      const v = this.modelValue[f.key]
      return v !== undefined ? v : [...(f.default || [])]
    },
    set(f, v) {
      this.$emit('update:modelValue', { ...this.modelValue, [f.key]: (v || []).filter((x) => !String(x).startsWith('\u0000')) })
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
    <v-list-subheader v-if="words[g.kind] || g.title">{{ words[g.kind] || g.title }}</v-list-subheader>
    <v-treeview v-if="isTree(g)" :items="g.items" item-title="title" item-value="target" item-children="children"
      open-all activatable active-strategy="single-independent" density="compact" :active="[]"
      @update:activated="(v) => v && v.length && pick(v[0])" data-link-tree>
      <template #prepend="{ item }"><v-icon :icon="item.icon || 'mdi-file-document-outline'" size="small"
        :class="{ 'text-medium-emphasis': item.context }"></v-icon></template>
      <template #title="{ item }">
        <div :class="{ 'text-medium-emphasis': item.context }" :data-link-option="item.target">
          <div class="text-body-2">{{ item.title }}</div>
          <div class="text-caption text-medium-emphasis text-truncate">{{ item.target }}</div>
        </div>
      </template>
    </v-treeview>
    <template v-else>
      <v-list-item v-for="it in g.items" :key="it.target" :title="it.title" :subtitle="it.sub && it.sub !== g.title ? it.sub : it.target"
        :prepend-icon="it.icon || g.icon || 'mdi-file-document-outline'" @click="pick(it.target)" :data-link-option="it.target"></v-list-item>
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

// searchMixin is the search of a picker: its finder (a computed of the
// component), q, its facets and their values, its groups; runSearch() asks
// again, a moment after (now: at once) — not search(): the pickers' prop of
// their search is called so —; skip(q) says a text is not searched (an
// address typed).
export const searchMixin = {
  components: { MdkitLinkFacets: LinkFacets, MdkitLinkOptions: LinkOptions },
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
        this.groups = (q.length >= 2 || this.filtered ? res.groups : []).filter((g) => g && g.items && g.items.length)
        this.loading = false
      }
      if (now) run()
      else this.timer = setTimeout(run, 250)
    },
  },
}
