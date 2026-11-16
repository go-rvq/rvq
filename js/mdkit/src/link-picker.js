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
import { searchMixin } from './link-search.js'
import { validIcon } from './link-icon.js'

// hrefOf is the address a link of target has: a path of the site as it is; a
// record path, its address (resolve: the cards'); '', none.
export async function hrefOf(target, resolve = metaResolver()) {
  if (!target) return ''
  if (!target.startsWith('admin:')) return target
  if (!resolve) return ''
  const cards = await resolve([target]).catch(() => ({}))
  return (cards && cards[target] && cards[target].href) || ''
}

// externalHref says href is an address of another site (http or https, not
// of this one's origin).
export function externalHref(href, origin = (typeof location !== 'undefined' ? location.origin : '')) {
  const h = String(href || '').trim()
  if (!/^https?:\/\//i.test(h)) return false
  return !origin || !(h === origin || h.startsWith(origin + '/'))
}

// metaTitle is the reader of the title of a page of another site the
// page's meta mdkit-link-cards says (its data-title: the admin's
// LinkTitleEvent): given an address, its title, '' of none; null, no meta.
// It waits ms at most (the server waits less: one try).
export function metaTitle(doc = (typeof document !== 'undefined' ? document : null), ms = 5000) {
  const meta = doc && doc.querySelector('meta[name="mdkit-link-cards"]')
  const url = meta && meta.getAttribute('data-title')
  if (!url) return null
  return async (href) => {
    const ctl = typeof AbortController !== 'undefined' ? new AbortController() : null
    const timer = ctl && setTimeout(() => ctl.abort(), ms)
    try {
      const res = await fetch(url, { method: 'POST', credentials: 'same-origin', body: new URLSearchParams({ url: href }), signal: ctl && ctl.signal })
      if (!res.ok) return ''
      const out = await res.json().catch(() => ({}))
      const data = out && out.data !== undefined ? out.data : out
      return (data && typeof data.title === 'string' && data.title) || ''
    } catch (e) {
      return ''
    } finally {
      if (timer) clearTimeout(timer)
    }
  }
}

export const LinkForm = {
  name: 'MdkitLinkForm',
  mixins: [searchMixin],
  props: {
    href: { type: String, default: '' },
    target: { type: String, default: '' },
    // targets: whether it may open in a new tab (the rich text's)
    targets: { type: Boolean, default: true },
    // linked: a link is there (it may be taken out)
    linked: { type: Boolean, default: false },
    // text: the text of the link (the one chosen); texts: whether it is
    // asked — empty, a page's or a post's chosen is its title, else the
    // address
    text: { type: String, default: '' },
    texts: { type: Boolean, default: true },
    words: { type: Object, default: () => wordsOf('link') },
    search: { type: Function, default: null },
    // titleOf: the title of a page of another site (a promise); by default
    // the page's meta's (metaTitle)
    titleOf: { type: Function, default: null },
    // icon: the icon of the link there; icons: whether it is asked (a text
    // link's: the icon before it) — a file chosen gives its own
    icon: { type: String, default: '' },
    icons: { type: Boolean, default: false },
  },
  emits: ['apply', 'remove', 'cancel'],
  data() {
    return { url: this.href || '', blank: this.target === '_blank', busy: false, title: '', label: this.text || '', iconName: this.icon || '', picked: '', lastPicked: '' }
  },
  computed: {
    finder() { return this.search || config.cardSearch || metaSearch() },
  },
  // the facets of the search, asked at once
  mounted() { this.resetSearch() },
  watch: {
    href(v) { this.url = v || '' },
    text(v) { this.label = v || '' },
    target(v) { this.blank = v === '_blank' },
  },
  template: `
<div data-link-form>
  <v-text-field v-model="url" autofocus density="compact" variant="outlined" hide-details prepend-inner-icon="mdi-link-variant"
    :label="words.address || 'Address'" placeholder="https://… /…" @keydown.enter.prevent="apply" data-link-href></v-text-field>
  <v-text-field v-if="texts" v-model="label" density="compact" variant="outlined" hide-details class="mt-3"
    prepend-inner-icon="mdi-format-text" :label="words.text || 'Link text'"
    :placeholder="title || words.textHint || 'Empty: the title of the page or post'" persistent-placeholder
    @keydown.enter.prevent="apply" data-link-text></v-text-field>
  <v-text-field v-if="icons" v-model.trim="iconName" density="compact" variant="outlined" hide-details clearable class="mt-3"
    :prepend-inner-icon="validIcon(iconName) || 'mdi-shape-outline'" :label="words.icon || 'Icon'" placeholder="mdi-…"
    data-link-icon></v-text-field>
  <v-text-field v-if="finder" v-model="q" density="compact" variant="outlined" hide-details clearable class="mt-3"
    prepend-inner-icon="mdi-magnify" :loading="loading" :label="words.search || 'Search a page or a record'"
    @update:model-value="runSearch()" data-link-q></v-text-field>
  <mdkit-link-facets v-if="finder" :facets="facets" :words="words" v-model="values"></mdkit-link-facets>
  <mdkit-link-results v-if="groups.length" :groups="groups" :words="words" @pick="pick"></mdkit-link-results>
  <div v-else-if="searched && !loading" class="text-body-2 text-medium-emphasis mt-2" data-link-none>
    {{ words.none || 'Nothing found.' }}
  </div>
  <v-switch v-if="targets" v-model="blank" color="primary" density="compact" hide-details class="mt-2"
    :label="words.newTab || 'Open in a new tab'" data-link-blank></v-switch>
  <div class="d-flex align-center ga-2 mt-3">
    <v-btn v-if="linked" variant="text" color="error" prepend-icon="mdi-link-variant-off" @click="$emit('remove')" data-link-remove>
      {{ words.remove || 'Remove the link' }}</v-btn>
    <v-spacer></v-spacer>
    <v-btn variant="text" @click="$emit('cancel')">{{ words.cancel || 'Cancel' }}</v-btn>
    <v-btn color="primary" variant="flat" :disabled="!url.trim() || busy" :loading="busy" @click="apply" data-link-apply>{{ words.apply || 'Apply' }}</v-btn>
  </div>
</div>`,
  methods: {
    // pick makes the address the option's (a record's: its address in the
    // admin)
    async pick(it) {
      this.busy = true
      this.url = (await hrefOf(it.target)) || this.url
      this.busy = false
      this.groups = []
      this.q = ''
      this.title = it.title
      // its icon (a file's: of its extension), in place of none or of the
      // one a choice before put — not of one written
      this.picked = validIcon(it.icon)
      if (this.picked && (!this.iconName || this.iconName === this.lastPicked)) this.iconName = this.picked
      this.lastPicked = this.picked
    },
    validIcon,
    async apply() {
      const href = this.url.trim()
      if (!href || this.busy) return
      // the text: the one typed; empty, the title of the page or post chosen
      // — of an address of another site, its page's title (asked of the
      // server: once, a short wait) —, else the address
      let title = this.title || ''
      if (!(this.label || '').trim() && !title && externalHref(href)) {
        const ask = this.titleOf || config.linkTitle || metaTitle()
        if (ask) {
          this.busy = true
          title = (await ask(href).catch(() => '')) || ''
          this.busy = false
        }
      }
      const text = (this.label || '').trim() || title || href
      this.$emit('apply', { href, target: this.targets && this.blank ? '_blank' : null, title, text,
        icon: this.icons ? validIcon(this.iconName) || null : null, pickedIcon: this.picked || '' })
    },
  },
}

export default {
  name: 'MdkitLinkPicker',
  components: { MdkitLinkForm: LinkForm },
  props: {
    words: { type: Object, default: () => wordsOf('link') },
    targets: { type: Boolean, default: true },
    // icons: the icon of the link asked (a text link's)
    icons: { type: Boolean, default: false },
  },
  data: () => ({ show: false, href: '', target: '', text: '', icon: '', linked: false, done: null, key: 0 }),
  template: `
<v-dialog v-model="show" max-width="560" scrollable data-link-picker>
  <v-card>
    <v-card-title class="d-flex align-center ga-2"><v-icon icon="mdi-link-variant"></v-icon>{{ words.title || 'Link' }}</v-card-title>
    <v-card-text>
      <mdkit-link-form :key="key" :href="href" :target="target" :text="text" :targets="targets" :linked="linked" :words="words"
        :icon="icon" :icons="icons" @apply="finish" @remove="finish(null)" @cancel="show = false"></mdkit-link-form>
    </v-card-text>
  </v-card>
</v-dialog>`,
  methods: {
    // open shows the picker with the link there (attrs: href, target), done
    // given the one chosen ({href, target, title}), or null: taken out
    open(attrs, done) {
      this.href = (attrs && attrs.href) || ''
      this.target = (attrs && attrs.target) || ''
      this.text = (attrs && attrs.text) || ''
      this.icon = (attrs && attrs.icon) || ''
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
