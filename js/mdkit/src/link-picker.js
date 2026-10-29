// The link of a text, chosen — an address typed, or a page of the site and,
// in the admin, a record, found as one types (the search of the link cards:
// the page's meta mdkit-link-cards' data-search) —, whether it opens in a new
// tab (`targets`: the rich text's), or the link taken out.
//   - mdkit-link-form: the form (in a dialog of the page's: the admin's
//     TipTap's), its events apply ({href, target}), remove, cancel;
//   - mdkit-link-picker: it in a dialog of its own (open(attrs, done): done
//     given {href, target}, or null to take it out) — the editors of mdkit.
// A record (admin:…:) chosen is its address in the admin (the cards' answer:
// its href).
import { config, wordsOf } from './config.js'
import { metaResolver } from './link-cards.js'
import { metaSearch } from './link-card-picker.js'

// hrefOf is the address a link of target has: a path of the site as it is; a
// record path, its address (resolve: the cards'); '', none.
export async function hrefOf(target, resolve = metaResolver()) {
  if (!target) return ''
  if (!target.startsWith('admin:')) return target
  if (!resolve) return ''
  const cards = await resolve([target]).catch(() => ({}))
  return (cards && cards[target] && cards[target].href) || ''
}

export const LinkForm = {
  name: 'MdkitLinkForm',
  props: {
    href: { type: String, default: '' },
    target: { type: String, default: '' },
    // targets: whether it may open in a new tab (the rich text's)
    targets: { type: Boolean, default: true },
    // linked: a link is there (it may be taken out)
    linked: { type: Boolean, default: false },
    words: { type: Object, default: () => wordsOf('link') },
    search: { type: Function, default: null },
  },
  emits: ['apply', 'remove', 'cancel'],
  data() {
    return { url: this.href || '', blank: this.target === '_blank', q: '', groups: [], loading: false, seq: 0, timer: null, busy: false, title: '' }
  },
  computed: {
    finder() { return this.search || config.cardSearch || metaSearch() },
  },
  watch: {
    href(v) { this.url = v || '' },
    target(v) { this.blank = v === '_blank' },
  },
  template: `
<div data-link-form>
  <v-text-field v-model="url" autofocus density="compact" variant="outlined" hide-details prepend-inner-icon="mdi-link-variant"
    :label="words.address || 'Address'" placeholder="https://… /…" @keydown.enter.prevent="apply" data-link-href></v-text-field>
  <v-text-field v-if="finder" v-model="q" density="compact" variant="outlined" hide-details clearable class="mt-3"
    prepend-inner-icon="mdi-magnify" :loading="loading" :label="words.search || 'Search a page or a record'"
    @update:model-value="typed" data-link-q></v-text-field>
  <v-list v-if="groups.length" density="compact" class="mt-1" style="max-height:16rem;overflow:auto" data-link-options>
    <template v-for="(g, i) in groups" :key="i">
      <v-list-subheader v-if="words[g.kind] || g.title">{{ words[g.kind] || g.title }}</v-list-subheader>
      <v-list-item v-for="it in g.items" :key="it.target" :title="it.title" :subtitle="it.target"
        :prepend-icon="g.icon || 'mdi-file-document-outline'" @click="pick(it)" :data-link-option="it.target"></v-list-item>
    </template>
  </v-list>
  <div v-else-if="(q || '').trim().length >= 2 && !loading" class="text-body-2 text-medium-emphasis mt-2" data-link-none>
    {{ words.none || 'Nothing found.' }}
  </div>
  <v-switch v-if="targets" v-model="blank" color="primary" density="compact" hide-details class="mt-2"
    :label="words.newTab || 'Open in a new tab'" data-link-blank></v-switch>
  <div class="d-flex align-center ga-2 mt-3">
    <v-btn v-if="linked" variant="text" color="error" prepend-icon="mdi-link-variant-off" @click="$emit('remove')" data-link-remove>
      {{ words.remove || 'Remove the link' }}</v-btn>
    <v-spacer></v-spacer>
    <v-btn variant="text" @click="$emit('cancel')">{{ words.cancel || 'Cancel' }}</v-btn>
    <v-btn color="primary" variant="flat" :disabled="!url.trim() || busy" @click="apply" data-link-apply>{{ words.apply || 'Apply' }}</v-btn>
  </div>
</div>`,
  methods: {
    // typed searches what is typed, a moment after the last key
    typed() {
      clearTimeout(this.timer)
      const q = (this.q || '').trim()
      if (q.length < 2 || !this.finder) {
        this.groups = []
        this.loading = false
        return
      }
      const n = ++this.seq
      this.loading = true
      this.timer = setTimeout(async () => {
        let groups = []
        try {
          groups = (await this.finder(q)) || []
        } catch (e) {
          groups = []
        }
        if (n !== this.seq) return
        this.groups = groups.filter((g) => g && g.items && g.items.length)
        this.loading = false
      }, 250)
    },
    // pick makes the address the option's (a record's: its address in the
    // admin)
    async pick(it) {
      this.busy = true
      this.url = (await hrefOf(it.target)) || this.url
      this.busy = false
      this.groups = []
      this.q = ''
      this.title = it.title
    },
    apply() {
      const href = this.url.trim()
      if (!href) return
      this.$emit('apply', { href, target: this.targets && this.blank ? '_blank' : null, title: this.title || '' })
    },
  },
}

export default {
  name: 'MdkitLinkPicker',
  components: { MdkitLinkForm: LinkForm },
  props: {
    words: { type: Object, default: () => wordsOf('link') },
    targets: { type: Boolean, default: true },
  },
  data: () => ({ show: false, href: '', target: '', linked: false, done: null, key: 0 }),
  template: `
<v-dialog v-model="show" max-width="560" scrollable data-link-picker>
  <v-card>
    <v-card-title class="d-flex align-center ga-2"><v-icon icon="mdi-link-variant"></v-icon>{{ words.title || 'Link' }}</v-card-title>
    <v-card-text>
      <mdkit-link-form :key="key" :href="href" :target="target" :targets="targets" :linked="linked" :words="words"
        @apply="finish" @remove="finish(null)" @cancel="show = false"></mdkit-link-form>
    </v-card-text>
  </v-card>
</v-dialog>`,
  methods: {
    // open shows the picker with the link there (attrs: href, target), done
    // given the one chosen ({href, target, title}), or null: taken out
    open(attrs, done) {
      this.href = (attrs && attrs.href) || ''
      this.target = (attrs && attrs.target) || ''
      this.linked = !!this.href
      this.done = done
      this.key++
      this.show = true
    },
    finish(link) {
      this.show = false
      if (this.done) this.done(link)
    },
  },
}
