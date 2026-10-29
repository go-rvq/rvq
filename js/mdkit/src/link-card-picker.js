// mdkit-link-card-picker: what a link card links, chosen — a Vue component,
// a dialog (open(done)) —: the pages of the site and, in the admin, its
// records, found as one types (the page's search: the meta
// mdkit-link-cards' data-search, or `search`), or an address typed (a path of
// the site, /…; a record path, admin:…:). Chosen, its target given to done.
import { config, wordsOf } from './config.js'
import { cardTarget } from './link-cards.js'

// metaSearch is the search of the page's meta (mdkit-link-cards
// data-search: a POST of q answering the groups of options — {data},
// {groups}, or them —); null, none.
export function metaSearch(doc = (typeof document !== 'undefined' ? document : null)) {
  const meta = doc && doc.querySelector('meta[name="mdkit-link-cards"]')
  const url = meta && meta.getAttribute('data-search')
  if (!url) return null
  const locale = meta.getAttribute('data-locale') || ''
  return async (q) => {
    const body = new URLSearchParams({ q })
    if (locale) body.append('locale', locale)
    const res = await fetch(url, { method: 'POST', credentials: 'same-origin', body })
    if (!res.ok) return []
    const out = await res.json().catch(() => ({}))
    return (out && (out.data || out.groups)) || (Array.isArray(out) ? out : [])
  }
}

export default {
  name: 'MdkitLinkCardPicker',
  props: {
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('card') },
    // search: the groups of options matching q (a promise); by default the
    // page's meta's
    search: { type: Function, default: null },
  },
  data: () => ({ show: false, q: '', groups: [], loading: false, done: null, seq: 0, timer: null }),
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
        @update:model-value="typedQ" @keydown.enter.prevent="enter" data-link-card-q></v-text-field>
      <v-list v-if="typed || groups.length" density="compact" class="mt-2" data-link-card-options>
        <v-list-item v-if="typed" prepend-icon="mdi-link-variant" :title="(words.useAddress || 'Use the address') + ' ' + typed"
          @click="pick(typed)" data-link-card-typed></v-list-item>
        <template v-for="(g, i) in groups" :key="i">
          <v-list-subheader v-if="words[g.kind] || g.title">{{ words[g.kind] || g.title }}</v-list-subheader>
          <v-list-item v-for="it in g.items" :key="it.target" :prepend-icon="g.icon || 'mdi-file-document-outline'"
            :title="it.title" :subtitle="it.sub && it.sub !== g.title ? it.sub : it.target" @click="pick(it.target)"
            :data-link-card-option="it.target"></v-list-item>
        </template>
      </v-list>
      <div v-else-if="(q || '').trim().length >= 2 && !loading" class="text-body-2 text-medium-emphasis mt-3" data-link-card-none>
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
      this.q = ''
      this.groups = []
      this.done = done
      this.show = true
    },
    // typedQ searches what is typed, a moment after the last key
    typedQ() {
      clearTimeout(this.timer)
      const q = (this.q || '').trim()
      if (q.length < 2 || !this.finder || this.typed) {
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
