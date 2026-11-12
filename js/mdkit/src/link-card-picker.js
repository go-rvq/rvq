// mdkit-link-card-picker: what a link card links, chosen — a Vue component,
// a dialog (open(done)) —: the pages of the site and, in the admin, its
// records, found as one types (the page's search: the meta
// mdkit-link-cards' data-search, or `search`), or an address typed (a path of
// the site, /…; a record path, admin:…:). Chosen, its target given to done.
import { config, wordsOf } from './config.js'
import { cardTarget } from './link-cards.js'
import { searchMixin } from './link-search.js'

// queryLocale is the locale of the address of the page (?locale=…); '', none.
function queryLocale() {
  try {
    return (typeof location !== 'undefined' && new URLSearchParams(location.search).get('locale')) || ''
  } catch (e) {
    return ''
  }
}

// metaSearch is the search of the page's meta (mdkit-link-cards
// data-search: a POST of q answering the groups of options — {data},
// {groups}, or them —; given the values of the facets, {key: [values]},
// they are posted too (f_<key>, an empty one f_<key>="") and the facets
// asked: {groups, facets}); null, none.
export function metaSearch(doc = (typeof document !== 'undefined' ? document : null)) {
  const meta = doc && doc.querySelector('meta[name="mdkit-link-cards"]')
  const url = meta && meta.getAttribute('data-search')
  if (!url) return null
  // the language of the text: the meta's, else the page's (its locale of
  // the query: the admin's listing); none, the server's
  const locale = meta.getAttribute('data-locale') || queryLocale()
  return async (q, values) => {
    const body = new URLSearchParams({ q })
    if (locale) body.append('locale', locale)
    if (values) {
      body.append('facets', '1')
      for (const [k, vs] of Object.entries(values)) {
        if (!vs.length) body.append('f_' + k, '')
        for (const v of vs) body.append('f_' + k, v)
      }
    }
    const res = await fetch(url, { method: 'POST', credentials: 'same-origin', body })
    if (!res.ok) return values ? { groups: [], facets: null } : []
    const out = await res.json().catch(() => ({}))
    const data = out && out.data !== undefined ? out.data : out
    if (values && data && !Array.isArray(data) && Array.isArray(data.groups)) return data
    return (data && (Array.isArray(data) ? data : data.groups)) || []
  }
}

export default {
  name: 'MdkitLinkCardPicker',
  mixins: [searchMixin],
  props: {
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('card') },
    // search: the groups of options matching q — and the values of the
    // facets: {groups, facets} — (a promise); by default the page's meta's
    search: { type: Function, default: null },
  },
  data: () => ({ show: false, done: null }),
  computed: {
    // typed is the address typed, when it is one a card may link
    typed() { return cardTarget((this.q || '').trim()) },
    finder() { return this.search || config.cardSearch || metaSearch() },
  },
  template: `
<v-dialog v-model="show" max-width="620" scrollable data-link-card-picker>
  <v-card>
    <v-card-title class="d-flex align-center ga-2">
      <v-icon icon="mdi-card-text-outline"></v-icon>{{ words.title || 'Link card' }}
    </v-card-title>
    <v-card-text>
      <v-text-field v-model="q" autofocus density="compact" variant="outlined" hide-details clearable
        prepend-inner-icon="mdi-magnify" :loading="loading"
        :placeholder="words.search || 'Search, or type an address of the site (/…)'"
        @update:model-value="runSearch()" @keydown.enter.prevent="enter" data-link-card-q></v-text-field>
      <mdkit-link-facets v-if="!typed" :facets="facets" v-model="values"></mdkit-link-facets>
      <v-list v-if="typed || groups.length" density="compact" class="mt-2" data-link-card-options>
        <v-list-item v-if="typed" prepend-icon="mdi-link-variant" :title="(words.useAddress || 'Use the address') + ' ' + typed"
          @click="pick(typed)" data-link-card-typed></v-list-item>
        <mdkit-link-options :groups="groups" :words="words" @pick="(it) => pick(it.target)"></mdkit-link-options>
      </v-list>
      <div v-else-if="searched && !loading" class="text-body-2 text-medium-emphasis mt-3" data-link-card-none>
        {{ words.none || 'Nothing found.' }}
      </div>
      <div class="text-caption text-medium-emphasis mt-3">{{ words.hint || 'The card shows the title, the cover and the summary, in the language of the page.' }}</div>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="text" @click="show = false">{{ words.cancel || 'Cancel' }}</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>`,
  methods: {
    // open shows the picker, done given the target chosen
    open(done) {
      this.done = done
      this.show = true
      this.resetSearch()
    },
    // skip: an address typed is not searched
    skip() { return !!this.typed },
    // enter takes the address typed, else the first option
    enter() {
      if (this.typed) return this.pick(this.typed)
      const first = this.groups[0] && this.groups[0].items[0]
      if (first) this.pick(first.target)
    },
    pick(target) {
      this.show = false
      if (this.done) this.done(target)
    },
  },
}
